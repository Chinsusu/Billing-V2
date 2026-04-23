package wallet

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListTopupRequestsQueryAddsReviewFilters(t *testing.T) {
	query, args, err := buildListTopupRequestsQuery(TopupRequestFilter{
		TenantID:       tenant.ID("tenant-1"),
		WalletID:       WalletID("wallet-1"),
		RequestedBy:    identity.UserID("account-1"),
		DisplayID:      90001,
		PaymentMethod:  PaymentMethodBankTransfer,
		Status:         TopupStatusUnderReview,
		AmountMinMinor: int64Ptr(100),
		AmountMaxMinor: int64Ptr(900),
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"topup.tenant_id = $1",
		"topup.display_id = $2",
		"topup.wallet_id = $3",
		"topup.requested_by = $4",
		"topup.payment_method = $5",
		"topup.status = $6",
		"topup.amount_minor >= $7",
		"topup.amount_minor <= $8",
		"LIMIT $9",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 9 || args[8] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildGetTopupRequestQueryAddsRequesterScope(t *testing.T) {
	query, args, err := buildGetTopupRequestQuery(TopupRequestLookup{
		ID:          TopupRequestID("topup-1"),
		TenantID:    tenant.ID("tenant-1"),
		RequestedBy: identity.UserID("account-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"topup.topup_request_id = $1", "topup.tenant_id = $2", "topup.requested_by = $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListTopupRequestsQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListTopupRequestsQuery(TopupRequestFilter{
		TenantID: tenant.ID("tenant-1"),
		Status:   TopupStatus("bad"),
	})
	if !errors.Is(err, ErrTopupStatusInvalid) {
		t.Fatalf("expected status error, got %v", err)
	}
}
