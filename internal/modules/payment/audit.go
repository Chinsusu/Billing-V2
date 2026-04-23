package payment

import (
	"context"
	"encoding/json"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const walletPaymentAuditAction = "invoice.wallet_paid"

type AuditAppender interface {
	Append(ctx context.Context, input audit.AppendInput) (audit.Log, error)
}

func (service *Service) appendWalletPaymentAudit(ctx context.Context, input PayInvoiceFromWalletInput, result WalletInvoicePayment) error {
	if service.audit == nil || result.LedgerEntry.ID.Empty() {
		return nil
	}
	record := result.Invoice.Invoice
	after := invoiceAuditStatus{Status: record.Status}
	if !record.PaidAt.IsZero() {
		after.PaidAt = record.PaidAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   record.TenantID,
		ActorID:    audit.ActorID(input.ActorID),
		ActorType:  audit.ActorTypeUser,
		Action:     walletPaymentAuditAction,
		TargetType: "invoice",
		TargetID:   audit.TargetID(record.ID),
		BeforeSnapshotRedacted: paymentAuditJSON(invoiceAuditStatus{
			Status: result.PreviousInvoiceStatus,
		}),
		AfterSnapshotRedacted: paymentAuditJSON(after),
		MetadataRedacted: paymentAuditJSON(walletPaymentAuditMetadata{
			InvoiceDisplayID:     record.DisplayID,
			PaymentTransactionID: result.Transaction.ID,
			TransactionDisplayID: result.Transaction.DisplayID,
			WalletID:             input.WalletID,
			LedgerEntryID:        result.LedgerEntry.ID,
			AmountMinor:          record.TotalMinor,
			Currency:             record.Currency,
		}),
		CorrelationID: audit.CorrelationID(record.ID),
	})
	return err
}

type invoiceAuditStatus struct {
	Status invoice.Status `json:"status"`
	PaidAt string         `json:"paid_at,omitempty"`
}

type walletPaymentAuditMetadata struct {
	InvoiceDisplayID     int64                `json:"invoice_display_id"`
	PaymentTransactionID TransactionID        `json:"payment_transaction_id"`
	TransactionDisplayID int64                `json:"transaction_display_id"`
	WalletID             wallet.WalletID      `json:"wallet_id"`
	LedgerEntryID        wallet.LedgerEntryID `json:"ledger_entry_id"`
	AmountMinor          int64                `json:"amount_minor"`
	Currency             string               `json:"currency"`
}

func paymentAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
