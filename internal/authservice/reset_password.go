package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
	"github.com/phoenix-industries/phoenix-server/pkg/validate"
)

type resetPasswordData struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func (s *Service) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var data resetPasswordData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httputil.Error(nil, "invalid request body", http.StatusBadRequest).WriteJSON(w)
		return
	}
	defer r.Body.Close()
	if data.Password == "" {
		httputil.ErrorBadRequest().WriteJSON(w)
		return
	}
	if data.NewPassword == "" {
		httputil.ErrorBadRequest().WriteJSON(w)
		return
	}

	if err := validate.Password(data.NewPassword); err != nil {
		httputil.Error(nil, err.Error(), http.StatusBadRequest).WriteJSON(w)
		return
	}

	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		httputil.ErrorUnauthorized().WriteJSON(w)
		s.logger.ErrorContext(r.Context(), "error occured in reset password handler", "error", err, "userID", userID)
		return
	}

	ctx := r.Context()
	res := AuthResponse{}
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		user, err := models.UserGetByID(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return httputil.Error(nil, "invalid credentials", http.StatusUnauthorized)
		}

		if valid, err := s.auth.VerifyPassword(data.Password, user.Password); err != nil {
			return httputil.Error(err, "request failed", http.StatusInternalServerError)
		} else if !valid {
			return httputil.Error(nil, "invalid credentials", http.StatusUnauthorized)
		}

		hash, err := s.auth.HashPassword(data.NewPassword)
		if err != nil {
			return httputil.Error(err, "failed to hash password", http.StatusInternalServerError)
		}
		user.Password = hash

		if err := models.UserUpdate(ctx, tx, user); err != nil {
			return httputil.Error(err, "failed to update user", http.StatusInternalServerError)
		}

		sessionID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.Error(err, "failed to generate id", http.StatusInternalServerError)
		}

		refreshToken, err := s.auth.GenerateToken()
		if err != nil {
			return httputil.Error(err, "failed to generate token", http.StatusInternalServerError)
		}

		session := models.UserSession{
			ID:        sessionID,
			UserID:    user.ID,
			Token:     refreshToken,
			IPAddress: httputil.IP(r),
			UserAgent: httputil.UserAgent(r),
		}
		if err := models.UserSessionInsert(ctx, tx, &session); err != nil {
			return httputil.Error(err, "failed to create session", http.StatusInternalServerError)
		}

		client := httputil.Client(r)
		accessToken, err := s.auth.GenerateJWT(user.ID, client, user.Role)
		if err != nil {
			return httputil.Error(err, "failed to generate jwt", http.StatusInternalServerError)
		}

		res.TokenType = auth.TokenType
		res.AccessToken = accessToken
		res.RefreshToken = refreshToken
		res.ExpiresAt = session.ExpiresAt.Unix()

		return nil
	})
	if err != nil {
		if httpErr := httputil.CastError(err); httpErr != nil {
			s.logger.ErrorContext(ctx, "error occured in reset password handler", "error", httpErr.Wrap())
			httpErr.WriteJSON(w)
			return
		}
		s.logger.ErrorContext(ctx, "error occured in reset password handler", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", auth.TokenPrefix+res.AccessToken)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		httputil.Error(err, "failed to return response", http.StatusInternalServerError).WriteJSON(w)
	}
}
