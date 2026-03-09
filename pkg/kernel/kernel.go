// Package kernel manages the core dependencies and services.
package kernel

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Service interface {
	Name() string
	Register(k *Kernel) error
}

type Kernel struct {
	services map[string]Service
	db       *database.Database
	mux      *http.ServeMux
	auth     *auth.Auth
	logger   *slog.Logger
}

func NewKernel(auth *auth.Auth, logger *slog.Logger) *Kernel {
	return &Kernel{
		services: map[string]Service{},
		mux:      http.NewServeMux(),
		auth:     auth,
		logger:   logger,
	}
}

func (k *Kernel) Database() *database.Database {
	return k.db
}

func (k *Kernel) Mux() *http.ServeMux {
	return k.mux
}

func (k *Kernel) Auth() *auth.Auth {
	return k.auth
}

func (k *Kernel) Logger() *slog.Logger {
	return k.logger
}

func (k *Kernel) Run(services ...Service) error {
	for _, service := range services {
		name := service.Name()
		if _, ok := k.services[name]; ok {
			return fmt.Errorf("service already registered: %s", name)
		}
		k.services[name] = service
		if err := service.Register(k); err != nil {
			return err
		}
		k.logger.Info("service registered", "name", name)
	}
	return nil
}
