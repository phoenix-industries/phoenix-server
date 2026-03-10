package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
)

type HTTPError struct {
	msg  string
	err  error
	code int
}

func Error(err error, msg string, code int) *HTTPError {
	return &HTTPError{
		err:  err,
		msg:  msg,
		code: code,
	}
}

func CastError(err error) *HTTPError {
	if e, ok := err.(*HTTPError); ok {
		return e
	}
	return nil
}

func (e *HTTPError) Error() string {
	return e.msg
}

func (e *HTTPError) Unwrap() error {
	return e.err
}

func (e *HTTPError) Wrap() error {
	if e.err == nil {
		return errors.New(e.msg)
	}
	return fmt.Errorf("%s: %w", e.msg, e.err)
}

func (e *HTTPError) Write(w http.ResponseWriter) {
	w.WriteHeader(e.code)
	w.Write([]byte(e.msg))
}

func (e *HTTPError) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.code)
	return json.NewEncoder(w).Encode(map[string]string{"error": e.msg})
}

func ErrorBadRequest() *HTTPError {
	return Error(nil, ErrInvalidRequest.Error(), http.StatusBadRequest)
}

func ErrorUnauthorized() *HTTPError {
	return Error(nil, ErrUnauthorized.Error(), http.StatusUnauthorized)
}
