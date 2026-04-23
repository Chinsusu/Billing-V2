package wallet

import (
	"context"
	"encoding/json"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
)

const (
	topupAuditActionApproved = "wallet.topup.approved"
	topupAuditActionRejected = "wallet.topup.rejected"
)

type AuditAppender interface {
	Append(ctx context.Context, input audit.AppendInput) (audit.Log, error)
}

func (service *Service) appendTopupReviewAudit(
	ctx context.Context,
	action string,
	before TopupRequest,
	after TopupRequest,
	reviewer identity.UserID,
) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:               after.TenantID,
		ActorID:                audit.ActorID(reviewer),
		ActorType:              audit.ActorTypeUser,
		Action:                 action,
		TargetType:             "topup_request",
		TargetID:               audit.TargetID(after.ID),
		BeforeSnapshotRedacted: walletAuditJSON(topupAuditStatus{Status: before.Status}),
		AfterSnapshotRedacted: walletAuditJSON(topupAuditStatus{
			Status:        after.Status,
			LedgerEntryID: after.LedgerEntryID,
		}),
		MetadataRedacted: walletAuditJSON(topupAuditMetadata{
			DisplayID:   after.DisplayID,
			WalletID:    after.WalletID,
			AmountMinor: after.AmountMinor,
			Currency:    after.Currency,
		}),
		CorrelationID: audit.CorrelationID(after.ID),
	})
	return err
}

type topupAuditStatus struct {
	Status        TopupStatus   `json:"status"`
	LedgerEntryID LedgerEntryID `json:"ledger_entry_id,omitempty"`
}

type topupAuditMetadata struct {
	DisplayID   int64    `json:"display_id"`
	WalletID    WalletID `json:"wallet_id"`
	AmountMinor int64    `json:"amount_minor"`
	Currency    string   `json:"currency"`
}

func walletAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
