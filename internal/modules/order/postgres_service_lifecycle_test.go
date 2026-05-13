package order

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestTransitionServiceLifecycleArgsNormalizeAndValidate(t *testing.T) {
	args, err := transitionServiceLifecycleArgs(TransitionServiceLifecycleInput{
		ID:               " service-1 ",
		TenantID:         tenant.ID(" tenant-1 "),
		ActorID:          " admin-1 ",
		Action:           ServiceLifecycleActionSuspend,
		FromStatus:       ServiceStatusActive,
		ToStatus:         ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonManualAdmin,
		Reason:           "abuse ticket AB-1",
	})
	if err != nil {
		t.Fatalf("expected lifecycle args: %v", err)
	}
	if len(args) != 11 {
		t.Fatalf("expected 11 args, got %d", len(args))
	}
	if args[0] != ServiceID("service-1") ||
		args[1] != tenant.ID("tenant-1") ||
		args[7] != ServiceEventSuspended {
		t.Fatalf("unexpected lifecycle args: %#v", args)
	}
	reason, ok := args[5].(sql.NullString)
	if !ok || !reason.Valid || reason.String != string(SuspensionReasonManualAdmin) {
		t.Fatalf("unexpected suspension reason arg: %#v", args[5])
	}
}

func TestTransitionServiceLifecycleArgsRejectsMissingReason(t *testing.T) {
	_, err := transitionServiceLifecycleArgs(TransitionServiceLifecycleInput{
		ID:               "service-1",
		TenantID:         tenant.ID("tenant-1"),
		ActorID:          audit.ActorID("admin-1"),
		Action:           ServiceLifecycleActionSuspend,
		FromStatus:       ServiceStatusActive,
		ToStatus:         ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonManualAdmin,
	})
	if !errors.Is(err, ErrServiceLifecycleReasonMissing) {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestTransitionServiceLifecycleSQLScopesTenantAndExpectedStatus(t *testing.T) {
	for _, clause := range []string{
		"UPDATE service_instances",
		"service_instance_id = $1",
		"tenant_id = $2",
		"status = $3",
		"term_end = $9::timestamptz",
		"billing_status = $10",
		"suspension_reason = $11",
		"RETURNING",
	} {
		if !strings.Contains(transitionServiceLifecycleSQL, clause) {
			t.Fatalf("expected %q in lifecycle SQL: %s", clause, transitionServiceLifecycleSQL)
		}
	}
}

func TestTransitionServiceLifecycleSQLEmitsLifecycleOutboxEvent(t *testing.T) {
	for _, clause := range []string{
		"WITH updated AS",
		"INSERT INTO outbox_events",
		ServiceAggregateType,
		"'from_status', $3::text",
		"'to_status', status",
		"'display_id', display_id",
		"FROM updated",
	} {
		if !strings.Contains(transitionServiceLifecycleSQL, clause) {
			t.Fatalf("expected %q in lifecycle SQL: %s", clause, transitionServiceLifecycleSQL)
		}
	}
}
