package support

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrStoreMissing                = errors.New("support store missing")
	ErrAuthorizerMissing           = errors.New("support authorizer missing")
	ErrSuspenderMissing            = errors.New("support service suspender missing")
	ErrTicketIDMissing             = errors.New("support ticket id missing")
	ErrTicketSubjectMissing        = errors.New("support ticket subject missing")
	ErrTicketCategoryInvalid       = errors.New("support ticket category invalid")
	ErrTicketPriorityInvalid       = errors.New("support ticket priority invalid")
	ErrTicketStatusInvalid         = errors.New("support ticket status invalid")
	ErrTicketNoteBodyMissing       = errors.New("support ticket note body missing")
	ErrTicketNoteVisibilityInvalid = errors.New("support ticket note visibility invalid")
	ErrAbuseCaseIDMissing          = errors.New("abuse case id missing")
	ErrAbuseCaseTypeInvalid        = errors.New("abuse case type invalid")
	ErrAbuseSeverityInvalid        = errors.New("abuse severity invalid")
	ErrAbuseStatusInvalid          = errors.New("abuse status invalid")
	ErrAbuseReportSourceInvalid    = errors.New("abuse report source invalid")
	ErrEvidenceSummaryMissing      = errors.New("abuse evidence summary missing")
	ErrAbuseReasonMissing          = errors.New("abuse action reason missing")
	ErrAbuseServiceMismatch        = errors.New("abuse service mismatch")
	ErrCorrelationIDMissing        = errors.New("support correlation id missing")
	ErrRequesterMismatch           = errors.New("support requester mismatch")
	ErrSupportAccessDenied         = errors.New("support access denied")
	ErrSupportTicketNotFound       = errors.New("support ticket not found")
	ErrAbuseCaseNotFound           = errors.New("abuse case not found")
)

type TicketID string
type TicketNoteID string
type AbuseCaseID string
type ReferenceType string
type ReferenceID string
type CorrelationID string

type TicketCategory string

const (
	TicketCategoryBilling               TicketCategory = "billing"
	TicketCategoryTopup                 TicketCategory = "topup"
	TicketCategoryOrder                 TicketCategory = "order"
	TicketCategoryProvisioning          TicketCategory = "provisioning"
	TicketCategoryServiceAccess         TicketCategory = "service_access"
	TicketCategoryCredential            TicketCategory = "credential"
	TicketCategoryRenewalExpiry         TicketCategory = "renewal_expiry"
	TicketCategorySuspensionTermination TicketCategory = "suspension_termination"
	TicketCategoryProviderIssue         TicketCategory = "provider_issue"
	TicketCategoryAbuseTakedown         TicketCategory = "abuse_takedown"
	TicketCategoryAccountLogin          TicketCategory = "account_login"
	TicketCategoryResellerSetup         TicketCategory = "reseller_setup"
	TicketCategoryFeatureRequest        TicketCategory = "feature_request"
	TicketCategoryOther                 TicketCategory = "other"
)

func (category TicketCategory) Valid() bool {
	switch category {
	case TicketCategoryBilling, TicketCategoryTopup, TicketCategoryOrder, TicketCategoryProvisioning,
		TicketCategoryServiceAccess, TicketCategoryCredential, TicketCategoryRenewalExpiry,
		TicketCategorySuspensionTermination, TicketCategoryProviderIssue, TicketCategoryAbuseTakedown,
		TicketCategoryAccountLogin, TicketCategoryResellerSetup, TicketCategoryFeatureRequest, TicketCategoryOther:
		return true
	default:
		return false
	}
}

type TicketPriority string

const (
	TicketPriorityP0 TicketPriority = "p0"
	TicketPriorityP1 TicketPriority = "p1"
	TicketPriorityP2 TicketPriority = "p2"
	TicketPriorityP3 TicketPriority = "p3"
	TicketPriorityP4 TicketPriority = "p4"
)

func (priority TicketPriority) Valid() bool {
	switch priority {
	case TicketPriorityP0, TicketPriorityP1, TicketPriorityP2, TicketPriorityP3, TicketPriorityP4:
		return true
	default:
		return false
	}
}

type TicketStatus string

const (
	TicketStatusOpen              TicketStatus = "open"
	TicketStatusWaitingOnCustomer TicketStatus = "waiting_on_customer"
	TicketStatusWaitingOnSupport  TicketStatus = "waiting_on_support"
	TicketStatusResolved          TicketStatus = "resolved"
	TicketStatusClosed            TicketStatus = "closed"
)

func (status TicketStatus) Valid() bool {
	switch status {
	case TicketStatusOpen, TicketStatusWaitingOnCustomer, TicketStatusWaitingOnSupport, TicketStatusResolved, TicketStatusClosed:
		return true
	default:
		return false
	}
}

