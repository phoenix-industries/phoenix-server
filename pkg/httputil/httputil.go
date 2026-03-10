// Package httputil contains http utilities.
package httputil

import (
	"net/http"
	"strings"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
)

func UserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func IP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func Client(r *http.Request) string {
	return r.Header.Get("X-Client")
}

func GetAccessToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", ErrUnauthorized
	}
	if !strings.HasPrefix(header, auth.TokenPrefix) {
		return "", ErrUnauthorized
	}
	token := header[len(auth.TokenPrefix):]
	if token == "" {
		return "", ErrUnauthorized
	}
	return token, nil
}
