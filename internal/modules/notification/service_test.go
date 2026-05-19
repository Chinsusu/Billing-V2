package notification

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestServiceQueueCreatesRedactedNotificationAndDeliveryJob(t *testing.T) {
	store := &fakeNotificationStore{}
	queue := &fakeNotificationJobQueue{}
	service := NewServiceWithJobs(store, queue)

	created, err := service.Queue(context.Background(), QueueInput{
		TenantID:        "11111111-1111-1111-1111-111111111111",
		RecipientUserID: "22222222-2222-2222-2222-222222222222",
		Channel:         ChannelDashboard,
		TemplateKey:     EventServiceSuspended,
		EventType:       EventServiceSuspended,
		Priority:        PriorityHigh,
		PayloadJSON:     json.RawMessage(`{"service_display_id":61001,"root_password":"secret","safe":"value"}`),
		ReferenceType:   "service",
		ReferenceID:     "33333333-3333-3333-3333-333333333333",
	})
	if err != nil {
		t.Fatalf("expected queued notification: %v", err)
	}
	if created.ID == "" || queue.input.Type != DeliveryJobType || queue.input.ReferenceType != DeliveryReferenceType {
		t.Fatalf("expected delivery job for notification: notification=%+v job=%+v", created, queue.input)
	}
	if strings.Contains(string(store.input.PayloadRedacted), "secret") || !strings.Contains(string(store.input.PayloadRedacted), redactedValue) {
		t.Fatalf("expected redacted stored payload: %s", string(store.input.PayloadRedacted))
	}
	if strings.Contains(string(queue.input.PayloadJSON), "secret") || !strings.Contains(string(queue.input.PayloadJSON), redactedValue) {
		t.Fatalf("expected redacted job payload: %s", string(queue.input.PayloadJSON))
	}
	if queue.input.Priority != 30 {
		t.Fatalf("expected high priority delivery job, got %+v", queue.input)
	}
}

func TestServiceQueueUsesTelegramDeliveryJobTypeForTelegramChannel(t *testing.T) {
	store := &fakeNotificationStore{}
	queue := &fakeNotificationJobQueue{}
	service := NewServiceWithJobs(store, queue)

	_, err := service.Queue(context.Background(), QueueInput{
		TenantID:       "11111111-1111-1111-1111-111111111111",
		RecipientGroup: RecipientGroupAdminOps,
		Channel:        ChannelTelegram,
		TemplateKey:    EventProvisioningManualReview,
		EventType:      EventProvisioningManualReview,
		Priority:       PriorityCritical,
		PayloadJSON:    json.RawMessage(`{"job_display_id":82001}`),
		ReferenceType:  "job",
		ReferenceID:    "33333333-3333-3333-3333-333333333333",
	})
	if err != nil {
		t.Fatalf("expected queued telegram notification: %v", err)
	}
	if queue.input.Type != TelegramDeliveryJobType {
		t.Fatalf("expected telegram delivery job type, got %+v", queue.input)
	}
	if queue.input.Priority != 10 {
		t.Fatalf("expected critical telegram priority, got %+v", queue.input)
	}
}

func TestServiceQueueRejectsInvalidPayload(t *testing.T) {
	service := NewService(&fakeNotificationStore{})

	_, err := service.Queue(context.Background(), QueueInput{
		TenantID:        "11111111-1111-1111-1111-111111111111",
		RecipientUserID: "22222222-2222-2222-2222-222222222222",
		Channel:         ChannelDashboard,
		TemplateKey:     EventServiceExpired,
		EventType:       EventServiceExpired,
		PayloadJSON:     json.RawMessage(`{`),
		ReferenceType:   "service",
		ReferenceID:     "33333333-3333-3333-3333-333333333333",
	})
	if err != ErrPayloadInvalid {
		t.Fatalf("expected invalid payload error, got %v", err)
	}
}

type fakeNotificationStore struct {
	input CreateInput
}

func (store *fakeNotificationStore) CreateNotification(ctx context.Context, input CreateInput) (Notification, error) {
	store.input = input
	return Notification{
		ID:              "44444444-4444-4444-4444-444444444444",
		DisplayID:       10001,
		TenantID:        input.TenantID,
		RecipientUserID: input.RecipientUserID,
		RecipientGroup:  input.RecipientGroup,
		Channel:         input.Channel,
		TemplateKey:     input.TemplateKey,
		EventType:       input.EventType,
		Priority:        input.Priority,
		PayloadRedacted: input.PayloadRedacted,
		ReferenceType:   input.ReferenceType,
		ReferenceID:     input.ReferenceID,
		DedupeKey:       input.DedupeKey,
		Status:          StatusQueued,
		CorrelationID:   input.CorrelationID,
	}, nil
}

type fakeNotificationJobQueue struct {
	input jobs.CreateJobInput
}

func (queue *fakeNotificationJobQueue) CreateJob(ctx context.Context, input jobs.CreateJobInput) (jobs.Job, error) {
	queue.input = input
	return jobs.Job{ID: "job-1"}, nil
}
