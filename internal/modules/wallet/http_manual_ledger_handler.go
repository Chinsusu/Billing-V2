package wallet

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const walletIdempotencyKeyHeader = "Idempotency-Key"

type createWalletRefundBody struct {
	WalletID      WalletID      `json:"wallet_id"`
	AmountMinor   int64         `json:"amount_minor"`
	Currency      string        `json:"currency"`
	ReferenceType ReferenceType `json:"reference_type"`
	ReferenceID   ReferenceID   `json:"reference_id"`
	Reason        string        `json:"reason"`
	CorrelationID CorrelationID `json:"correlation_id"`
}

type createWalletAdjustmentBody struct {
	WalletID      WalletID      `json:"wallet_id"`
	Direction     Direction     `json:"direction"`
	AmountMinor   int64         `json:"amount_minor"`
	Currency      string        `json:"currency"`
	ReferenceType ReferenceType `json:"reference_type"`
	ReferenceID   ReferenceID   `json:"reference_id"`
	Reason        string        `json:"reason"`
	CorrelationID CorrelationID `json:"correlation_id"`
}

func (handler *HTTPHandler) adminWalletRefundsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPost: handler.tenantRoute(handler.handleCreateAdminWalletRefund, handler.options.AdminAdjustmentMiddleware),
	})
}

func (handler *HTTPHandler) adminWalletAdjustmentsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPost: handler.tenantRoute(handler.handleCreateAdminWalletAdjustment, handler.options.AdminAdjustmentMiddleware),
	})
}

func (handler *HTTPHandler) handleCreateAdminWalletRefund(w http.ResponseWriter, r *http.Request) {
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
	var body createWalletRefundBody
	if !decodeManualLedgerJSON(w, r, &body) {
		return
	}
	entry, err := handler.service.CreateWalletRefund(r.Context(), CreateWalletRefundInput{
		TenantID:       tenantID,
		WalletID:       body.WalletID,
		AmountMinor:    body.AmountMinor,
		Currency:       body.Currency,
		ReferenceType:  body.ReferenceType,
		ReferenceID:    body.ReferenceID,
		IdempotencyKey: IdempotencyKey(strings.TrimSpace(r.Header.Get(walletIdempotencyKeyHeader))),
		CreatedBy:      actor.ID,
		Reason:         body.Reason,
		CorrelationID:  body.CorrelationID,
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newLedgerEntryResponse(entry))
}

func (handler *HTTPHandler) handleCreateAdminWalletAdjustment(w http.ResponseWriter, r *http.Request) {
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
	var body createWalletAdjustmentBody
	if !decodeManualLedgerJSON(w, r, &body) {
		return
	}
	entry, err := handler.service.CreateWalletAdjustment(r.Context(), CreateWalletAdjustmentInput{
		TenantID:       tenantID,
		WalletID:       body.WalletID,
		Direction:      body.Direction,
		AmountMinor:    body.AmountMinor,
		Currency:       body.Currency,
		ReferenceType:  body.ReferenceType,
		ReferenceID:    body.ReferenceID,
		IdempotencyKey: IdempotencyKey(strings.TrimSpace(r.Header.Get(walletIdempotencyKeyHeader))),
		CreatedBy:      actor.ID,
		Reason:         body.Reason,
		CorrelationID:  body.CorrelationID,
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newLedgerEntryResponse(entry))
}

func decodeManualLedgerJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("body", "request.body_invalid", "Request body must be valid JSON."),
		})
		return false
	}
	return true
}
