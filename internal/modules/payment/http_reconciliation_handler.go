package payment

import (
	"net/http"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const adminPaymentReconciliationPrefix = "/admin/payment-reconciliation/"

func (handler *HTTPHandler) adminReconciliationsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminPaymentReconciliations, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminReconciliationRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminPaymentReconciliation, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) handleListAdminPaymentReconciliations(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := actorFromContext(w, r); !ok {
		return
	}
	filter, page, ok := reconciliationFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	records, err := handler.service.ListPaymentReconciliations(r.Context(), filter)
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newPaymentReconciliationResponses(records), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminPaymentReconciliation(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := actorFromContext(w, r); !ok {
		return
	}
	transactionID, ok := paymentReconciliationTransactionIDFromPath(w, r)
	if !ok {
		return
	}
	record, err := handler.service.GetPaymentReconciliation(r.Context(), ReconciliationLookup{
		TenantID:      tenantID,
		TransactionID: transactionID,
	})
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newPaymentReconciliationResponse(record))
}

func reconciliationFilterFromRequest(w http.ResponseWriter, r *http.Request) (ReconciliationFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := ReconciliationFilter{Limit: page.Limit}
	query := r.URL.Query()
	if accountUserID := identity.UserID(strings.TrimSpace(query.Get("account_user_id"))); accountUserID != "" {
		filter.AccountUserID = accountUserID
	}
	if displayID, present, ok := paymentPositiveInt64Query(w, r, "display_id"); !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	status := TransactionStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writePaymentError(w, r, ErrStatusInvalid)
			return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	filter.Provider = strings.TrimSpace(query.Get("provider"))
	if invoiceID := invoice.InvoiceID(strings.TrimSpace(query.Get("invoice_id"))); invoiceID != "" {
		filter.InvoiceID = invoiceID
	}
	if invoiceDisplayID, present, ok := paymentPositiveInt64Query(w, r, "invoice_display_id"); !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.InvoiceDisplayID = invoiceDisplayID
	}
	if walletID := wallet.WalletID(strings.TrimSpace(query.Get("wallet_id"))); !walletID.Empty() {
		filter.WalletID = walletID
	}
	if walletDisplayID, present, ok := paymentPositiveInt64Query(w, r, "wallet_display_id"); !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.WalletDisplayID = walletDisplayID
	}
	amountMin, amountMax, ok := paymentAmountRangeQuery(w, r)
	if !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.AmountMinMinor = amountMin
	filter.AmountMaxMinor = amountMax
	createdFrom, ok := reconciliationTimeFromRequest(w, r, "created_from")
	if !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	}
	createdTo, ok := reconciliationTimeFromRequest(w, r, "created_to")
	if !ok {
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.CreatedFrom = createdFrom
	filter.CreatedTo = createdTo
	if !createdFrom.IsZero() && !createdTo.IsZero() && createdTo.Before(createdFrom) {
		writePaymentError(w, r, ErrCreatedTimeWindowInvalid)
		return ReconciliationFilter{}, httpserver.CursorPageRequest{}, false
	}
	return filter, page, true
}

func reconciliationTimeFromRequest(w http.ResponseWriter, r *http.Request, field string) (time.Time, bool) {
	value := strings.TrimSpace(r.URL.Query().Get(field))
	if value == "" {
		return time.Time{}, true
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField(field, "payment.created_time_invalid", "Created time must be RFC3339."),
		})
		return time.Time{}, false
	}
	return parsed, true
}

func paymentReconciliationTransactionIDFromPath(w http.ResponseWriter, r *http.Request) (TransactionID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, adminPaymentReconciliationPrefix))
	if value == "" || strings.Contains(value, "/") {
		writePaymentError(w, r, ErrTransactionIDMissing)
		return "", false
	}
	return TransactionID(value), true
}
