package identity

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/middleware"
)

type AdminReadHTTPService interface {
	ListAdminTenants(ctx context.Context, filter tenant.ListTenantsFilter) ([]tenant.TenantSummary, error)
	ListAdminUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error)
}

type AdminReadRouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type AdminReadHTTPHandlerOptions struct {
	AdminMiddleware AdminReadRouteMiddleware
}

type AdminReadHTTPHandler struct {
	service AdminReadHTTPService
	options AdminReadHTTPHandlerOptions
}

func NewAdminReadHTTPHandler(service AdminReadHTTPService) *AdminReadHTTPHandler {
	return NewAdminReadHTTPHandlerWithOptions(service, AdminReadHTTPHandlerOptions{})
}

func NewAdminReadHTTPHandlerWithOptions(service AdminReadHTTPService, options AdminReadHTTPHandlerOptions) *AdminReadHTTPHandler {
	return &AdminReadHTTPHandler{service: service, options: options}
}

func (handler *AdminReadHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/tenants", middleware.RequireMethod(http.MethodGet, handler.adminTenantRoute(handler.handleListAdminTenants)))
	mux.HandleFunc("/admin/accounts", middleware.RequireMethod(http.MethodGet, handler.adminTenantRoute(func(w http.ResponseWriter, r *http.Request) {
		handler.handleListAdminUsers(w, r, "")
	})))
	mux.HandleFunc("/admin/customers", middleware.RequireMethod(http.MethodGet, handler.adminTenantRoute(func(w http.ResponseWriter, r *http.Request) {
		handler.handleListAdminUsers(w, r, UserTypeClient)
	})))
}

func (handler *AdminReadHTTPHandler) adminTenantRoute(next http.HandlerFunc) http.HandlerFunc {
	return identityTenantContext(requireIdentityTenantContext(applyAdminReadMiddleware(next, handler.options.AdminMiddleware)))
}

func identityTenantContext(next http.HandlerFunc) http.HandlerFunc {
	handler := tenant.HeaderContextMiddleware(http.HandlerFunc(next))
	return handler.ServeHTTP
}

func requireIdentityTenantContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := identityTenantIDFromContext(w, r); !ok {
			return
		}
		next(w, r)
	}
}

func applyAdminReadMiddleware(next http.HandlerFunc, routeMiddleware AdminReadRouteMiddleware) http.HandlerFunc {
	if routeMiddleware == nil {
		return next
	}
	return routeMiddleware(next)
}

func (handler *AdminReadHTTPHandler) handleListAdminTenants(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := adminTenantFilterFromRequest(w, r)
	if !ok {
		return
	}
	records, err := handler.service.ListAdminTenants(r.Context(), filter)
	if err != nil {
		writeAdminReadError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newAdminTenantResponses(records), httpserver.NewPage(page.Limit, ""))
}

func (handler *AdminReadHTTPHandler) handleListAdminUsers(w http.ResponseWriter, r *http.Request, defaultType UserType) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := adminUserFilterFromRequest(w, r, defaultType)
	if !ok {
		return
	}
	records, err := handler.service.ListAdminUsers(r.Context(), filter)
	if err != nil {
		writeAdminReadError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newAdminAccountResponses(records), httpserver.NewPage(page.Limit, ""))
}

func (handler *AdminReadHTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeAdminReadError(w, r, ErrAdminReadStoreMissing)
		return false
	}
	return true
}

