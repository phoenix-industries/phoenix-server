package authservice

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type logoutData struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) HandleLogout(w http.ResponseWriter, r *http.Request) *httputil.Response {
	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ResponseFromError(err)
	}
	if userID == "" {
		return httputil.ErrUnauthorized.Response()
	}

	var data logoutData
	if err := httputil.BodyJSON(r, &data); err != nil {
		return httputil.NewStatusError(nil, "invalid request body", http.StatusBadRequest).Response()
	}
	if data.RefreshToken == "" {
		return httputil.ErrUnauthorized.Response()
	}

	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		session, err := models.UserSessionGetByToken(ctx, tx, data.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if session.UserID != userID {
			return httputil.ErrUnauthorized
		}
		if err := models.UserSessionDeleteByID(ctx, tx, session.ID); err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}
		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusOK, nil)
}
