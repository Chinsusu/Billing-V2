package notification

import (
	"strings"
	"testing"
	"time"
)

func TestPasswordResetEventDoesNotPersistTokenMaterial(t *testing.T) {
	event := PasswordResetEvent(PasswordResetEventInput{
		TenantID:      "11111111-1111-1111-1111-111111111111",
		UserID:        "22222222-2222-2222-2222-222222222222",
		Email:         "client@example.com",
		ExpiresAt:     time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC),
		CorrelationID: "22222222-2222-2222-2222-222222222222",
	})

	if event.EventType != EventPasswordReset || event.TemplateKey != EventPasswordReset || event.Channel != ChannelEmail {
		t.Fatalf("unexpected password reset event: %+v", event)
	}
	if strings.Contains(strings.ToLower(string(event.PayloadJSON)), "token") ||
		strings.Contains(strings.ToLower(string(event.PayloadJSON)), "reset_link") {
		t.Fatalf("password reset payload must not contain token material: %s", string(event.PayloadJSON))
	}
}

func TestProvisioningManualReviewEventTargetsAdminOps(t *testing.T) {
	event := ProvisioningManualReviewEvent(ProvisioningEventInput{
		TenantID:           "11111111-1111-1111-1111-111111111111",
		JobID:              "33333333-3333-3333-3333-333333333333",
		JobDisplayID:       53001,
		OrderDisplayID:     42001,
		ErrorCode:          "PROVIDER_TIMEOUT_REQUEST_KNOWN",
		ManualReviewReason: "provider request outcome is unknown",
	})

	if event.EventType != EventProvisioningManualReview || event.RecipientGroup != RecipientGroupAdminOps ||
		event.Priority != PriorityCritical {
		t.Fatalf("unexpected provisioning event: %+v", event)
	}
	if err := event.Normalize().Validate(); err != nil {
		t.Fatalf("expected valid provisioning notification event: %v", err)
	}
}

func TestServiceLifecycleEventUsesLifecycleSpecificPriority(t *testing.T) {
	event := ServiceLifecycleEvent(ServiceLifecycleEventInput{
		TenantID:         "11111111-1111-1111-1111-111111111111",
		UserID:           "22222222-2222-2222-2222-222222222222",
		ServiceID:        "44444444-4444-4444-4444-444444444444",
		ServiceDisplayID: 61001,
		DedupeWindow:     "final",
		EventType:        EventServiceTerminated,
		Action:           "terminate",
		Status:           "terminated",
		Reason:           "expired grace period ended",
	})

	if event.Priority != PriorityCritical || event.DedupeKey != "service_lifecycle:44444444-4444-4444-4444-444444444444:service.terminated:final" {
		t.Fatalf("unexpected lifecycle event: %+v", event)
	}
}
