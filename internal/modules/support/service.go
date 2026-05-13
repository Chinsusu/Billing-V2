package support

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ServiceSuspender interface {
	TransitionServiceLifecycle(ctx context.Context, input order.TransitionServiceLifecycleInput) (order.ServiceInstance, error)
}

type Service struct {
	store      Store
	authorizer rbac.Authorizer
	audit      AuditAppender
	suspender  ServiceSuspender
}

func NewService(store Store, authorizer rbac.Authorizer) *Service {
	return &Service{store: store, authorizer: authorizer}
}

func NewServiceWithAudit(store Store, authorizer rbac.Authorizer, audit AuditAppender) *Service {
	return &Service{store: store, authorizer: authorizer, audit: audit}
}

func NewServiceWithAuditAndSuspender(
	store Store,
	authorizer rbac.Authorizer,
	audit AuditAppender,
	suspender ServiceSuspender,
) *Service {
	return &Service{store: store, authorizer: authorizer, audit: audit, suspender: suspender}
}

func (service *Service) CreateTicket(ctx context.Context, input CreateTicketInput) (Ticket, error) {
	if err := service.ready(); err != nil {
		return Ticket{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Ticket{}, err
	}
	if err := ensureActorTenantAccess(ctx, input.Actor, input.TenantID); err != nil {
		return Ticket{}, err
	}
	if !isOwnClientTicket(input.Actor, input.TenantID, input.RequesterUserID) {
		if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionTicketManage, rbac.RiskMedium, "create support ticket"); err != nil {
			return Ticket{}, err
		}
	}
	ticket, err := service.store.CreateSupportTicket(ctx, input)
	if err != nil {
		return Ticket{}, err
	}
	var initialNote TicketNote
	if input.InitialNote != "" {
		initialNote, err = service.store.CreateSupportTicketNote(ctx, CreateTicketNoteInput{
			TicketID:     ticket.ID,
			TenantID:     ticket.TenantID,
			AuthorID:     input.Actor.ID,
			Visibility:   NoteVisibilityPublic,
			BodyRedacted: input.InitialNote,
		})
		if err != nil {
			return Ticket{}, err
		}
	}
	if err := service.appendTicketCreatedAudit(ctx, input, ticket); err != nil {
		return Ticket{}, err
	}
	if initialNote.ID != "" {
		if err := service.appendTicketNoteCreatedAudit(ctx, AddTicketNoteInput{
			Actor:         input.Actor,
			TicketID:      ticket.ID,
			TenantID:      ticket.TenantID,
			Visibility:    NoteVisibilityPublic,
			BodyRedacted:  input.InitialNote,
			CorrelationID: input.CorrelationID,
		}, initialNote); err != nil {
			return Ticket{}, err
		}
	}
	return ticket, nil
}

func (service *Service) AddTicketNote(ctx context.Context, input AddTicketNoteInput) (TicketNote, error) {
	if err := service.ready(); err != nil {
		return TicketNote{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TicketNote{}, err
	}
	if err := ensureActorTenantAccess(ctx, input.Actor, input.TenantID); err != nil {
		return TicketNote{}, err
	}
	if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionTicketManage, rbac.RiskMedium, "add support ticket note"); err != nil {
		return TicketNote{}, err
	}
	note, err := service.store.CreateSupportTicketNote(ctx, CreateTicketNoteInput{
		TicketID:     input.TicketID,
		TenantID:     input.TenantID,
		AuthorID:     input.Actor.ID,
		Visibility:   input.Visibility,
		BodyRedacted: input.BodyRedacted,
	})
	if err != nil {
		return TicketNote{}, err
	}
	if err := service.appendTicketNoteCreatedAudit(ctx, input, note); err != nil {
		return TicketNote{}, err
	}
	return note, nil
}

func (service *Service) CreateRiskFlag(ctx context.Context, input CreateRiskFlagInput) (RiskFlag, error) {
	if err := service.ready(); err != nil {
		return RiskFlag{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return RiskFlag{}, err
	}
	if err := ensureActorTenantAccess(ctx, input.Actor, input.TenantID); err != nil {
		return RiskFlag{}, err
	}
	if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionRiskFlagCreate, rbac.RiskHigh, "create risk flag"); err != nil {
		return RiskFlag{}, err
	}
	flag, err := service.store.CreateRiskFlag(ctx, input)
	if err != nil {
		return RiskFlag{}, err
	}
	if err := service.appendRiskFlagCreatedAudit(ctx, input, flag); err != nil {
		return RiskFlag{}, err
	}
	return flag, nil
}

