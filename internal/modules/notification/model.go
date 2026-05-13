package notification

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrStoreMissing          = errors.New("notification store missing")
	ErrNotificationIDMissing = errors.New("notification id missing")
	ErrChannelInvalid        = errors.New("notification channel invalid")
	ErrStatusInvalid         = errors.New("notification status invalid")
	ErrPriorityInvalid       = errors.New("notification priority invalid")
	ErrTemplateKeyMissing    = errors.New("notification template key missing")
	ErrEventTypeMissing      = errors.New("notification event type missing")
	ErrRecipientMissing      = errors.New("notification recipient missing")
	ErrPayloadInvalid        = errors.New("notification payload invalid")
	ErrDedupeKeyMissing      = errors.New("notification dedupe key missing")
	ErrCorrelationIDMissing  = errors.New("notification correlation id missing")
	ErrNotificationNotFound  = errors.New("notification not found")
)

type ID string
type Channel string
type Status string
type Priority string
type ReferenceType string
type ReferenceID string
type RecipientGroup string
type CorrelationID string

const (
	ChannelEmail     Channel = "email"
	ChannelDashboard Channel = "dashboard"
	ChannelTelegram  Channel = "telegram"
	ChannelWebhook   Channel = "webhook"
)

func (channel Channel) Valid() bool {
	switch channel {
	case ChannelEmail, ChannelDashboard, ChannelTelegram, ChannelWebhook:
		return true
	default:
		return false
	}
}

