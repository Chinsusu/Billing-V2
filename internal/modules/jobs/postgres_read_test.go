package jobs

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListJobsQueryAddsFilters(t *testing.T) {
	query, args, err := buildListJobsQuery(Filter{
		TenantID:        "tenant_1",
		DisplayID:       81001,
		Type:            "provider.provision",
		Status:          StatusFailedRetryable,
		ReferenceType:   "order",
		ReferenceID:     "order_1",
		SourceID:        "source_1",
		SourceDisplayID: 10002,
		Limit:           25,
	})
	if err != nil {
		t.Fatalf("expected list query: %v", err)
	}
	for _, clause := range []string{
		"FROM jobs",
		"tenant_id = $1",
		"display_id = $2",
		"job_type = $3",
		"status = $4",
		"reference_type = $5",
		"reference_id = $6",
		"source_id = $7",
		"source.source_id = jobs.source_id",
		"source.display_id = $8",
		"ORDER BY created_at DESC",
		"LIMIT $9",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 9 || args[0] != tenant.ID("tenant_1") || args[8] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildGetJobQueryScopesTenant(t *testing.T) {
	query, args, err := buildGetJobQuery(Lookup{ID: "job_1", TenantID: "tenant_1"})
	if err != nil {
		t.Fatalf("expected get query: %v", err)
	}
	for _, clause := range []string{"job_id = $1", "tenant_id = $2"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[0] != ID("job_1") || args[1] != tenant.ID("tenant_1") {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListJobsQueryRejectsMissingTenant(t *testing.T) {
	_, _, err := buildListJobsQuery(Filter{Status: StatusQueued})

	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestBuildJobSummaryQueryScopesTenantAndType(t *testing.T) {
	query, args, err := buildJobSummaryQuery(SummaryFilter{TenantID: "tenant_1", Type: "provider.provision"})
	if err != nil {
		t.Fatalf("expected summary query: %v", err)
	}
	for _, clause := range []string{
		"WITH scoped AS",
		"tenant_id = $1",
		"job_type = $2",
		"COUNT(*) FILTER (WHERE status = 'queued')",
		"COUNT(*) FILTER (WHERE status = 'failed_retryable')",
		"MIN(created_at) FILTER (WHERE status = 'queued')",
		"latest_failure AS",
		"status IN ('failed_retryable', 'failed_terminal', 'manual_review')",
		"LEFT JOIN latest_failure ON TRUE",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[0] != tenant.ID("tenant_1") || args[1] != Type("provider.provision") {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListAttemptsQueryScopesThroughParentJob(t *testing.T) {
	query, args, err := buildListAttemptsQuery(AttemptFilter{JobID: "job_1", TenantID: "tenant_1", Limit: 15})
	if err != nil {
		t.Fatalf("expected attempts query: %v", err)
	}
	for _, clause := range []string{
		"FROM job_attempts attempt",
		"JOIN jobs job ON job.job_id = attempt.job_id",
		"attempt.job_id = $1",
		"job.tenant_id = $2",
		"ORDER BY attempt.attempt_number DESC",
		"LIMIT $3",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 || args[0] != ID("job_1") || args[1] != tenant.ID("tenant_1") || args[2] != 15 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildAttemptJobVisibleQueryScopesTenant(t *testing.T) {
	query, args, err := buildAttemptJobVisibleQuery(AttemptFilter{JobID: "job_1", TenantID: "tenant_1", Limit: 15})
	if err != nil {
		t.Fatalf("expected parent job query: %v", err)
	}
	for _, clause := range []string{
		"SELECT 1",
		"FROM jobs",
		"job_id = $1",
		"tenant_id = $2",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[0] != ID("job_1") || args[1] != tenant.ID("tenant_1") {
		t.Fatalf("unexpected args: %#v", args)
	}
}
