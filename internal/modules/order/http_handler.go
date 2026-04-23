package order

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
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
	ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error)
	GetOrder(ctx context.Context, lookup OrderLookup) (Order, error)
	TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error)
	ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error)
	GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	AdminMiddleware         RouteMiddleware
	AdminManageMiddleware   RouteMiddleware
	AdminServiceMiddleware  RouteMiddleware
	ClientMiddleware        RouteMiddleware
	ClientServiceMiddleware RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

const (
	adminOrderPrefix    = "/admin/orders/"
	clientOrderPrefix   = "/client/orders/"
	adminServicePrefix  = "/admin/services/"
	clientServicePrefix = "/client/services/"
)

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
	mux.HandleFunc("/admin/orders", handler.adminOrdersRoute)
	mux.HandleFunc("/admin/orders/", handler.adminOrderRoute)
	mux.HandleFunc("/admin/services", handler.adminServicesRoute)
	mux.HandleFunc("/admin/services/", handler.adminServiceRoute)
	mux.HandleFunc("/client/orders", handler.clientOrdersRoute)
	mux.HandleFunc("/client/orders/", handler.clientOrderRoute)
	mux.HandleFunc("/client/services", handler.clientServicesRoute)
	mux.HandleFunc("/client/services/", handler.clientServiceRoute)
}

func (handler *HTTPHandler) adminOrdersRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminOrders, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminOrderRoute(w http.ResponseWriter, r *http.Request) {
	if isAdminOrderStatusPath(r.URL.Path) {
		dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPatch: handler.tenantRoute(handler.handleTransitionAdminOrderStatus, handler.options.AdminManageMiddleware),
		})
		return
	}
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminOrder, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) clientOrdersRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet:  handler.tenantRoute(handler.handleListClientOrders, handler.options.ClientMiddleware),
		http.MethodPost: middleware.RequireMethod(http.MethodPost, handler.tenantRoute(handler.handleCreateClientOrder, handler.options.ClientMiddleware)),
	})
}

func (handler *HTTPHandler) clientOrderRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetClientOrder, handler.options.ClientMiddleware),
	})
}

func dispatchOrderMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", orderAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func orderAllowHeader(methods map[string]http.HandlerFunc) string {
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

func (handler *HTTPHandler) handleListAdminOrders(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := orderFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	orders, err := handler.service.ListOrders(r.Context(), filter)
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newOrderResponses(orders), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminOrder(w http.ResponseWriter, r *http.Request) {
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
	orderID, ok := adminOrderIDFromPath(w, r)
	if !ok {
		return
	}
	order, err := handler.service.GetOrder(r.Context(), OrderLookup{
		ID:       orderID,
		TenantID: tenantID,
	})
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newOrderResponse(order))
}

func (handler *HTTPHandler) handleListClientOrders(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := orderFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.BuyerUserID = actor.ID
	orders, err := handler.service.ListOrders(r.Context(), filter)
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newOrderResponses(orders), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientOrder(w http.ResponseWriter, r *http.Request) {
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
	orderID, ok := orderIDFromPath(w, r)
	if !ok {
		return
	}
	order, err := handler.service.GetOrder(r.Context(), OrderLookup{
		ID:          orderID,
		TenantID:    tenantID,
		BuyerUserID: actor.ID,
	})
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newOrderResponse(order))
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

func adminOrderIDFromPath(w http.ResponseWriter, r *http.Request) (OrderID, bool) {
	return orderIDFromPrefix(w, r, adminOrderPrefix)
}

func orderIDFromPath(w http.ResponseWriter, r *http.Request) (OrderID, bool) {
	return orderIDFromPrefix(w, r, clientOrderPrefix)
}

func orderIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (OrderID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrOrderIDMissing)
		return "", false
	}
	return OrderID(value), true
}

func orderFilterFromRequest(w http.ResponseWriter, r *http.Request) (OrderFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return OrderFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := OrderFilter{Limit: page.Limit}
	query := r.URL.Query()
	buyerUserID := identity.UserID(strings.TrimSpace(query.Get("buyer_user_id")))
	if buyerUserID != "" {
		filter.BuyerUserID = buyerUserID
	}
	orderStatus := OrderStatus(strings.TrimSpace(query.Get("status")))
	if orderStatus != "" {
		if !orderStatus.Valid() {
			writeOrderError(w, r, ErrOrderStatusInvalid)
			return OrderFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.OrderStatus = orderStatus
	}
	billingStatus := BillingStatus(strings.TrimSpace(query.Get("billing_status")))
	if billingStatus != "" {
		if !billingStatus.Valid() {
			writeOrderError(w, r, ErrBillingStatusInvalid)
			return OrderFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.BillingStatus = billingStatus
	}
	return filter, page, true
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

func writeOrderError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := orderValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrOrderNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "order.not_found", "Order was not found.")
	case errors.Is(err, ErrServiceNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "service.not_found", "Service instance was not found.")
	case errors.Is(err, ErrOrderStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "order.status_conflict", "Order status no longer matches the expected value.")
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
	case errors.Is(err, ErrOrderIDMissing):
		return validationField("order_id", "order.order_id_missing", "Order id is required."), true
	case errors.Is(err, ErrServiceIDMissing):
		return validationField("service_id", "service.service_id_missing", "Service id is required."), true
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
	case errors.Is(err, ErrServiceStatusInvalid):
		return validationField("status", "service.status_invalid", "Service status is invalid."), true
	case errors.Is(err, ErrStatusTransitionInvalid):
		return validationField("to_status", "order.status_transition_invalid", "Order status change is not allowed."), true
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
