package payment

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type HTTPService interface {
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error)
	PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error)
	ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error)
	GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	AdminMiddleware    RouteMiddleware
	ResellerMiddleware RouteMiddleware
	ClientMiddleware   RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

const (
	IdempotencyKeyHeader = "Idempotency-Key"

	adminTransactionPrefix          = "/admin/transactions/"
	clientTransactionPrefix         = "/client/transactions/"
	clientInvoiceWalletPaymentsPath = "/client/invoice-wallet-payments"
	maxJSONBodyBytes                = 1 << 20
)

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/payment-reconciliation", handler.adminReconciliationsRoute)
	mux.HandleFunc("/admin/payment-reconciliation/", handler.adminReconciliationRoute)
	mux.HandleFunc("/admin/transactions", handler.adminTransactionsRoute)
	mux.HandleFunc("/admin/transactions/", handler.adminTransactionRoute)
	mux.HandleFunc("/reseller/transactions", handler.resellerTransactionsRoute)
	mux.HandleFunc(clientInvoiceWalletPaymentsPath, handler.clientInvoiceWalletPaymentsRoute)
	mux.HandleFunc("/client/transactions", handler.clientTransactionsRoute)
	mux.HandleFunc("/client/transactions/", handler.clientTransactionRoute)
}

func (handler *HTTPHandler) adminTransactionsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminTransactions, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminTransactionRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminTransaction, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) resellerTransactionsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminTransactions, handler.options.ResellerMiddleware),
	})
}

func (handler *HTTPHandler) clientTransactionsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListClientTransactions, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) clientTransactionRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetClientTransaction, handler.options.ClientMiddleware),
	})
}

func dispatchPaymentMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", paymentAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func paymentAllowHeader(methods map[string]http.HandlerFunc) string {
	allowed := make([]string, 0, len(methods))
	for method := range methods {
		allowed = append(allowed, method)
	}
	sort.Strings(allowed)
	return strings.Join(allowed, ", ")
}

func (handler *HTTPHandler) tenantRoute(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	return tenantContext(requireTenantContext(applyRouteMiddleware(next, routeMiddleware)))
}

func tenantContext(next http.HandlerFunc) http.HandlerFunc {
	handler := tenant.HeaderContextMiddleware(http.HandlerFunc(next))
	return handler.ServeHTTP
}

func requireTenantContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := tenantIDFromContext(w, r); !ok {
			return
		}
		next(w, r)
	}
}

func applyRouteMiddleware(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	if routeMiddleware == nil {
		return next
	}
	return routeMiddleware(next)
}

func (handler *HTTPHandler) handleListAdminTransactions(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := transactionFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	transactions, err := handler.service.ListTransactions(r.Context(), filter)
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTransactionResponses(transactions), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminTransaction(w http.ResponseWriter, r *http.Request) {
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
	transactionID, ok := adminTransactionIDFromPath(w, r)
	if !ok {
		return
	}
	transaction, err := handler.service.GetTransaction(r.Context(), TransactionLookup{ID: transactionID, TenantID: tenantID})
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTransactionResponse(transaction))
}

func (handler *HTTPHandler) handleListClientTransactions(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := transactionFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.AccountUserID = actor.ID
	transactions, err := handler.service.ListTransactions(r.Context(), filter)
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTransactionResponses(transactions), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientTransaction(w http.ResponseWriter, r *http.Request) {
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
	transactionID, ok := clientTransactionIDFromPath(w, r)
	if !ok {
		return
	}
	transaction, err := handler.service.GetTransaction(r.Context(), TransactionLookup{
		ID:            transactionID,
		TenantID:      tenantID,
		AccountUserID: actor.ID,
	})
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTransactionResponse(transaction))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writePaymentError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func transactionFilterFromRequest(w http.ResponseWriter, r *http.Request) (TransactionFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := TransactionFilter{Limit: page.Limit}
	query := r.URL.Query()
	accountUserID := identity.UserID(strings.TrimSpace(query.Get("account_user_id")))
	if accountUserID != "" {
		filter.AccountUserID = accountUserID
	}
	if accountDisplayID, present, ok := paymentPositiveInt64Query(w, r, "account_display_id"); !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.AccountDisplayID = accountDisplayID
	}
	if displayID, present, ok := paymentPositiveInt64Query(w, r, "display_id"); !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	orderID := order.OrderID(strings.TrimSpace(query.Get("order_id")))
	if orderID != "" {
		filter.OrderID = orderID
	}
	if orderDisplayID, present, ok := paymentPositiveInt64Query(w, r, "order_display_id"); !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.OrderDisplayID = orderDisplayID
	}
	invoiceID := invoice.InvoiceID(strings.TrimSpace(query.Get("invoice_id")))
	if invoiceID != "" {
		filter.InvoiceID = invoiceID
	}
	if invoiceDisplayID, present, ok := paymentPositiveInt64Query(w, r, "invoice_display_id"); !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.InvoiceDisplayID = invoiceDisplayID
	}
	transactionType := TransactionType(strings.TrimSpace(query.Get("type")))
	if transactionType != "" {
		if !transactionType.Valid() {
			writePaymentError(w, r, ErrTypeInvalid)
			return TransactionFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Type = transactionType
	}
	status := TransactionStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writePaymentError(w, r, ErrStatusInvalid)
			return TransactionFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	amountMin, amountMax, ok := paymentAmountRangeQuery(w, r)
	if !ok {
		return TransactionFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.AmountMinMinor = amountMin
	filter.AmountMaxMinor = amountMax
	return filter, page, true
}

func adminTransactionIDFromPath(w http.ResponseWriter, r *http.Request) (TransactionID, bool) {
	return transactionIDFromPrefix(w, r, adminTransactionPrefix)
}

func clientTransactionIDFromPath(w http.ResponseWriter, r *http.Request) (TransactionID, bool) {
	return transactionIDFromPrefix(w, r, clientTransactionPrefix)
}

func transactionIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (TransactionID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writePaymentError(w, r, ErrTransactionIDMissing)
		return "", false
	}
	return TransactionID(value), true
}

func tenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writePaymentError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writePaymentError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func pageFromRequest(w http.ResponseWriter, r *http.Request) (httpserver.CursorPageRequest, bool) {
	page, err := httpserver.ParseCursorPage(r)
	if err == nil {
		return page, true
	}
	switch {
	case errors.Is(err, httpserver.ErrPageLimitTooLarge):
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{validationField("limit", "request.limit_too_large", "Limit is too large.")})
	default:
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{validationField("limit", "request.limit_invalid", "Limit must be a positive number.")})
	}
	return httpserver.CursorPageRequest{}, false
}

func writePaymentError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := paymentValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrTransactionNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "payment.transaction_not_found", "Payment transaction was not found.")
	case errors.Is(err, invoice.ErrInvoiceNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "invoice.not_found", "Invoice was not found.")
	case errors.Is(err, order.ErrOrderNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "order.not_found", "Order was not found.")
	case errors.Is(err, wallet.ErrWalletNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "wallet.not_found", "Wallet was not found.")
	case errors.Is(err, ErrInvoiceNotPayable):
		httpserver.WriteError(w, r, http.StatusConflict, "payment.invoice_not_payable", "Invoice is not payable.")
	case errors.Is(err, ErrIdempotencyConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "payment.idempotency_conflict", "Idempotency key conflicts with another payment.")
	case errors.Is(err, order.ErrOrderStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "order.status_conflict", "Order status changed before payment completed.")
	case errors.Is(err, order.ErrProvisioningSourceNotFound):
		httpserver.WriteError(w, r, http.StatusConflict, "order.provisioning_source_not_found", "No active provisioning source is available for this order.")
	case errors.Is(err, ErrWalletCurrencyMismatch):
		httpserver.WriteError(w, r, http.StatusConflict, "payment.wallet_currency_mismatch", "Wallet currency does not match invoice currency.")
	case errors.Is(err, wallet.ErrInsufficientBalance):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.insufficient_balance", "Wallet balance is insufficient.")
	case errors.Is(err, invoice.ErrInvoiceStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "invoice.status_conflict", "Invoice status changed before payment completed.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "payment.service_unavailable", "Payment service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "payment.operation_failed", "Payment operation failed.")
	}
}

func paymentValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrAccountIDMissing):
		return validationField("account_user_id", "payment.account_missing", "Account user is required."), true
	case errors.Is(err, ErrTransactionIDMissing):
		return validationField("transaction_id", "payment.transaction_id_missing", "Transaction id is required."), true
	case errors.Is(err, invoice.ErrInvoiceIDMissing):
		return validationField("invoice_id", "invoice.invoice_id_missing", "Invoice id is required."), true
	case errors.Is(err, wallet.ErrWalletIDMissing):
		return validationField("wallet_id", "wallet.wallet_id_missing", "Wallet id is required."), true
	case errors.Is(err, ErrTypeInvalid):
		return validationField("type", "payment.type_invalid", "Transaction type is invalid."), true
	case errors.Is(err, ErrStatusInvalid):
		return validationField("status", "payment.status_invalid", "Transaction status is invalid."), true
	case errors.Is(err, ErrCurrencyMissing):
		return validationField("currency", "payment.currency_missing", "Currency is required."), true
	case errors.Is(err, ErrCurrencyInvalid):
		return validationField("currency", "payment.currency_invalid", "Currency is invalid."), true
	case errors.Is(err, ErrAmountInvalid):
		return validationField("amount_minor", "payment.amount_invalid", "Amount must be greater than zero."), true
	case errors.Is(err, ErrIdempotencyKeyMissing):
		return validationField("idempotency_key", "payment.idempotency_key_missing", "Idempotency key is required."), true
	case errors.Is(err, ErrCreatedTimeInvalid):
		return validationField("created_at", "payment.created_time_invalid", "Created time is invalid."), true
	case errors.Is(err, ErrCreatedTimeWindowInvalid):
		return validationField("created_at", "payment.created_time_window_invalid", "Created time window is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
