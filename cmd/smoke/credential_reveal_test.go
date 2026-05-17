package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRevealTargetCredentialUsesCookieOnlySession(t *testing.T) {
	fixture := targetCredentialRevealFixture{
		CredentialID:     "credential_1",
		EncryptedPayload: "encrypted-fixture",
		ServiceDisplayID: 43001,
		ExpectedPayload:  targetCredentialRevealPayload,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Actor-Id") != "" || r.Header.Get("X-Tenant-Id") != "" {
			t.Fatalf("credential reveal smoke should not send dev actor headers")
		}
		cookie, err := r.Cookie("billing_session")
		if err != nil || cookie.Value != "session-token" {
			t.Fatalf("expected session cookie")
		}
		if !strings.HasSuffix(r.URL.Path, "/credentials/credential_1/reveal") {
			t.Fatalf("unexpected reveal path")
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"credential_1","credential_type":"recovery_code","masked_hint":"Recovery code / ****","status":"active","payload":{"username":"target-smoke-user","password":"target-credential-smoke-secret","host":"target-smoke.invalid"},"revealed_at":"2026-05-17T00:00:00Z","reveal_expires_message":"shown once"},"request_id":"req_test"}`))
	}))
	defer server.Close()

	record, err := revealTargetCredential(context.Background(), server.Client(), server.URL, &http.Cookie{Name: "billing_session", Value: "session-token"}, fixture)
	if err != nil {
		t.Fatalf("revealTargetCredential returned error: %v", err)
	}
	if record.Type != targetCredentialRevealType || record.MaskedHint != targetCredentialRevealMaskedHint {
		t.Fatalf("unexpected reveal metadata")
	}
}

func TestValidateTargetCredentialRevealResponseRejectsSensitiveMetadataWithoutLeakingPayload(t *testing.T) {
	fixture := targetCredentialRevealFixture{
		CredentialID:     "credential_1",
		EncryptedPayload: "encrypted-fixture",
		ServiceDisplayID: 43001,
		ExpectedPayload:  targetCredentialRevealPayload,
	}
	header := http.Header{}
	header.Set("Cache-Control", "no-store")
	header.Set("Pragma", "no-cache")
	body := []byte(`{"data":{"id":"credential_1","credential_type":"recovery_code","masked_hint":"Recovery code / ****","status":"active","payload":{"username":"target-smoke-user","password":"target-credential-smoke-secret","host":"target-smoke.invalid"},"encrypted_payload":"encrypted-fixture","revealed_at":"2026-05-17T00:00:00Z","reveal_expires_message":"shown once"},"request_id":"req_test"}`)

	_, err := validateTargetCredentialRevealResponse(header, body, fixture)
	if err == nil {
		t.Fatal("expected sensitive response metadata to fail")
	}
	for _, leaked := range []string{"target-credential-smoke-secret", "encrypted-fixture"} {
		if strings.Contains(err.Error(), leaked) {
			t.Fatalf("error leaked credential material: %v", err)
		}
	}
}

func TestValidateTargetCredentialRevealDBEvidenceRejectsAuditLeaksWithoutValues(t *testing.T) {
	err := validateTargetCredentialRevealDBEvidence(targetCredentialRevealDBEvidence{
		LastRevealedByClient: true,
		RateLimitAttempts:    1,
		AuditCount:           1,
		AuditDisplayID:       70001,
		AuditHasDisplayID:    true,
		AuditLeakedSecret:    true,
	})
	if err == nil {
		t.Fatal("expected audit leak evidence to fail")
	}
	if strings.Contains(err.Error(), "target-credential-smoke-secret") {
		t.Fatalf("error leaked credential material: %v", err)
	}
}
