package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type serviceInstanceResponse struct {
	ID                 ServiceID                   `json:"id"`
	DisplayID          int64                       `json:"display_id"`
	TenantID           tenant.ID                   `json:"tenant_id"`
	OrderID            OrderID                     `json:"order_id"`
	TenantPlanID       catalog.TenantPlanID        `json:"tenant_plan_id"`
	ProviderSourceID   catalog.ProviderSourceID    `json:"provider_source_id"`
	ExternalResourceID provider.ExternalResourceID `json:"external_resource_id"`
	Status             ServiceStatus               `json:"status"`
	BillingStatus      BillingStatus               `json:"billing_status"`
	SuspensionReason   SuspensionReason            `json:"suspension_reason,omitempty"`
	TermStart          time.Time                   `json:"term_start"`
	TermEnd            time.Time                   `json:"term_end"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

func newServiceInstanceResponse(service ServiceInstance) serviceInstanceResponse {
	return serviceInstanceResponse{
		ID:                 service.ID,
		DisplayID:          service.DisplayID,
		TenantID:           service.TenantID,
		OrderID:            service.OrderID,
		TenantPlanID:       service.TenantPlanID,
		ProviderSourceID:   service.ProviderSourceID,
		ExternalResourceID: service.ExternalResourceID,
		Status:             service.Status,
		BillingStatus:      service.BillingStatus,
		SuspensionReason:   service.SuspensionReason,
		TermStart:          service.TermStart,
		TermEnd:            service.TermEnd,
		CreatedAt:          service.CreatedAt,
		UpdatedAt:          service.UpdatedAt,
	}
}

func newServiceInstanceResponses(services []ServiceInstance) []serviceInstanceResponse {
	responses := make([]serviceInstanceResponse, 0, len(services))
	for _, service := range services {
		responses = append(responses, newServiceInstanceResponse(service))
	}
	return responses
}
