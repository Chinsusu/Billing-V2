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
		"status = 'active'",
		"billing_status = 'paid'",
		"status = 'expired'",
		"billing_status = 'overdue'",
		"status = 'suspended'",
		"suspension_reason = 'expiry'",
		"expected_billing_status",
		"expected_suspension_reason",
		"term_end <= $2",
		"LIMIT $3",
	} {
		if !strings.Contains(listDueServiceLifecycleActionsSQL, clause) {
			t.Fatalf("expected %q in due SQL: %s", clause, listDueServiceLifecycleActionsSQL)
		}
	}
}
