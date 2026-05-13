package support

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateTicketAllowsClientOwnTicketWithoutManagePermission(t *testing.T) {
	store := &fakeSupportStore{}
	auditLog := &fakeSupportAudit{}
	service := NewServiceWithAudit(store, nil, auditLog)
	actor := identity.NewActor("user_1", "tenant_1", identity.ActorTypeClient)

	ticket, err := service.CreateTicket(testTenantContext("tenant_1"), CreateTicketInput{
		Actor:           actor,
		TenantID:        "tenant_1",
		RequesterUserID: "user_1",
		Category:        TicketCategoryAccountLogin,
		Priority:        TicketPriorityP2,
		Subject:         "Cannot log in",
		InitialNote:     "redacted browser details",
		CorrelationID:   "11111111-1111-1111-1111-111111111111",
	})

	if err != nil {
		t.Fatalf("expected own ticket to be allowed: %v", err)
	}
	if ticket.ID != "ticket_1" || store.createdTicket.RequesterUserID != "user_1" {
		t.Fatalf("unexpected ticket/store state: ticket=%#v input=%#v", ticket, store.createdTicket)
	}
	if len(auditLog.inputs) != 2 || auditLog.inputs[0].Action != auditActionSupportTicketCreated ||
		auditLog.inputs[1].Action != auditActionSupportNoteCreated {
		t.Fatalf("expected ticket audit, got %#v", auditLog.inputs)
	}
	if strings.Contains(string(auditLog.inputs[0].MetadataRedacted), "Cannot log in") {
		t.Fatalf("audit metadata must not contain support subject: %s", auditLog.inputs[0].MetadataRedacted)
	}
	if store.createdNote.BodyRedacted != "redacted browser details" {
		t.Fatalf("expected initial note to be stored redacted, got %#v", store.createdNote)
	}
}

func TestCreateTicketDeniesClientTicketForAnotherRequester(t *testing.T) {
	store := &fakeSupportStore{}
	service := NewService(store, rbac.StaticAuthorizer{Permissions: rbac.NewPermissionSet()})
	actor := identity.NewActor("user_1", "tenant_1", identity.ActorTypeClient)

	_, err := service.CreateTicket(testTenantContext("tenant_1"), CreateTicketInput{
		Actor:           actor,
		TenantID:        "tenant_1",
		RequesterUserID: "user_2",
		Category:        TicketCategoryBilling,
		Priority:        TicketPriorityP3,
		Subject:         "Billing question",
		CorrelationID:   "22222222-2222-2222-2222-222222222222",
	})

	if !errors.Is(err, rbac.ErrPermissionDenied) {
		t.Fatalf("expected permission denial, got %v", err)
	}
	if store.createdTicket.Subject != "" {
		t.Fatalf("ticket should not be stored on denied access: %#v", store.createdTicket)
	}
}

func TestCreateAbuseCaseRequiresPermissionAndAuditsWithoutEvidence(t *testing.T) {
	store := &fakeSupportStore{}
	auditLog := &fakeSupportAudit{}
	service := NewServiceWithAudit(store, rbac.StaticAuthorizer{
		Permissions: rbac.NewPermissionSet(rbac.PermissionAbuseCaseManage),
	}, auditLog)
	actor := identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)

	abuseCase, err := service.CreateAbuseCase(testTenantContext("tenant_1"), CreateAbuseCaseInput{
		Actor:                   actor,
		TenantID:                "tenant_1",
		UserID:                  "user_1",
		ServiceID:               "service_1",
		CaseType:                AbuseCaseTypePhishing,
		Severity:                AbuseSeverityHigh,
		ReportSource:            AbuseReportSourceProvider,
		EvidenceSummaryRedacted: "redacted provider evidence with token abc",
		CorrelationID:           "33333333-3333-3333-3333-333333333333",
	})

	if err != nil {
		t.Fatalf("expected abuse case creation: %v", err)
	}
	if abuseCase.ID != "abuse_1" {
		t.Fatalf("unexpected abuse case: %#v", abuseCase)
	}
	if len(auditLog.inputs) != 1 || auditLog.inputs[0].Action != auditActionAbuseCaseCreated {
		t.Fatalf("expected abuse audit, got %#v", auditLog.inputs)
	}
	if strings.Contains(string(auditLog.inputs[0].MetadataRedacted), "token abc") ||
		strings.Contains(string(auditLog.inputs[0].AfterSnapshotRedacted), "token abc") {
		t.Fatalf("audit output must not include evidence: %#v", auditLog.inputs[0])
	}
}

