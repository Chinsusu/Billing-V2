package order

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const (
	serviceSuspendSuffix   = "/suspend"
	serviceUnsuspendSuffix = "/unsuspend"
	serviceTerminateSuffix = "/terminate"
)

type serviceLifecycleRequest struct {
	FromStatus       ServiceStatus    `json:"from_status"`
	BillingStatus    BillingStatus    `json:"billing_status"`
	SuspensionReason SuspensionReason `json:"suspension_reason"`
	Reason           string           `json:"reason"`
	NotifyClient     bool             `json:"notify_client"`
}

type serviceLifecycleRoute struct {
	Action                  ServiceLifecycleAction
	ToStatus                ServiceStatus
	DefaultBillingStatus    BillingStatus
	DefaultSuspensionReason SuspensionReason
}

func isServiceLifecyclePath(path string, prefix string, suffix string) bool {
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	value := strings.TrimPrefix(path, prefix)
	if !strings.HasSuffix(value, suffix) {
		return false
	}
	serviceID := strings.TrimSuffix(value, suffix)
	return serviceID != "" && !strings.Contains(serviceID, "/")
}

func serviceIDFromLifecyclePath(w http.ResponseWriter, r *http.Request, prefix string, suffix string) (ServiceID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	value = strings.TrimSuffix(value, suffix)
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrServiceIDMissing)
		return "", false
	}
	return ServiceID(value), true
}

func (handler *HTTPHandler) handleSuspendAdminService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, adminServicePrefix, serviceSuspendSuffix, serviceLifecycleRoute{
		Action:                  ServiceLifecycleActionSuspend,
		ToStatus:                ServiceStatusSuspended,
		DefaultSuspensionReason: SuspensionReasonManualAdmin,
	})
}

func (handler *HTTPHandler) handleSuspendResellerService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, resellerServicePrefix, serviceSuspendSuffix, serviceLifecycleRoute{
		Action:                  ServiceLifecycleActionSuspend,
		ToStatus:                ServiceStatusSuspended,
		DefaultSuspensionReason: SuspensionReasonManualReseller,
	})
}

func (handler *HTTPHandler) handleUnsuspendAdminService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, adminServicePrefix, serviceUnsuspendSuffix, serviceLifecycleRoute{
		Action:               ServiceLifecycleActionUnsuspend,
		ToStatus:             ServiceStatusActive,
		DefaultBillingStatus: BillingStatusPaid,
	})
}

func (handler *HTTPHandler) handleUnsuspendResellerService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, resellerServicePrefix, serviceUnsuspendSuffix, serviceLifecycleRoute{
		Action:               ServiceLifecycleActionUnsuspend,
		ToStatus:             ServiceStatusActive,
		DefaultBillingStatus: BillingStatusPaid,
	})
}

func (handler *HTTPHandler) handleTerminateAdminService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, adminServicePrefix, serviceTerminateSuffix, serviceLifecycleRoute{
		Action:   ServiceLifecycleActionTerminate,
		ToStatus: ServiceStatusTerminated,
	})
}

func (handler *HTTPHandler) handleTerminateResellerService(w http.ResponseWriter, r *http.Request) {
	handler.handleServiceLifecycle(w, r, resellerServicePrefix, serviceTerminateSuffix, serviceLifecycleRoute{
		Action:   ServiceLifecycleActionTerminate,
		ToStatus: ServiceStatusTerminated,
	})
}

func (handler *HTTPHandler) handleServiceLifecycle(w http.ResponseWriter, r *http.Request, prefix string, suffix string, route serviceLifecycleRoute) {
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
	serviceID, ok := serviceIDFromLifecyclePath(w, r, prefix, suffix)
	if !ok {
		return
	}
	var request serviceLifecycleRequest
	if !decodeOrderJSON(w, r, &request) {
		return
	}
	record, err := handler.service.TransitionServiceLifecycle(r.Context(), request.toInput(serviceID, tenantID, actor.ID, route))
	if err != nil {
		writeOrderError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newServiceInstanceResponse(record))
}

func (request serviceLifecycleRequest) toInput(serviceID ServiceID, tenantID tenant.ID, actorID identity.UserID, route serviceLifecycleRoute) TransitionServiceLifecycleInput {
	billingStatus := request.BillingStatus
	if billingStatus == "" {
		billingStatus = route.DefaultBillingStatus
	}
	suspensionReason := request.SuspensionReason
	if suspensionReason == "" {
		suspensionReason = route.DefaultSuspensionReason
	}
	return TransitionServiceLifecycleInput{
		ID:               serviceID,
		TenantID:         tenantID,
		ActorID:          audit.ActorID(actorID),
		ActorType:        audit.ActorTypeUser,
		Action:           route.Action,
		FromStatus:       request.FromStatus,
		ToStatus:         route.ToStatus,
		BillingStatus:    billingStatus,
		SuspensionReason: suspensionReason,
		Reason:           request.Reason,
	}
}
