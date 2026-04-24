package main

import (
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestAddRunSummary(t *testing.T) {
	total := jobs.RunSummary{Claimed: 1, Succeeded: 1}
	addRunSummary(&total, jobs.RunSummary{
		Claimed:        2,
		Succeeded:      1,
		Retried:        1,
		ManualReview:   1,
		TerminalFailed: 1,
		Cancelled:      1,
	})

	if total.Claimed != 3 ||
		total.Succeeded != 2 ||
		total.Retried != 1 ||
		total.ManualReview != 1 ||
		total.TerminalFailed != 1 ||
		total.Cancelled != 1 {
		t.Fatalf("unexpected total summary: %+v", total)
	}
}

func TestProvisioningJobFailureIsActionable(t *testing.T) {
	err := provisioningJobFailure(42001, provisioningJobSmokeRecord{
		DisplayID:                81001,
		Status:                   "manual_review",
		AttemptCount:             2,
		LastErrorCode:            "provider_timeout",
		LastErrorMessageRedacted: "provider timed out",
	}, jobs.RunSummary{Claimed: 1, ManualReview: 1})

	message := err.Error()
	for _, expected := range []string{"81001", "42001", "manual_review", "provider_timeout", "provider timed out"} {
		if !strings.Contains(message, expected) {
			t.Fatalf("expected failure message to contain %q, got %s", expected, message)
		}
	}
}
