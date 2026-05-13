package support

import (
	"context"
	"encoding/json"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
)

const (
	auditActionSupportTicketCreated = "support.ticket.created"
	auditActionSupportNoteCreated   = "support.ticket.note.created"
	auditActionRiskFlagCreated      = "risk.flag.created"
	auditActionAbuseCaseCreated     = "abuse.case.created"
	auditActionAbuseServiceSuspend  = "abuse.service.suspended"
)

type AuditAppender interface {
	Append(ctx context.Context, input audit.AppendInput) (audit.Log, error)
}

type ticketAuditMetadata struct {
	DisplayID       int64          `json:"display_id"`
	RequesterUserID string         `json:"requester_user_id,omitempty"`
	Category        TicketCategory `json:"category,omitempty"`
	Priority        TicketPriority `json:"priority,omitempty"`
	ReferenceType   ReferenceType  `json:"reference_type,omitempty"`
	ReferenceID     ReferenceID    `json:"reference_id,omitempty"`
}

type ticketNoteAuditMetadata struct {
	DisplayID  int64          `json:"display_id"`
	TicketID   TicketID       `json:"ticket_id"`
	Visibility NoteVisibility `json:"visibility"`
}

type riskFlagAuditMetadata struct {
	DisplayID int64         `json:"display_id"`
	FlagType  RiskFlagType  `json:"flag_type"`
	Severity  AbuseSeverity `json:"severity"`
	UserID    string        `json:"user_id,omitempty"`
	ServiceID string        `json:"service_id,omitempty"`
	OrderID   string        `json:"order_id,omitempty"`
}

type abuseCaseAuditMetadata struct {
	DisplayID        int64             `json:"display_id"`
	CaseType         AbuseCaseType     `json:"case_type"`
	Severity         AbuseSeverity     `json:"severity"`
	ReportSource     AbuseReportSource `json:"report_source"`
	UserID           string            `json:"user_id,omitempty"`
	ServiceID        order.ServiceID   `json:"service_id,omitempty"`
	ProviderSourceID string            `json:"provider_source_id,omitempty"`
}

type abuseStatusAudit struct {
	Status      AbuseStatus `json:"status"`
	ActionTaken string      `json:"action_taken,omitempty"`
}

func (service *Service) appendTicketCreatedAudit(ctx context.Context, input CreateTicketInput, ticket Ticket) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   ticket.TenantID,
		ActorID:    audit.ActorID(input.Actor.ID),
		ActorType:  supportAuditActorType(input.Actor),
		Action:     auditActionSupportTicketCreated,
		TargetType: "support_ticket",
		TargetID:   audit.TargetID(ticket.ID),
		AfterSnapshotRedacted: supportAuditJSON(struct {
			Status TicketStatus `json:"status"`
		}{Status: ticket.Status}),
		MetadataRedacted: supportAuditJSON(ticketAuditMetadata{
			DisplayID:       ticket.DisplayID,
			RequesterUserID: string(ticket.RequesterUserID),
			Category:        ticket.Category,
			Priority:        ticket.Priority,
			ReferenceType:   ticket.ReferenceType,
			ReferenceID:     ticket.ReferenceID,
		}),
		CorrelationID: audit.CorrelationID(input.CorrelationID),
	})
	return err
}

func (service *Service) appendTicketNoteCreatedAudit(ctx context.Context, input AddTicketNoteInput, note TicketNote) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   note.TenantID,
		ActorID:    audit.ActorID(input.Actor.ID),
		ActorType:  supportAuditActorType(input.Actor),
		Action:     auditActionSupportNoteCreated,
		TargetType: "support_ticket_note",
		TargetID:   audit.TargetID(note.ID),
		MetadataRedacted: supportAuditJSON(ticketNoteAuditMetadata{
			DisplayID:  note.DisplayID,
			TicketID:   note.TicketID,
			Visibility: note.Visibility,
		}),
		CorrelationID: audit.CorrelationID(input.CorrelationID),
	})
	return err
}

func (service *Service) appendRiskFlagCreatedAudit(ctx context.Context, input CreateRiskFlagInput, flag RiskFlag) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   flag.TenantID,
		ActorID:    audit.ActorID(input.Actor.ID),
		ActorType:  supportAuditActorType(input.Actor),
		Action:     auditActionRiskFlagCreated,
		TargetType: "risk_flag",
		TargetID:   audit.TargetID(flag.ID),
		AfterSnapshotRedacted: supportAuditJSON(struct {
			Status RiskFlagStatus `json:"status"`
		}{Status: flag.Status}),
		MetadataRedacted: supportAuditJSON(riskFlagAuditMetadata{
			DisplayID: flag.DisplayID,
			FlagType:  flag.FlagType,
			Severity:  flag.Severity,
			UserID:    string(flag.UserID),
			ServiceID: string(flag.ServiceID),
			OrderID:   string(flag.OrderID),
		}),
		CorrelationID: audit.CorrelationID(input.CorrelationID),
	})
	return err
}

func (service *Service) appendAbuseCaseCreatedAudit(ctx context.Context, input CreateAbuseCaseInput, abuseCase AbuseCase) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:              abuseCase.TenantID,
		ActorID:               audit.ActorID(input.Actor.ID),
		ActorType:             supportAuditActorType(input.Actor),
		Action:                auditActionAbuseCaseCreated,
		TargetType:            "abuse_case",
		TargetID:              audit.TargetID(abuseCase.ID),
		AfterSnapshotRedacted: supportAuditJSON(abuseStatusAudit{Status: abuseCase.Status}),
		MetadataRedacted: supportAuditJSON(abuseCaseAuditMetadata{
			DisplayID:        abuseCase.DisplayID,
			CaseType:         abuseCase.CaseType,
			Severity:         abuseCase.Severity,
			ReportSource:     abuseCase.ReportSource,
			UserID:           string(abuseCase.UserID),
			ServiceID:        abuseCase.ServiceID,
			ProviderSourceID: string(abuseCase.ProviderSourceID),
		}),
		CorrelationID: audit.CorrelationID(input.CorrelationID),
	})
	return err
}

func (service *Service) appendAbuseServiceSuspendedAudit(
	ctx context.Context,
	input SuspendServiceForAbuseInput,
	before AbuseCase,
	after AbuseCase,
) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:               after.TenantID,
		ActorID:                audit.ActorID(input.Actor.ID),
		ActorType:              supportAuditActorType(input.Actor),
		Action:                 auditActionAbuseServiceSuspend,
		TargetType:             "abuse_case",
		TargetID:               audit.TargetID(after.ID),
		BeforeSnapshotRedacted: supportAuditJSON(abuseStatusAudit{Status: before.Status}),
		AfterSnapshotRedacted: supportAuditJSON(abuseStatusAudit{
			Status:      after.Status,
			ActionTaken: "service_suspended",
		}),
		MetadataRedacted: supportAuditJSON(abuseCaseAuditMetadata{
			DisplayID:        after.DisplayID,
			CaseType:         after.CaseType,
			Severity:         after.Severity,
			ReportSource:     after.ReportSource,
			UserID:           string(after.UserID),
			ServiceID:        input.ServiceID,
			ProviderSourceID: string(after.ProviderSourceID),
		}),
		CorrelationID: audit.CorrelationID(after.CorrelationID),
	})
	return err
}

func supportAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
