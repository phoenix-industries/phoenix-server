package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	OK         bool   `json:"ok"`
	StatusCode int    `json:"-"`
	Data       any    `json:"data,omitempty"`
	err        error  `json:"-"`
	Error      string `json:"error,omitempty"`
}

func NewResponse(statusCode int, data any, err error) *Response {
	errmsg := ""
	if err != nil {
		errmsg = err.Error()
	}
	return &Response{
		OK:         err == nil,
		StatusCode: statusCode,
		Data:       data,
		err:        err,
		Error:      errmsg,
	}
}

func NewResponseOK(statusCode int, data any) *Response {
	return NewResponse(statusCode, data, nil)
}

func NewResponseError(statusCode int, err error) *Response {
	return NewResponse(statusCode, nil, err)
}

func ResponseFromError(err error) *Response {
	var statusErr *StatusError
	if errors.As(err, &statusErr) {
		return NewResponse(statusErr.StatusCode(), nil, err)
	}
	return NewResponse(http.StatusInternalServerError, nil, err)
}

func (r *Response) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	println("status code", r.StatusCode)
	w.WriteHeader(r.StatusCode)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}
	return nil
}

func (r *Response) EmbededError() error {
	if r.err == nil {
		return nil
	}
	return r.err
}
