package jobs

import (
	"strings"
	"testing"
	"time"
)

func TestRetryJobSQLAllowsOnlyRetryableOrManualReview(t *testing.T) {
	for _, clause := range []string{
		"UPDATE jobs",
		"tenant_id = $2",
		"status IN ('failed_retryable', 'manual_review')",
		"RETURNING",
	} {
		if !strings.Contains(retryJobSQL, clause) {
			t.Fatalf("expected %q in retry SQL: %s", clause, retryJobSQL)
		}
	}
}

func TestManualReviewAndCancelSQLAvoidActiveWorkers(t *testing.T) {
	for _, sql := range []string{markManualReviewJobSQL, cancelJobSQL} {
		for _, blocked := range []string{"'claimed'", "'running'", "'succeeded'"} {
			if strings.Contains(sql, blocked) {
				t.Fatalf("did not expect active/succeeded state %s in SQL: %s", blocked, sql)
			}
		}
		if !strings.Contains(sql, "tenant_id = $2") {
			t.Fatalf("expected tenant scope in SQL: %s", sql)
		}
	}
}

func TestRetryNextAttemptDefaultsToNow(t *testing.T) {
	now := time.Date(2026, 4, 24, 3, 0, 0, 0, time.UTC)

	if got := retryNextAttemptAt(RetryJobInput{Now: now}); !got.Equal(now) {
		t.Fatalf("expected retry time to default to now, got %v", got)
	}
}
