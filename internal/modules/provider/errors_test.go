package provider

import "testing"

func TestDefaultRetrySafetyMapsTaxonomy(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		want RetrySafety
	}{
		{name: "auth", code: ErrorAuthFailed, want: RetrySafetyDoNotRetry},
		{name: "rate limit", code: ErrorRateLimited, want: RetrySafetySafeRetry},
		{name: "network before send", code: ErrorNetworkBeforeSend, want: RetrySafetySafeRetry},
		{name: "timeout unknown", code: ErrorTimeoutUnknown, want: RetrySafetyUnsafeRetry},
		{name: "partial success", code: ErrorPartialSuccess, want: RetrySafetyManualReviewRequired},
		{name: "credential missing", code: ErrorCredentialMissing, want: RetrySafetyManualReviewRequired},
		{name: "unsupported", code: ErrorCapabilityNotSupported, want: RetrySafetyDoNotRetry},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultRetrySafety(tt.code); got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestAdapterErrorAllowsRetrySafetyOverride(t *testing.T) {
	err := AdapterError{
		Code:   ErrorMaintenance,
		Safety: RetrySafetyManualReviewRequired,
	}

	if got := err.RetrySafety(); got != RetrySafetyManualReviewRequired {
		t.Fatalf("expected override, got %s", got)
	}
}

func TestStatusForError(t *testing.T) {
	if got := StatusForError(NewError(ErrorCapabilityNotSupported, "unsupported")); got != OperationStatusCapabilityNotSupported {
		t.Fatalf("expected capability status, got %s", got)
	}
	if got := StatusForError(NewError(ErrorPartialSuccess, "missing credential")); got != OperationStatusPartialSuccess {
		t.Fatalf("expected partial status, got %s", got)
	}
	if got := StatusForError(NewError(ErrorTimeoutUnknown, "request may have reached provider")); got != OperationStatusUnknown {
		t.Fatalf("expected unknown status, got %s", got)
	}
}
