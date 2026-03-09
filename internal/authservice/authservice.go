// Package authservice provides the authentication service.
package authservice

import (
	"log/slog"
	"net/http"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
	"github.com/phoenix-industries/phoenix-server/pkg/kernel"
)

type Service struct {
	db     *database.Database
	auth   *auth.Auth
	logger *slog.Logger
}

func New() kernel.Service {
	return &Service{}
}

func (s *Service) Name() string {
	return "auth_service"
}

func (s *Service) Register(k *kernel.Kernel) error {
	s.db = k.Database()
	s.auth = k.Auth()
	s.logger = k.Logger()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", s.RegisterHandler)
	mux.HandleFunc("/login", s.LoginHandler)
	mux.HandleFunc("/logout", s.LogoutHandler)

	middleware := httputil.NewMiddleware(s.logger)
	handler := httputil.ChainMiddlewares(
		middleware.Logging,
		middleware.Recovery,
	)(mux)

	k.Mux().Handle("/auth", handler)

	return nil
}

type AuthResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
