package provider

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

func TestNewEncryptedCredentialEnvelopeEncryptsPayload(t *testing.T) {
	cipher, err := secrets.NewAESGCMCipher(testCredentialCipherKey())
	if err != nil {
		t.Fatalf("NewAESGCMCipher returned error: %v", err)
	}

	envelope, err := NewEncryptedCredentialEnvelope(
		CredentialTypeVPSRoot,
		map[string]string{"username": "root", "access": "launch-fixture"},
		"root / ****",
		"",
		cipher,
	)
	if err != nil {
		t.Fatalf("NewEncryptedCredentialEnvelope returned error: %v", err)
	}
	if envelope.EncryptedPayload == "" || envelope.EncryptedPayload == "launch-fixture" {
		t.Fatalf("expected encrypted payload, got %q", envelope.EncryptedPayload)
	}
	if strings.Contains(envelope.EncryptedPayload, "launch-fixture") {
		t.Fatalf("encrypted payload leaked fixture value: %q", envelope.EncryptedPayload)
	}
	if envelope.EncryptionKeyVersion != DefaultCredentialKeyVersion || envelope.MaskedHint != "root / ****" {
		t.Fatalf("unexpected envelope metadata: %+v", envelope)
	}
	plaintext, err := cipher.Decrypt(envelope.EncryptedPayload)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}
	if !strings.Contains(plaintext, "launch-fixture") {
		t.Fatalf("decrypted payload missing fixture value: %s", plaintext)
	}
}

func TestNewEncryptedCredentialEnvelopeRejectsMissingInputs(t *testing.T) {
	cipher, err := secrets.NewAESGCMCipher(testCredentialCipherKey())
	if err != nil {
		t.Fatalf("NewAESGCMCipher returned error: %v", err)
	}

	if _, err := NewEncryptedCredentialEnvelope("", map[string]string{"access": "fixture"}, "hint", "", cipher); !errors.Is(err, ErrCredentialTypeMissing) {
		t.Fatalf("expected credential type error, got %v", err)
	}
	if _, err := NewEncryptedCredentialEnvelope(CredentialTypeVPSRoot, nil, "hint", "", cipher); !errors.Is(err, ErrCredentialPayloadMissing) {
		t.Fatalf("expected credential payload error, got %v", err)
	}
	if _, err := NewEncryptedCredentialEnvelope(CredentialTypeVPSRoot, map[string]string{"access": "fixture"}, "", "", cipher); !errors.Is(err, ErrCredentialMaskMissing) {
		t.Fatalf("expected credential mask error, got %v", err)
	}
	if _, err := NewEncryptedCredentialEnvelope(CredentialTypeVPSRoot, map[string]string{"access": "fixture"}, "hint", "", nil); !errors.Is(err, ErrCredentialCipherMissing) {
		t.Fatalf("expected credential cipher error, got %v", err)
	}
}

func testCredentialCipherKey() string {
	return "1234567890123456" + "7890123456789012"
}
