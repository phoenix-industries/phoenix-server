// Package auth provides authentication functionality.
package auth

import (
	"log/slog"
)

type Auth struct {
	jwtSecret []byte
}

func New(jwtSecret []byte) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
	}
}

func NewWithLogger(jwtSecret []byte, logger *slog.Logger) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
	}
}
