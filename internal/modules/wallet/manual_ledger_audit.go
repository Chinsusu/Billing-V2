package wallet

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
)

const (
	walletAuditActionRefundCreated     = "wallet.refund.created"
	walletAuditActionAdjustmentCreated = "wallet.adjustment.created"
)

func (service *Service) appendManualLedgerAudit(ctx context.Context, action string, entry LedgerEntry) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   entry.TenantID,
		ActorID:    audit.ActorID(entry.CreatedBy),
		ActorType:  audit.ActorTypeUser,
		Action:     action,
		TargetType: "wallet_ledger_entry",
		TargetID:   audit.TargetID(entry.ID),
		AfterSnapshotRedacted: walletAuditJSON(manualLedgerAuditStatus{
			Status:            entry.Status,
			BalanceAfterMinor: entry.BalanceAfterMinor,
		}),
		MetadataRedacted: walletAuditJSON(manualLedgerAuditMetadata{
			WalletID:       entry.WalletID,
			DisplayID:      entry.DisplayID,
			Direction:      entry.Direction,
			EntryType:      entry.EntryType,
			AmountMinor:    entry.AmountMinor,
			Currency:       entry.Currency,
			ReferenceType:  entry.ReferenceType,
			ReferenceID:    entry.ReferenceID,
			IdempotencyKey: entry.IdempotencyKey,
			Reason:         entry.Reason,
		}),
		CorrelationID: audit.CorrelationID(entry.CorrelationID),
	})
	return err
}

type manualLedgerAuditStatus struct {
	Status            LedgerStatus `json:"status"`
	BalanceAfterMinor int64        `json:"balance_after_minor"`
}

type manualLedgerAuditMetadata struct {
	WalletID       WalletID       `json:"wallet_id"`
	DisplayID      int64          `json:"display_id"`
	Direction      Direction      `json:"direction"`
	EntryType      EntryType      `json:"entry_type"`
	AmountMinor    int64          `json:"amount_minor"`
	Currency       string         `json:"currency"`
	ReferenceType  ReferenceType  `json:"reference_type"`
	ReferenceID    ReferenceID    `json:"reference_id"`
	IdempotencyKey IdempotencyKey `json:"idempotency_key"`
	Reason         string         `json:"reason"`
}
