package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ProvisioningJob struct {
	ID                  ProvisioningJobID
	DisplayID           int64
	OrderID             OrderID
	TenantID            tenant.ID
	ProviderSourceID    catalog.ProviderSourceID
	ProviderOperationID ProviderOperationID
	Status              ProvisioningStatus
	IdempotencyKey      IdempotencyKey
	AttemptNumber       int
	LastErrorCode       string
	LastErrorMessage    string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateProvisioningJobInput struct {
	OrderID             OrderID
	TenantID            tenant.ID
	ProviderSourceID    catalog.ProviderSourceID
	ProviderOperationID ProviderOperationID
	Status              ProvisioningStatus
	IdempotencyKey      IdempotencyKey
	AttemptNumber       int
}

type RecordProvisioningResultInput struct {
	OrderID             OrderID
	TenantID            tenant.ID
	ProviderSourceID    catalog.ProviderSourceID
	ProviderOperationID ProviderOperationID
	Status              ProvisioningStatus
	IdempotencyKey      IdempotencyKey
	AttemptNumber       int
	LastErrorCode       string
	LastErrorMessage    string
}

func (input CreateProvisioningJobInput) Normalize() CreateProvisioningJobInput {
	output := input
	output.ProviderOperationID = ProviderOperationID(trim(string(output.ProviderOperationID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	if output.Status == "" {
		output.Status = ProvisioningStatusQueued
	}
	if output.AttemptNumber == 0 {
		output.AttemptNumber = 1
	}
	return output
}

func (input CreateProvisioningJobInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if input.ProviderOperationID == "" {
		return ErrProviderOperationIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if input.AttemptNumber <= 0 {
		return ErrAttemptInvalid
	}
	if !input.Status.Valid() {
		return ErrProvisioningStatusInvalid
	}
	return nil
}

func (input RecordProvisioningResultInput) Normalize() RecordProvisioningResultInput {
	output := input
	output.ProviderOperationID = ProviderOperationID(trim(string(output.ProviderOperationID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.LastErrorCode = trim(output.LastErrorCode)
	output.LastErrorMessage = trim(output.LastErrorMessage)
	if output.AttemptNumber == 0 {
		output.AttemptNumber = 1
	}
	return output
}

func (input RecordProvisioningResultInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if input.ProviderOperationID == "" {
		return ErrProviderOperationIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if input.AttemptNumber <= 0 {
		return ErrAttemptInvalid
	}
	if !input.Status.Valid() {
		return ErrProvisioningStatusInvalid
	}
	return nil
}
