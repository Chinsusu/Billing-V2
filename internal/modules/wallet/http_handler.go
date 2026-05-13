package wallet

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	AdminMiddleware           RouteMiddleware
	AdminReviewMiddleware     RouteMiddleware
	AdminAdjustmentMiddleware RouteMiddleware
	ResellerMiddleware        RouteMiddleware
	ClientMiddleware          RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

const (
	adminWalletPrefix          = "/admin/wallets/"
	adminWalletRefundsPath     = "/admin/wallet-refunds"
	adminWalletAdjustmentsPath = "/admin/wallet-adjustments"
	adminTopupRequestPrefix    = "/admin/topup-requests/"
	resellerWalletPrefix       = "/reseller/wallets/"
	clientWalletPrefix         = "/client/wallets/"
	clientTopupRequestPrefix   = "/client/topup-requests/"
)

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/wallets", handler.adminWalletsRoute)
	mux.HandleFunc("/admin/wallets/", handler.adminWalletRoute)
	mux.HandleFunc(adminWalletRefundsPath, handler.adminWalletRefundsRoute)
	mux.HandleFunc(adminWalletAdjustmentsPath, handler.adminWalletAdjustmentsRoute)
	mux.HandleFunc("/admin/topup-requests", handler.adminTopupRequestsRoute)
	mux.HandleFunc("/admin/topup-requests/", handler.adminTopupRequestRoute)
	mux.HandleFunc("/reseller/wallets", handler.resellerWalletsRoute)
	mux.HandleFunc("/reseller/wallets/", handler.resellerWalletRoute)
	mux.HandleFunc("/reseller/topup-requests", handler.resellerTopupRequestsRoute)
	mux.HandleFunc("/client/wallets", handler.clientWalletsRoute)
	mux.HandleFunc("/client/wallets/", handler.clientWalletRoute)
	mux.HandleFunc("/client/topup-requests", handler.clientTopupRequestsRoute)
	mux.HandleFunc("/client/topup-requests/", handler.clientTopupRequestRoute)
}

func (handler *HTTPHandler) adminWalletsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminWallets, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminWalletRoute(w http.ResponseWriter, r *http.Request) {
	walletID, action, ok := walletPath(w, r, adminWalletPrefix)
	if !ok {
		return
	}
	switch action {
	case "":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleGetAdminWallet(w, r, walletID)
			}, handler.options.AdminMiddleware),
		})
	case "ledger":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleListAdminLedger(w, r, walletID)
			}, handler.options.AdminMiddleware),
		})
	default:
		writeWalletError(w, r, ErrWalletIDMissing)
	}
}

func (handler *HTTPHandler) resellerWalletsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminWallets, handler.options.ResellerMiddleware),
	})
}

func (handler *HTTPHandler) resellerWalletRoute(w http.ResponseWriter, r *http.Request) {
	walletID, action, ok := walletPath(w, r, resellerWalletPrefix)
	if !ok {
		return
	}
	switch action {
	case "":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleGetAdminWallet(w, r, walletID)
			}, handler.options.ResellerMiddleware),
		})
	case "ledger":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleListAdminLedger(w, r, walletID)
			}, handler.options.ResellerMiddleware),
		})
	default:
		writeWalletError(w, r, ErrWalletIDMissing)
	}
}

func (handler *HTTPHandler) clientWalletsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListClientWallets, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) clientWalletRoute(w http.ResponseWriter, r *http.Request) {
	walletID, action, ok := walletPath(w, r, clientWalletPrefix)
	if !ok {
		return
	}
	switch action {
	case "":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleGetClientWallet(w, r, walletID)
			}, handler.options.ClientMiddleware),
		})
	case "ledger":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleListClientLedger(w, r, walletID)
			}, handler.options.ClientMiddleware),
		})
	default:
		writeWalletError(w, r, ErrWalletIDMissing)
	}
}

func dispatchWalletMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", walletAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func walletAllowHeader(methods map[string]http.HandlerFunc) string {
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

