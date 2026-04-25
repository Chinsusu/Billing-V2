package catalog

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func (handler *HTTPHandler) adminCatalogProductsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet:  handler.adminReadRoute(handler.handleListProducts),
		http.MethodPost: handler.adminManageRoute(handler.handleCreateProduct),
	})
}

func (handler *HTTPHandler) adminCatalogPlansRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet:  handler.adminReadRoute(handler.handleListAdminPlans),
		http.MethodPost: handler.adminManageRoute(handler.handleCreatePlan),
	})
}

func (handler *HTTPHandler) adminCatalogProviderSourcesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchCatalogMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet:  handler.adminReadRoute(handler.handleListProviderSources),
		http.MethodPost: handler.adminManageRoute(handler.handleCreateProviderSource),
	})
}

func dispatchCatalogMethods(w http.ResponseWriter, r *http.Request, methods map[string]http.HandlerFunc) {
	if handler, ok := methods[r.Method]; ok {
		handler(w, r)
		return
	}
	w.Header().Set("Allow", catalogAllowHeader(methods))
	httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
}

func catalogAllowHeader(methods map[string]http.HandlerFunc) string {
	allowed := make([]string, 0, len(methods))
	for method := range methods {
		allowed = append(allowed, method)
	}
	if len(allowed) == 2 {
		if allowed[0] > allowed[1] {
			allowed[0], allowed[1] = allowed[1], allowed[0]
		}
	}
	return strings.Join(allowed, ", ")
}

func (handler *HTTPHandler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := productFilterFromRequest(w, r)
	if !ok {
		return
	}
	products, err := handler.service.ListProducts(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newProductResponses(products), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleListAdminPlans(w http.ResponseWriter, r *http.Request) {
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

func (handler *HTTPHandler) handleListProviderSources(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := providerSourceFilterFromRequest(w, r)
	if !ok {
		return
	}
	sources, err := handler.service.ListProviderSources(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newProviderSourceResponses(sources), httpserver.NewPage(page.Limit, ""))
}

func productFilterFromRequest(w http.ResponseWriter, r *http.Request) (ProductFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return ProductFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := ProductFilter{Limit: page.Limit}
	query := r.URL.Query()
	productType := ProductType(strings.TrimSpace(query.Get("product_type")))
	if productType != "" {
		if !productType.Valid() {
			writeFieldValidation(w, r, "product_type", "catalog.product_type_invalid", "Product type is invalid.")
			return ProductFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Type = productType
	}
	status := ProductStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeFieldValidation(w, r, "status", "catalog.status_invalid", "Status is invalid.")
			return ProductFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	return filter, page, true
}

func providerSourceFilterFromRequest(w http.ResponseWriter, r *http.Request) (ProviderSourceFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return ProviderSourceFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := ProviderSourceFilter{Limit: page.Limit}
	query := r.URL.Query()
	if displayID, present, ok := catalogPositiveInt64Query(w, r, "display_id"); !ok {
		return ProviderSourceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	sourceType := provider.Type(strings.TrimSpace(query.Get("source_type")))
	if sourceType != "" {
		if !providerTypeValid(sourceType) {
			writeFieldValidation(w, r, "source_type", "catalog.source_type_invalid", "Source type is invalid.")
			return ProviderSourceFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Type = sourceType
	}
	status := ProviderSourceStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeFieldValidation(w, r, "status", "catalog.status_invalid", "Status is invalid.")
			return ProviderSourceFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	return filter, page, true
}
