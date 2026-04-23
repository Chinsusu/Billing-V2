package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/middleware"
)

const (
	TenantHeader = "X-Tenant-Id"
	ActorHeader  = "X-Actor-Id"

	maxJSONBodyBytes = 1 << 20
)

type HTTPService interface {
	CreateProduct(ctx context.Context, input CreateProductInput) (Product, error)
	CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error)
	CreateProviderSource(ctx context.Context, input CreateProviderSourceInput) (ProviderSource, error)
	CreatePlanSource(ctx context.Context, input CreatePlanSourceInput) (PlanSource, error)
	CloneTenantProduct(ctx context.Context, input CreateTenantProductInput) (TenantProduct, error)
	CloneTenantPlan(ctx context.Context, input CreateTenantPlanInput) (TenantPlan, error)
	ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error)
	ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error)
}

type HTTPHandler struct {
	service HTTPService
}

func NewHTTPHandler(service HTTPService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/catalog/products", middleware.RequireMethod(http.MethodPost, handler.handleCreateProduct))
	mux.HandleFunc("/admin/catalog/plans", middleware.RequireMethod(http.MethodPost, handler.handleCreatePlan))
	mux.HandleFunc("/admin/catalog/provider-sources", middleware.RequireMethod(http.MethodPost, handler.handleCreateProviderSource))
	mux.HandleFunc("/admin/catalog/plan-sources", middleware.RequireMethod(http.MethodPost, handler.handleCreatePlanSource))
	mux.HandleFunc("/reseller/catalog/master-plans", middleware.RequireMethod(http.MethodGet, handler.handleListMasterPlans))
	mux.HandleFunc("/reseller/catalog/products/clone", middleware.RequireMethod(http.MethodPost, handler.handleCloneTenantProduct))
	mux.HandleFunc("/reseller/catalog/plans/clone", middleware.RequireMethod(http.MethodPost, handler.handleCloneTenantPlan))
	mux.HandleFunc("/reseller/catalog", middleware.RequireMethod(http.MethodGet, handler.handleListTenantCatalog))
	mux.HandleFunc("/client/catalog", middleware.RequireMethod(http.MethodGet, handler.handleListClientCatalog))
}

func (handler *HTTPHandler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request createProductRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	product, err := handler.service.CreateProduct(r.Context(), request.toInput(actorIDFromHeader(r)))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newProductResponse(product))
}

func (handler *HTTPHandler) handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request createPlanRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	plan, err := handler.service.CreatePlan(r.Context(), request.toInput())
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newPlanResponse(plan))
}

func (handler *HTTPHandler) handleCreateProviderSource(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request createProviderSourceRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	source, err := handler.service.CreateProviderSource(r.Context(), request.toInput())
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newProviderSourceResponse(source))
}

func (handler *HTTPHandler) handleCreatePlanSource(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request createPlanSourceRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	source, err := handler.service.CreatePlanSource(r.Context(), request.toInput())
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newPlanSourceResponse(source))
}

func (handler *HTTPHandler) handleListMasterPlans(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := masterPlanFilterFromRequest(w, r)
	if !ok {
		return
	}
	plans, err := handler.service.ListMasterPlans(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newPlanResponses(plans), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleCloneTenantProduct(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromHeader(w, r)
	if !ok {
		return
	}
	var request cloneTenantProductRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	product, err := handler.service.CloneTenantProduct(r.Context(), request.toInput(tenantID))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newTenantProductResponse(product))
}

func (handler *HTTPHandler) handleCloneTenantPlan(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromHeader(w, r)
	if !ok {
		return
	}
	var request cloneTenantPlanRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	plan, err := handler.service.CloneTenantPlan(r.Context(), request.toInput(tenantID))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newTenantPlanResponse(plan))
}

