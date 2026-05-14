package jobs

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateJobArgsNormalizeAndValidate(t *testing.T) {
	args, err := createJobArgs(CreateJobInput{
		TenantID:       tenant.ID(" tenant-1 "),
		Type:           Type(" provider.provision "),
		ReferenceType:  ReferenceType(" order "),
		ReferenceID:    ReferenceID(" order-1 "),
		SourceID:       SourceID(" source-1 "),
		IdempotencyKey: " provision-key-1 ",
		CorrelationID:  CorrelationID(" correlation-1 "),
	})
	if err != nil {
		t.Fatalf("expected create job args: %v", err)
	}
	if len(args) != 10 {
		t.Fatalf("expected 10 args, got %d", len(args))
	}
	if args[0] != tenant.ID("tenant-1") || args[1] != Type("provider.provision") ||
		args[6] != 100 || args[8] != 5 || args[5] != "{}" {
		t.Fatalf("unexpected create job args: %#v", args)
	}
}

func TestCreateJobArgsRejectsMissingTenant(t *testing.T) {
	_, err := createJobArgs(CreateJobInput{
		Type:           Type("provider.provision"),
		ReferenceType:  ReferenceType("order"),
		ReferenceID:    ReferenceID("order-1"),
		IdempotencyKey: "provision-key-1",
		CorrelationID:  CorrelationID("correlation-1"),
	})
	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestCreateJobSQLIsIdempotentByTenantTypeAndKey(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO jobs",
		"ON CONFLICT (tenant_id, job_type, idempotency_key) WHERE tenant_id IS NOT NULL",
		"RETURNING",
	} {
		if !strings.Contains(createJobSQL, clause) {
			t.Fatalf("expected %q in create job SQL: %s", clause, createJobSQL)
		}
	}
}

func TestClaimJobsQueryAddsTypeFilterOnlyWhenRequested(t *testing.T) {
	request := ClaimRequest{
		WorkerID: "worker_1",
		Limit:    10,
		LockFor:  time.Minute,
		Types:    []Type{"provider.provision", ""},
	}

	query, args := claimJobsQuery(request, fixedRunnerTime())

	if !strings.Contains(query, "job_type = ANY($5)") {
		t.Fatalf("expected type filter in query: %s", query)
	}
	if len(args) != 5 {
		t.Fatalf("expected type filter arg, got %d args", len(args))
	}
}

func TestClaimJobsQueryOmitsTypeFilterWhenEmpty(t *testing.T) {
	request := ClaimRequest{WorkerID: "worker_1", Limit: 10, LockFor: time.Minute}

	query, args := claimJobsQuery(request, fixedRunnerTime())

	if strings.Contains(query, "job_type = ANY") {
		t.Fatalf("did not expect type filter in query: %s", query)
	}
	if len(args) != 4 {
		t.Fatalf("expected base args, got %d args", len(args))
	}
}

func TestClaimJobsQueryQualifiesReturningColumns(t *testing.T) {
	request := ClaimRequest{WorkerID: "worker_1", Limit: 10, LockFor: time.Minute}

	query, _ := claimJobsQuery(request, fixedRunnerTime())

	if !strings.Contains(query, "RETURNING job.job_id, job.display_id") {
		t.Fatalf("expected qualified job returning columns: %s", query)
	}
}

func TestCompletionValidateRejectsClaimState(t *testing.T) {
	completion := Completion{Status: StatusClaimed}

	if err := completion.Validate(); err == nil {
		t.Fatal("expected invalid completion status")
	}
}

func TestOutboxCompletionValidateRejectsProcessingState(t *testing.T) {
	completion := OutboxCompletion{Status: OutboxStatusProcessing}

	if err := completion.Validate(); err == nil {
		t.Fatal("expected invalid outbox completion status")
	}
}