type NoteVisibility string

const (
	NoteVisibilityPublic   NoteVisibility = "public"
	NoteVisibilityInternal NoteVisibility = "internal"
)

func (visibility NoteVisibility) Valid() bool {
	switch visibility {
	case NoteVisibilityPublic, NoteVisibilityInternal:
		return true
	default:
		return false
	}
}

type AbuseCaseType string

const (
	AbuseCaseTypeSpam                   AbuseCaseType = "spam"
	AbuseCaseTypePhishing               AbuseCaseType = "phishing"
	AbuseCaseTypeMalware                AbuseCaseType = "malware"
	AbuseCaseTypeBotnet                 AbuseCaseType = "botnet"
	AbuseCaseTypeBruteForce             AbuseCaseType = "brute_force"
	AbuseCaseTypePortScanning           AbuseCaseType = "port_scanning"
	AbuseCaseTypeDDoS                   AbuseCaseType = "ddos"
	AbuseCaseTypeCopyright              AbuseCaseType = "copyright"
	AbuseCaseTypeCredentialTheft        AbuseCaseType = "credential_theft"
	AbuseCaseTypeProxyScrapingViolation AbuseCaseType = "proxy_scraping_violation"
	AbuseCaseTypePaymentFraud           AbuseCaseType = "payment_fraud"
	AbuseCaseTypeChargeback             AbuseCaseType = "chargeback"
	AbuseCaseTypeIllegalContent         AbuseCaseType = "illegal_content"
	AbuseCaseTypeAUPViolation           AbuseCaseType = "aup_violation"
	AbuseCaseTypeProviderTakedown       AbuseCaseType = "provider_takedown"
	AbuseCaseTypeOther                  AbuseCaseType = "other"
)

func (caseType AbuseCaseType) Valid() bool {
	switch caseType {
	case AbuseCaseTypeSpam, AbuseCaseTypePhishing, AbuseCaseTypeMalware, AbuseCaseTypeBotnet,
		AbuseCaseTypeBruteForce, AbuseCaseTypePortScanning, AbuseCaseTypeDDoS, AbuseCaseTypeCopyright,
		AbuseCaseTypeCredentialTheft, AbuseCaseTypeProxyScrapingViolation, AbuseCaseTypePaymentFraud,
		AbuseCaseTypeChargeback, AbuseCaseTypeIllegalContent, AbuseCaseTypeAUPViolation,
		AbuseCaseTypeProviderTakedown, AbuseCaseTypeOther:
		return true
	default:
		return false
	}
}

type AbuseSeverity string

const (
	AbuseSeverityLow      AbuseSeverity = "low"
	AbuseSeverityMedium   AbuseSeverity = "medium"
	AbuseSeverityHigh     AbuseSeverity = "high"
	AbuseSeverityCritical AbuseSeverity = "critical"
)

func (severity AbuseSeverity) Valid() bool {
	switch severity {
	case AbuseSeverityLow, AbuseSeverityMedium, AbuseSeverityHigh, AbuseSeverityCritical:
		return true
	default:
		return false
	}
}

func (severity AbuseSeverity) RequiresImmediateSuspension() bool {
	return severity == AbuseSeverityHigh || severity == AbuseSeverityCritical
}

type AbuseStatus string

const (
	AbuseStatusNew                  AbuseStatus = "new"
	AbuseStatusTriaging             AbuseStatus = "triaging"
	AbuseStatusAwaitingClientAction AbuseStatus = "awaiting_client_action"
	AbuseStatusSuspended            AbuseStatus = "suspended"
	AbuseStatusResolved             AbuseStatus = "resolved"
	AbuseStatusFalsePositive        AbuseStatus = "rejected_false_positive"
	AbuseStatusTerminated           AbuseStatus = "terminated"
	AbuseStatusEscalated            AbuseStatus = "escalated"
	AbuseStatusClosed               AbuseStatus = "closed"
)

func (status AbuseStatus) Valid() bool {
	switch status {
	case AbuseStatusNew, AbuseStatusTriaging, AbuseStatusAwaitingClientAction, AbuseStatusSuspended,
		AbuseStatusResolved, AbuseStatusFalsePositive, AbuseStatusTerminated, AbuseStatusEscalated, AbuseStatusClosed:
		return true
	default:
		return false
	}
}

type AbuseReportSource string

