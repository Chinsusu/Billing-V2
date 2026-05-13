package notification

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	EventPasswordReset              = "auth.password_reset"
	EventTopupApproved              = "wallet.topup.approved"
	EventTopupRejected              = "wallet.topup.rejected"
	EventProvisioningFailed         = "provisioning.failed"
	EventProvisioningManualReview   = "provisioning.manual_review"
	EventServiceExpiring            = "service.expiring"
	EventServiceExpired             = "service.expired"
	EventServiceSuspended           = "service.suspended"
	EventServiceTerminated          = "service.terminated"
	EventServiceLifecycleTransition = "service.lifecycle"

	RecipientGroupAdminOps = "admin_ops"
)

type PasswordResetEventInput struct {
	TenantID      tenant.ID
	UserID        identity.UserID
	Email         string
	ExpiresAt     time.Time
	CorrelationID CorrelationID
}

type TopupStatusEventInput struct {
	TenantID       tenant.ID
	UserID         identity.UserID
	TopupID        ReferenceID
	TopupDisplayID int64
	Status         string
	AmountMinor    int64
	Currency       string
	ReviewNote     string
	CorrelationID  CorrelationID
}

type ProvisioningEventInput struct {
	TenantID           tenant.ID
	RecipientGroup     RecipientGroup
	OrderID            ReferenceID
	OrderDisplayID     int64
	JobID              ReferenceID
	JobDisplayID       int64
	ErrorCode          string
	ManualReviewReason string
	CorrelationID      CorrelationID
}

type ServiceLifecycleEventInput struct {
	TenantID         tenant.ID
	UserID           identity.UserID
	ServiceID        ReferenceID
	ServiceDisplayID int64
	DedupeWindow     string
	EventType        string
	Action           string
	Status           string
	BillingStatus    string
	TermEnd          time.Time
	Reason           string
	CorrelationID    CorrelationID
}

func PasswordResetEvent(input PasswordResetEventInput) QueueInput {
	return QueueInput{
		TenantID:        input.TenantID,
		RecipientUserID: input.UserID,
		Channel:         ChannelEmail,
		TemplateKey:     EventPasswordReset,
		EventType:       EventPasswordReset,
		Priority:        PriorityHigh,
		PayloadJSON: marshalPayload(map[string]any{
			"email":      input.Email,
			"expires_at": input.ExpiresAt,
		}),
		ReferenceType: "user",
		ReferenceID:   ReferenceID(input.UserID),
		DedupeKey:     "password_reset:" + string(input.UserID) + ":" + input.ExpiresAt.UTC().Format(time.RFC3339),
		CorrelationID: input.CorrelationID,
	}
}

func TopupStatusEvent(input TopupStatusEventInput) QueueInput {
	eventType := EventTopupApproved
	priority := PriorityNormal
	if input.Status == "rejected" {
		eventType = EventTopupRejected
		priority = PriorityHigh
	}
	return QueueInput{
		TenantID:        input.TenantID,
		RecipientUserID: input.UserID,
		Channel:         ChannelDashboard,
		TemplateKey:     eventType,
		EventType:       eventType,
		Priority:        priority,
		PayloadJSON: marshalPayload(map[string]any{
			"topup_display_id": input.TopupDisplayID,
			"status":           input.Status,
			"amount_minor":     input.AmountMinor,
			"currency":         input.Currency,
			"review_note":      input.ReviewNote,
		}),
		ReferenceType: "topup_request",
		ReferenceID:   input.TopupID,
		DedupeKey:     "topup_status:" + string(input.TopupID) + ":" + input.Status,
		CorrelationID: input.CorrelationID,
	}
}

func ProvisioningFailedEvent(input ProvisioningEventInput) QueueInput {
	return provisioningEvent(EventProvisioningFailed, PriorityCritical, input)
}

func ProvisioningManualReviewEvent(input ProvisioningEventInput) QueueInput {
	return provisioningEvent(EventProvisioningManualReview, PriorityCritical, input)
}

func ServiceLifecycleEvent(input ServiceLifecycleEventInput) QueueInput {
	eventType := input.EventType
	if eventType == "" {
		eventType = EventServiceLifecycleTransition
	}
	priority := PriorityNormal
	if eventType == EventServiceExpiring || eventType == EventServiceSuspended {
		priority = PriorityHigh
	}
	if eventType == EventServiceTerminated {
		priority = PriorityCritical
	}
	return QueueInput{
		TenantID:        input.TenantID,
		RecipientUserID: input.UserID,
		Channel:         ChannelDashboard,
		TemplateKey:     eventType,
		EventType:       eventType,
		Priority:        priority,
		PayloadJSON: marshalPayload(map[string]any{
			"service_display_id": input.ServiceDisplayID,
			"action":             input.Action,
			"status":             input.Status,
			"billing_status":     input.BillingStatus,
			"term_end":           input.TermEnd,
			"reason":             input.Reason,
		}),
		ReferenceType: "service",
		ReferenceID:   input.ServiceID,
		DedupeKey:     serviceLifecycleDedupeKey(input, eventType),
		CorrelationID: input.CorrelationID,
	}
}

func provisioningEvent(eventType string, priority Priority, input ProvisioningEventInput) QueueInput {
	recipientGroup := input.RecipientGroup
	if recipientGroup == "" {
		recipientGroup = RecipientGroupAdminOps
	}
	return QueueInput{
		TenantID:       input.TenantID,
		RecipientGroup: recipientGroup,
		Channel:        ChannelDashboard,
		TemplateKey:    eventType,
		EventType:      eventType,
		Priority:       priority,
		PayloadJSON: marshalPayload(map[string]any{
			"order_display_id":     input.OrderDisplayID,
			"job_display_id":       input.JobDisplayID,
			"error_code":           input.ErrorCode,
			"manual_review_reason": input.ManualReviewReason,
		}),
		ReferenceType: "job",
		ReferenceID:   input.JobID,
		DedupeKey:     eventType + ":" + string(input.JobID),
		CorrelationID: input.CorrelationID,
	}
}

func marshalPayload(value any) json.RawMessage {
	body, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(body)
}

func serviceLifecycleDedupeKey(input ServiceLifecycleEventInput, eventType string) string {
	window := input.DedupeWindow
	if window == "" && !input.TermEnd.IsZero() {
		window = input.TermEnd.UTC().Format(time.RFC3339)
	}
	if window == "" {
		window = input.Action
	}
	return "service_lifecycle:" + string(input.ServiceID) + ":" + eventType + ":" + window
}
