package payment

import (
	"context"
	"sort"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

type ReconciliationReportStatus string

const (
	ReconciliationReportStatusBalanced   ReconciliationReportStatus = "balanced"
	ReconciliationReportStatusMismatched ReconciliationReportStatus = "mismatched"
)

type DailyReconciliationInput struct {
	TenantID tenant.ID
	Date     time.Time
}

type DailyReconciliationReport struct {
	TenantID    tenant.ID
	Date        string
	WindowFrom  time.Time
	WindowTo    time.Time
	GeneratedAt time.Time
	Status      ReconciliationReportStatus
	Wallets     WalletReconciliationSummary
	Invoices    InvoiceReconciliationSummary
	Payments    PaymentReconciliationSummary
}

type DailyReconciliationData struct {
	WalletsChecked             int
	WalletMismatches           []WalletBalanceMismatch
	InvoicesChecked            int
	InvoicePaymentMismatches   []InvoicePaymentMismatch
	PaymentsChecked            int
	DuplicatePaymentReferences []DuplicatePaymentReference
}

type WalletReconciliationSummary struct {
	Checked    int
	Balanced   int
	Mismatched int
	Mismatches []WalletBalanceMismatch
}

type WalletBalanceMismatch struct {
	WalletID              wallet.WalletID
	WalletDisplayID       int64
	Currency              string
	AvailableBalanceMinor int64
	LedgerBalanceMinor    int64
	DifferenceMinor       int64
	LastLedgerEntryID     wallet.LedgerEntryID
	LastLedgerDisplayID   int64
}

type InvoiceReconciliationSummary struct {
	Checked    int
	Mismatched int
	Mismatches []InvoicePaymentMismatch
}

type InvoicePaymentMismatch struct {
	InvoiceID                     invoice.InvoiceID
	InvoiceDisplayID              int64
	Status                        invoice.Status
	TotalMinor                    int64
	PostedPaymentTotalMinor       int64
	PostedPaymentTransactionCount int
	Reason                        string
}

type PaymentReconciliationSummary struct {
	Checked                 int
	DuplicateReferenceCount int
	DuplicateReferences     []DuplicatePaymentReference
}

type DuplicatePaymentReference struct {
	ReferenceType         string
	ReferenceID           invoice.InvoiceID
	ReferenceDisplayID    int64
	TransactionDisplayIDs []int64
	TransactionCount      int
	TotalAmountMinor      int64
}

type DailyReconciliationStore interface {
	GetDailyReconciliationData(ctx context.Context, input DailyReconciliationInput) (DailyReconciliationData, error)
}

func (input DailyReconciliationInput) Normalize() DailyReconciliationInput {
	output := input
	if !output.Date.IsZero() {
		date := output.Date.UTC()
		output.Date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	}
	return output
}

func (input DailyReconciliationInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.Date.IsZero() {
		return ErrCreatedTimeInvalid
	}
	return nil
}

func (input DailyReconciliationInput) WindowTo() time.Time {
	return input.Date.Add(24 * time.Hour)
}

func newDailyReconciliationReport(
	input DailyReconciliationInput,
	data DailyReconciliationData,
	generatedAt time.Time,
) DailyReconciliationReport {
	sortDailyReconciliationData(&data)
	status := ReconciliationReportStatusBalanced
	if len(data.WalletMismatches) > 0 ||
		len(data.InvoicePaymentMismatches) > 0 ||
		len(data.DuplicatePaymentReferences) > 0 {
		status = ReconciliationReportStatusMismatched
	}
	walletMismatchCount := len(data.WalletMismatches)
	invoiceMismatchCount := len(data.InvoicePaymentMismatches)
	return DailyReconciliationReport{
		TenantID:    input.TenantID,
		Date:        input.Date.Format("2006-01-02"),
		WindowFrom:  input.Date,
		WindowTo:    input.WindowTo(),
		GeneratedAt: generatedAt.UTC(),
		Status:      status,
		Wallets: WalletReconciliationSummary{
			Checked:    data.WalletsChecked,
			Balanced:   data.WalletsChecked - walletMismatchCount,
			Mismatched: walletMismatchCount,
			Mismatches: data.WalletMismatches,
		},
		Invoices: InvoiceReconciliationSummary{
			Checked:    data.InvoicesChecked,
			Mismatched: invoiceMismatchCount,
			Mismatches: data.InvoicePaymentMismatches,
		},
		Payments: PaymentReconciliationSummary{
			Checked:                 data.PaymentsChecked,
			DuplicateReferenceCount: len(data.DuplicatePaymentReferences),
			DuplicateReferences:     data.DuplicatePaymentReferences,
		},
	}
}

func sortDailyReconciliationData(data *DailyReconciliationData) {
	sort.Slice(data.WalletMismatches, func(left int, right int) bool {
		return data.WalletMismatches[left].WalletDisplayID < data.WalletMismatches[right].WalletDisplayID
	})
	sort.Slice(data.InvoicePaymentMismatches, func(left int, right int) bool {
		return data.InvoicePaymentMismatches[left].InvoiceDisplayID < data.InvoicePaymentMismatches[right].InvoiceDisplayID
	})
	sort.Slice(data.DuplicatePaymentReferences, func(left int, right int) bool {
		return data.DuplicatePaymentReferences[left].ReferenceDisplayID < data.DuplicatePaymentReferences[right].ReferenceDisplayID
	})
}
