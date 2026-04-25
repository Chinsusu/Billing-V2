package wallet

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type createTopupRequestBody struct {
	WalletID         WalletID      `json:"wallet_id"`
	AmountMinor      int64         `json:"amount_minor"`
	Currency         string        `json:"currency"`
	PaymentMethod    PaymentMethod `json:"payment_method"`
	PaymentReference string        `json:"payment_reference"`
}

type reviewTopupRequestBody struct {
	ReviewNote string `json:"review_note"`
}

type topupRequestResponse struct {
	ID                   TopupRequestID  `json:"id"`
	DisplayID            int64           `json:"display_id"`
	TenantID             tenant.ID       `json:"tenant_id"`
	WalletID             WalletID        `json:"wallet_id"`
	WalletDisplayID      int64           `json:"wallet_display_id,omitempty"`
	RequestedBy          identity.UserID `json:"requested_by"`
	RequestedByDisplayID int64           `json:"requested_by_display_id,omitempty"`
	AmountMinor          int64           `json:"amount_minor"`
	Currency             string          `json:"currency"`
	PaymentMethod        PaymentMethod   `json:"payment_method"`
	PaymentReference     string          `json:"payment_reference,omitempty"`
	Status               TopupStatus     `json:"status"`
	ReviewedBy           identity.UserID `json:"reviewed_by,omitempty"`
	ReviewedByDisplayID  int64           `json:"reviewed_by_display_id,omitempty"`
	ReviewedAt           *time.Time      `json:"reviewed_at,omitempty"`
	ReviewNote           string          `json:"review_note,omitempty"`
	LedgerEntryID        LedgerEntryID   `json:"ledger_entry_id,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

func newTopupRequestResponse(request TopupRequest) topupRequestResponse {
	return topupRequestResponse{
		ID:                   request.ID,
		DisplayID:            request.DisplayID,
		TenantID:             request.TenantID,
		WalletID:             request.WalletID,
		WalletDisplayID:      request.WalletDisplayID,
		RequestedBy:          request.RequestedBy,
		RequestedByDisplayID: request.RequestedByDisplayID,
		AmountMinor:          request.AmountMinor,
		Currency:             request.Currency,
		PaymentMethod:        request.PaymentMethod,
		PaymentReference:     request.PaymentReference,
		Status:               request.Status,
		ReviewedBy:           request.ReviewedBy,
		ReviewedByDisplayID:  request.ReviewedByDisplayID,
		ReviewedAt:           request.ReviewedAt,
		ReviewNote:           request.ReviewNote,
		LedgerEntryID:        request.LedgerEntryID,
		CreatedAt:            request.CreatedAt,
		UpdatedAt:            request.UpdatedAt,
	}
}

func newTopupRequestResponses(requests []TopupRequest) []topupRequestResponse {
	responses := make([]topupRequestResponse, 0, len(requests))
	for _, request := range requests {
		responses = append(responses, newTopupRequestResponse(request))
	}
	return responses
}
