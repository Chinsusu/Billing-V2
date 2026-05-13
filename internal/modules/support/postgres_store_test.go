package support

import (
	"strings"
	"testing"
)

func TestSupportSQLUsesTenantScopedTablesAndReturningColumns(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO support_tickets",
		"tenant_id",
		"requester_user_id",
		"correlation_id",
		"RETURNING " + supportTicketColumns,
	} {
		if !strings.Contains(createSupportTicketSQL, clause) {
			t.Fatalf("expected %q in support ticket SQL: %s", clause, createSupportTicketSQL)
		}
	}
}

func TestRiskFlagSQLStoresRedactedNotesAndScopedTargets(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO risk_flags",
		"service_instance_id",
		"order_id",
		"note_redacted",
		"created_by",
		"RETURNING " + riskFlagColumns,
	} {
		if !strings.Contains(createRiskFlagSQL, clause) {
			t.Fatalf("expected %q in risk flag SQL: %s", clause, createRiskFlagSQL)
		}
	}
}

func TestAbuseCaseSQLStoresRedactedEvidenceAndTenantLookup(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO abuse_cases",
		"evidence_summary_redacted",
		"provider_source_id",
		"correlation_id",
		"RETURNING " + abuseCaseColumns,
	} {
		if !strings.Contains(createAbuseCaseSQL, clause) {
			t.Fatalf("expected %q in abuse case SQL: %s", clause, createAbuseCaseSQL)
		}
	}
	for _, clause := range []string{
		"FROM abuse_cases",
		"abuse_case_id = $1",
		"tenant_id = $2",
	} {
		if !strings.Contains(getAbuseCaseSQL, clause) {
			t.Fatalf("expected %q in abuse lookup SQL: %s", clause, getAbuseCaseSQL)
		}
	}
}

func TestMarkAbuseCaseSuspendedSQLIsTenantScoped(t *testing.T) {
	for _, clause := range []string{
		"UPDATE abuse_cases",
		"status = 'suspended'",
		"action_taken = $3",
		"abuse_case_id = $1",
		"tenant_id = $2",
		"RETURNING " + abuseCaseColumns,
	} {
		if !strings.Contains(markAbuseCaseSuspendedSQL, clause) {
			t.Fatalf("expected %q in abuse suspend SQL: %s", clause, markAbuseCaseSuspendedSQL)
		}
	}
}
