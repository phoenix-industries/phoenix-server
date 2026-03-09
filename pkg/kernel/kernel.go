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
	Register(k *Kernel) (http.Handler, error)
}

type Kernel struct {
	mux      *http.ServeMux
	db       *database.Database
	auth     *auth.Auth
	logger   *slog.Logger
	services map[string]Service
}

func NewKernel(db *database.Database, auth *auth.Auth, logger *slog.Logger) *Kernel {
	return &Kernel{
		mux:      http.NewServeMux(),
		db:       db,
		auth:     auth,
		logger:   logger,
		services: map[string]Service{},
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
		handler, err := service.Register(k)
		if err != nil {
			return err
		}
		if handler != nil {
			k.mux.Handle(fmt.Sprintf("/%s/", name), http.StripPrefix(fmt.Sprintf("/%s", name), handler))
		}
		k.logger.Info("service registered", "name", name)
	}
	return nil
}
