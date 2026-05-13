package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

const (
	DeliveryJobType       jobs.Type          = "notification.deliver"
	DeliveryReferenceType jobs.ReferenceType = "notification"
)

type Service struct {
	store Store
	queue jobs.QueueStore
}

type DeliveryJobPayload struct {
	NotificationID        ID              `json:"notification_id"`
	NotificationDisplayID int64           `json:"notification_display_id"`
	EventType             string          `json:"event_type"`
	TemplateKey           string          `json:"template_key"`
	Channel               Channel         `json:"channel"`
	RecipientUserID       string          `json:"recipient_user_id,omitempty"`
	RecipientGroup        RecipientGroup  `json:"recipient_group,omitempty"`
	PayloadRedacted       json.RawMessage `json:"payload_redacted"`
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func NewServiceWithJobs(store Store, queue jobs.QueueStore) *Service {
	return &Service{store: store, queue: queue}
}

func (service *Service) Queue(ctx context.Context, input QueueInput) (Notification, error) {
	if service == nil || service.store == nil {
		return Notification{}, ErrStoreMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Notification{}, err
	}
	payloadRedacted, err := RedactPayload(input.PayloadJSON)
	if err != nil {
		return Notification{}, err
	}
	notification, err := service.store.CreateNotification(ctx, CreateInput{
		TenantID:        input.TenantID,
		RecipientUserID: input.RecipientUserID,
		RecipientGroup:  input.RecipientGroup,
		Channel:         input.Channel,
		TemplateKey:     input.TemplateKey,
		EventType:       input.EventType,
		Priority:        input.Priority,
		PayloadRedacted: payloadRedacted,
		ReferenceType:   input.ReferenceType,
		ReferenceID:     input.ReferenceID,
		DedupeKey:       input.DedupeKey,
		CorrelationID:   input.CorrelationID,
	})
	if err != nil {
		return Notification{}, err
	}
	if service.queue == nil {
		return notification, nil
	}
	if _, err := service.queue.CreateJob(ctx, deliveryJobInput(notification)); err != nil {
		return Notification{}, fmt.Errorf("queue notification delivery job: %w", err)
	}
	return notification, nil
}

func deliveryJobInput(notification Notification) jobs.CreateJobInput {
	payload := deliveryPayloadJSON(notification)
	return jobs.CreateJobInput{
		TenantID:       notification.TenantID,
		Type:           DeliveryJobType,
		ReferenceType:  DeliveryReferenceType,
		ReferenceID:    jobs.ReferenceID(notification.ID),
		PayloadJSON:    payload,
		Priority:       deliveryPriority(notification.Priority),
		IdempotencyKey: "notification:" + string(notification.ID) + ":" + string(notification.Channel),
		MaxAttempts:    5,
		CorrelationID:  jobs.CorrelationID(notification.CorrelationID),
	}
}

func deliveryPayloadJSON(notification Notification) json.RawMessage {
	body, err := json.Marshal(DeliveryJobPayload{
		NotificationID:        notification.ID,
		NotificationDisplayID: notification.DisplayID,
		EventType:             notification.EventType,
		TemplateKey:           notification.TemplateKey,
		Channel:               notification.Channel,
		RecipientUserID:       string(notification.RecipientUserID),
		RecipientGroup:        notification.RecipientGroup,
		PayloadRedacted:       notification.PayloadRedacted,
	})
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(body)
}

func deliveryPriority(priority Priority) int {
	switch priority {
	case PriorityCritical:
		return 10
	case PriorityHigh:
		return 30
	case PriorityNormal:
		return 70
	default:
		return 100
	}
}
