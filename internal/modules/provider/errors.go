package provider

import "fmt"

type RetrySafety string

const (
	RetrySafetySafeRetry            RetrySafety = "safe_retry"
	RetrySafetyUnsafeRetry          RetrySafety = "unsafe_retry"
	RetrySafetyDoNotRetry           RetrySafety = "do_not_retry"
	RetrySafetyManualReviewRequired RetrySafety = "manual_review_required"
)

type ErrorCode string

const (
	ErrorAuthFailed             ErrorCode = "PROVIDER_AUTH_FAILED"
	ErrorPermissionDenied       ErrorCode = "PROVIDER_PERMISSION_DENIED"
	ErrorAccountSuspended       ErrorCode = "PROVIDER_ACCOUNT_SUSPENDED"
	ErrorConfigInvalid          ErrorCode = "PROVIDER_CONFIG_INVALID"
	ErrorAdapterNotFound        ErrorCode = "PROVIDER_ADAPTER_NOT_FOUND"
	ErrorOutOfStock             ErrorCode = "PROVIDER_OUT_OF_STOCK"
	ErrorPlanNotFound           ErrorCode = "PROVIDER_PLAN_NOT_FOUND"
	ErrorRegionUnavailable      ErrorCode = "PROVIDER_REGION_UNAVAILABLE"
	ErrorCapabilityNotSupported ErrorCode = "PROVIDER_CAPABILITY_NOT_SUPPORTED"
	ErrorRateLimited            ErrorCode = "PROVIDER_RATE_LIMITED"
	ErrorTemporary              ErrorCode = "PROVIDER_TEMPORARY_ERROR"
	ErrorNetworkBeforeSend      ErrorCode = "PROVIDER_NETWORK_ERROR_BEFORE_SEND"
	ErrorMaintenance            ErrorCode = "PROVIDER_MAINTENANCE"
	ErrorTimeoutUnknown         ErrorCode = "PROVIDER_TIMEOUT_UNKNOWN"
	ErrorTimeoutRequestKnown    ErrorCode = "PROVIDER_TIMEOUT_REQUEST_KNOWN"
	ErrorPartialSuccess         ErrorCode = "PROVIDER_PARTIAL_SUCCESS"
	ErrorStateDrift             ErrorCode = "PROVIDER_STATE_DRIFT"
	ErrorResourceAlreadyExists  ErrorCode = "PROVIDER_RESOURCE_ALREADY_EXISTS"
	ErrorCredentialMissing      ErrorCode = "PROVIDER_CREDENTIAL_MISSING"
	ErrorCredentialInvalid      ErrorCode = "PROVIDER_CREDENTIAL_INVALID"
	ErrorResponseInvalid        ErrorCode = "PROVIDER_RESPONSE_INVALID"
)

type AdapterError struct {
	Code            ErrorCode
	MessageRedacted string
	Safety          RetrySafety
	Cause           error
}

func NewError(code ErrorCode, messageRedacted string) AdapterError {
	return AdapterError{Code: code, MessageRedacted: messageRedacted}
}

func (err AdapterError) Error() string {
	if err.MessageRedacted == "" {
		return string(err.Code)
	}
	return fmt.Sprintf("%s: %s", err.Code, err.MessageRedacted)
}

func (err AdapterError) Unwrap() error {
	return err.Cause
}

func (err AdapterError) RetrySafety() RetrySafety {
	if err.Safety != "" {
		return err.Safety
	}
	return DefaultRetrySafety(err.Code)
}

func DefaultRetrySafety(code ErrorCode) RetrySafety {
	switch code {
	case ErrorRateLimited,
		ErrorTemporary,
		ErrorNetworkBeforeSend,
		ErrorMaintenance:
		return RetrySafetySafeRetry
	case ErrorTimeoutUnknown:
		return RetrySafetyUnsafeRetry
	case ErrorTimeoutRequestKnown,
		ErrorPartialSuccess,
		ErrorStateDrift,
		ErrorCredentialMissing:
		return RetrySafetyManualReviewRequired
	case ErrorResourceAlreadyExists:
		return RetrySafetyManualReviewRequired
	case ErrorAuthFailed,
		ErrorPermissionDenied,
		ErrorAccountSuspended,
		ErrorConfigInvalid,
		ErrorAdapterNotFound,
		ErrorOutOfStock,
		ErrorPlanNotFound,
		ErrorRegionUnavailable,
		ErrorCapabilityNotSupported,
		ErrorCredentialInvalid,
		ErrorResponseInvalid:
		return RetrySafetyDoNotRetry
	default:
		return RetrySafetyManualReviewRequired
	}
}

func StatusForError(err AdapterError) OperationStatus {
	switch err.Code {
	case ErrorCapabilityNotSupported:
		return OperationStatusCapabilityNotSupported
	case ErrorPartialSuccess:
		return OperationStatusPartialSuccess
	case ErrorTimeoutUnknown, ErrorTimeoutRequestKnown, ErrorStateDrift:
		return OperationStatusUnknown
	}
	if err.RetrySafety() == RetrySafetyManualReviewRequired || err.RetrySafety() == RetrySafetyUnsafeRetry {
		return OperationStatusManualReviewRequired
	}
	return OperationStatusFailed
}
