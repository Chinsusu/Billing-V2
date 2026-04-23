package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Reservation struct {
	ID               ReservationID
	DisplayID        int64
	OrderID          OrderID
	TenantID         tenant.ID
	ProviderSourceID catalog.ProviderSourceID
	Status           ReservationStatus
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateReservationInput struct {
	OrderID          OrderID
	TenantID         tenant.ID
	ProviderSourceID catalog.ProviderSourceID
	Status           ReservationStatus
	ExpiresAt        time.Time
}

func (input CreateReservationInput) Normalize() CreateReservationInput {
	output := input
	if output.Status == "" {
		output.Status = ReservationStatusPendingReserve
	}
	return output
}

func (input CreateReservationInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if !input.Status.Valid() {
		return ErrReservationStatusInvalid
	}
	if input.ExpiresAt.IsZero() {
		return ErrReservationExpiryMissing
	}
	return nil
}