func adminTenantFilterFromRequest(w http.ResponseWriter, r *http.Request) (tenant.ListTenantsFilter, httpserver.CursorPageRequest, bool) {
	tenantID, ok := identityTenantIDFromContext(w, r)
	if !ok {
		return tenant.ListTenantsFilter{}, httpserver.CursorPageRequest{}, false
	}
	page, ok := adminReadPageFromRequest(w, r)
	if !ok {
		return tenant.ListTenantsFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := tenant.ListTenantsFilter{
		ScopeTenantID: tenantID,
		Limit:         page.Limit,
	}
	query := r.URL.Query()
	if parentID := tenant.ID(strings.TrimSpace(query.Get("parent_tenant_id"))); !parentID.Empty() {
		filter.ParentID = parentID
	}
	if tenantType := tenant.Type(strings.TrimSpace(query.Get("type"))); tenantType != "" {
		if !tenantType.Valid() {
			writeAdminReadError(w, r, tenant.ErrTenantTypeInvalid)
			return tenant.ListTenantsFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Type = tenantType
	}
	if status := tenant.Status(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeAdminReadError(w, r, tenant.ErrTenantStatusInvalid)
			return tenant.ListTenantsFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	displayID, present, ok := adminReadDisplayIDFromRequest(w, r)
	if !ok {
		return tenant.ListTenantsFilter{}, httpserver.CursorPageRequest{}, false
	}
	if present {
		filter.DisplayID = displayID
	}
	return filter, page, true
}

func adminUserFilterFromRequest(w http.ResponseWriter, r *http.Request, defaultType UserType) (UserListFilter, httpserver.CursorPageRequest, bool) {
	tenantID, ok := identityTenantIDFromContext(w, r)
	if !ok {
		return UserListFilter{}, httpserver.CursorPageRequest{}, false
	}
	page, ok := adminReadPageFromRequest(w, r)
	if !ok {
		return UserListFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := UserListFilter{
		TenantID: tenantID,
		Type:     defaultType,
		Limit:    page.Limit,
		Email:    strings.TrimSpace(r.URL.Query().Get("email")),
	}
	query := r.URL.Query()
	if userType := UserType(strings.TrimSpace(query.Get("type"))); userType != "" {
		if !userType.Valid() {
			writeAdminReadError(w, r, ErrUserTypeInvalid)
			return UserListFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Type = userType
	}
	if filter.Type != "" && !filter.Type.Valid() {
		writeAdminReadError(w, r, ErrUserTypeInvalid)
		return UserListFilter{}, httpserver.CursorPageRequest{}, false
	}
	if status := UserStatus(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeAdminReadError(w, r, ErrUserStatusInvalid)
			return UserListFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	displayID, present, ok := adminReadDisplayIDFromRequest(w, r)
	if !ok {
		return UserListFilter{}, httpserver.CursorPageRequest{}, false
	}
	if present {
		filter.DisplayID = displayID
	}
	return filter, page, true
}

func adminReadPageFromRequest(w http.ResponseWriter, r *http.Request) (httpserver.CursorPageRequest, bool) {
	page, err := httpserver.ParseCursorPage(r)
	if err != nil {
		writeAdminReadError(w, r, err)
		return httpserver.CursorPageRequest{}, false
	}
	return page, true
}

func adminReadDisplayIDFromRequest(w http.ResponseWriter, r *http.Request) (int64, bool, bool) {
	value, present, err := httpserver.ParseOptionalPositiveInt64Query(r, "display_id")
	if err != nil {
		writeAdminReadError(w, r, err)
		return 0, false, false
	}
	return value, present, true
}

func identityTenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeAdminReadError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func writeAdminReadError(w http.ResponseWriter, r *http.Request, err error) {
	if field, ok := adminReadValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		httpserver.WriteError(w, r, http.StatusBadRequest, "tenant.context_missing", "Tenant context is required.")
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		httpserver.WriteError(w, r, http.StatusForbidden, "tenant.context_invalid", "Tenant context is invalid.")
	case errors.Is(err, ErrAdminReadStoreMissing),
		errors.Is(err, tenant.ErrStoreExecutorMissing),
		errors.Is(err, ErrUserStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "identity.admin_read_unavailable", "Account read service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "identity.admin_read_failed", "Account read operation failed.")
	}
}

func adminReadValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, httpserver.ErrPageLimitInvalid), errors.Is(err, httpserver.ErrPageLimitTooLarge):
		return adminReadValidationFieldValue("limit", "pagination.limit_invalid", "Limit is invalid."), true
	case errors.Is(err, httpserver.ErrQueryIntegerInvalid):
		return adminReadValidationFieldValue("display_id", "query.display_id_invalid", "Display id is invalid."), true
	case errors.Is(err, tenant.ErrTenantTypeInvalid):
		return adminReadValidationFieldValue("type", "tenant.type_invalid", "Tenant type is invalid."), true
	case errors.Is(err, tenant.ErrTenantStatusInvalid):
		return adminReadValidationFieldValue("status", "tenant.status_invalid", "Tenant status is invalid."), true
	case errors.Is(err, ErrUserTypeInvalid):
		return adminReadValidationFieldValue("type", "identity.user_type_invalid", "User type is invalid."), true
	case errors.Is(err, ErrUserStatusInvalid):
		return adminReadValidationFieldValue("status", "identity.user_status_invalid", "User status is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func adminReadValidationFieldValue(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
