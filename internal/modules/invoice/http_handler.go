package invoice

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type HTTPService interface {
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error)
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
	adminInvoicePrefix  = "/admin/invoices/"
	clientInvoicePrefix = "/client/invoices/"
)

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/invoices", handler.adminInvoicesRoute)
	mux.HandleFunc("/admin/invoices/", handler.adminInvoiceRoute)
	mux.HandleFunc("/reseller/invoices", handler.resellerInvoicesRoute)
	mux.HandleFunc("/client/invoices", handler.clientInvoicesRoute)
	mux.HandleFunc("/client/invoices/", handler.clientInvoiceRoute)
}

func (handler *HTTPHandler) adminInvoicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchInvoiceMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminInvoices, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminInvoiceRoute(w http.ResponseWriter, r *http.Request) {
	dispatchInvoiceMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminInvoice, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) resellerInvoicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchInvoiceMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminInvoices, handler.options.ResellerMiddleware),
	})
}

func (handler *HTTPHandler) clientInvoicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchInvoiceMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListClientInvoices, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) clientInvoiceRoute(w http.ResponseWriter, r *http.Request) {
	dispatchInvoiceMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetClientInvoice, handler.options.ClientMiddleware),
	})
}

func dispatchInvoiceMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", invoiceAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func invoiceAllowHeader(methods map[string]http.HandlerFunc) string {
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

func (handler *HTTPHandler) handleListAdminInvoices(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := invoiceFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	invoices, err := handler.service.ListInvoices(r.Context(), filter)
	if err != nil {
		writeInvoiceError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newInvoiceResponses(invoices), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminInvoice(w http.ResponseWriter, r *http.Request) {
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
	invoiceID, ok := adminInvoiceIDFromPath(w, r)
	if !ok {
		return
	}
	detail, err := handler.service.GetInvoice(r.Context(), InvoiceLookup{ID: invoiceID, TenantID: tenantID})
	if err != nil {
		writeInvoiceError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newInvoiceDetailResponse(detail))
}

func (handler *HTTPHandler) handleListClientInvoices(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	filter, page, ok := invoiceFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.BuyerUserID = actor.ID
	invoices, err := handler.service.ListInvoices(r.Context(), filter)
	if err != nil {
		writeInvoiceError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newInvoiceResponses(invoices), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientInvoice(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	invoiceID, ok := clientInvoiceIDFromPath(w, r)
	if !ok {
		return
	}
	detail, err := handler.service.GetInvoice(r.Context(), InvoiceLookup{
		ID: invoiceID, TenantID: tenantID, BuyerUserID: actor.ID,
	})
	if err != nil {
		writeInvoiceError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newInvoiceDetailResponse(detail))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeInvoiceError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func adminInvoiceIDFromPath(w http.ResponseWriter, r *http.Request) (InvoiceID, bool) {
	return invoiceIDFromPrefix(w, r, adminInvoicePrefix)
}

func clientInvoiceIDFromPath(w http.ResponseWriter, r *http.Request) (InvoiceID, bool) {
	return invoiceIDFromPrefix(w, r, clientInvoicePrefix)
}

func invoiceIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (InvoiceID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeInvoiceError(w, r, ErrInvoiceIDMissing)
		return "", false
	}
	return InvoiceID(value), true
}

func invoiceFilterFromRequest(w http.ResponseWriter, r *http.Request) (InvoiceFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := InvoiceFilter{Limit: page.Limit}
	query := r.URL.Query()
	if buyerUserID := identity.UserID(strings.TrimSpace(query.Get("buyer_user_id"))); buyerUserID != "" {
		filter.BuyerUserID = buyerUserID
	}
	if buyerDisplayID, present, ok := invoicePositiveInt64Query(w, r, "buyer_display_id"); !ok {
		return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.BuyerDisplayID = buyerDisplayID
	}
	if displayID, present, ok := invoicePositiveInt64Query(w, r, "display_id"); !ok {
		return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	if orderID := order.OrderID(strings.TrimSpace(query.Get("order_id"))); orderID != "" {
		filter.OrderID = orderID
	}
	if orderDisplayID, present, ok := invoicePositiveInt64Query(w, r, "order_display_id"); !ok {
		return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.OrderDisplayID = orderDisplayID
	}
	if status := Status(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeInvoiceError(w, r, ErrStatusInvalid)
			return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	amountMin, amountMax, ok := invoiceAmountRangeQuery(w, r)
	if !ok {
		return InvoiceFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.AmountMinMinor = amountMin
	filter.AmountMaxMinor = amountMax
	return filter, page, true
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
		writeInvoiceError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeInvoiceError(w, r, err)
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

func writeInvoiceError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := invoiceValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrInvoiceNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "invoice.not_found", "Invoice was not found.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "invoice.service_unavailable", "Invoice service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "invoice.operation_failed", "Invoice operation failed.")
	}
}

func invoiceValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrInvoiceIDMissing):
		return validationField("invoice_id", "invoice.invoice_id_missing", "Invoice id is required."), true
	case errors.Is(err, ErrBuyerIDMissing):
		return validationField("buyer_user_id", "invoice.buyer_missing", "Buyer is required."), true
	case errors.Is(err, ErrStatusInvalid):
		return validationField("status", "invoice.status_invalid", "Invoice status is invalid."), true
	case errors.Is(err, ErrAmountInvalid):
		return validationField("amount_minor", "invoice.amount_invalid", "Money amount must not be negative."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
