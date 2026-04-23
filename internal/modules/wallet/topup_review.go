package wallet

import (
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ApproveTopupRequestInput struct {
	ID         TopupRequestID
	TenantID   tenant.ID
	ReviewedBy identity.UserID
	ReviewNote string
}

type RejectTopupRequestInput struct {
	ID         TopupRequestID
	TenantID   tenant.ID
	ReviewedBy identity.UserID
	ReviewNote string
}

func (input ApproveTopupRequestInput) Normalize() ApproveTopupRequestInput {
	output := input
	output.ReviewNote = trim(output.ReviewNote)
	return output
}

func (input ApproveTopupRequestInput) Validate() error {
	if input.ID.Empty() {
		return ErrTopupRequestIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ReviewedBy == "" {
		return identity.ErrActorIDMissing
	}
	return nil
}

func (input RejectTopupRequestInput) Normalize() RejectTopupRequestInput {
	output := input
	output.ReviewNote = trim(output.ReviewNote)
	return output
}

func (input RejectTopupRequestInput) Validate() error {
	if input.ID.Empty() {
		return ErrTopupRequestIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ReviewedBy == "" {
		return identity.ErrActorIDMissing
	}
	if input.ReviewNote == "" {
		return ErrReviewNoteMissing
	}
	return nil
}

func approveLedgerInput(request TopupRequest, reviewer identity.UserID) PostLedgerEntryInput {
	referenceID := ReferenceID(request.ID)
	return PostLedgerEntryInput{
		WalletID:       request.WalletID,
		TenantID:       request.TenantID,
		Direction:      DirectionCredit,
		AmountMinor:    request.AmountMinor,
		Currency:       request.Currency,
		EntryType:      EntryTypeTopup,
		ReferenceType:  ReferenceType("topup_request"),
		ReferenceID:    referenceID,
		IdempotencyKey: IdempotencyKey(fmt.Sprintf("topup_request:%s:approve", request.ID)),
		CreatedBy:      reviewer,
		Reason:         fmt.Sprintf("Top-up request %d approved", request.DisplayID),
		CorrelationID:  CorrelationID(request.ID),
	}
}

func reviewableTopupStatus(status TopupStatus) bool {
	return status == TopupStatusSubmitted || status == TopupStatusUnderReview
}
