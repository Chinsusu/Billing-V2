package order

import (
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCredentialRevealRateLimitInputNormalizesAndValidates(t *testing.T) {
	windowStart := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	input := CredentialRevealRateLimitInput{
		TenantID:    tenant.ID(" tenant_1 "),
		ActorID:     identity.UserID(" actor_1 "),
		ServiceID:   ServiceID(" service_1 "),
		WindowStart: windowStart,
	}.Normalize()

	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid rate limit input: %v", err)
	}
	if input.TenantID != tenant.ID("tenant_1") ||
		input.ActorID != identity.UserID("actor_1") ||
		input.ServiceID != ServiceID("service_1") {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
}

func TestIncrementCredentialRevealRateLimitSQLUpsertsCounter(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO service_credential_reveal_rate_limits",
		"tenant_id, actor_id, service_instance_id, window_start",
		"ON CONFLICT (tenant_id, actor_id, service_instance_id, window_start) DO UPDATE",
		"attempt_count = service_credential_reveal_rate_limits.attempt_count + 1",
		"RETURNING",
	} {
		if !strings.Contains(incrementCredentialRevealRateLimitSQL, clause) {
			t.Fatalf("expected %q in reveal rate limit SQL: %s", clause, incrementCredentialRevealRateLimitSQL)
		}
	}
}
