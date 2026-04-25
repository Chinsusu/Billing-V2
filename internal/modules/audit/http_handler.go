package audit

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type HTTPService interface {
	ListLogs(ctx context.Context, filter Filter) ([]Log, error)
	GetLog(ctx context.Context, lookup Lookup) (Log, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	AdminMiddleware RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

const adminAuditLogPrefix = "/admin/audit-logs/"

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/audit-logs", handler.adminAuditLogsRoute)
	mux.HandleFunc("/admin/audit-logs/", handler.adminAuditLogRoute)
}

func (handler *HTTPHandler) adminAuditLogsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchAuditMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminAuditLogs, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminAuditLogRoute(w http.ResponseWriter, r *http.Request) {
	dispatchAuditMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminAuditLog, handler.options.AdminMiddleware),
	})
}

func dispatchAuditMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", auditAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func auditAllowHeader(methods map[string]http.HandlerFunc) string {
	allowed := make([]string, 0, len(methods))
	for method := range methods {
		allowed = append(allowed, method)
	}
	sort.Strings(allowed)
	return strings.Join(allowed, ", ")
}

func (handler *HTTPHandler) tenantRoute(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	return auditTenantContext(requireAuditTenantContext(applyAuditRouteMiddleware(next, routeMiddleware)))
}

func auditTenantContext(next http.HandlerFunc) http.HandlerFunc {
	handler := tenant.HeaderContextMiddleware(http.HandlerFunc(next))
	return handler.ServeHTTP
}

func requireAuditTenantContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := auditTenantIDFromContext(w, r); !ok {
			return
		}
		next(w, r)
	}
}

func applyAuditRouteMiddleware(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	if routeMiddleware == nil {
		return next
	}
	return routeMiddleware(next)
}

func (handler *HTTPHandler) handleListAdminAuditLogs(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := auditTenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := auditActorFromContext(w, r); !ok {
		return
	}
	filter, page, ok := auditFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	logs, err := handler.service.ListLogs(r.Context(), filter)
	if err != nil {
		writeAuditError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newLogSummaryResponses(logs), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminAuditLog(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := auditTenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := auditActorFromContext(w, r); !ok {
		return
	}
	logID, ok := auditLogIDFromPath(w, r)
	if !ok {
		return
	}
	record, err := handler.service.GetLog(r.Context(), Lookup{ID: logID, TenantID: tenantID})
	if err != nil {
		writeAuditError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newLogDetailResponse(record))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeAuditError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func auditFilterFromRequest(w http.ResponseWriter, r *http.Request) (Filter, httpserver.CursorPageRequest, bool) {
	page, ok := auditPageFromRequest(w, r)
	if !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	}
	filter := Filter{Limit: page.Limit}
	query := r.URL.Query()
	filter.ActorID = ActorID(strings.TrimSpace(query.Get("actor_id")))
	actorType := ActorType(strings.TrimSpace(query.Get("actor_type")))
	if actorType != "" {
		if !actorType.Valid() {
			writeAuditError(w, r, ErrActorTypeInvalid)
			return Filter{}, httpserver.CursorPageRequest{}, false
		}
		filter.ActorType = actorType
	}
	if displayID, present, ok := auditPositiveInt64Query(w, r, "display_id"); !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	if actorDisplayID, present, ok := auditPositiveInt64Query(w, r, "actor_display_id"); !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.ActorDisplayID = actorDisplayID
	}
	filter.Action = strings.TrimSpace(query.Get("action"))
	filter.TargetType = strings.TrimSpace(query.Get("target_type"))
	filter.TargetID = TargetID(strings.TrimSpace(query.Get("target_id")))
	if targetDisplayID, present, ok := auditPositiveInt64Query(w, r, "target_display_id"); !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.TargetDisplayID = targetDisplayID
	}
	createdFrom, ok := auditTimeFromRequest(w, r, "created_from")
	if !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	}
	createdTo, ok := auditTimeFromRequest(w, r, "created_to")
	if !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	}
	filter.CreatedFrom = createdFrom
	filter.CreatedTo = createdTo
	if !createdFrom.IsZero() && !createdTo.IsZero() && createdTo.Before(createdFrom) {
		writeAuditError(w, r, ErrCreatedWindowInvalid)
		return Filter{}, httpserver.CursorPageRequest{}, false
	}
	return filter, page, true
}

func auditTimeFromRequest(w http.ResponseWriter, r *http.Request, field string) (time.Time, bool) {
	value := strings.TrimSpace(r.URL.Query().Get(field))
	if value == "" {
		return time.Time{}, true
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			auditValidationField(field, "audit.created_time_invalid", "Created time must be RFC3339."),
		})
		return time.Time{}, false
	}
	return parsed, true
}

func auditLogIDFromPath(w http.ResponseWriter, r *http.Request) (ID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, adminAuditLogPrefix))
	if value == "" || strings.Contains(value, "/") {
		writeAuditError(w, r, ErrAuditLogNotFound)
		return "", false
	}
	return ID(value), true
}

func auditTenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeAuditError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func auditActorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeAuditError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func auditPageFromRequest(w http.ResponseWriter, r *http.Request) (httpserver.CursorPageRequest, bool) {
	page, err := httpserver.ParseCursorPage(r)
	if err == nil {
		return page, true
	}
	switch {
	case errors.Is(err, httpserver.ErrPageLimitTooLarge):
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{auditValidationField("limit", "request.limit_too_large", "Limit is too large.")})
	default:
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{auditValidationField("limit", "request.limit_invalid", "Limit must be a positive number.")})
	}
	return httpserver.CursorPageRequest{}, false
}

func writeAuditError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := auditErrorField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrAuditLogNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "audit.log_not_found", "Audit log was not found.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrAuditStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "audit.service_unavailable", "Audit service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit.operation_failed", "Audit operation failed.")
	}
}

func auditErrorField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return auditValidationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return auditValidationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrActorTypeInvalid):
		return auditValidationField("actor_type", "audit.actor_type_invalid", "Actor type is invalid."), true
	case errors.Is(err, ErrCreatedTimeInvalid):
		return auditValidationField("created_at", "audit.created_time_invalid", "Created time is invalid."), true
	case errors.Is(err, ErrCreatedWindowInvalid):
		return auditValidationField("created_at", "audit.created_time_window_invalid", "Created time window is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func auditValidationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
