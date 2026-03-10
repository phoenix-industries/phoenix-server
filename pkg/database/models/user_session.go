package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type UserSession struct {
	Model
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

func (s *UserSession) Validate() error {
	if s.UserID == "" {
		return errors.New("user_id is required")
	}
	if s.Token == "" {
		return errors.New("token is required")
	}
	if s.IPAddress == "" {
		return errors.New("ip_address is required")
	}
	if s.UserAgent == "" {
		return errors.New("user_agent is required")
	}
	return nil
}

func UserSessionInsert(ctx context.Context, db database.DB, userSession *UserSession) error {
	if userSession.ID == "" {
		return errors.New("id is not set")
	}
	if userSession.ExpiresAt.IsZero() {
		userSession.ExpiresAt = time.Now().Add(auth.DefaultSessionDuration)
	}
	stmt := `
		INSERT INTO user_sessions
		(id, user_id, token, ip_address, user_agent, expires_at)
		VALUES
		($1, $2, $3, $4, $5, $6)
	`
	if _, err := db.Exec(ctx, stmt, userSession.ID, userSession.UserID, userSession.Token, userSession.IPAddress, userSession.UserAgent, userSession.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func UserSessionGetByID(ctx context.Context, db database.DB, id string) (*UserSession, error) {
	stmt := `SELECT * FROM user_sessions WHERE id = $1 and deleted_at is null`
	var userSession UserSession
	if err := pgxscan.Get(ctx, db, &userSession, stmt, id); err != nil {
		return nil, err
	}
	return &userSession, nil
}

func UserSessionGetByToken(ctx context.Context, db database.DB, token string) (*UserSession, error) {
	stmt := `SELECT * FROM user_sessions WHERE token = $1 and deleted_at is null`
	var userSession UserSession
	if err := pgxscan.Get(ctx, db, &userSession, stmt, token); err != nil {
		return nil, err
	}
	return &userSession, nil
}

func UserSessionGetByUserID(ctx context.Context, db database.DB, userID string) (*UserSession, error) {
	stmt := `SELECT * FROM user_sessions WHERE user_id = $1 and deleted_at is null`
	var userSession UserSession
	if err := pgxscan.Get(ctx, db, &userSession, stmt, userID); err != nil {
		return nil, err
	}
	return &userSession, nil
}

func UserSessionDeleteByID(ctx context.Context, db database.DB, id string) error {
	stmt := `UPDATE user_sessions SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 and deleted_at is null`
	_, err := db.Exec(ctx, stmt, id)
	return err
}
