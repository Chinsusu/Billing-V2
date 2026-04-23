package order

import (
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type TransitionOrderStatusInput struct {
	ID            OrderID
	TenantID      tenant.ID
	ActorID       identity.UserID
	FromStatus    OrderStatus
	ToStatus      OrderStatus
	BillingStatus BillingStatus
}

func (input TransitionOrderStatusInput) Normalize() TransitionOrderStatusInput {
	return TransitionOrderStatusInput{
		ID:            OrderID(trim(string(input.ID))),
		TenantID:      tenant.ID(trim(string(input.TenantID))),
		ActorID:       identity.UserID(trim(string(input.ActorID))),
		FromStatus:    OrderStatus(trim(string(input.FromStatus))),
		ToStatus:      OrderStatus(trim(string(input.ToStatus))),
		BillingStatus: BillingStatus(trim(string(input.BillingStatus))),
	}
}

func (input TransitionOrderStatusInput) Validate() error {
	if input.ID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ActorID == "" {
		return identity.ErrActorIDMissing
	}
	if !input.FromStatus.Valid() || !input.ToStatus.Valid() {
		return ErrOrderStatusInvalid
	}
	if !input.BillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if !CanTransitionOrder(input.FromStatus, input.ToStatus) {
		return ErrStatusTransitionInvalid
	}
	return nil
}
