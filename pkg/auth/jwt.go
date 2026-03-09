package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

const (
	defaultIssuer    = "phoenix"
	defaultExpiresAt = 10 * time.Minute
)

type Claims struct {
	jwt.RegisteredClaims
	Role database.UserRole `json:"role"`
}

func (c *Claims) Valid() error {
	if c.Role == "" {
		return errors.New("invalid token: role is empty")
	}
	return nil
}

func JWTSecretFromEnv() ([]byte, error) {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		return nil, errors.New("jwt secret is not set in env `AUTH_JWT_SECRET`")
	}
	return []byte(secret), nil
}

// GenerateToken generates a JWT token with the HS256 method.
func (a *Auth) GenerateToken(subject, audience string, role database.UserRole) (string, error) {
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultIssuer,
			Subject:   subject,
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultExpiresAt)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := t.SignedString(a.jwtSecret)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return ss, nil
}

// ParseToken parses a JWT token.
func (a *Auth) ParseToken(token string) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, a.jwtParseKeyFunc, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	} else if !t.Valid {
		return nil, errors.New("invalid token")
	} else if _, ok := t.Claims.(*Claims); !ok {
		return nil, errors.New("invalid token claims")
	}

	issuer, err := t.Claims.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("invalid token issuer: %w", err)
	} else if issuer != defaultIssuer {
		return nil, errors.New("invalid token issuer")
	}

	return t, nil
}

func (a *Auth) jwtParseKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("invalid token")
	}
	return a.jwtSecret, nil
}
