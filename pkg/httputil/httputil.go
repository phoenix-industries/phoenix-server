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

func GetUserID(auther *auth.Auth, r *http.Request) (string, error) {
	token, err := GetAccessToken(r)
	if err != nil {
		return "", err
	}
	jwt, err := auther.ParseJWT(token)
	if err != nil {
		return "", err
	}
	userID, err := jwt.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return userID, nil
}

func GetUserRole(auther *auth.Auth, r *http.Request) (auth.Role, error) {
	token, err := GetAccessToken(r)
	if err != nil {
		return "", err
	}
	jwt, err := auther.ParseJWT(token)
	if err != nil {
		return "", err
	}
	claims, err := auther.GetJWTClaims(jwt)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

func CheckAuth(auther *auth.Auth, r *http.Request) error {
	_, err := GetAccessToken(r)
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}
