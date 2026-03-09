// Package kernel provides core dependencies for the services.
package kernel

import (
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
	mux      *http.ServeMux
	db       *database.Database
	auth     *auth.Auth
	logger   *slog.Logger
}

func NewKernel(port string, auth *auth.Auth, logger *slog.Logger) *Kernel {
	return &Kernel{
		mux:    http.NewServeMux(),
		auth:   auth,
		logger: logger,
	}
}

func (k *Kernel) Mux() *http.ServeMux {
	return k.mux
}

func (k *Kernel) Logger() *slog.Logger {
	return k.logger
}

func (k *Kernel) Run(services ...Service) error {
	for _, service := range services {
		name := service.Name()
		if _, ok := k.services[name]; ok {
			continue
		}
		k.services[name] = service
		if err := service.Register(k); err != nil {
			return err
		}
		k.logger.Info("service registered", "name", name)
	}
	return nil
}
