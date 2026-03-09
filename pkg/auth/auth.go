// Package auth provides authentication functionality.
package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

type Auth struct {
	jwtSecret             []byte
	passwordSigningKey    []byte
	passwordSaltSeparator []byte
}

func New(jwtSecret []byte) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
	}
}

func NewFromEnv() (*Auth, error) {
	jwtSecret, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_JWT_SECRET"))
	if err != nil {
		return nil, errors.New("failed to decode jwt secret")
	}
	signingKey, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_PASSWORD_SIGNING_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode signer key: %v", err)
	}
	saltSep, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_PASSWORD_SALT_SEPARATOR"))
	if err != nil {
		return nil, fmt.Errorf("hasher.init: failed to decode salt separator: %v", err)
	}
	return &Auth{
		jwtSecret:             []byte(jwtSecret),
		passwordSigningKey:    signingKey,
		passwordSaltSeparator: saltSep,
	}, nil
}
