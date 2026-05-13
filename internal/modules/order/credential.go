package order

import (
	"context"
	"encoding/json"
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
	ErrCredentialIDMissing            = errors.New("credential id missing")
	ErrCredentialTypeMissing          = errors.New("credential type missing")
	ErrCredentialTypeInvalid          = errors.New("credential type invalid")
	ErrCredentialPayloadMissing       = errors.New("credential encrypted payload missing")
	ErrCredentialKeyVersionMissing    = errors.New("credential encryption key version missing")
	ErrCredentialAlgorithmMissing     = errors.New("credential encryption algorithm missing")
	ErrCredentialMaskedHintMissing    = errors.New("credential masked hint missing")
	ErrCredentialStatusInvalid        = errors.New("credential status invalid")
	ErrCredentialNotFound             = errors.New("credential not found")
	ErrCredentialStoreMissing         = errors.New("credential store missing")
	ErrCredentialCipherMissing        = errors.New("credential cipher missing")
	ErrCredentialDecryptFailed        = errors.New("credential decrypt failed")
	ErrCredentialRevealRateLimited    = errors.New("credential reveal rate limited")
	ErrCredentialRevealLimiterMissing = errors.New("credential reveal limiter missing")
	ErrCredentialRevealTimeInvalid    = errors.New("credential reveal time invalid")
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

type ServiceCredentialFilter struct {
	TenantID  tenant.ID
	ServiceID ServiceID
	Status    CredentialStatus
}

type ServiceCredentialLookup struct {
	ID        CredentialID
	TenantID  tenant.ID
	ServiceID ServiceID
}

type MarkServiceCredentialRevealedInput struct {
	ID         CredentialID
	TenantID   tenant.ID
	ServiceID  ServiceID
	ActorID    identity.UserID
	RevealedAt time.Time
}

type RevealServiceCredentialInput struct {
	TenantID     tenant.ID
	ServiceID    ServiceID
	CredentialID CredentialID
	ActorID      identity.UserID
	BuyerUserID  identity.UserID
	ClientIP     string
	UserAgent    string
	Reason       string
}

type RevealServiceCredentialResult struct {
	Credential           ServiceCredential
	Payload              json.RawMessage
	RevealedAt           time.Time
	RevealExpiresMessage string
}

type ServiceCredentialStore interface {
	ListServiceCredentials(ctx context.Context, filter ServiceCredentialFilter) ([]ServiceCredential, error)
	GetServiceCredential(ctx context.Context, lookup ServiceCredentialLookup) (ServiceCredential, error)
	MarkServiceCredentialRevealed(ctx context.Context, input MarkServiceCredentialRevealedInput) (ServiceCredential, error)
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

func (filter ServiceCredentialFilter) Normalize() ServiceCredentialFilter {
	return ServiceCredentialFilter{
		TenantID:  tenant.ID(trim(string(filter.TenantID))),
		ServiceID: ServiceID(trim(string(filter.ServiceID))),
		Status:    filter.Status,
	}
}

func (filter ServiceCredentialFilter) Validate() error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrCredentialStatusInvalid
	}
	return nil
}

func (lookup ServiceCredentialLookup) Normalize() ServiceCredentialLookup {
	return ServiceCredentialLookup{
		ID:        CredentialID(trim(string(lookup.ID))),
		TenantID:  tenant.ID(trim(string(lookup.TenantID))),
		ServiceID: ServiceID(trim(string(lookup.ServiceID))),
	}
}

func (lookup ServiceCredentialLookup) Validate() error {
	if lookup.ID.Empty() {
		return ErrCredentialIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if lookup.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	return nil
}

func (input MarkServiceCredentialRevealedInput) Normalize() MarkServiceCredentialRevealedInput {
	return MarkServiceCredentialRevealedInput{
		ID:         CredentialID(trim(string(input.ID))),
		TenantID:   tenant.ID(trim(string(input.TenantID))),
		ServiceID:  ServiceID(trim(string(input.ServiceID))),
		ActorID:    identity.UserID(trim(string(input.ActorID))),
		RevealedAt: input.RevealedAt,
	}
}

func (input MarkServiceCredentialRevealedInput) Validate() error {
	if input.ID.Empty() {
		return ErrCredentialIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if trim(string(input.ActorID)) == "" {
		return identity.ErrActorIDMissing
	}
	if input.RevealedAt.IsZero() {
		return ErrCredentialRevealTimeInvalid
	}
	return nil
}

func (input RevealServiceCredentialInput) Normalize() RevealServiceCredentialInput {
	return RevealServiceCredentialInput{
		TenantID:     tenant.ID(trim(string(input.TenantID))),
		ServiceID:    ServiceID(trim(string(input.ServiceID))),
		CredentialID: CredentialID(trim(string(input.CredentialID))),
		ActorID:      identity.UserID(trim(string(input.ActorID))),
		BuyerUserID:  identity.UserID(trim(string(input.BuyerUserID))),
		ClientIP:     trim(input.ClientIP),
		UserAgent:    trim(input.UserAgent),
		Reason:       trim(input.Reason),
	}
}

func (input RevealServiceCredentialInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if input.CredentialID.Empty() {
		return ErrCredentialIDMissing
	}
	if trim(string(input.ActorID)) == "" {
		return identity.ErrActorIDMissing
	}
	return nil
}
