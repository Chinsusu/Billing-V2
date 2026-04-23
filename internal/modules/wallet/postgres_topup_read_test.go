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
		TenantID:      tenant.ID("tenant-1"),
		WalletID:      WalletID("wallet-1"),
		RequestedBy:   identity.UserID("account-1"),
		PaymentMethod: PaymentMethodBankTransfer,
		Status:        TopupStatusUnderReview,
		Limit:         25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"topup.tenant_id = $1",
		"topup.wallet_id = $2",
		"topup.requested_by = $3",
		"topup.payment_method = $4",
		"topup.status = $5",
		"LIMIT $6",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 6 || args[5] != 25 {
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
