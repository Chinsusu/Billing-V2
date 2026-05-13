package notification

import (
	"strings"
	"testing"
)

func TestRedactPayloadRemovesSensitiveFields(t *testing.T) {
	payload, err := RedactPayload([]byte(`{
		"service_display_id": 42001,
		"credential": {"username": "root", "password": "secret"},
		"nested": [{"api_key": "key-1"}, {"safe": "value"}],
		"reset_token": "reset-token"
	}`))
	if err != nil {
		t.Fatalf("expected payload redaction: %v", err)
	}
	body := string(payload)
	for _, leaked := range []string{"secret", "key-1", "reset-token"} {
		if strings.Contains(body, leaked) {
			t.Fatalf("payload leaked %q: %s", leaked, body)
		}
	}
	if !strings.Contains(body, redactedValue) || !strings.Contains(body, `"safe":"value"`) {
		t.Fatalf("expected redacted sensitive values and preserved safe values: %s", body)
	}
}

func TestQueueInputNormalizeDefaultsDedupeAndCorrelation(t *testing.T) {
	input := QueueInput{
		TenantID:        " 11111111-1111-1111-1111-111111111111 ",
		RecipientUserID: " 22222222-2222-2222-2222-222222222222 ",
		Channel:         ChannelDashboard,
		TemplateKey:     EventServiceExpired,
		EventType:       EventServiceExpired,
		ReferenceType:   " service ",
		ReferenceID:     " 33333333-3333-3333-3333-333333333333 ",
	}.Normalize()

	if input.Priority != PriorityNormal {
		t.Fatalf("expected default priority, got %q", input.Priority)
	}
	if input.CorrelationID != "33333333-3333-3333-3333-333333333333" {
		t.Fatalf("expected reference correlation id, got %q", input.CorrelationID)
	}
	if input.DedupeKey != "service.expired:dashboard:22222222-2222-2222-2222-222222222222:service:33333333-3333-3333-3333-333333333333" {
		t.Fatalf("unexpected dedupe key: %q", input.DedupeKey)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid queue input: %v", err)
	}
}
