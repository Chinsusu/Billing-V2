package notification

import (
	"context"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestLocalDeliveryHandlerMarksNotificationSent(t *testing.T) {
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	store := &fakeNotificationDeliveryStore{}
	handler := NewLocalDeliveryHandler(store)
	handler.Now = func() time.Time { return now }

	completion, err := handler.Handle(context.Background(), jobs.Job{
		ID:            "job-1",
		Type:          DeliveryJobType,
		ReferenceType: DeliveryReferenceType,
		ReferenceID:   "44444444-4444-4444-4444-444444444444",
	})
	if err != nil {
		t.Fatalf("expected local delivery success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded || store.sentID != "44444444-4444-4444-4444-444444444444" ||
		!store.sentAt.Equal(now) {
		t.Fatalf("unexpected completion/store: %+v store=%+v", completion, store)
	}
}

func TestLocalDeliveryHandlerFailsInvalidJobTerminally(t *testing.T) {
	handler := NewLocalDeliveryHandler(&fakeNotificationDeliveryStore{})

	completion, err := handler.Handle(context.Background(), jobs.Job{
		ID:            "job-1",
		Type:          DeliveryJobType,
		ReferenceType: "order",
		ReferenceID:   "11111111-1111-1111-1111-111111111111",
	})
	if err != nil {
		t.Fatalf("expected terminal completion, got error: %v", err)
	}
	if completion.Status != jobs.StatusFailedTerminal || completion.LastErrorCode != "notification_job_invalid" {
		t.Fatalf("unexpected invalid job completion: %+v", completion)
	}
}

type fakeNotificationDeliveryStore struct {
	sentID ID
	sentAt time.Time
}

func (store *fakeNotificationDeliveryStore) MarkNotificationSent(ctx context.Context, id ID, sentAt time.Time) (Notification, error) {
	store.sentID = id
	store.sentAt = sentAt
	return Notification{ID: id, Status: StatusSent}, nil
}

func (store *fakeNotificationDeliveryStore) MarkNotificationFailed(ctx context.Context, id ID, failedAt time.Time, errorCode string, errorMessageRedacted string) (Notification, error) {
	return Notification{ID: id, Status: StatusFailed}, nil
}
