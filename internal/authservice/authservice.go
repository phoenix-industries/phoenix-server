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
	return "auth"
}

func (s *Service) Register(k *kernel.Kernel) (http.Handler, error) {
	s.db = k.Database()
	s.auth = k.Auth()
	s.logger = k.Logger().WithGroup(s.Name())

	mux := http.NewServeMux()
	mux.HandleFunc("/", httputil.NotFoundHandler(s.logger))
	mux.HandleFunc("POST /register", s.HandleRegister)
	mux.HandleFunc("POST /login", s.HandleLogin)
	mux.HandleFunc("POST /logout", s.HandleLogout)
	mux.HandleFunc("POST /refresh", s.HandleRefresh)

	return mux, nil
}

type AuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
