package notification

import (
	"strings"
	"testing"
)

func TestCreateNotificationSQLIsIdempotentByTenantChannelDedupe(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO notifications",
		"ON CONFLICT (tenant_id, channel, dedupe_key)",
		"payload_redacted",
		"RETURNING",
	} {
		if !strings.Contains(createNotificationSQL, clause) {
			t.Fatalf("expected %q in create notification SQL: %s", clause, createNotificationSQL)
		}
	}
}
