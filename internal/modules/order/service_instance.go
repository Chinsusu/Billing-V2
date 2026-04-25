package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ServiceInstance struct {
	ID                      ServiceID
	DisplayID               int64
	TenantID                tenant.ID
	OrderID                 OrderID
	OrderDisplayID          int64
	BuyerDisplayID          int64
	TenantPlanID            catalog.TenantPlanID
	ProviderSourceID        catalog.ProviderSourceID
	ProviderSourceDisplayID int64
	ExternalResourceID      provider.ExternalResourceID
	Status                  ServiceStatus
	BillingStatus           BillingStatus
	SuspensionReason        SuspensionReason
	TermStart               time.Time
	TermEnd                 time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CreateServiceInstanceInput struct {
	TenantID           tenant.ID
	OrderID            OrderID
	TenantPlanID       catalog.TenantPlanID
	ProviderSourceID   catalog.ProviderSourceID
	ExternalResourceID provider.ExternalResourceID
	Status             ServiceStatus
	BillingStatus      BillingStatus
	SuspensionReason   SuspensionReason
	TermStart          time.Time
	TermEnd            time.Time
}

func (input CreateServiceInstanceInput) Normalize() CreateServiceInstanceInput {
	output := input
	output.ExternalResourceID = provider.ExternalResourceID(trim(string(output.ExternalResourceID)))
	if output.Status == "" {
		output.Status = ServiceStatusActive
	}
	if output.BillingStatus == "" {
		output.BillingStatus = BillingStatusPaid
	}
	return output
}

func (input CreateServiceInstanceInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantPlanID.Empty() {
		return ErrTenantPlanIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if input.ExternalResourceID == "" {
		return ErrExternalResourceIDMissing
	}
	if !input.Status.Valid() {
		return ErrServiceStatusInvalid
	}
	if !input.BillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if input.SuspensionReason != "" && !input.SuspensionReason.Valid() {
		return ErrSuspensionReasonInvalid
	}
	if input.Status == ServiceStatusSuspended && input.SuspensionReason == "" {
		return ErrSuspensionReasonInvalid
	}
	if input.TermStart.IsZero() || input.TermEnd.IsZero() || !input.TermEnd.After(input.TermStart) {
		return ErrTermWindowInvalid
	}
	return nil
}
