package payment

import (
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestDailyReconciliationSQLUsesTenantAndWindow(t *testing.T) {
	for name, query := range map[string]string{
		"wallet_summary":      dailyWalletSummarySQL,
		"wallet_mismatch":     dailyWalletMismatchSQL,
		"invoice_checked":     dailyInvoiceCheckedSQL,
		"invoice_mismatch":    dailyInvoiceMismatchSQL,
		"payment_checked":     dailyPaymentCheckedSQL,
		"duplicate_reference": dailyDuplicatePaymentReferenceSQL,
	} {
		if !strings.Contains(query, "tenant_id = $1") && !strings.Contains(query, "wallet.tenant_id = $1") {
			t.Fatalf("%s query missing tenant scope: %s", name, query)
		}
	}
	for _, query := range []string{dailyInvoiceCheckedSQL, dailyInvoiceMismatchSQL, dailyPaymentCheckedSQL, dailyDuplicatePaymentReferenceSQL} {
		if !strings.Contains(query, ">= $2") || !strings.Contains(query, "< $3") {
			t.Fatalf("expected window bounds in query: %s", query)
		}
	}
	if !strings.Contains(dailyWalletMismatchSQL, "CASE WHEN direction = 'credit' THEN amount_minor ELSE -amount_minor END") {
		t.Fatalf("wallet mismatch query must recompute ledger balance: %s", dailyWalletMismatchSQL)
	}
	if !strings.Contains(dailyInvoiceMismatchSQL, "paid_invoice_amount_mismatch") ||
		!strings.Contains(dailyDuplicatePaymentReferenceSQL, "HAVING COUNT(*) > 1") {
		t.Fatalf("expected mismatch and duplicate checks")
	}
}

func TestDailyReconciliationInputNormalizeUsesUTCDay(t *testing.T) {
	input := DailyReconciliationInput{
		TenantID: tenant.ID("tenant-1"),
		Date:     time.Date(2026, 5, 13, 8, 30, 0, 0, time.FixedZone("ICT", 7*60*60)),
	}.Normalize()
	expected := time.Date(2026, 5, 13, 0, 0, 0, 0, time.UTC)
	if !input.Date.Equal(expected) || !input.WindowTo().Equal(expected.Add(24*time.Hour)) {
		t.Fatalf("unexpected normalized window: %+v to %s", input, input.WindowTo())
	}
}
