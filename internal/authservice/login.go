package authservice

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
	"github.com/phoenix-industries/phoenix-server/pkg/validate"
)

const (
	errInvalidCredentials = "invalid credentials"
)

type loginData struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data loginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httputil.Error(nil, "invalid request body", http.StatusBadRequest).WriteJSON(w)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	res := AuthResponse{}
	err := s.db.InTx(ctx, func(tx pgx.Tx) error {
		var user *models.User
		if validate.IsEmail(data.Identifier) {
			u, err := models.UserGetByEmail(ctx, s.db.Pool(), data.Identifier)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return httputil.Error(nil, errInvalidCredentials, http.StatusUnauthorized)
				}
				return httputil.Error(err, "request failed", http.StatusInternalServerError)
			}
			user = u
		} else if validate.IsPhoneNumber(data.Identifier) {
			u, err := models.UserGetByPhone(ctx, s.db.Pool(), data.Identifier)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return httputil.Error(nil, errInvalidCredentials, http.StatusUnauthorized)
				}
				return httputil.Error(err, "request failed", http.StatusInternalServerError)
			}
			user = u
		} else {
			return httputil.Error(nil, errInvalidCredentials, http.StatusUnauthorized)
		}
		if user == nil {
			return httputil.Error(nil, errInvalidCredentials, http.StatusUnauthorized)
		}

		if valid, err := s.auth.VerifyPassword(data.Password, user.Password); err != nil {
			return httputil.Error(err, "request failed", http.StatusInternalServerError)
		} else if !valid {
			return httputil.Error(nil, errInvalidCredentials, http.StatusUnauthorized)
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
		jwt, err := s.auth.GenerateJWT(user.ID, client, user.Role)
		if err != nil {
			return httputil.Error(err, "failed to generate jwt", http.StatusInternalServerError)
		}

		res.TokenType = auth.TokenType
		res.AccessToken = jwt
		res.RefreshToken = refreshToken
		res.ExpiresAt = session.ExpiresAt.Unix()

		return nil
	})
	if err != nil {
		if httpErr := httputil.CastError(err); httpErr != nil {
			s.logger.ErrorContext(ctx, "error occured in login handler", "error", httpErr.Wrap())
			httpErr.WriteJSON(w)
			return
		}
		s.logger.ErrorContext(ctx, "error occured in login handler", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", auth.TokenPrefix+res.AccessToken)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&res)
}
