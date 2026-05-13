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
	Quantity         int
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

type ReserveInventoryInput struct {
	OrderID          OrderID
	TenantID         tenant.ID
	ProviderSourceID catalog.ProviderSourceID
	Quantity         int
	ExpiresAt        time.Time
}

type ExpireReservationsInput struct {
	TenantID tenant.ID
	Now      time.Time
}

const DefaultReservationTTL = 5 * time.Minute

func (input ReserveInventoryInput) Normalize() ReserveInventoryInput {
	output := input
	if output.Quantity == 0 {
		output.Quantity = 1
	}
	if !output.ExpiresAt.IsZero() {
		output.ExpiresAt = output.ExpiresAt.UTC()
	}
	return output
}

func (input ReserveInventoryInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if input.Quantity <= 0 {
		return ErrReservationQuantityInvalid
	}
	if input.ExpiresAt.IsZero() {
		return ErrReservationExpiryMissing
	}
	return nil
}

func (input ExpireReservationsInput) Normalize() ExpireReservationsInput {
	output := input
	if !output.Now.IsZero() {
		output.Now = output.Now.UTC()
	}
	return output
}

func (input ExpireReservationsInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.Now.IsZero() {
		return ErrReservationExpiryMissing
	}
	return nil
}
