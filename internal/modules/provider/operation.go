package provider

import (
	"context"
	"time"
)

// OperationStatus is the outcome of an adapter operation.
type OperationStatus string

const (
	StatusSuccess              OperationStatus = "success"
	StatusFailed               OperationStatus = "failed"
	StatusPartialSuccess       OperationStatus = "partial_success"
	StatusUnknown              OperationStatus = "unknown"
	StatusManualReviewRequired OperationStatus = "manual_review_required"
	StatusCapabilityNotSupported OperationStatus = "capability_not_supported"
)

// HealthStatus is the health of a provider/source.
type HealthStatus string

const (
	HealthHealthy   HealthStatus = "healthy"
	HealthDegraded  HealthStatus = "degraded"
	HealthDown      HealthStatus = "down"
	HealthUnknown   HealthStatus = "unknown"
)

// StockStatus is the stock availability for a plan/source.
type StockStatus string

const (
	StockAvailable  StockStatus = "available"
	StockOutOfStock StockStatus = "out_of_stock"
	StockUnknown    StockStatus = "unknown"
)

// OperationContext carries the required metadata for any adapter call.
// All fields come from the provisioning job; adapters must not read config or
// secrets independently.
type OperationContext struct {
	OperationID        string
	TenantID           string
	SourceID           string
	ProviderAccountID  string
	ActorID            string
	IdempotencyKey     string
	CorrelationID      string
	AttemptNumber      int
	RequestTimeout     time.Duration
	CapabilitySnapshot CapabilityProfile
}

// OperationResult is the normalized output of every adapter operation.
// Raw provider payloads must be stored out-of-band; this struct must never
// contain plaintext credentials or secrets.
type OperationResult struct {
	Status              OperationStatus
	ExternalRequestID   string
	ExternalResourceID  string
	ServiceIdentifier   string
	// CredentialPayload holds encrypted credential bytes ready for DB storage.
	// The adapter encrypts before returning; plaintext must not leave the adapter.
	CredentialPayload   []byte
	ProviderStatus      string
	RetrySafety         RetrySafety
	ErrorCode           ErrorCode
	ErrorMessage        string // redacted — no secrets
	RawResponseRef      string // opaque reference to object storage, not inline
	ObservedAt          time.Time
}

// HealthResult is returned by CheckHealth.
type HealthResult struct {
	Status  HealthStatus
	Message string
}

// StockResult is returned by CheckStock.
type StockResult struct {
	Status        StockStatus
	CapacityCount *int // nil when unknown
}

// ProvisionInput bundles plan and source data needed to provision a resource.
// Adapters must not perform wallet debits or modify orders.
type ProvisionInput struct {
	OpCtx              OperationContext
	PlanSpecSnapshot   map[string]string // key-value plan spec, no secrets
	SourceConfigRef    string            // opaque reference; actual secret loaded by secret loader
}

// StatusInput is used for getStatus calls.
type StatusInput struct {
	OpCtx             OperationContext
	ExternalResourceID string
}

// SuspendInput is used for suspend/unsuspend calls.
type SuspendInput struct {
	OpCtx              OperationContext
	ExternalResourceID string
	Reason             string
}

// TerminateInput is used for terminate calls.
type TerminateInput struct {
	OpCtx              OperationContext
	ExternalResourceID string
}

// RenewInput is used for renew calls.
type RenewInput struct {
	OpCtx              OperationContext
	ExternalResourceID string
}

// ResetPasswordInput is used for password reset calls.
type ResetPasswordInput struct {
	OpCtx              OperationContext
	ExternalResourceID string
}

// Adapter is the interface every provider adapter must implement.
// Adapters normalize provider-specific behavior; they do not debit wallets,
// decide permissions, or self-retry outside the job policy.
type Adapter interface {
	// ProviderType returns the stable string identifier for this adapter.
	ProviderType() string

	// CapabilityProfile returns the declared capabilities for a source.
	CapabilityProfile(ctx context.Context, sourceID string) (CapabilityProfile, error)

	// CheckHealth tests provider reachability and credential validity.
	CheckHealth(ctx context.Context, opCtx OperationContext) (HealthResult, error)

	// CheckStock checks available capacity for the plan/source combination.
	CheckStock(ctx context.Context, opCtx OperationContext) (StockResult, error)

	// Provision creates a resource on the provider. Result status may be
	// unknown or partial_success; callers must handle manual_review_required.
	Provision(ctx context.Context, input ProvisionInput) (OperationResult, error)

	// GetStatus fetches the current state of a provisioned resource.
	GetStatus(ctx context.Context, input StatusInput) (OperationResult, error)

	// Suspend suspends an active resource. Returns CAPABILITY_NOT_SUPPORTED if
	// the provider does not support the operation.
	Suspend(ctx context.Context, input SuspendInput) (OperationResult, error)

	// Unsuspend restores a suspended resource.
	Unsuspend(ctx context.Context, input SuspendInput) (OperationResult, error)

	// Terminate permanently destroys a resource. Callers must not retry blindly
	// if the result is unknown.
	Terminate(ctx context.Context, input TerminateInput) (OperationResult, error)

	// Renew extends the resource lifecycle when the provider requires an API call.
	Renew(ctx context.Context, input RenewInput) (OperationResult, error)

	// ResetPassword rotates the credential for a resource.
	ResetPassword(ctx context.Context, input ResetPasswordInput) (OperationResult, error)
}
