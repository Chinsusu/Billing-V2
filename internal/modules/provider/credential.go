package provider

import (
	"encoding/json"
	"errors"
	"strings"
)

const DefaultCredentialKeyVersion = "v1"

var (
	ErrCredentialCipherMissing  = errors.New("credential cipher missing")
	ErrCredentialTypeMissing    = errors.New("credential type missing")
	ErrCredentialPayloadMissing = errors.New("credential payload missing")
	ErrCredentialMaskMissing    = errors.New("credential masked hint missing")
)

type CredentialType string

const (
	CredentialTypeVPSRoot      CredentialType = "vps_root"
	CredentialTypeProxyAuth    CredentialType = "proxy_auth"
	CredentialTypeSSHKey       CredentialType = "ssh_key"
	CredentialTypeConsoleURL   CredentialType = "console_url"
	CredentialTypeAPIToken     CredentialType = "api_token"
	CredentialTypeRecoveryCode CredentialType = "recovery_code"
)

type CredentialCipher interface {
	Encrypt(plaintext string) (string, error)
}

func NewEncryptedCredentialEnvelope(
	credentialType CredentialType,
	payload map[string]string,
	maskedHint string,
	keyVersion string,
	cipher CredentialCipher,
) (CredentialEnvelope, error) {
	if credentialType == "" {
		return CredentialEnvelope{}, ErrCredentialTypeMissing
	}
	if len(payload) == 0 {
		return CredentialEnvelope{}, ErrCredentialPayloadMissing
	}
	if strings.TrimSpace(maskedHint) == "" {
		return CredentialEnvelope{}, ErrCredentialMaskMissing
	}
	if cipher == nil {
		return CredentialEnvelope{}, ErrCredentialCipherMissing
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return CredentialEnvelope{}, err
	}
	encrypted, err := cipher.Encrypt(string(body))
	if err != nil {
		return CredentialEnvelope{}, err
	}
	if keyVersion = strings.TrimSpace(keyVersion); keyVersion == "" {
		keyVersion = DefaultCredentialKeyVersion
	}
	return CredentialEnvelope{
		Type:                 credentialType,
		EncryptedPayload:     encrypted,
		EncryptionKeyVersion: keyVersion,
		MaskedHint:           strings.TrimSpace(maskedHint),
	}, nil
}

func (envelope CredentialEnvelope) EncryptedPayloadValue() string {
	if value := strings.TrimSpace(envelope.EncryptedPayload); value != "" {
		return value
	}
	return strings.TrimSpace(envelope.EncryptedPayloadRef)
}

func (envelope CredentialEnvelope) HasEncryptedPayload() bool {
	return envelope.EncryptedPayloadValue() != ""
}
