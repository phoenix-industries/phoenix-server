// Package httputil contains http utilities.
package httputil

import "net/http"

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
