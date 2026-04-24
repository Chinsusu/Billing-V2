package checkout

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const maxJSONBodyBytes = 1 << 20

type HTTPService interface {
	CheckoutOrder(ctx context.Context, input CheckoutOrderInput) (invoice.InvoiceDetail, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	ClientMiddleware RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

type checkoutOrderRequest struct {
	OrderID order.OrderID `json:"order_id"`
}

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/client/checkouts", handler.clientCheckoutsRoute)
}

func (handler *HTTPHandler) clientCheckoutsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCheckoutMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPost: handler.tenantRoute(handler.handleCreateClientCheckout, handler.options.ClientMiddleware),
	})
}

func dispatchCheckoutMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", checkoutAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func checkoutAllowHeader(methods map[string]http.HandlerFunc) string {
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

func (handler *HTTPHandler) handleCreateClientCheckout(w http.ResponseWriter, r *http.Request) {
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
	var request checkoutOrderRequest
	if !decodeCheckoutJSON(w, r, &request) {
		return
	}
	detail, err := handler.service.CheckoutOrder(r.Context(), CheckoutOrderInput{
		TenantID:       tenantID,
		BuyerUserID:    actor.ID,
		OrderID:        request.OrderID,
		IdempotencyKey: idempotencyKeyFromHeader(r),
	})
	if err != nil {
		writeCheckoutError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newInvoiceDetailResponse(detail))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeCheckoutError(w, r, invoice.ErrServiceStoreMissing)
		return false
	}
	return true
}

func decodeCheckoutJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
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

func tenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeCheckoutError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeCheckoutError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func idempotencyKeyFromHeader(r *http.Request) invoice.IdempotencyKey {
	return invoice.IdempotencyKey(strings.TrimSpace(r.Header.Get(IdempotencyKeyHeader)))
}

func writeCheckoutError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := checkoutValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, order.ErrOrderNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "order.not_found", "Order was not found.")
	case errors.Is(err, invoice.ErrOrderNotCheckoutable):
		httpserver.WriteError(w, r, http.StatusConflict, "checkout.order_not_checkoutable", "Order is not ready for checkout.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, invoice.ErrServiceStoreMissing), errors.Is(err, invoice.ErrStoreExecutorMissing), errors.Is(err, invoice.ErrOrderReaderMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "checkout.service_unavailable", "Checkout service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "checkout.operation_failed", "Checkout operation failed.")
	}
}

func checkoutValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, invoice.ErrBuyerIDMissing):
		return validationField("actor_id", "checkout.buyer_missing", "Buyer actor is required."), true
	case errors.Is(err, order.ErrOrderIDMissing):
		return validationField("order_id", "order.order_id_missing", "Order id is required."), true
	case errors.Is(err, invoice.ErrIdempotencyKeyMissing):
		return validationField("idempotency_key", "checkout.idempotency_key_missing", "Idempotency key is required."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
