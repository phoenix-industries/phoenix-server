// Package server provides HTTP server.
package server

import (
	"log/slog"
	"net/http"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
)

type Server struct {
	mux    *http.ServeMux
	auth   *auth.Auth
	logger *slog.Logger
	port   string
}

func NewServer(port string, auth *auth.Auth, logger *slog.Logger) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		auth:   auth,
		logger: logger,
		port:   port,
	}
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) Logger() *slog.Logger {
	return s.logger
}

type Router interface {
	Register(mux *http.ServeMux) error
}

func (s *Server) Register(r Router) error {
	return r.Register(s.mux)
}

func (s *Server) Run() error {
	s.logger.Info("Starting server", "port", s.port)
	if err := http.ListenAndServe(s.port, s.mux); err != nil {
		return err
	}
	return nil
}
