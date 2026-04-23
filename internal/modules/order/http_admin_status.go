package order

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const adminOrderStatusSuffix = "/status"

func isAdminOrderStatusPath(path string) bool {
	if !strings.HasPrefix(path, adminOrderPrefix) {
		return false
	}
	value := strings.TrimPrefix(path, adminOrderPrefix)
	if !strings.HasSuffix(value, adminOrderStatusSuffix) {
		return false
	}
	orderID := strings.TrimSuffix(value, adminOrderStatusSuffix)
	return orderID != "" && !strings.Contains(orderID, "/")
}

func (handler *HTTPHandler) handleTransitionAdminOrderStatus(w http.ResponseWriter, r *http.Request) {
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
	orderID, ok := adminOrderIDFromStatusPath(w, r)
	if !ok {
		return
	}
	var request transitionOrderStatusRequest
	if !decodeOrderJSON(w, r, &request) {
		return
	}
	order, err := handler.service.TransitionOrderStatus(r.Context(), request.toInput(orderID, tenantID))
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newOrderResponse(order))
}

func adminOrderIDFromStatusPath(w http.ResponseWriter, r *http.Request) (OrderID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, adminOrderPrefix))
	value = strings.TrimSuffix(value, adminOrderStatusSuffix)
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrOrderIDMissing)
		return "", false
	}
	return OrderID(value), true
}

func (request transitionOrderStatusRequest) toInput(orderID OrderID, tenantID tenant.ID) TransitionOrderStatusInput {
	return TransitionOrderStatusInput{
		ID:            orderID,
		TenantID:      tenantID,
		FromStatus:    request.FromStatus,
		ToStatus:      request.ToStatus,
		BillingStatus: request.BillingStatus,
	}
}
