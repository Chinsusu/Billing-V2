package audit

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListLogsQueryAddsFilters(t *testing.T) {
	createdFrom := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	createdTo := createdFrom.Add(24 * time.Hour)
	query, args, err := buildListLogsQuery(Filter{
		TenantID:        tenant.ID("tenant-1"),
		ActorID:         ActorID("actor-1"),
		ActorDisplayID:  10002,
		ActorType:       ActorTypeUser,
		DisplayID:       70001,
		Action:          "invoice.paid",
		TargetType:      "invoice",
		TargetID:        TargetID("target-1"),
		TargetDisplayID: 44001,
		CreatedFrom:     createdFrom,
		CreatedTo:       createdTo,
		Limit:           25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"tenant_id = $1",
		"actor_display_id",
		"target_display_id",
		"actor_id = $2",
		"actor.display_id = $3",
		"actor_type = $4",
		"display_id = $5",
		"action = $6",
		"target_type = $7",
		"target_id = $8",
		"inv.display_id = $9",
		"ord.display_id = $9",
		"job.display_id = $9",
		"topup.display_id = $9",
		"svc.display_id = $9",
		"source.display_id = $9",
		"created_at >= $10",
		"created_at <= $11",
		"LIMIT $12",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 12 || args[11] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListLogsQueryRejectsBadWindow(t *testing.T) {
	_, _, err := buildListLogsQuery(Filter{
		TenantID:    tenant.ID("tenant-1"),
		CreatedFrom: time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC),
		CreatedTo:   time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrCreatedWindowInvalid) {
		t.Fatalf("expected created window error, got %v", err)
	}
}

func TestBuildGetLogQueryUsesTenantScope(t *testing.T) {
	query, args, err := buildGetLogQuery(Lookup{ID: ID("audit-1"), TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"audit_id = $1", "tenant_id = $2"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 {
		t.Fatalf("unexpected args: %#v", args)
	}
}
