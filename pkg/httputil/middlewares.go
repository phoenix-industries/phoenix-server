package httputil

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/phoenix-industries/phoenix-server/internal/buildinfo"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
)

type Middleware func(http.Handler) http.Handler

type RequestID string

const RequestIDKey RequestID = "request_id"

func ChainMiddlewares(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for _, m := range middlewares {
			next = m(next)
		}
		return next
	}
}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID, err := auth.GenerateID()
			if err != nil {
				logger.Error("failed to generate id", "error", err)
				http.Error(w, "failed to generate id", http.StatusInternalServerError)
				return
			}
			r = r.WithContext(context.WithValue(ctx, RequestIDKey, requestID))
			w.Header().Set("X-Request-ID", requestID)
			if buildinfo.DevMode() {
				start := time.Now()
				next.ServeHTTP(w, r)
				duration := time.Since(start)
				logger.Info("request handled", "request_id", requestID, "duration", duration, "client", Client(r), "user_agent", UserAgent(r), "ip", IP(r))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RecoveryMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := r.Context().Value(RequestIDKey)
					logger.Error("request panic", "request_id", requestID, "error", err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: improve and add rate limiting
func AuthGuardMiddleware(auth *auth.Auth) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := CheckAuth(auth, r); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