func TestCreateRiskFlagRequiresPermission(t *testing.T) {
	store := &fakeSupportStore{}
	service := NewService(store, rbac.StaticAuthorizer{Permissions: rbac.NewPermissionSet()})
	actor := identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)

	_, err := service.CreateRiskFlag(testTenantContext("tenant_1"), CreateRiskFlagInput{
		Actor:         actor,
		TenantID:      "tenant_1",
		UserID:        "user_1",
		FlagType:      RiskFlagTypeProviderRisk,
		Severity:      AbuseSeverityMedium,
		CorrelationID: "44444444-4444-4444-4444-444444444444",
	})

	if !errors.Is(err, rbac.ErrPermissionDenied) {
		t.Fatalf("expected risk flag permission denial, got %v", err)
	}
	if store.createdRiskFlag.ID != "" {
		t.Fatalf("risk flag should not be stored on denied access: %#v", store.createdRiskFlag)
	}
}

func TestSuspendServiceForAbuseRequiresServiceSuspendPermission(t *testing.T) {
	store := &fakeSupportStore{abuseCase: testAbuseCase()}
	suspender := &fakeServiceSuspender{}
	service := NewServiceWithAuditAndSuspender(store, rbac.StaticAuthorizer{
		Permissions: rbac.NewPermissionSet(rbac.PermissionAbuseCaseManage),
	}, nil, suspender)
	actor := identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)

	_, err := service.SuspendServiceForAbuse(testTenantContext("tenant_1"), SuspendServiceForAbuseInput{
		Actor:      actor,
		TenantID:   "tenant_1",
		CaseID:     "abuse_1",
		ServiceID:  "service_1",
		FromStatus: order.ServiceStatusActive,
		Reason:     "confirmed provider takedown",
	})

	if !errors.Is(err, rbac.ErrPermissionDenied) {
		t.Fatalf("expected service suspend permission denial, got %v", err)
	}
	if suspender.input.ID != "" {
		t.Fatalf("suspender should not be called: %#v", suspender.input)
	}
}

func TestSuspendServiceForAbuseSuspendsServiceAndMarksCase(t *testing.T) {
	store := &fakeSupportStore{abuseCase: testAbuseCase()}
	auditLog := &fakeSupportAudit{}
	suspender := &fakeServiceSuspender{}
	service := NewServiceWithAuditAndSuspender(store, rbac.StaticAuthorizer{
		Permissions: rbac.NewPermissionSet(rbac.PermissionAbuseCaseManage, rbac.PermissionServiceSuspend),
	}, auditLog, suspender)
	actor := identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)

	abuseCase, err := service.SuspendServiceForAbuse(testTenantContext("tenant_1"), SuspendServiceForAbuseInput{
		Actor:      actor,
		TenantID:   "tenant_1",
		CaseID:     "abuse_1",
		ServiceID:  "service_1",
		FromStatus: order.ServiceStatusActive,
		Reason:     "confirmed provider takedown",
	})

	if err != nil {
		t.Fatalf("expected abuse suspension: %v", err)
	}
	if suspender.input.Action != order.ServiceLifecycleActionSuspend ||
		suspender.input.SuspensionReason != order.SuspensionReasonAbuse ||
		suspender.input.Reason != "confirmed provider takedown" {
		t.Fatalf("unexpected lifecycle input: %#v", suspender.input)
	}
	if abuseCase.Status != AbuseStatusSuspended ||
		store.markSuspended.ActionTaken != "service_suspended: confirmed provider takedown" {
		t.Fatalf("unexpected suspended case: case=%#v mark=%#v", abuseCase, store.markSuspended)
	}
	if len(auditLog.inputs) != 1 || auditLog.inputs[0].Action != auditActionAbuseServiceSuspend {
		t.Fatalf("expected abuse suspension audit, got %#v", auditLog.inputs)
	}
}

func testTenantContext(id tenant.ID) context.Context {
	return tenant.WithContext(context.Background(), tenant.NewContext(id))
}

func testAbuseCase() AbuseCase {
	return AbuseCase{
		ID:                      "abuse_1",
		DisplayID:               10001,
		TenantID:                "tenant_1",
		UserID:                  "user_1",
		ServiceID:               "service_1",
		CaseType:                AbuseCaseTypePhishing,
		Severity:                AbuseSeverityHigh,
		ReportSource:            AbuseReportSourceProvider,
		Status:                  AbuseStatusNew,
		EvidenceSummaryRedacted: "redacted provider evidence",
		CreatedBy:               "admin_1",
		CorrelationID:           "55555555-5555-5555-5555-555555555555",
	}
}

