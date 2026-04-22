package provider

// ErrorCode is a normalized provider error code.
type ErrorCode string

const (
	// Auth and configuration errors — do not retry.
	ErrAuthFailed          ErrorCode = "PROVIDER_AUTH_FAILED"
	ErrPermissionDenied    ErrorCode = "PROVIDER_PERMISSION_DENIED"
	ErrAccountSuspended    ErrorCode = "PROVIDER_ACCOUNT_SUSPENDED"
	ErrConfigInvalid       ErrorCode = "PROVIDER_CONFIG_INVALID"
	ErrAdapterNotFound     ErrorCode = "PROVIDER_ADAPTER_NOT_FOUND"

	// Capacity and product mapping — do not retry.
	ErrOutOfStock              ErrorCode = "PROVIDER_OUT_OF_STOCK"
	ErrPlanNotFound            ErrorCode = "PROVIDER_PLAN_NOT_FOUND"
	ErrRegionUnavailable       ErrorCode = "PROVIDER_REGION_UNAVAILABLE"
	ErrCapabilityNotSupported  ErrorCode = "PROVIDER_CAPABILITY_NOT_SUPPORTED"

	// Transient runtime — safe to retry.
	ErrRateLimited             ErrorCode = "PROVIDER_RATE_LIMITED"
	ErrTemporaryError          ErrorCode = "PROVIDER_TEMPORARY_ERROR"
	ErrNetworkBeforeSend       ErrorCode = "PROVIDER_NETWORK_ERROR_BEFORE_SEND"
	ErrMaintenance             ErrorCode = "PROVIDER_MAINTENANCE"

	// Uncertain and dangerous — manual review or status lookup required.
	ErrTimeoutUnknown          ErrorCode = "PROVIDER_TIMEOUT_UNKNOWN"
	ErrTimeoutRequestKnown     ErrorCode = "PROVIDER_TIMEOUT_REQUEST_KNOWN"
	ErrPartialSuccess          ErrorCode = "PROVIDER_PARTIAL_SUCCESS"
	ErrStateDrift              ErrorCode = "PROVIDER_STATE_DRIFT"
	ErrResourceAlreadyExists   ErrorCode = "PROVIDER_RESOURCE_ALREADY_EXISTS"

	// Credential-specific.
	ErrCredentialMissing        ErrorCode = "PROVIDER_CREDENTIAL_MISSING"
	ErrCredentialInvalid        ErrorCode = "PROVIDER_CREDENTIAL_INVALID"
	ErrCredentialRotationFailed ErrorCode = "PROVIDER_CREDENTIAL_ROTATION_FAILED"
)

// RetrySafety describes whether a failed operation may be safely retried.
type RetrySafety string

const (
	RetrySafe           RetrySafety = "safe_retry"
	RetryUnsafe         RetrySafety = "unsafe_retry"
	RetryDoNot          RetrySafety = "do_not_retry"
	RetryManualReview   RetrySafety = "manual_review_required"
)

// retrySafetyFor returns the canonical retry safety for a given error code.
func retrySafetyFor(code ErrorCode) RetrySafety {
	switch code {
	case ErrAuthFailed, ErrPermissionDenied, ErrAccountSuspended,
		ErrConfigInvalid, ErrAdapterNotFound, ErrOutOfStock,
		ErrPlanNotFound, ErrCapabilityNotSupported:
		return RetryDoNot

	case ErrRateLimited, ErrTemporaryError, ErrNetworkBeforeSend, ErrMaintenance:
		return RetrySafe

	case ErrTimeoutUnknown:
		return RetryUnsafe

	case ErrTimeoutRequestKnown, ErrPartialSuccess, ErrStateDrift,
		ErrCredentialMissing, ErrCredentialInvalid, ErrCredentialRotationFailed:
		return RetryManualReview

	case ErrRegionUnavailable:
		// Policy-dependent; default safe (can choose another region).
		return RetrySafe

	case ErrResourceAlreadyExists:
		// Safe if the caller maps the existing resource.
		return RetrySafe

	default:
		return RetryManualReview
	}
}

// AdapterError wraps a normalized error code with a redacted human-readable message.
// Raw provider errors must never be surfaced here.
type AdapterError struct {
	Code    ErrorCode
	Message string // redacted — must not contain secrets or raw provider payloads
}

func (e *AdapterError) Error() string {
	return string(e.Code) + ": " + e.Message
}

// RetrySafety returns the retry classification for this error.
func (e *AdapterError) RetrySafety() RetrySafety {
	return retrySafetyFor(e.Code)
}
