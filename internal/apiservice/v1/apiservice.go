// Package apiservice provides the API service.
package apiservice

import (
	"context"
	"log/slog"

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

// must implement kernel.Service
var _ kernel.Service = (*Service)(nil)

func New() *Service {
	return &Service{}
}

func (s *Service) Name() string {
	return "api/v1"
}

func (s *Service) Register(ctx context.Context, env *kernel.Env) error {
	name := s.Name()
	s.db = env.Database()
	s.auth = env.Auth()
	s.logger = env.Logger().WithGroup(name)

	pub := env.Router().Group("/" + name)
	r := env.Router().Group("/" + name)
	r.Use(httputil.AuthGuardMiddleware(s.auth))

	r.HandleFunc("GET /users", s.HandleGetUsers)
	r.HandleFunc("GET /users/{id}", s.HandleGetUserByID)
	r.HandleFunc("PATCH /users/{id}", s.HandleUpdateUser)

	r.HandleFunc("POST /products/categories", s.HandleCreateProductCategory)
	pub.HandleFunc("GET /products/categories", s.HandleListProductCategories)
	pub.HandleFunc("GET /products/categories/{id}", s.HandleGetProductCategoryByID)
	r.HandleFunc("PUT /products/categories/{id}", s.HandleUpdateProductCategory)
	r.HandleFunc("DELETE /products/categories/{id}", s.HandleDeleteProductCategory)

	r.HandleFunc("POST /products", s.HandleCreateProduct)
	r.HandleFunc("POST /products/buy", s.HandleBuyProduct)
	pub.HandleFunc("GET /products", s.HandleListProducts)
	pub.HandleFunc("GET /products/{id}", s.HandleGetProductByID)
	r.HandleFunc("PATCH /products/{id}", s.HandleUpdateProduct)
	r.HandleFunc("DELETE /products/{id}", s.HandleDeleteProduct)

	return nil
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return nil
}
