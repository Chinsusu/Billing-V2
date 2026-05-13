package order

import (
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	DefaultCredentialEncryptionAlgorithm  = "aes-256-gcm"
	DefaultCredentialEncryptionKeyVersion = "v1"
	DefaultCredentialMaskedHint           = "Encrypted credential"
)

var (
	ErrCredentialIDMissing         = errors.New("credential id missing")
	ErrCredentialTypeMissing       = errors.New("credential type missing")
	ErrCredentialTypeInvalid       = errors.New("credential type invalid")
	ErrCredentialPayloadMissing    = errors.New("credential encrypted payload missing")
	ErrCredentialKeyVersionMissing = errors.New("credential encryption key version missing")
	ErrCredentialAlgorithmMissing  = errors.New("credential encryption algorithm missing")
	ErrCredentialMaskedHintMissing = errors.New("credential masked hint missing")
	ErrCredentialStatusInvalid     = errors.New("credential status invalid")
)

type CredentialID string
type CredentialType string
type CredentialStatus string

const (
	CredentialTypeVPSRoot      CredentialType = "vps_root"
	CredentialTypeProxyAuth    CredentialType = "proxy_auth"
	CredentialTypeSSHKey       CredentialType = "ssh_key"
	CredentialTypeConsoleURL   CredentialType = "console_url"
	CredentialTypeAPIToken     CredentialType = "api_token"
	CredentialTypeRecoveryCode CredentialType = "recovery_code"
)

const (
	CredentialStatusActive  CredentialStatus = "active"
	CredentialStatusRotated CredentialStatus = "rotated"
	CredentialStatusRevoked CredentialStatus = "revoked"
	CredentialStatusExpired CredentialStatus = "expired"
)

type ServiceCredential struct {
	ID                   CredentialID
	TenantID             tenant.ID
	ServiceID            ServiceID
	Type                 CredentialType
	EncryptedPayload     string
	EncryptionKeyVersion string
	EncryptionAlgorithm  string
	SecretVersion        string
	MaskedHint           string
	Status               CredentialStatus
	LastRevealedAt       time.Time
	LastRevealedBy       identity.UserID
	RotatedAt            time.Time
	RotatedBy            identity.UserID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CreateServiceCredentialInput struct {
	TenantID             tenant.ID
	ServiceID            ServiceID
	Type                 CredentialType
	EncryptedPayload     string
	EncryptionKeyVersion string
	EncryptionAlgorithm  string
	SecretVersion        string
	MaskedHint           string
	Status               CredentialStatus
}

func (id CredentialID) Empty() bool { return trim(string(id)) == "" }

func (credentialType CredentialType) Valid() bool {
	switch credentialType {
	case CredentialTypeVPSRoot, CredentialTypeProxyAuth, CredentialTypeSSHKey, CredentialTypeConsoleURL, CredentialTypeAPIToken, CredentialTypeRecoveryCode:
		return true
	default:
		return false
	}
}

func (status CredentialStatus) Valid() bool {
	switch status {
	case CredentialStatusActive, CredentialStatusRotated, CredentialStatusRevoked, CredentialStatusExpired:
		return true
	default:
		return false
	}
}

func (input CreateServiceCredentialInput) Normalize() CreateServiceCredentialInput {
	output := input
	output.EncryptedPayload = trim(output.EncryptedPayload)
	output.EncryptionKeyVersion = trim(output.EncryptionKeyVersion)
	output.EncryptionAlgorithm = trim(output.EncryptionAlgorithm)
	output.SecretVersion = trim(output.SecretVersion)
	output.MaskedHint = trim(output.MaskedHint)
	if output.EncryptionKeyVersion == "" {
		output.EncryptionKeyVersion = DefaultCredentialEncryptionKeyVersion
	}
	if output.EncryptionAlgorithm == "" {
		output.EncryptionAlgorithm = DefaultCredentialEncryptionAlgorithm
	}
	if output.MaskedHint == "" {
		output.MaskedHint = DefaultCredentialMaskedHint
	}
	if output.Status == "" {
		output.Status = CredentialStatusActive
	}
	return output
}

func (input CreateServiceCredentialInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if input.Type == "" {
		return ErrCredentialTypeMissing
	}
	if !input.Type.Valid() {
		return ErrCredentialTypeInvalid
	}
	if input.EncryptedPayload == "" {
		return ErrCredentialPayloadMissing
	}
	if input.EncryptionKeyVersion == "" {
		return ErrCredentialKeyVersionMissing
	}
	if input.EncryptionAlgorithm == "" {
		return ErrCredentialAlgorithmMissing
	}
	if input.MaskedHint == "" {
		return ErrCredentialMaskedHintMissing
	}
	if !input.Status.Valid() {
		return ErrCredentialStatusInvalid
	}
	return nil
}
