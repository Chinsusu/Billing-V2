package order

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const serviceRenewSuffix = "/renew"

type clientServiceRenewalRequest struct {
	WalletID   wallet.WalletID `json:"wallet_id"`
	FromStatus ServiceStatus   `json:"from_status"`
	Reason     string          `json:"reason"`
}

type clientServiceRenewalResponse struct {
	Service            serviceInstanceResponse             `json:"service"`
	Invoice            clientServiceRenewalInvoiceResponse `json:"invoice"`
	PaymentTransaction clientServiceRenewalPaymentResponse `json:"payment_transaction"`
	Ledger             clientServiceRenewalLedgerResponse  `json:"ledger"`
	AmountMinor        int64                               `json:"amount_minor"`
	Currency           string                              `json:"currency"`
	Renewed            bool                                `json:"renewed"`
}

type clientServiceRenewalInvoiceResponse struct {
	ID         string `json:"id"`
	DisplayID  int64  `json:"display_id"`
	Status     string `json:"status"`
	TotalMinor int64  `json:"total_minor"`
	Currency   string `json:"currency"`
}

type clientServiceRenewalPaymentResponse struct {
	ID        string `json:"id"`
	DisplayID int64  `json:"display_id"`
	Status    string `json:"status"`
}

type clientServiceRenewalLedgerResponse struct {
	ID        wallet.LedgerEntryID `json:"id"`
	DisplayID int64                `json:"display_id"`
	WalletID  wallet.WalletID      `json:"wallet_id"`
	EntryType string               `json:"entry_type"`
}

func isServiceRenewPath(path string, prefix string) bool {
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	value := strings.TrimPrefix(path, prefix)
	if !strings.HasSuffix(value, serviceRenewSuffix) {
		return false
	}
	serviceID := strings.TrimSuffix(value, serviceRenewSuffix)
	return serviceID != "" && !strings.Contains(serviceID, "/")
}

func serviceIDFromRenewPath(w http.ResponseWriter, r *http.Request, prefix string) (ServiceID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	value = strings.TrimSuffix(value, serviceRenewSuffix)
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrServiceIDMissing)
		return "", false
	}
	return ServiceID(value), true
}

func (handler *HTTPHandler) handleRenewClientService(w http.ResponseWriter, r *http.Request) {
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
	serviceID, ok := serviceIDFromRenewPath(w, r, clientServicePrefix)
	if !ok {
		return
	}
	var request clientServiceRenewalRequest
	if !decodeOrderJSON(w, r, &request) {
		return
	}
	result, err := handler.service.RenewClientService(r.Context(), request.toInput(tenantID, actor.ID, serviceID, idempotencyKeyFromHeader(r)))
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newClientServiceRenewalResponse(result))
}

func (request clientServiceRenewalRequest) toInput(tenantID tenant.ID, actorID identity.UserID, serviceID ServiceID, idempotencyKey IdempotencyKey) ClientServiceRenewalInput {
	return ClientServiceRenewalInput{
		TenantID:       tenantID,
		BuyerUserID:    actorID,
		ServiceID:      serviceID,
		WalletID:       request.WalletID,
		ActorID:        actorID,
		FromStatus:     request.FromStatus,
		Reason:         request.Reason,
		IdempotencyKey: idempotencyKey,
	}
}

func newClientServiceRenewalResponse(result ClientServiceRenewal) clientServiceRenewalResponse {
	return clientServiceRenewalResponse{
		Service: newServiceInstanceResponse(result.Service),
		Invoice: clientServiceRenewalInvoiceResponse{
			ID:         result.InvoiceID,
			DisplayID:  result.InvoiceDisplayID,
			Status:     "paid",
			TotalMinor: result.AmountMinor,
			Currency:   result.Currency,
		},
		PaymentTransaction: clientServiceRenewalPaymentResponse{
			ID:        result.PaymentTransactionID,
			DisplayID: result.PaymentTransactionDisplay,
			Status:    "posted",
		},
		Ledger: clientServiceRenewalLedgerResponse{
			ID:        result.LedgerEntryID,
			DisplayID: result.LedgerEntryDisplayID,
			WalletID:  result.WalletID,
			EntryType: "purchase",
		},
		AmountMinor: result.AmountMinor,
		Currency:    result.Currency,
		Renewed:     result.Renewed,
	}
}
