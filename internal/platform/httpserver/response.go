package httpserver

import (
	"encoding/json"
	"net/http"
)

type SuccessEnvelope struct {
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

type ErrorEnvelope struct {
	Error     ErrorBody `json:"error"`
	RequestID string    `json:"request_id"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteSuccess(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	writeJSON(w, statusCode, SuccessEnvelope{
		Data:      data,
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func WriteError(w http.ResponseWriter, r *http.Request, statusCode int, code string, message string) {
	writeJSON(w, statusCode, ErrorEnvelope{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}
