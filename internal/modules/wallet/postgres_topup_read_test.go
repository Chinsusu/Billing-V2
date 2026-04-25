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
		TenantID:             tenant.ID("tenant-1"),
		WalletID:             WalletID("wallet-1"),
		WalletDisplayID:      70004,
		RequestedBy:          identity.UserID("account-1"),
		RequestedByDisplayID: 10002,
		DisplayID:            90001,
		PaymentMethod:        PaymentMethodBankTransfer,
		Status:               TopupStatusUnderReview,
		AmountMinMinor:       int64Ptr(100),
		AmountMaxMinor:       int64Ptr(900),
		Limit:                25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"topup.tenant_id = $1",
		"topup.display_id = $2",
		"topup.wallet_id = $3",
		"wallet.wallet_id = topup.wallet_id",
		"wallet.display_id = $4",
		"topup.requested_by = $5",
		"requester.user_id = topup.requested_by",
		"requester.display_id = $6",
		"topup.payment_method = $7",
		"topup.status = $8",
		"topup.amount_minor >= $9",
		"topup.amount_minor <= $10",
		"LIMIT $11",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 11 || args[10] != 25 {
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
