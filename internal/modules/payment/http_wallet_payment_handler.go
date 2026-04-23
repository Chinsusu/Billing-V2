package payment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func (handler *HTTPHandler) clientInvoiceWalletPaymentsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPost: handler.tenantRoute(handler.handleCreateClientInvoiceWalletPayment, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) handleCreateClientInvoiceWalletPayment(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	actor, ok := actorFromContext(w, r)
	if !ok {
		return
	}
	var request createInvoiceWalletPaymentRequest
	if !decodePaymentJSON(w, r, &request) {
		return
	}
	input := request.toInput(tenantID, actor.ID, idempotencyKeyFromHeader(r)).Normalize()
	if err := input.Validate(); err != nil {
		writePaymentError(w, r, err)
		return
	}
	result, err := handler.service.PayInvoiceFromWallet(r.Context(), input)
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newInvoiceWalletPaymentResponse(result))
}

func decodePaymentJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must be valid JSON.")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must contain one JSON object.")
		return false
	}
	return true
}

func idempotencyKeyFromHeader(r *http.Request) IdempotencyKey {
	return IdempotencyKey(strings.TrimSpace(r.Header.Get(IdempotencyKeyHeader)))
}
