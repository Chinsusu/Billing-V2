package payment

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type dailyReconciliationReportResponse struct {
	TenantID    tenant.ID                            `json:"tenant_id"`
	Date        string                               `json:"date"`
	WindowFrom  time.Time                            `json:"window_from"`
	WindowTo    time.Time                            `json:"window_to"`
	GeneratedAt time.Time                            `json:"generated_at"`
	Status      ReconciliationReportStatus           `json:"status"`
	Wallets     walletReconciliationResponse         `json:"wallets"`
	Invoices    invoiceReconciliationResponse        `json:"invoices"`
	Payments    paymentReconciliationSummaryResponse `json:"payments"`
}

type walletReconciliationResponse struct {
	Checked    int                             `json:"checked"`
	Balanced   int                             `json:"balanced"`
	Mismatched int                             `json:"mismatched"`
	Mismatches []walletBalanceMismatchResponse `json:"mismatches"`
}

type walletBalanceMismatchResponse struct {
	WalletDisplayID       int64  `json:"wallet_display_id"`
	Currency              string `json:"currency"`
	AvailableBalanceMinor int64  `json:"available_balance_minor"`
	LedgerBalanceMinor    int64  `json:"ledger_balance_minor"`
	DifferenceMinor       int64  `json:"difference_minor"`
	LastLedgerDisplayID   int64  `json:"last_ledger_display_id,omitempty"`
}

type invoiceReconciliationResponse struct {
	Checked    int                              `json:"checked"`
	Mismatched int                              `json:"mismatched"`
	Mismatches []invoicePaymentMismatchResponse `json:"mismatches"`
}

type invoicePaymentMismatchResponse struct {
	InvoiceDisplayID              int64          `json:"invoice_display_id"`
	Status                        invoice.Status `json:"status"`
	TotalMinor                    int64          `json:"total_minor"`
	PostedPaymentTotalMinor       int64          `json:"posted_payment_total_minor"`
	PostedPaymentTransactionCount int            `json:"posted_payment_transaction_count"`
	Reason                        string         `json:"reason"`
}

type paymentReconciliationSummaryResponse struct {
	Checked                 int                                 `json:"checked"`
	DuplicateReferenceCount int                                 `json:"duplicate_reference_count"`
	DuplicateReferences     []duplicatePaymentReferenceResponse `json:"duplicate_references"`
}

type duplicatePaymentReferenceResponse struct {
	ReferenceType         string  `json:"reference_type"`
	ReferenceDisplayID    int64   `json:"reference_display_id"`
	TransactionDisplayIDs []int64 `json:"transaction_display_ids"`
	TransactionCount      int     `json:"transaction_count"`
	TotalAmountMinor      int64   `json:"total_amount_minor"`
}

func newDailyReconciliationReportResponse(report DailyReconciliationReport) dailyReconciliationReportResponse {
	return dailyReconciliationReportResponse{
		TenantID:    report.TenantID,
		Date:        report.Date,
		WindowFrom:  report.WindowFrom,
		WindowTo:    report.WindowTo,
		GeneratedAt: report.GeneratedAt,
		Status:      report.Status,
		Wallets: walletReconciliationResponse{
			Checked:    report.Wallets.Checked,
			Balanced:   report.Wallets.Balanced,
			Mismatched: report.Wallets.Mismatched,
			Mismatches: newWalletBalanceMismatchResponses(report.Wallets.Mismatches),
		},
		Invoices: invoiceReconciliationResponse{
			Checked:    report.Invoices.Checked,
			Mismatched: report.Invoices.Mismatched,
			Mismatches: newInvoicePaymentMismatchResponses(report.Invoices.Mismatches),
		},
		Payments: paymentReconciliationSummaryResponse{
			Checked:                 report.Payments.Checked,
			DuplicateReferenceCount: report.Payments.DuplicateReferenceCount,
			DuplicateReferences:     newDuplicatePaymentReferenceResponses(report.Payments.DuplicateReferences),
		},
	}
}

func newWalletBalanceMismatchResponses(mismatches []WalletBalanceMismatch) []walletBalanceMismatchResponse {
	responses := make([]walletBalanceMismatchResponse, 0, len(mismatches))
	for _, mismatch := range mismatches {
		responses = append(responses, walletBalanceMismatchResponse{
			WalletDisplayID:       mismatch.WalletDisplayID,
			Currency:              mismatch.Currency,
			AvailableBalanceMinor: mismatch.AvailableBalanceMinor,
			LedgerBalanceMinor:    mismatch.LedgerBalanceMinor,
			DifferenceMinor:       mismatch.DifferenceMinor,
			LastLedgerDisplayID:   mismatch.LastLedgerDisplayID,
		})
	}
	return responses
}

func newInvoicePaymentMismatchResponses(mismatches []InvoicePaymentMismatch) []invoicePaymentMismatchResponse {
	responses := make([]invoicePaymentMismatchResponse, 0, len(mismatches))
	for _, mismatch := range mismatches {
		responses = append(responses, invoicePaymentMismatchResponse{
			InvoiceDisplayID:              mismatch.InvoiceDisplayID,
			Status:                        mismatch.Status,
			TotalMinor:                    mismatch.TotalMinor,
			PostedPaymentTotalMinor:       mismatch.PostedPaymentTotalMinor,
			PostedPaymentTransactionCount: mismatch.PostedPaymentTransactionCount,
			Reason:                        mismatch.Reason,
		})
	}
	return responses
}

func newDuplicatePaymentReferenceResponses(duplicates []DuplicatePaymentReference) []duplicatePaymentReferenceResponse {
	responses := make([]duplicatePaymentReferenceResponse, 0, len(duplicates))
	for _, duplicate := range duplicates {
		responses = append(responses, duplicatePaymentReferenceResponse{
			ReferenceType:         duplicate.ReferenceType,
			ReferenceDisplayID:    duplicate.ReferenceDisplayID,
			TransactionDisplayIDs: duplicate.TransactionDisplayIDs,
			TransactionCount:      duplicate.TransactionCount,
			TotalAmountMinor:      duplicate.TotalAmountMinor,
		})
	}
	return responses
}
