// Package authservice provides the authentication service.
package authservice

import (
	"context"
	"log/slog"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/kernel"
)

var _ kernel.Service = (*Service)(nil)

type Service struct {
	db     *database.Database
	auth   *auth.Auth
	logger *slog.Logger
}

func New() *Service {
	return &Service{}
}

func (s *Service) Name() string {
	return "auth"
}

func (s *Service) Register(ctx context.Context, env *kernel.Env) error {
	name := s.Name()
	s.db = env.Database()
	s.auth = env.Auth()
	s.logger = env.Logger().WithGroup(name)

	r := env.Router().Group("/" + name)

	r.HandleFunc("POST /register", s.HandleRegister)
	r.HandleFunc("POST /login", s.HandleLogin)
	r.HandleFunc("POST /logout", s.HandleLogout)
	r.HandleFunc("POST /refresh", s.HandleRefresh)

	return nil
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return nil
}

type AuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