const (
	AbuseReportSourceProvider           AbuseReportSource = "provider"
	AbuseReportSourceDatacenter         AbuseReportSource = "datacenter"
	AbuseReportSourceEmailDesk          AbuseReportSource = "email_abuse_desk"
	AbuseReportSourceLegal              AbuseReportSource = "legal"
	AbuseReportSourcePaymentProcessor   AbuseReportSource = "payment_processor"
	AbuseReportSourceInternalMonitoring AbuseReportSource = "internal_monitoring"
	AbuseReportSourceClient             AbuseReportSource = "client"
	AbuseReportSourceReseller           AbuseReportSource = "reseller"
	AbuseReportSourceThirdParty         AbuseReportSource = "third_party"
	AbuseReportSourceOther              AbuseReportSource = "other"
)

func (source AbuseReportSource) Valid() bool {
	switch source {
	case AbuseReportSourceProvider, AbuseReportSourceDatacenter, AbuseReportSourceEmailDesk,
		AbuseReportSourceLegal, AbuseReportSourcePaymentProcessor, AbuseReportSourceInternalMonitoring,
		AbuseReportSourceClient, AbuseReportSourceReseller, AbuseReportSourceThirdParty, AbuseReportSourceOther:
		return true
	default:
		return false
	}
}

type Ticket struct {
	ID              TicketID
	DisplayID       int64
	TenantID        tenant.ID
	RequesterUserID identity.UserID
	CreatedBy       identity.UserID
	AssignedUserID  identity.UserID
	Category        TicketCategory
	Priority        TicketPriority
	Status          TicketStatus
	Subject         string
	ReferenceType   ReferenceType
	ReferenceID     ReferenceID
	CorrelationID   CorrelationID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TicketNote struct {
	ID           TicketNoteID
	DisplayID    int64
	TicketID     TicketID
	TenantID     tenant.ID
	AuthorID     identity.UserID
	Visibility   NoteVisibility
	BodyRedacted string
	CreatedAt    time.Time
}

type AbuseCase struct {
	ID                      AbuseCaseID
	DisplayID               int64
	TenantID                tenant.ID
	UserID                  identity.UserID
	ServiceID               order.ServiceID
	ProviderSourceID        catalog.ProviderSourceID
	CaseType                AbuseCaseType
	Severity                AbuseSeverity
	ReportSource            AbuseReportSource
	Status                  AbuseStatus
	EvidenceSummaryRedacted string
	DeadlineAt              time.Time
	AssignedOwnerID         identity.UserID
	ActionTaken             string
	FinalResolution         string
	CreatedBy               identity.UserID
	CorrelationID           CorrelationID
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CreateTicketInput struct {
	Actor           identity.Actor
	TenantID        tenant.ID
	RequesterUserID identity.UserID
	Category        TicketCategory
	Priority        TicketPriority
	Subject         string
	ReferenceType   ReferenceType
	ReferenceID     ReferenceID
	InitialNote     string
	CorrelationID   CorrelationID
}

type CreateAbuseCaseInput struct {
	Actor                   identity.Actor
	TenantID                tenant.ID
	UserID                  identity.UserID
	ServiceID               order.ServiceID
	ProviderSourceID        catalog.ProviderSourceID
	CaseType                AbuseCaseType
	Severity                AbuseSeverity
	ReportSource            AbuseReportSource
	EvidenceSummaryRedacted string
	DeadlineAt              time.Time
	AssignedOwnerID         identity.UserID
	CorrelationID           CorrelationID
}

type SuspendServiceForAbuseInput struct {
	Actor      identity.Actor
	TenantID   tenant.ID
	CaseID     AbuseCaseID
	ServiceID  order.ServiceID
	FromStatus order.ServiceStatus
	Reason     string
}

type Store interface {
	CreateSupportTicket(ctx context.Context, input CreateTicketInput) (Ticket, error)
	CreateSupportTicketNote(ctx context.Context, input CreateTicketNoteInput) (TicketNote, error)
	CreateRiskFlag(ctx context.Context, input CreateRiskFlagInput) (RiskFlag, error)
	CreateAbuseCase(ctx context.Context, input CreateAbuseCaseInput) (AbuseCase, error)
	GetAbuseCase(ctx context.Context, tenantID tenant.ID, id AbuseCaseID) (AbuseCase, error)
	MarkAbuseCaseSuspended(ctx context.Context, input MarkAbuseCaseSuspendedInput) (AbuseCase, error)
}

type CreateTicketNoteInput struct {
	TicketID     TicketID
	TenantID     tenant.ID
	AuthorID     identity.UserID
	Visibility   NoteVisibility
	BodyRedacted string
}

type MarkAbuseCaseSuspendedInput struct {
	ID          AbuseCaseID
	TenantID    tenant.ID
	ActionTaken string
}

func trim(value string) string {
	return strings.TrimSpace(value)
}
