package jobs

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type HTTPService interface {
	ListJobs(ctx context.Context, filter Filter) ([]Job, error)
	GetJob(ctx context.Context, lookup Lookup) (Job, error)
}

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type HTTPHandlerOptions struct {
	AdminMiddleware    RouteMiddleware
	ResellerMiddleware RouteMiddleware
}

type HTTPHandler struct {
	service HTTPService
	options HTTPHandlerOptions
}

const (
	adminJobPrefix    = "/admin/jobs/"
	resellerJobPrefix = "/reseller/jobs/"
)

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{})
}

func NewHTTPHandlerWithOptions(service HTTPService, options HTTPHandlerOptions) *HTTPHandler {
	return &HTTPHandler{service: service, options: options}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/jobs", handler.adminJobsRoute)
	mux.HandleFunc("/admin/jobs/", handler.adminJobRoute)
	mux.HandleFunc("/reseller/jobs", handler.resellerJobsRoute)
	mux.HandleFunc("/reseller/jobs/", handler.resellerJobRoute)
}

func (handler *HTTPHandler) adminJobsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchJobMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListJobs, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) adminJobRoute(w http.ResponseWriter, r *http.Request) {
	dispatchJobMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminJob, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) resellerJobsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchJobMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListJobs, handler.options.ResellerMiddleware),
	})
}

func (handler *HTTPHandler) resellerJobRoute(w http.ResponseWriter, r *http.Request) {
	dispatchJobMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetResellerJob, handler.options.ResellerMiddleware),
	})
}

func dispatchJobMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", jobAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func jobAllowHeader(methods map[string]http.HandlerFunc) string {
	allowed := make([]string, 0, len(methods))
	for method := range methods {
		allowed = append(allowed, method)
	}
	sort.Strings(allowed)
	return strings.Join(allowed, ", ")
}

func (handler *HTTPHandler) tenantRoute(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	return jobTenantContext(requireJobTenantContext(applyJobRouteMiddleware(next, routeMiddleware)))
}

func jobTenantContext(next http.HandlerFunc) http.HandlerFunc {
	handler := tenant.HeaderContextMiddleware(http.HandlerFunc(next))
	return handler.ServeHTTP
}

func requireJobTenantContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := jobTenantIDFromContext(w, r); !ok {
			return
		}
		next(w, r)
	}
}

func applyJobRouteMiddleware(next http.HandlerFunc, routeMiddleware RouteMiddleware) http.HandlerFunc {
	if routeMiddleware == nil {
		return next
	}
	return routeMiddleware(next)
}

func (handler *HTTPHandler) handleListJobs(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := jobTenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := jobActorFromContext(w, r); !ok {
		return
	}
	filter, page, ok := jobFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	jobs, err := handler.service.ListJobs(r.Context(), filter)
	if err != nil {
		writeJobError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newJobResponses(jobs), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminJob(w http.ResponseWriter, r *http.Request) {
	handler.handleGetJob(w, r, adminJobPrefix)
}

func (handler *HTTPHandler) handleGetResellerJob(w http.ResponseWriter, r *http.Request) {
	handler.handleGetJob(w, r, resellerJobPrefix)
}

func (handler *HTTPHandler) handleGetJob(w http.ResponseWriter, r *http.Request, prefix string) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := jobTenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := jobActorFromContext(w, r); !ok {
		return
	}
	jobID, ok := jobIDFromPrefix(w, r, prefix)
	if !ok {
		return
	}
	job, err := handler.service.GetJob(r.Context(), Lookup{ID: jobID, TenantID: tenantID})
	if err != nil {
		writeJobError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newJobResponse(job))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeJobError(w, r, ErrServiceStoreMissing)
		return false
	}
	return true
}

func jobFilterFromRequest(w http.ResponseWriter, r *http.Request) (Filter, httpserver.CursorPageRequest, bool) {
	page, ok := jobPageFromRequest(w, r)
	if !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	}
	filter := Filter{Limit: page.Limit}
	query := r.URL.Query()
	if displayID, present, ok := jobPositiveInt64Query(w, r, "display_id"); !ok {
		return Filter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	filter.Type = Type(strings.TrimSpace(query.Get("job_type")))
	status := Status(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeJobError(w, r, ErrStatusInvalid)
			return Filter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	filter.ReferenceType = ReferenceType(strings.TrimSpace(query.Get("reference_type")))
	filter.ReferenceID = ReferenceID(strings.TrimSpace(query.Get("reference_id")))
	filter.SourceID = SourceID(strings.TrimSpace(query.Get("source_id")))
	return filter, page, true
}

func jobIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (ID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeJobError(w, r, ErrJobIDMissing)
		return "", false
	}
	return ID(value), true
}

func jobTenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeJobError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func jobActorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeJobError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func jobPageFromRequest(w http.ResponseWriter, r *http.Request) (httpserver.CursorPageRequest, bool) {
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

func writeJobError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := jobValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrJobNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "job.not_found", "Job was not found.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing), errors.Is(err, ErrStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "job.service_unavailable", "Job service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "job.operation_failed", "Job operation failed.")
	}
}

func jobValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrJobIDMissing):
		return validationField("job_id", "job.job_id_missing", "Job id is required."), true
	case errors.Is(err, ErrStatusInvalid):
		return validationField("status", "job.status_invalid", "Job status is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{Field: field, Code: code, Message: message}
}