func (service *Service) CreateAbuseCase(ctx context.Context, input CreateAbuseCaseInput) (AbuseCase, error) {
	if err := service.ready(); err != nil {
		return AbuseCase{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return AbuseCase{}, err
	}
	if err := ensureActorTenantAccess(ctx, input.Actor, input.TenantID); err != nil {
		return AbuseCase{}, err
	}
	if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionAbuseCaseManage, rbac.RiskHigh, "create abuse case"); err != nil {
		return AbuseCase{}, err
	}
	abuseCase, err := service.store.CreateAbuseCase(ctx, input)
	if err != nil {
		return AbuseCase{}, err
	}
	if err := service.appendAbuseCaseCreatedAudit(ctx, input, abuseCase); err != nil {
		return AbuseCase{}, err
	}
	return abuseCase, nil
}

func (service *Service) SuspendServiceForAbuse(ctx context.Context, input SuspendServiceForAbuseInput) (AbuseCase, error) {
	if err := service.ready(); err != nil {
		return AbuseCase{}, err
	}
	if service.suspender == nil {
		return AbuseCase{}, ErrSuspenderMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return AbuseCase{}, err
	}
	if err := ensureActorTenantAccess(ctx, input.Actor, input.TenantID); err != nil {
		return AbuseCase{}, err
	}
	if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionAbuseCaseManage, rbac.RiskHigh, input.Reason); err != nil {
		return AbuseCase{}, err
	}
	if err := service.authorize(ctx, input.Actor, input.TenantID, rbac.PermissionServiceSuspend, rbac.RiskHigh, input.Reason); err != nil {
		return AbuseCase{}, err
	}
	before, err := service.store.GetAbuseCase(ctx, input.TenantID, input.CaseID)
	if err != nil {
		return AbuseCase{}, err
	}
	if before.ServiceID != input.ServiceID {
		return AbuseCase{}, ErrAbuseServiceMismatch
	}
	if before.Status == AbuseStatusSuspended {
		return before, nil
	}
	if _, err := service.suspender.TransitionServiceLifecycle(ctx, order.TransitionServiceLifecycleInput{
		ID:               input.ServiceID,
		TenantID:         input.TenantID,
		ActorID:          audit.ActorID(input.Actor.ID),
		ActorType:        supportAuditActorType(input.Actor),
		Action:           order.ServiceLifecycleActionSuspend,
		FromStatus:       input.FromStatus,
		ToStatus:         order.ServiceStatusSuspended,
		SuspensionReason: order.SuspensionReasonAbuse,
		Reason:           input.Reason,
	}); err != nil {
		return AbuseCase{}, err
	}
	after, err := service.store.MarkAbuseCaseSuspended(ctx, MarkAbuseCaseSuspendedInput{
		ID:          input.CaseID,
		TenantID:    input.TenantID,
		ActionTaken: "service_suspended: " + input.Reason,
	})
	if err != nil {
		return AbuseCase{}, err
	}
	if err := service.appendAbuseServiceSuspendedAudit(ctx, input, before, after); err != nil {
		return AbuseCase{}, err
	}
	return after, nil
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrStoreMissing
	}
	return nil
}

func (service *Service) authorize(
	ctx context.Context,
	actor identity.Actor,
	tenantID tenant.ID,
	permission rbac.Permission,
	risk rbac.RiskLevel,
	reason string,
) error {
	if service.authorizer == nil {
		return ErrAuthorizerMissing
	}
	return service.authorizer.Check(ctx, rbac.CheckRequest{
		Actor:            actor,
		Permission:       permission,
		ResourceTenantID: tenantID,
		Risk:             risk,
		Reason:           reason,
	})
}

func ensureActorTenantAccess(ctx context.Context, actor identity.Actor, tenantID tenant.ID) error {
	tenantContext, err := tenant.RequireAccess(ctx, tenantID)
	if err != nil {
		return err
	}
	if !actor.IsSystem && actor.TenantID != tenantContext.ActorTenantID {
		return tenant.ErrTenantMismatch
	}
	return nil
}

func isOwnClientTicket(actor identity.Actor, tenantID tenant.ID, requesterUserID identity.UserID) bool {
	return actor.Type == identity.ActorTypeClient && actor.TenantID == tenantID && actor.ID == requesterUserID
}
