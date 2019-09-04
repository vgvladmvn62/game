package server

import (
	"encoding/json"
	"net/http"
)

// APIError returned from endpoints in json
type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewAPIError is just a constructor. Nothing extraordinary
func NewAPIError(msg string, code int) APIError {
	return APIError{msg, code}
}

// Send sends a marshaled Error with given code
func (e APIError) Send(w http.ResponseWriter) error {
	out, _ := json.Marshal(e) // TODO: can there be error?

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)
	_, err := w.Write(out)
	return err
}
