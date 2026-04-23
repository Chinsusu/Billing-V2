package payment

import (
	"context"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const defaultReconciliationListLimit = 100
const maxReconciliationListLimit = 500

type ReconciliationFilter struct {
	TenantID    tenant.ID
	Status      TransactionStatus
	Provider    string
	InvoiceID   invoice.InvoiceID
	WalletID    wallet.WalletID
	CreatedFrom time.Time
	CreatedTo   time.Time
	Limit       int
}

type ReconciliationLookup struct {
	TenantID      tenant.ID
	TransactionID TransactionID
}

type PaymentReconciliation struct {
	Transaction Transaction
	Provider    string
	Invoice     ReconciliationInvoice
	Ledger      ReconciliationLedger
}

type ReconciliationInvoice struct {
	ID         invoice.InvoiceID
	DisplayID  int64
	Status     invoice.Status
	TotalMinor int64
	PaidAt     time.Time
}

func (record ReconciliationInvoice) Empty() bool {
	return record.ID.Empty()
}

type ReconciliationLedger struct {
	ID                wallet.LedgerEntryID
	DisplayID         int64
	WalletID          wallet.WalletID
	WalletDisplayID   int64
	Direction         wallet.Direction
	EntryType         wallet.EntryType
	Status            wallet.LedgerStatus
	BalanceAfterMinor int64
}

func (record ReconciliationLedger) Empty() bool {
	return record.ID.Empty()
}

type ReconciliationStore interface {
	ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error)
	GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error)
}

func normalizeReconciliationFilter(filter ReconciliationFilter) ReconciliationFilter {
	output := filter
	output.Provider = trim(output.Provider)
	if output.Limit <= 0 {
		output.Limit = defaultReconciliationListLimit
	}
	if output.Limit > maxReconciliationListLimit {
		output.Limit = maxReconciliationListLimit
	}
	return output
}

func validateReconciliationFilter(filter ReconciliationFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrStatusInvalid
	}
	if !filter.CreatedFrom.IsZero() && !filter.CreatedTo.IsZero() && filter.CreatedTo.Before(filter.CreatedFrom) {
		return ErrCreatedTimeWindowInvalid
	}
	return nil
}

func validateReconciliationLookup(lookup ReconciliationLookup) error {
	if lookup.TransactionID.Empty() {
		return ErrTransactionIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}
