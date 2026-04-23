package catalog

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const (
	adminCatalogProductPrefix        = "/admin/catalog/products/"
	adminCatalogPlanPrefix           = "/admin/catalog/plans/"
	adminCatalogProviderSourcePrefix = "/admin/catalog/provider-sources/"
)

func (handler *HTTPHandler) adminCatalogProductRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPatch: handler.adminRoute(handler.handleUpdateProductStatus),
	})
}

func (handler *HTTPHandler) adminCatalogPlanRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPatch: handler.adminRoute(handler.handleUpdatePlanStatus),
	})
}

func (handler *HTTPHandler) adminCatalogProviderSourceRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodPatch: handler.adminRoute(handler.handleUpdateProviderSourceStatus),
	})
}

func (handler *HTTPHandler) handleUpdateProductStatus(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	id, ok := productIDFromPath(w, r)
	if !ok {
		return
	}
	var request updateProductStatusRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	product, err := handler.service.UpdateProductStatus(r.Context(), request.toInput(id))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newProductResponse(product))
}

func (handler *HTTPHandler) handleUpdatePlanStatus(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	id, ok := planIDFromPath(w, r)
	if !ok {
		return
	}
	var request updatePlanStatusRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	plan, err := handler.service.UpdatePlanStatus(r.Context(), request.toInput(id))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newPlanResponse(plan))
}

func (handler *HTTPHandler) handleUpdateProviderSourceStatus(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	id, ok := providerSourceIDFromPath(w, r)
	if !ok {
		return
	}
	var request updateProviderSourceStatusRequest
	if !decodeCatalogJSON(w, r, &request) {
		return
	}
	source, err := handler.service.UpdateProviderSourceStatus(r.Context(), request.toInput(id))
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newProviderSourceResponse(source))
}

func productIDFromPath(w http.ResponseWriter, r *http.Request) (ProductID, bool) {
	value, ok := catalogPathValue(w, r, adminCatalogProductPrefix, "product_id", "catalog.product_id_missing", "Product id is required.")
	return ProductID(value), ok
}

func planIDFromPath(w http.ResponseWriter, r *http.Request) (PlanID, bool) {
	value, ok := catalogPathValue(w, r, adminCatalogPlanPrefix, "plan_id", "catalog.plan_id_missing", "Plan id is required.")
	return PlanID(value), ok
}

func providerSourceIDFromPath(w http.ResponseWriter, r *http.Request) (ProviderSourceID, bool) {
	value, ok := catalogPathValue(w, r, adminCatalogProviderSourcePrefix, "source_id", "catalog.source_id_missing", "Source id is required.")
	return ProviderSourceID(value), ok
}

func catalogPathValue(w http.ResponseWriter, r *http.Request, prefix string, field string, code string, message string) (string, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeFieldValidation(w, r, field, code, message)
		return "", false
	}
	return value, true
}