type fakeSupportStore struct {
	createdTicket   CreateTicketInput
	createdNote     CreateTicketNoteInput
	createdRiskFlag RiskFlag
	createdAbuse    CreateAbuseCaseInput
	abuseCase       AbuseCase
	markSuspended   MarkAbuseCaseSuspendedInput
}

func (store *fakeSupportStore) CreateSupportTicket(ctx context.Context, input CreateTicketInput) (Ticket, error) {
	store.createdTicket = input
	return Ticket{
		ID:              "ticket_1",
		DisplayID:       10001,
		TenantID:        input.TenantID,
		RequesterUserID: input.RequesterUserID,
		CreatedBy:       input.Actor.ID,
		Category:        input.Category,
		Priority:        input.Priority,
		Status:          TicketStatusOpen,
		Subject:         input.Subject,
		CorrelationID:   input.CorrelationID,
	}, nil
}

func (store *fakeSupportStore) CreateSupportTicketNote(ctx context.Context, input CreateTicketNoteInput) (TicketNote, error) {
	store.createdNote = input
	return TicketNote{
		ID:           "note_1",
		DisplayID:    10002,
		TicketID:     input.TicketID,
		TenantID:     input.TenantID,
		AuthorID:     input.AuthorID,
		Visibility:   input.Visibility,
		BodyRedacted: input.BodyRedacted,
	}, nil
}

func (store *fakeSupportStore) CreateRiskFlag(ctx context.Context, input CreateRiskFlagInput) (RiskFlag, error) {
	store.createdRiskFlag = RiskFlag{
		ID:            "risk_1",
		DisplayID:     10003,
		TenantID:      input.TenantID,
		UserID:        input.UserID,
		ServiceID:     input.ServiceID,
		OrderID:       input.OrderID,
		FlagType:      input.FlagType,
		Severity:      input.Severity,
		Status:        RiskFlagStatusOpen,
		NoteRedacted:  input.NoteRedacted,
		CreatedBy:     input.Actor.ID,
		CorrelationID: input.CorrelationID,
	}
	return store.createdRiskFlag, nil
}

func (store *fakeSupportStore) CreateAbuseCase(ctx context.Context, input CreateAbuseCaseInput) (AbuseCase, error) {
	store.createdAbuse = input
	return AbuseCase{
		ID:                      "abuse_1",
		DisplayID:               10001,
		TenantID:                input.TenantID,
		UserID:                  input.UserID,
		ServiceID:               input.ServiceID,
		ProviderSourceID:        input.ProviderSourceID,
		CaseType:                input.CaseType,
		Severity:                input.Severity,
		ReportSource:            input.ReportSource,
		Status:                  AbuseStatusNew,
		EvidenceSummaryRedacted: input.EvidenceSummaryRedacted,
		CreatedBy:               input.Actor.ID,
		CorrelationID:           input.CorrelationID,
	}, nil
}

func (store *fakeSupportStore) GetAbuseCase(ctx context.Context, tenantID tenant.ID, id AbuseCaseID) (AbuseCase, error) {
	if store.abuseCase.ID == "" {
		return AbuseCase{}, ErrAbuseCaseNotFound
	}
	return store.abuseCase, nil
}

func (store *fakeSupportStore) MarkAbuseCaseSuspended(ctx context.Context, input MarkAbuseCaseSuspendedInput) (AbuseCase, error) {
	store.markSuspended = input
	abuseCase := store.abuseCase
	abuseCase.Status = AbuseStatusSuspended
	abuseCase.ActionTaken = input.ActionTaken
	return abuseCase, nil
}

type fakeSupportAudit struct {
	inputs []audit.AppendInput
}

func (auditLog *fakeSupportAudit) Append(ctx context.Context, input audit.AppendInput) (audit.Log, error) {
	auditLog.inputs = append(auditLog.inputs, input)
	return audit.Log{}, nil
}

type fakeServiceSuspender struct {
	input order.TransitionServiceLifecycleInput
}

func (suspender *fakeServiceSuspender) TransitionServiceLifecycle(ctx context.Context, input order.TransitionServiceLifecycleInput) (order.ServiceInstance, error) {
	suspender.input = input
	return order.ServiceInstance{
		ID:               input.ID,
		TenantID:         input.TenantID,
		Status:           input.ToStatus,
		SuspensionReason: input.SuspensionReason,
	}, nil
}
