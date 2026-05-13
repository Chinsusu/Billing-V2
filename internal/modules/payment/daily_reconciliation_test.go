package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildDailyReconciliationReportReturnsBalancedCleanData(t *testing.T) {
	service := newDailyReconciliationTestService(DailyReconciliationData{
		WalletsChecked:  2,
		InvoicesChecked: 3,
		PaymentsChecked: 4,
	})

	report, err := service.BuildDailyReconciliationReport(context.Background(), dailyReconciliationInput())
	if err != nil {
		t.Fatalf("expected report: %v", err)
	}
	if report.Status != ReconciliationReportStatusBalanced {
		t.Fatalf("expected balanced report, got %+v", report)
	}
	if report.Wallets.Balanced != 2 || report.Wallets.Mismatched != 0 ||
		report.Invoices.Mismatched != 0 || report.Payments.DuplicateReferenceCount != 0 {
		t.Fatalf("unexpected clean report: %+v", report)
	}
	if report.Date != "2026-05-13" ||
		!report.WindowTo.Equal(time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC)) ||
		!report.GeneratedAt.Equal(dailyReconciliationNow()) {
		t.Fatalf("unexpected report timestamps: %+v", report)
	}
}

func TestBuildDailyReconciliationReportFlagsWalletMismatch(t *testing.T) {
	service := newDailyReconciliationTestService(DailyReconciliationData{
		WalletsChecked: 2,
		WalletMismatches: []WalletBalanceMismatch{{
			WalletID:              "wallet-1",
			WalletDisplayID:       70001,
			Currency:              "USD",
			AvailableBalanceMinor: 900,
			LedgerBalanceMinor:    1000,
			DifferenceMinor:       -100,
			LastLedgerEntryID:     "ledger-1",
			LastLedgerDisplayID:   71001,
		}},
	})

	report, err := service.BuildDailyReconciliationReport(context.Background(), dailyReconciliationInput())
	if err != nil {
		t.Fatalf("expected report: %v", err)
	}
	if report.Status != ReconciliationReportStatusMismatched ||
		report.Wallets.Balanced != 1 ||
		report.Wallets.Mismatched != 1 {
		t.Fatalf("expected wallet mismatch report, got %+v", report)
	}
}

func TestBuildDailyReconciliationReportFlagsInvoicePaymentMismatch(t *testing.T) {
	service := newDailyReconciliationTestService(DailyReconciliationData{
		InvoicesChecked: 2,
		InvoicePaymentMismatches: []InvoicePaymentMismatch{{
			InvoiceID:                     "invoice-1",
			InvoiceDisplayID:              44001,
			Status:                        invoice.StatusPaid,
			TotalMinor:                    1200,
			PostedPaymentTotalMinor:       1000,
			PostedPaymentTransactionCount: 1,
			Reason:                        "paid_invoice_amount_mismatch",
		}},
	})

	report, err := service.BuildDailyReconciliationReport(context.Background(), dailyReconciliationInput())
	if err != nil {
		t.Fatalf("expected report: %v", err)
	}
	if report.Status != ReconciliationReportStatusMismatched ||
		report.Invoices.Mismatched != 1 ||
		report.Invoices.Mismatches[0].InvoiceDisplayID != 44001 {
		t.Fatalf("expected invoice mismatch report, got %+v", report)
	}
}

func TestBuildDailyReconciliationReportFlagsDuplicatePaymentReference(t *testing.T) {
	service := newDailyReconciliationTestService(DailyReconciliationData{
		PaymentsChecked: 2,
		DuplicatePaymentReferences: []DuplicatePaymentReference{{
			ReferenceType:         "invoice",
			ReferenceID:           "invoice-1",
			ReferenceDisplayID:    44001,
			TransactionDisplayIDs: []int64{51001, 51002},
			TransactionCount:      2,
			TotalAmountMinor:      2400,
		}},
	})

	report, err := service.BuildDailyReconciliationReport(context.Background(), dailyReconciliationInput())
	if err != nil {
		t.Fatalf("expected report: %v", err)
	}
	if report.Status != ReconciliationReportStatusMismatched ||
		report.Payments.DuplicateReferenceCount != 1 ||
		report.Payments.DuplicateReferences[0].TransactionDisplayIDs[1] != 51002 {
		t.Fatalf("expected duplicate payment report, got %+v", report)
	}
}

func TestBuildDailyReconciliationReportRejectsMissingDate(t *testing.T) {
	service := newDailyReconciliationTestService(DailyReconciliationData{})

	_, err := service.BuildDailyReconciliationReport(context.Background(), DailyReconciliationInput{TenantID: tenant.ID("tenant-1")})
	if !errors.Is(err, ErrCreatedTimeInvalid) {
		t.Fatalf("expected date error, got %v", err)
	}
}

func dailyReconciliationInput() DailyReconciliationInput {
	return DailyReconciliationInput{
		TenantID: tenant.ID("tenant-1"),
		Date:     time.Date(2026, 5, 13, 8, 30, 0, 0, time.FixedZone("ICT", 7*60*60)),
	}
}

func dailyReconciliationNow() time.Time {
	return time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
}

func newDailyReconciliationTestService(data DailyReconciliationData) *Service {
	service := NewService(&fakeDailyReconciliationStore{data: data})
	service.now = dailyReconciliationNow
	return service
}

type fakeDailyReconciliationStore struct {
	data DailyReconciliationData
}

func (store *fakeDailyReconciliationStore) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	return Transaction{}, nil
}

func (store *fakeDailyReconciliationStore) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	return nil, nil
}

func (store *fakeDailyReconciliationStore) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	return Transaction{}, nil
}

func (store *fakeDailyReconciliationStore) GetDailyReconciliationData(ctx context.Context, input DailyReconciliationInput) (DailyReconciliationData, error) {
	return store.data, nil
}

func TestParseInt64CSV(t *testing.T) {
	values, err := parseInt64CSV("51001, 51002")
	if err != nil {
		t.Fatalf("expected parsed values: %v", err)
	}
	if len(values) != 2 || values[0] != 51001 || values[1] != 51002 {
		t.Fatalf("unexpected parsed values: %#v", values)
	}
	if _, err := parseInt64CSV("51001,bad"); err == nil {
		t.Fatalf("expected invalid display id to fail")
	}
}
