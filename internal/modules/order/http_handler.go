package order

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/middleware"
)

const (
	TenantHeader         = tenant.HeaderTenantID
	IdempotencyKeyHeader = "Idempotency-Key"

	maxJSONBodyBytes = 1 << 20
)

type HTTPService interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	ClientMiddleware RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		options: options,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/client/orders", middleware.RequireMethod(http.MethodPost, handler.tenantRoute(handler.handleCreateClientOrder, handler.options.ClientMiddleware)))
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

func (handler *HTTPHandler) handleCreateClientOrder(w http.ResponseWriter, r *http.Request) {
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
	var request createOrderRequest
	if !decodeOrderJSON(w, r, &request) {
		return
	}
	order, err := handler.service.CreateOrder(r.Context(), request.toInput(tenantID, actor.ID, idempotencyKeyFromHeader(r)))
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newOrderResponse(order))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeOrderError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func decodeOrderJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
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
		writeOrderError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeOrderError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func idempotencyKeyFromHeader(r *http.Request) IdempotencyKey {
	return IdempotencyKey(strings.TrimSpace(r.Header.Get(IdempotencyKeyHeader)))
}

func writeOrderError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := orderValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "order.service_unavailable", "Order service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "order.operation_failed", "Order operation failed.")
	}
}

func orderValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrBuyerIDMissing):
		return validationField("actor_id", "order.buyer_missing", "Buyer actor is required."), true
	case errors.Is(err, ErrTenantPlanIDMissing):
		return validationField("tenant_plan_id", "order.tenant_plan_id_missing", "Tenant plan id is required."), true
	case errors.Is(err, ErrIdempotencyKeyMissing):
		return validationField("idempotency_key", "order.idempotency_key_missing", "Idempotency key is required."), true
	case errors.Is(err, ErrCurrencyMissing):
		return validationField("currency", "order.currency_missing", "Currency is required."), true
	case errors.Is(err, ErrCurrencyInvalid):
		return validationField("currency", "order.currency_invalid", "Currency is invalid."), true
	case errors.Is(err, ErrAmountInvalid):
		return validationField("amount_minor", "order.amount_invalid", "Money amount must not be negative."), true
	case errors.Is(err, ErrQuantityInvalid):
		return validationField("quantity", "order.quantity_invalid", "Quantity must be greater than zero."), true
	case errors.Is(err, ErrOrderStatusInvalid):
		return validationField("order_status", "order.status_invalid", "Order status is invalid."), true
	case errors.Is(err, ErrBillingStatusInvalid):
		return validationField("billing_status", "order.billing_status_invalid", "Billing status is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{
		Field:   field,
		Code:    code,
		Message: message,
	}
}
