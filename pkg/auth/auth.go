// Package auth provides authentication functionality.
package auth

import (
	"log/slog"
)

type Auth struct {
	jwtSecret []byte
	logger    *slog.Logger
}

func New(jwtSecret []byte) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
		logger:    slog.Default(),
	}
}

func NewWithLogger(jwtSecret []byte, logger *slog.Logger) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}
