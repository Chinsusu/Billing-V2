package order

import (
	"net"
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
	if isServiceCredentialRevealPath(r.URL.Path, adminServicePrefix) {
		dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleRevealAdminServiceCredential, handler.options.AdminCredentialMiddleware),
		})
		return
	}
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminService, handler.options.AdminServiceMiddleware),
	})
}

func (handler *HTTPHandler) resellerServicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminServices, handler.options.ResellerServiceMiddleware),
	})
}

func (handler *HTTPHandler) resellerServiceRoute(w http.ResponseWriter, r *http.Request) {
	if isServiceCredentialRevealPath(r.URL.Path, resellerServicePrefix) {
		dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleRevealResellerServiceCredential, handler.options.ResellerCredentialMiddleware),
		})
		return
	}
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetResellerService, handler.options.ResellerServiceMiddleware),
	})
}

func (handler *HTTPHandler) clientServicesRoute(w http.ResponseWriter, r *http.Request) {
	dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListClientServices, handler.options.ClientServiceMiddleware),
	})
}

func (handler *HTTPHandler) clientServiceRoute(w http.ResponseWriter, r *http.Request) {
	if isServiceCredentialRevealPath(r.URL.Path, clientServicePrefix) {
		dispatchOrderMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleRevealClientServiceCredential, handler.options.ClientCredentialMiddleware),
		})
		return
	}
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
	credentials, ok := handler.serviceCredentialsForResponse(w, r, service)
	if !ok {
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponseWithCredentials(service, credentials))
}

func (handler *HTTPHandler) handleGetResellerService(w http.ResponseWriter, r *http.Request) {
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
	serviceID, ok := resellerServiceIDFromPath(w, r)
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
	credentials, ok := handler.serviceCredentialsForResponse(w, r, service)
	if !ok {
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponseWithCredentials(service, credentials))
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
	credentials, ok := handler.serviceCredentialsForResponse(w, r, service)
	if !ok {
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponseWithCredentials(service, credentials))
}

func (handler *HTTPHandler) handleRevealAdminServiceCredential(w http.ResponseWriter, r *http.Request) {
	handler.handleRevealServiceCredential(w, r, adminServicePrefix, false)
}

func (handler *HTTPHandler) handleRevealResellerServiceCredential(w http.ResponseWriter, r *http.Request) {
	handler.handleRevealServiceCredential(w, r, resellerServicePrefix, false)
}

func (handler *HTTPHandler) handleRevealClientServiceCredential(w http.ResponseWriter, r *http.Request) {
	handler.handleRevealServiceCredential(w, r, clientServicePrefix, true)
}

func (handler *HTTPHandler) handleRevealServiceCredential(w http.ResponseWriter, r *http.Request, prefix string, clientOwned bool) {
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
	serviceID, credentialID, ok := serviceCredentialRevealIDsFromPath(w, r, prefix)
	if !ok {
		return
	}
	var request credentialRevealRequest
	if !decodeOptionalOrderJSON(w, r, &request) {
		return
	}
	revealer, ok := handler.service.(HTTPCredentialRevealer)
	if !ok {
		writeOrderError(w, r, ErrCredentialStoreMissing)
		return
	}
	input := request.toInput(tenantID, serviceID, credentialID, actor.ID, clientIPFromRequest(r), r.UserAgent())
	if clientOwned {
		input.BuyerUserID = actor.ID
	}
	result, err := revealer.RevealServiceCredential(r.Context(), input)
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	httpserver.WriteSuccess(w, r, http.StatusOK, newCredentialRevealResponse(result))
}

func (handler *HTTPHandler) serviceCredentialsForResponse(w http.ResponseWriter, r *http.Request, service ServiceInstance) ([]ServiceCredential, bool) {
	reader, ok := handler.service.(HTTPCredentialReader)
	if !ok {
		return nil, true
	}
	credentials, err := reader.ListServiceCredentials(r.Context(), ServiceCredentialFilter{
		TenantID:  service.TenantID,
		ServiceID: service.ID,
		Status:    CredentialStatusActive,
	})
	if err != nil {
		writeOrderError(w, r, err)
		return nil, false
	}
	return credentials, true
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
	if buyerDisplayID, present, ok := orderPositiveInt64Query(w, r, "buyer_display_id"); !ok {
		return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.BuyerDisplayID = buyerDisplayID
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
	if sourceDisplayID, present, ok := orderPositiveInt64Query(w, r, "provider_source_display_id"); !ok {
		return ServiceInstanceFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.ProviderSourceDisplayID = sourceDisplayID
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

func resellerServiceIDFromPath(w http.ResponseWriter, r *http.Request) (ServiceID, bool) {
	return serviceIDFromPrefix(w, r, resellerServicePrefix)
}

func serviceIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (ServiceID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrServiceIDMissing)
		return "", false
	}
	return ServiceID(value), true
}

func isServiceCredentialRevealPath(path string, prefix string) bool {
	_, _, ok := serviceCredentialRevealIDs(path, prefix)
	return ok
}

func serviceCredentialRevealIDsFromPath(w http.ResponseWriter, r *http.Request, prefix string) (ServiceID, CredentialID, bool) {
	serviceID, credentialID, ok := serviceCredentialRevealIDs(r.URL.Path, prefix)
	if !ok {
		writeOrderError(w, r, ErrCredentialIDMissing)
		return "", "", false
	}
	return serviceID, credentialID, true
}

func serviceCredentialRevealIDs(path string, prefix string) (ServiceID, CredentialID, bool) {
	value := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.Split(value, "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] != "credentials" || parts[2] == "" || parts[3] != "reveal" {
		return "", "", false
	}
	return ServiceID(parts[0]), CredentialID(parts[2]), true
}

func clientIPFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		if first, _, ok := strings.Cut(forwardedFor, ","); ok {
			return strings.TrimSpace(first)
		}
		return forwardedFor
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
