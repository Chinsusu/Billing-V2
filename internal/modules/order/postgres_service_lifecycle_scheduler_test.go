package order

import (
	"strings"
	"testing"
	"time"
)

func TestListDueServiceLifecycleActionsArgsUseGraceCutoff(t *testing.T) {
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	args, err := listDueServiceLifecycleActionsArgs(ListDueServiceLifecycleActionsInput{
		Now:         now,
		Limit:       25,
		GracePeriod: 72 * time.Hour,
	})
	if err != nil {
		t.Fatalf("expected due args: %v", err)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[0] != now || args[1] != now.Add(-72*time.Hour) || args[2] != 25 {
		t.Fatalf("unexpected due args: %#v", args)
	}
}

func TestListDueServiceLifecycleActionsSQLGuardsStatusesAndGrace(t *testing.T) {
	for _, clause := range []string{
		"svc.provider_source_id",
		"source.source_type",
		"svc.external_resource_id",
		"JOIN provider_sources source",
		"svc.status = 'active'",
		"svc.billing_status = 'paid'",
		"svc.status = 'expired'",
		"svc.billing_status = 'overdue'",
		"svc.status = 'suspended'",
		"svc.suspension_reason = 'expiry'",
		"expected_billing_status",
		"expected_suspension_reason",
		"svc.term_end <= $2",
		"LIMIT $3",
	} {
		if !strings.Contains(listDueServiceLifecycleActionsSQL, clause) {
			t.Fatalf("expected %q in due SQL: %s", clause, listDueServiceLifecycleActionsSQL)
		}
	}
}
