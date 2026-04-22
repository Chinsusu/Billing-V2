package httpserver

import (
	"encoding/json"
	"net/http"
)

type SuccessEnvelope struct {
	Data      any    `json:"data"`
	Page      *Page  `json:"page,omitempty"`
	RequestID string `json:"request_id"`
}

type ErrorEnvelope struct {
	Error     ErrorBody `json:"error"`
	RequestID string    `json:"request_id"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details any               `json:"details,omitempty"`
	Fields  []ValidationField `json:"fields,omitempty"`
}

type ValidationField struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteSuccess(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	writeJSON(w, statusCode, SuccessEnvelope{
		Data:      data,
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func WriteList(w http.ResponseWriter, r *http.Request, statusCode int, data any, page Page) {
	writeJSON(w, statusCode, SuccessEnvelope{
		Data:      data,
		Page:      &page,
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func WriteError(w http.ResponseWriter, r *http.Request, statusCode int, code string, message string) {
	WriteErrorWithDetails(w, r, statusCode, code, message, nil)
}

func WriteErrorWithDetails(w http.ResponseWriter, r *http.Request, statusCode int, code string, message string, details any) {
	writeJSON(w, statusCode, ErrorEnvelope{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func WriteValidationError(w http.ResponseWriter, r *http.Request, fields []ValidationField) {
	writeJSON(w, http.StatusBadRequest, ErrorEnvelope{
		Error: ErrorBody{
			Code:    "validation.failed",
			Message: "Request validation failed.",
			Fields:  append([]ValidationField(nil), fields...),
		},
		RequestID: RequestIDFromContext(r.Context()),
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}
