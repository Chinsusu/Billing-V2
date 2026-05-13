package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type serviceInstanceResponse struct {
	ID                      ServiceID                   `json:"id"`
	DisplayID               int64                       `json:"display_id"`
	TenantID                tenant.ID                   `json:"tenant_id"`
	OrderID                 OrderID                     `json:"order_id"`
	OrderDisplayID          int64                       `json:"order_display_id,omitempty"`
	BuyerDisplayID          int64                       `json:"buyer_display_id,omitempty"`
	TenantPlanID            catalog.TenantPlanID        `json:"tenant_plan_id"`
	ProviderSourceID        catalog.ProviderSourceID    `json:"provider_source_id"`
	ProviderSourceDisplayID int64                       `json:"provider_source_display_id,omitempty"`
	ExternalResourceID      provider.ExternalResourceID `json:"external_resource_id"`
	Status                  ServiceStatus               `json:"status"`
	BillingStatus           BillingStatus               `json:"billing_status"`
	SuspensionReason        SuspensionReason            `json:"suspension_reason,omitempty"`
	TermStart               time.Time                   `json:"term_start"`
	TermEnd                 time.Time                   `json:"term_end"`
	Credentials             []serviceCredentialResponse `json:"credentials,omitempty"`
	CreatedAt               time.Time                   `json:"created_at"`
	UpdatedAt               time.Time                   `json:"updated_at"`
}

func newServiceInstanceResponse(service ServiceInstance) serviceInstanceResponse {
	return newServiceInstanceResponseWithCredentials(service, nil)
}

func newServiceInstanceResponseWithCredentials(service ServiceInstance, credentials []ServiceCredential) serviceInstanceResponse {
	response := serviceInstanceResponse{
		ID:                      service.ID,
		DisplayID:               service.DisplayID,
		TenantID:                service.TenantID,
		OrderID:                 service.OrderID,
		OrderDisplayID:          service.OrderDisplayID,
		BuyerDisplayID:          service.BuyerDisplayID,
		TenantPlanID:            service.TenantPlanID,
		ProviderSourceID:        service.ProviderSourceID,
		ProviderSourceDisplayID: service.ProviderSourceDisplayID,
		ExternalResourceID:      service.ExternalResourceID,
		Status:                  service.Status,
		BillingStatus:           service.BillingStatus,
		SuspensionReason:        service.SuspensionReason,
		TermStart:               service.TermStart,
		TermEnd:                 service.TermEnd,
		CreatedAt:               service.CreatedAt,
		UpdatedAt:               service.UpdatedAt,
	}
	if len(credentials) > 0 {
		response.Credentials = newServiceCredentialResponses(credentials)
	}
	return response
}

func newServiceInstanceResponses(services []ServiceInstance) []serviceInstanceResponse {
	responses := make([]serviceInstanceResponse, 0, len(services))
	for _, service := range services {
		responses = append(responses, newServiceInstanceResponse(service))
	}
	return responses
}
