// Package apiservice provides the API service.
package apiservice

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
	return "api/v1"
}

func (s *Service) Register(k *kernel.Kernel) (http.Handler, error) {
	s.db = k.Database()
	s.auth = k.Auth()
	s.logger = k.Logger().WithGroup(s.Name())

	mux := http.NewServeMux()
	mux.HandleFunc("/", httputil.NotFoundHandler(s.logger))
	mux.HandleFunc("GET /users", s.HandleGetUsers)
	mux.HandleFunc("GET /users/{id}", s.HandleGetUserByID)
	mux.HandleFunc("PATCH /users/{id}", s.HandleUpdateUser)

	handler := httputil.ChainMiddlewares(httputil.AuthGuardMiddleware(s.auth))(mux)
	return handler, nil
}
