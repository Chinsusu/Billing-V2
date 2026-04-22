package jobs

import (
	"strings"
	"testing"
	"time"
)

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
