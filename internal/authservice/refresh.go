package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type refreshData struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var data refreshData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httputil.Error(nil, "invalid request body", http.StatusBadRequest).WriteJSON(w)
		return
	}
	if data.RefreshToken == "" {
		httputil.ErrorBadRequest().WriteJSON(w)
		return
	}

	res := AuthResponse{}
	ctx := r.Context()
	err := s.db.InTx(ctx, func(tx pgx.Tx) error {
		session, err := models.UserSessionGetByToken(ctx, tx, data.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if time.Now().After(session.ExpiresAt.Add(auth.DefaultJWTDuration)) {
			return httputil.ErrorUnauthorized()
		}
		user, err := models.UserGetByID(ctx, tx, session.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		client := httputil.Client(r)
		accessToken, err := s.auth.GenerateJWT(user.ID, client, user.Role)
		if err != nil {
			return httputil.Error(err, "failed to generate jwt", http.StatusInternalServerError)
		}

		res.TokenType = auth.TokenType
		res.AccessToken = accessToken
		res.RefreshToken = session.Token
		res.ExpiresAt = session.ExpiresAt.Unix()

		return nil
	})
	if err != nil {
		if httpErr := httputil.CastError(err); httpErr != nil {
			s.logger.ErrorContext(ctx, "error occured in refresh handler", "error", httpErr.Wrap())
			httpErr.WriteJSON(w)
			return
		}
		s.logger.ErrorContext(ctx, "error occured in refresh handler", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", auth.TokenPrefix+res.AccessToken)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		httputil.Error(err, "failed to return response", http.StatusInternalServerError).WriteJSON(w)
	}
}
