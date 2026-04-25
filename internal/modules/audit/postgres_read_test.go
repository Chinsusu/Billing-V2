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
		TenantID:    tenant.ID("tenant-1"),
		ActorID:     ActorID("actor-1"),
		ActorType:   ActorTypeUser,
		DisplayID:   70001,
		Action:      "invoice.paid",
		TargetType:  "invoice",
		TargetID:    TargetID("target-1"),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Limit:       25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"tenant_id = $1",
		"actor_display_id",
		"target_display_id",
		"actor_id = $2",
		"actor_type = $3",
		"display_id = $4",
		"action = $5",
		"target_type = $6",
		"target_id = $7",
		"created_at >= $8",
		"created_at <= $9",
		"LIMIT $10",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 10 || args[9] != 25 {
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
