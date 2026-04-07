package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type logoutData struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
	token, err := httputil.GetAccessToken(r)
	if err != nil {
		httputil.ErrorBadRequest().WriteJSON(w)
		return
	}
	jwt, err := s.auth.ParseJWT(token)
	if err != nil {
		httputil.ErrorBadRequest().WriteJSON(w)
		s.logger.ErrorContext(r.Context(), "token parsing error occured in logout handler", "error", err, "token", token)
		return
	}
	userID, err := jwt.Claims.GetSubject()
	if err != nil {
		httputil.ErrorBadRequest().WriteJSON(w)
		s.logger.ErrorContext(r.Context(), "error occured in logout handler while geting the subject claim from jwt token", "error", err, "token", token)
		return
	}
	if userID == "" {
		httputil.ErrorUnauthorized().WriteJSON(w)
		s.logger.ErrorContext(r.Context(), "invalid token spotted in logout handler", "token", token)
		return
	}

	var data logoutData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httputil.Error(nil, "invalid request body", http.StatusBadRequest).WriteJSON(w)
		return
	}
	if data.RefreshToken == "" {
		httputil.ErrorBadRequest().WriteJSON(w)
		return
	}

	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		session, err := models.UserSessionGetByToken(ctx, tx, data.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if session.UserID != userID {
			return httputil.ErrorUnauthorized()
		}
		if err := models.UserSessionDeleteByID(ctx, tx, session.ID); err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}
		return nil
	})
	if err != nil {
		if httpErr := httputil.CastError(err); httpErr != nil {
			s.logger.ErrorContext(ctx, "error occured in logout handler", "error", httpErr.Wrap())
			httpErr.WriteJSON(w)
			return
		}
		s.logger.ErrorContext(ctx, "error occured in logout handler", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
}