const (
	StatusQueued    Status = "queued"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

func (status Status) Valid() bool {
	switch status {
	case StatusQueued, StatusSent, StatusFailed, StatusCancelled:
		return true
	default:
		return false
	}
}

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

func (priority Priority) Valid() bool {
	switch priority {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

type Notification struct {
	ID                       ID
	DisplayID                int64
	TenantID                 tenant.ID
	RecipientUserID          identity.UserID
	RecipientGroup           RecipientGroup
	Channel                  Channel
	TemplateKey              string
	EventType                string
	Priority                 Priority
	PayloadRedacted          json.RawMessage
	ReferenceType            ReferenceType
	ReferenceID              ReferenceID
	DedupeKey                string
	Status                   Status
	LastErrorCode            string
	LastErrorMessageRedacted string
	CorrelationID            CorrelationID
	SentAt                   time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type QueueInput struct {
	TenantID        tenant.ID
	RecipientUserID identity.UserID
	RecipientGroup  RecipientGroup
	Channel         Channel
	TemplateKey     string
	EventType       string
	Priority        Priority
	PayloadJSON     json.RawMessage
	ReferenceType   ReferenceType
	ReferenceID     ReferenceID
	DedupeKey       string
	CorrelationID   CorrelationID
}

type CreateInput struct {
	TenantID        tenant.ID
	RecipientUserID identity.UserID
	RecipientGroup  RecipientGroup
	Channel         Channel
	TemplateKey     string
	EventType       string
	Priority        Priority
	PayloadRedacted json.RawMessage
	ReferenceType   ReferenceType
	ReferenceID     ReferenceID
	DedupeKey       string
	CorrelationID   CorrelationID
}

type Store interface {
	CreateNotification(ctx context.Context, input CreateInput) (Notification, error)
}

type DeliveryStore interface {
	MarkNotificationSent(ctx context.Context, id ID, sentAt time.Time) (Notification, error)
	MarkNotificationFailed(ctx context.Context, id ID, failedAt time.Time, errorCode string, errorMessageRedacted string) (Notification, error)
}

func (input QueueInput) Normalize() QueueInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.RecipientUserID = identity.UserID(trim(string(output.RecipientUserID)))
	output.RecipientGroup = RecipientGroup(trim(string(output.RecipientGroup)))
	output.Channel = Channel(trim(string(output.Channel)))
	output.TemplateKey = trim(output.TemplateKey)
	output.EventType = trim(output.EventType)
	output.Priority = Priority(trim(string(output.Priority)))
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.DedupeKey = trim(output.DedupeKey)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	output.PayloadJSON = defaultPayload(output.PayloadJSON)
	if output.Priority == "" {
		output.Priority = PriorityNormal
	}
	if output.CorrelationID == "" {
		output.CorrelationID = defaultCorrelationID(output)
	}
	if output.DedupeKey == "" {
		output.DedupeKey = defaultDedupeKey(output)
	}
	return output
}

func (input QueueInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.RecipientUserID == "" && input.RecipientGroup == "" {
		return ErrRecipientMissing
	}
	if !input.Channel.Valid() {
		return ErrChannelInvalid
	}
	if input.TemplateKey == "" {
		return ErrTemplateKeyMissing
	}
	if input.EventType == "" {
		return ErrEventTypeMissing
	}
	if !input.Priority.Valid() {
		return ErrPriorityInvalid
	}
	if input.DedupeKey == "" {
		return ErrDedupeKeyMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	if _, err := RedactPayload(input.PayloadJSON); err != nil {
		return err
	}
	return nil
}

func (input CreateInput) Normalize() CreateInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.RecipientUserID = identity.UserID(trim(string(output.RecipientUserID)))
	output.RecipientGroup = RecipientGroup(trim(string(output.RecipientGroup)))
	output.Channel = Channel(trim(string(output.Channel)))
	output.TemplateKey = trim(output.TemplateKey)
	output.EventType = trim(output.EventType)
	output.Priority = Priority(trim(string(output.Priority)))
	output.PayloadRedacted = defaultPayload(output.PayloadRedacted)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.DedupeKey = trim(output.DedupeKey)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateInput) Validate() error {
	queueInput := QueueInput{
		TenantID:        input.TenantID,
		RecipientUserID: input.RecipientUserID,
		RecipientGroup:  input.RecipientGroup,
		Channel:         input.Channel,
		TemplateKey:     input.TemplateKey,
		EventType:       input.EventType,
		Priority:        input.Priority,
		PayloadJSON:     input.PayloadRedacted,
		ReferenceType:   input.ReferenceType,
		ReferenceID:     input.ReferenceID,
		DedupeKey:       input.DedupeKey,
		CorrelationID:   input.CorrelationID,
	}
	return queueInput.Validate()
}

func (notification Notification) Validate() error {
	if notification.ID == "" {
		return ErrNotificationIDMissing
	}
	input := CreateInput{
		TenantID:        notification.TenantID,
		RecipientUserID: notification.RecipientUserID,
		RecipientGroup:  notification.RecipientGroup,
		Channel:         notification.Channel,
		TemplateKey:     notification.TemplateKey,
		EventType:       notification.EventType,
		Priority:        notification.Priority,
		PayloadRedacted: notification.PayloadRedacted,
		ReferenceType:   notification.ReferenceType,
		ReferenceID:     notification.ReferenceID,
		DedupeKey:       notification.DedupeKey,
		CorrelationID:   notification.CorrelationID,
	}
	if err := input.Validate(); err != nil {
		return err
	}
	if !notification.Status.Valid() {
		return ErrStatusInvalid
	}
	return nil
}

func trim(value string) string {
	return strings.TrimSpace(value)
}

func defaultPayload(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return append(json.RawMessage(nil), value...)
}

func defaultCorrelationID(input QueueInput) CorrelationID {
	if input.ReferenceID != "" {
		return CorrelationID(input.ReferenceID)
	}
	if input.RecipientUserID != "" {
		return CorrelationID(input.RecipientUserID)
	}
	return CorrelationID(input.TenantID)
}

func defaultDedupeKey(input QueueInput) string {
	recipient := string(input.RecipientGroup)
	if input.RecipientUserID != "" {
		recipient = string(input.RecipientUserID)
	}
	parts := []string{
		string(input.EventType),
		string(input.Channel),
		recipient,
		string(input.ReferenceType),
		string(input.ReferenceID),
	}
	return strings.Join(parts, ":")
}
