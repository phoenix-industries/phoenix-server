package chatservice

import (
	"context"
	"log/slog"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/chat"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
	"github.com/phoenix-industries/phoenix-server/pkg/kernel"
)

type Service struct {
	db     *database.Database
	auth   *auth.Auth
	logger *slog.Logger
	server *chat.Server
}

// must implement kernel.Service
var _ kernel.Service = (*Service)(nil)

func New() *Service {
	return &Service{}
}

func (s *Service) Name() string {
	return "chat/v1"
}

func (s *Service) Register(ctx context.Context, env *kernel.Env) error {
	name := s.Name()
	s.db = env.Database()
	s.auth = env.Auth()
	s.logger = env.Logger().WithGroup(name)
	s.server = chat.NewServer(s.db, s.logger.WithGroup(name), chat.DefaultConfig())

	r := env.Router().Group("/" + name)
	r.Use(httputil.AuthGuardMiddleware(s.auth))

	r.HandleFunc("POST /direct", s.HandleSendDirect)
	r.HandleFunc("GET /rooms", s.HandleListRooms)
	r.HandleFuncNative("/connect", s.HandleConnect)

	return nil
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return nil
}