func (handler *HTTPHandler) handleListTenantCatalog(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := tenantCatalogFilterFromRequest(w, r)
	if !ok {
		return
	}
	catalog, err := handler.service.ListTenantCatalog(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTenantCatalogResponse(catalog), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleListClientCatalog(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := tenantCatalogFilterFromRequest(w, r)
	if !ok {
		return
	}
	catalog, err := handler.service.ListTenantCatalog(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTenantCatalogPublicResponse(catalog), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeCatalogError(w, r, ErrCatalogServiceStoreMissing)
		return false
	}
	return true
}

func decodeCatalogJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
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

func actorIDFromHeader(r *http.Request) UserID {
	return UserID(strings.TrimSpace(r.Header.Get(ActorHeader)))
}

func tenantIDFromHeader(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantID := tenant.ID(strings.TrimSpace(r.Header.Get(TenantHeader)))
	if tenantID.Empty() {
		writeCatalogError(w, r, tenant.ErrTenantIDMissing)
		return "", false
	}
	return tenantID, true
}

func masterPlanFilterFromRequest(w http.ResponseWriter, r *http.Request) (MasterPlanFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return MasterPlanFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := MasterPlanFilter{Limit: page.Limit}
	query := r.URL.Query()
	productType := ProductType(strings.TrimSpace(query.Get("product_type")))
	if productType != "" {
		if !productType.Valid() {
			writeFieldValidation(w, r, "product_type", "catalog.product_type_invalid", "Product type is invalid.")
			return MasterPlanFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.ProductType = productType
	}
	status := PlanStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeFieldValidation(w, r, "status", "catalog.status_invalid", "Status is invalid.")
			return MasterPlanFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	return filter, page, true
}

func tenantCatalogFilterFromRequest(w http.ResponseWriter, r *http.Request) (TenantCatalogFilter, httpserver.CursorPageRequest, bool) {
	tenantID, ok := tenantIDFromHeader(w, r)
	if !ok {
		return TenantCatalogFilter{}, httpserver.CursorPageRequest{}, false
	}
	page, ok := pageFromRequest(w, r)
	if !ok {
		return TenantCatalogFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := TenantCatalogFilter{
		TenantID: tenantID,
		Limit:    page.Limit,
	}
	query := r.URL.Query()
	productType := ProductType(strings.TrimSpace(query.Get("product_type")))
	if productType != "" {
		if !productType.Valid() {
			writeFieldValidation(w, r, "product_type", "catalog.product_type_invalid", "Product type is invalid.")
			return TenantCatalogFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.ProductType = productType
	}
	visibility := TenantPlanVisibility(strings.TrimSpace(query.Get("visibility")))
	if visibility != "" {
		if !visibility.Valid() {
			writeFieldValidation(w, r, "visibility", "catalog.visibility_invalid", "Visibility is invalid.")
			return TenantCatalogFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Visibility = visibility
	}
	status := TenantPlanStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeFieldValidation(w, r, "status", "catalog.status_invalid", "Status is invalid.")
			return TenantCatalogFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
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
		writeFieldValidation(w, r, "limit", "request.limit_too_large", "Limit is too large.")
	default:
		writeFieldValidation(w, r, "limit", "request.limit_invalid", "Limit must be a positive number.")
	}
	return httpserver.CursorPageRequest{}, false
}

func writeCatalogError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := catalogValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	if catalogNotFound(err) {
		httpserver.WriteError(w, r, http.StatusNotFound, "catalog.not_found", "Catalog record was not found.")
		return
	}
	if errors.Is(err, ErrCatalogServiceStoreMissing) || errors.Is(err, ErrCatalogStoreExecutorMissing) {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "catalog.service_unavailable", "Catalog service is unavailable.")
		return
	}
	httpserver.WriteError(w, r, http.StatusInternalServerError, "catalog.operation_failed", "Catalog operation failed.")
}

func writeFieldValidation(w http.ResponseWriter, r *http.Request, field string, code string, message string) {
	httpserver.WriteValidationError(w, r, []httpserver.ValidationField{{
		Field:   field,
		Code:    code,
		Message: message,
	}})
}

func catalogValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrProductIDMissing):
		return validationField("product_id", "catalog.product_id_missing", "Product id is required."), true
	case errors.Is(err, ErrProductTypeInvalid):
		return validationField("product_type", "catalog.product_type_invalid", "Product type is invalid."), true
	case errors.Is(err, ErrProductNameMissing):
		return validationField("name", "catalog.name_missing", "Name is required."), true
	case errors.Is(err, ErrProductStatusInvalid), errors.Is(err, ErrPlanStatusInvalid), errors.Is(err, ErrSourceStatusInvalid),
		errors.Is(err, ErrPlanSourceStatus), errors.Is(err, ErrTenantProductStatus), errors.Is(err, ErrTenantPlanStatus):
		return validationField("status", "catalog.status_invalid", "Status is invalid."), true
	case errors.Is(err, ErrPlanIDMissing):
		return validationField("plan_id", "catalog.plan_id_missing", "Plan id is required."), true
	case errors.Is(err, ErrPlanCodeMissing):
		return validationField("plan_code", "catalog.plan_code_missing", "Plan code is required."), true
	case errors.Is(err, ErrPlanNameMissing):
		return validationField("name", "catalog.name_missing", "Name is required."), true
	case errors.Is(err, ErrBillingCycleInvalid):
		return validationField("billing_cycle_type", "catalog.billing_cycle_invalid", "Billing cycle is invalid."), true
	case errors.Is(err, ErrBillingCycleValue):
		return validationField("billing_cycle_value", "catalog.billing_cycle_value_invalid", "Billing cycle value is invalid."), true
	case errors.Is(err, ErrCurrencyMissing):
		return validationField("currency", "catalog.currency_missing", "Currency is required."), true
	case errors.Is(err, ErrCurrencyInvalid):
		return validationField("currency", "catalog.currency_invalid", "Currency is invalid."), true
	case errors.Is(err, ErrMoneyAmountInvalid):
		return validationField("amount_minor", "catalog.amount_invalid", "Money amount must not be negative."), true
	case errors.Is(err, ErrVersionInvalid):
		return validationField("version", "catalog.version_invalid", "Version must be greater than zero."), true
	case errors.Is(err, ErrSourceIDMissing):
		return validationField("source_id", "catalog.source_id_missing", "Source id is required."), true
	case errors.Is(err, ErrSourceTypeInvalid):
		return validationField("source_type", "catalog.source_type_invalid", "Source type is invalid."), true
	case errors.Is(err, ErrSourceNameMissing):
		return validationField("name", "catalog.name_missing", "Name is required."), true
	case errors.Is(err, ErrInventoryModeInvalid):
		return validationField("inventory_mode", "catalog.inventory_mode_invalid", "Inventory mode is invalid."), true
	case errors.Is(err, ErrRiskLevelInvalid):
		return validationField("risk_level", "catalog.risk_level_invalid", "Risk level is invalid."), true
	case errors.Is(err, ErrPlanSourcePriority):
		return validationField("priority", "catalog.priority_invalid", "Priority must be greater than zero."), true
	case errors.Is(err, ErrTenantProductIDMissing):
		return validationField("tenant_product_id", "catalog.tenant_product_id_missing", "Tenant product id is required."), true
	case errors.Is(err, ErrTenantPlanIDMissing):
		return validationField("tenant_plan_id", "catalog.tenant_plan_id_missing", "Tenant plan id is required."), true
	case errors.Is(err, ErrTenantPlanVisibility):
		return validationField("visibility", "catalog.visibility_invalid", "Visibility is invalid."), true
	case errors.Is(err, ErrCreatedByMissing):
		return validationField("actor_id", "catalog.actor_missing", "Actor id is required."), true
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

func catalogNotFound(err error) bool {
	return errors.Is(err, ErrProductNotFound) ||
		errors.Is(err, ErrPlanNotFound) ||
		errors.Is(err, ErrProviderSourceNotFound) ||
		errors.Is(err, ErrPlanSourceNotFound) ||
		errors.Is(err, ErrTenantProductNotFound) ||
		errors.Is(err, ErrTenantPlanNotFound)
}
