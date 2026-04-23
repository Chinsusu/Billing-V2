package order

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func (handler *HTTPHandler) adminServicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminServices, handler.options.AdminServiceMiddleware),
	})
}

func (handler *HTTPHandler) adminServiceRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminService, handler.options.AdminServiceMiddleware),
	})
}

func (handler *HTTPHandler) clientServicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListClientServices, handler.options.ClientServiceMiddleware),
	})
}

func (handler *HTTPHandler) clientServiceRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetClientService, handler.options.ClientServiceMiddleware),
	})
}

func (handler *HTTPHandler) handleListAdminServices(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := serviceInstanceFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	services, err := handler.service.ListServiceInstances(r.Context(), filter)
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newServiceInstanceResponses(services), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminService(w http.ResponseWriter, r *http.Request) {
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
	serviceID, ok := adminServiceIDFromPath(w, r)
	if !ok {
		return
	}
	service, err := handler.service.GetServiceInstance(r.Context(), ServiceInstanceLookup{
		ID:       serviceID,
		TenantID: tenantID,
	})
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponse(service))
}

func (handler *HTTPHandler) handleListClientServices(w http.ResponseWriter, r *http.Request) {
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
	filter, page, ok := serviceInstanceFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.BuyerUserID = actor.ID
	services, err := handler.service.ListServiceInstances(r.Context(), filter)
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newServiceInstanceResponses(services), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientService(w http.ResponseWriter, r *http.Request) {
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
	serviceID, ok := clientServiceIDFromPath(w, r)
	if !ok {
		return
	}
	service, err := handler.service.GetServiceInstance(r.Context(), ServiceInstanceLookup{
		ID:          serviceID,
		TenantID:    tenantID,
		BuyerUserID: actor.ID,
	})
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponse(service))
}

func serviceInstanceFilterFromRequest(w http.ResponseWriter, r *http.Request) (ServiceInstanceFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := ServiceInstanceFilter{Limit: page.Limit}
	query := r.URL.Query()
	buyerUserID := identity.UserID(strings.TrimSpace(query.Get("buyer_user_id")))
	if buyerUserID != "" {
		filter.BuyerUserID = buyerUserID
	}
	if displayID, present, ok := orderPositiveInt64Query(w, r, "display_id"); !ok {
		return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	orderID := OrderID(strings.TrimSpace(query.Get("order_id")))
	if orderID != "" {
		filter.OrderID = orderID
	}
	if orderDisplayID, present, ok := orderPositiveInt64Query(w, r, "order_display_id"); !ok {
		return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.OrderDisplayID = orderDisplayID
	}
	status := ServiceStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeOrderError(w, r, ErrServiceStatusInvalid)
			return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	return filter, page, true
}

func adminServiceIDFromPath(w http.ResponseWriter, r *http.Request) (ServiceID, bool) {
	return serviceIDFromPrefix(w, r, adminServicePrefix)
}

func clientServiceIDFromPath(w http.ResponseWriter, r *http.Request) (ServiceID, bool) {
	return serviceIDFromPrefix(w, r, clientServicePrefix)
}

func serviceIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (ServiceID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrServiceIDMissing)
		return "", false
	}
	return ServiceID(value), true
}
