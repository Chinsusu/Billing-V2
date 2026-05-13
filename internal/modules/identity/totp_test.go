package identity

import (
	"testing"
	"time"
)

func TestTOTPCodeVerifiesCurrentWindow(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	code, err := TOTPCodeAt("JBSWY3DPEHPK3PXP", now)
	if err != nil {
		t.Fatalf("TOTPCodeAt returned error: %v", err)
	}
	ok, err := VerifyTOTPCode("JBSWY3DPEHPK3PXP", code, now)
	if err != nil {
		t.Fatalf("VerifyTOTPCode returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected code to verify")
	}
}

func TestTOTPCodeRejectsMalformedCode(t *testing.T) {
	if _, err := VerifyTOTPCode("JBSWY3DPEHPK3PXP", "12ab56", time.Now()); err != ErrTwoFactorCodeInvalid {
		t.Fatalf("expected invalid code error, got %v", err)
	}
}