func (handler *HTTPHandler) handleListAdminWallets(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := walletFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	wallets, err := handler.service.ListWallets(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newWalletResponses(wallets), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminWallet(w http.ResponseWriter, r *http.Request, walletID WalletID) {
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
	wallet, err := handler.service.GetWallet(r.Context(), WalletLookup{ID: walletID, TenantID: tenantID})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newWalletResponse(wallet))
}

func (handler *HTTPHandler) handleListAdminLedger(w http.ResponseWriter, r *http.Request, walletID WalletID) {
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
	filter, page, ok := ledgerFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.WalletID = walletID
	entries, err := handler.service.ListLedgerEntries(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newLedgerEntryResponses(entries), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleListClientWallets(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := walletFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.OwnerType = OwnerTypeUser
	filter.OwnerID = UserOwnerID(actor.ID)
	wallets, err := handler.service.ListWallets(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newWalletResponses(wallets), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientWallet(w http.ResponseWriter, r *http.Request, walletID WalletID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	wallet, err := handler.service.GetWallet(r.Context(), WalletLookup{
		ID: walletID, TenantID: tenantID, OwnerType: OwnerTypeUser, OwnerID: UserOwnerID(actor.ID),
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newWalletResponse(wallet))
}

func (handler *HTTPHandler) handleListClientLedger(w http.ResponseWriter, r *http.Request, walletID WalletID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	if _, err := handler.service.GetWallet(r.Context(), WalletLookup{
		ID: walletID, TenantID: tenantID, OwnerType: OwnerTypeUser, OwnerID: UserOwnerID(actor.ID),
	}); err != nil {
		writeWalletError(w, r, err)
		return
	}
	filter, page, ok := ledgerFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.WalletID = walletID
	entries, err := handler.service.ListLedgerEntries(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newLedgerEntryResponses(entries), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeWalletError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func walletPath(w http.ResponseWriter, r *http.Request, prefix string) (WalletID, string, bool) {
	value := strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
	parts := strings.Split(value, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeWalletError(w, r, ErrWalletIDMissing)
		return "", "", false
	}
	if len(parts) == 1 {
		return WalletID(parts[0]), "", true
	}
	if len(parts) == 2 && parts[1] == "ledger" {
		return WalletID(parts[0]), "ledger", true
	}
	return "", "invalid", true
}

func clientTenantActor(w http.ResponseWriter, r *http.Request) (tenant.ID, identity.Actor, bool) {
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return "", identity.Actor{}, false
	}
	actor, ok := actorFromContext(w, r)
	if !ok {
		return "", identity.Actor{}, false
	}
	return tenantID, actor, true
}

func tenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeWalletError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeWalletError(w, r, err)
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

func writeWalletError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := walletValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrWalletNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "wallet.not_found", "Wallet was not found.")
	case errors.Is(err, ErrLedgerEntryNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "wallet.ledger_not_found", "Wallet ledger entry was not found.")
	case errors.Is(err, ErrTopupRequestNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "wallet.topup_not_found", "Wallet top-up request was not found.")
	case errors.Is(err, ErrTopupStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.topup_status_conflict", "Top-up request status does not allow this action.")
	case errors.Is(err, ErrInsufficientBalance):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.insufficient_balance", "Wallet balance is insufficient.")
	case errors.Is(err, ErrWalletCurrencyMismatch):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.currency_mismatch", "Wallet currency does not match the requested currency.")
	case errors.Is(err, ErrWalletStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.status_conflict", "Wallet status does not allow this action.")
	case errors.Is(err, ErrIdempotencyConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.idempotency_conflict", "Idempotency key conflicts with another wallet ledger entry.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "wallet.service_unavailable", "Wallet service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "wallet.operation_failed", "Wallet operation failed.")
	}
}

func walletValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrWalletIDMissing):
		return validationField("wallet_id", "wallet.wallet_id_missing", "Wallet id is required."), true
	case errors.Is(err, ErrTopupRequestIDMissing):
		return validationField("topup_request_id", "wallet.topup_request_id_missing", "Top-up request id is required."), true
	case errors.Is(err, ErrOwnerTypeInvalid):
		return validationField("owner_type", "wallet.owner_type_invalid", "Wallet owner type is invalid."), true
	case errors.Is(err, ErrStatusInvalid):
		return validationField("status", "wallet.status_invalid", "Wallet status is invalid."), true
	case errors.Is(err, ErrDirectionInvalid):
		return validationField("direction", "wallet.direction_invalid", "Ledger direction is invalid."), true
	case errors.Is(err, ErrEntryTypeInvalid):
		return validationField("entry_type", "wallet.entry_type_invalid", "Ledger entry type is invalid."), true
	case errors.Is(err, ErrLedgerStatusInvalid):
		return validationField("status", "wallet.ledger_status_invalid", "Ledger status is invalid."), true
	case errors.Is(err, ErrTopupStatusInvalid):
		return validationField("status", "wallet.topup_status_invalid", "Top-up status is invalid."), true
	case errors.Is(err, ErrPaymentMethodInvalid):
		return validationField("payment_method", "wallet.payment_method_invalid", "Payment method is invalid."), true
	case errors.Is(err, ErrReviewNoteMissing):
		return validationField("review_note", "wallet.review_note_missing", "Review note is required."), true
	case errors.Is(err, ErrReasonMissing):
		return validationField("reason", "wallet.reason_missing", "Reason is required."), true
	case errors.Is(err, ErrAmountInvalid):
		return validationField("amount_minor", "wallet.amount_invalid", "Amount must be greater than zero."), true
	case errors.Is(err, ErrCurrencyMissing), errors.Is(err, ErrCurrencyInvalid):
		return validationField("currency", "wallet.currency_invalid", "Currency must be a three-letter code."), true
	case errors.Is(err, ErrReferenceTypeMissing):
		return validationField("reference_type", "wallet.reference_type_missing", "Reference type is required."), true
	case errors.Is(err, ErrReferenceIDMissing):
		return validationField("reference_id", "wallet.reference_id_missing", "Reference id is required."), true
	case errors.Is(err, ErrIdempotencyKeyMissing):
		return validationField("idempotency_key", "wallet.idempotency_key_missing", "Idempotency key is required."), true
	case errors.Is(err, ErrCorrelationIDMissing):
		return validationField("correlation_id", "wallet.correlation_id_missing", "Correlation id is required."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
