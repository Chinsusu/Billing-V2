package provider

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrOperationIDMissing    = errors.New("operation id missing")
	ErrTenantMissing         = errors.New("tenant id missing")
	ErrSourceMissing         = errors.New("source id missing")
	ErrIdempotencyKeyMissing = errors.New("idempotency key missing")
	ErrCorrelationIDMissing  = errors.New("correlation id missing")
	ErrAttemptInvalid        = errors.New("attempt number invalid")
)

type OperationID string
type SourceID string
type AccountID string
type ActorID string
type IdempotencyKey string
type CorrelationID string
type ExternalRequestID string
type ExternalResourceID string
type ServiceIdentifier string
type RawResponseReference string

type SourceSnapshot struct {
	ProviderType Type
	Name         string
	Region       string
	PlanKey      string
}

type OperationContext struct {
	OperationID            OperationID
	TenantID               tenant.ID
	SourceID               SourceID
	ProviderAccountID      AccountID
	ActorOrSystemID        ActorID
	IdempotencyKey         IdempotencyKey
	CorrelationID          CorrelationID
	RequestTimeout         time.Duration
	Deadline               time.Time
	AttemptNumber          int
	CapabilitySnapshot     CapabilityProfile
	ProviderSourceSnapshot SourceSnapshot
}

func (operation OperationContext) Validate() error {
	if operation.OperationID == "" {
		return ErrOperationIDMissing
	}
	if operation.TenantID.Empty() {
		return ErrTenantMissing
	}
	if operation.SourceID == "" {
		return ErrSourceMissing
	}
	if operation.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if operation.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	if operation.AttemptNumber <= 0 {
		return ErrAttemptInvalid
	}
	return nil
}

type OperationStatus string

const (
	OperationStatusSuccess                OperationStatus = "success"
	OperationStatusFailed                 OperationStatus = "failed"
	OperationStatusPartialSuccess         OperationStatus = "partial_success"
	OperationStatusUnknown                OperationStatus = "unknown"
	OperationStatusManualReviewRequired   OperationStatus = "manual_review_required"
	OperationStatusCapabilityNotSupported OperationStatus = "capability_not_supported"
)

type CredentialEnvelope struct {
	Type                 CredentialType
	EncryptedPayload     string
	EncryptedPayloadRef  string
	EncryptionKeyVersion string
	SecretVersion        string
	MaskedHint           string
}

type OperationResult struct {
	Status               OperationStatus
	ExternalRequestID    ExternalRequestID
	ExternalResourceID   ExternalResourceID
	ServiceIdentifier    ServiceIdentifier
	Credential           CredentialEnvelope
	ProviderStatus       string
	RetrySafety          RetrySafety
	ErrorCode            ErrorCode
	ErrorMessageRedacted string
	RawResponseReference RawResponseReference
	ObservedAt           time.Time
}

func SuccessResult(observedAt time.Time) OperationResult {
	return OperationResult{
		Status:      OperationStatusSuccess,
		RetrySafety: RetrySafetyDoNotRetry,
		ObservedAt:  observedAt,
	}
}

func ResultFromError(err AdapterError, observedAt time.Time) OperationResult {
	return OperationResult{
		Status:               StatusForError(err),
		RetrySafety:          err.RetrySafety(),
		ErrorCode:            err.Code,
		ErrorMessageRedacted: err.MessageRedacted,
		ObservedAt:           observedAt,
	}
}

type HealthStatus string

const (
	HealthStatusHealthy  HealthStatus = "healthy"
	HealthStatusDegraded HealthStatus = "degraded"
	HealthStatusDown     HealthStatus = "down"
	HealthStatusUnknown  HealthStatus = "unknown"
)

type StockStatus string

const (
	StockStatusAvailable  StockStatus = "available"
	StockStatusOutOfStock StockStatus = "out_of_stock"
	StockStatusUnknown    StockStatus = "unknown"
)

type HealthResult struct {
	HealthStatus HealthStatus
	Result       OperationResult
}

type StockResult struct {
	StockStatus   StockStatus
	CapacityCount int
	Result        OperationResult
}

type HealthRequest struct{}

type StockRequest struct {
	PlanKey string
	Region  string
}

type ProvisionRequest struct {
	PlanKey       string
	Hostname      string
	InputSnapshot map[string]string
}

type ResourceRequest struct {
	ExternalResourceID ExternalResourceID
	Reason             string
}

type Adapter interface {
	ProviderType() Type
	CapabilityProfile() CapabilityProfile
	CheckHealth(ctx context.Context, operation OperationContext, request HealthRequest) (HealthResult, error)
	CheckStock(ctx context.Context, operation OperationContext, request StockRequest) (StockResult, error)
	Provision(ctx context.Context, operation OperationContext, request ProvisionRequest) (OperationResult, error)
	GetStatus(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	Suspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	Unsuspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	Terminate(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	Renew(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	ResetPassword(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
	ChangeIP(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error)
}
