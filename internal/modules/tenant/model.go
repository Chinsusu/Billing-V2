package tenant

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrTenantNotFound         = errors.New("tenant not found")
	ErrTenantNameMissing      = errors.New("tenant name missing")
	ErrTenantSlugMissing      = errors.New("tenant slug missing")
	ErrTenantTypeInvalid      = errors.New("tenant type invalid")
	ErrTenantStatusInvalid    = errors.New("tenant status invalid")
	ErrCurrencyMissing        = errors.New("tenant currency missing")
	ErrDomainNotFound         = errors.New("tenant domain not found")
	ErrDomainMissing          = errors.New("tenant domain missing")
	ErrDomainTypeInvalid      = errors.New("tenant domain type invalid")
	ErrDomainStatusInvalid    = errors.New("tenant domain status invalid")
	ErrDomainTLSStatusInvalid = errors.New("tenant domain tls status invalid")
)

const TypeDirectStore Type = "direct_store"

type Status string

const (
	StatusActive       Status = "active"
	StatusSuspended    Status = "suspended"
	StatusDisabled     Status = "disabled"
	StatusPendingSetup Status = "pending_setup"
)

func (tenantType Type) Valid() bool {
	switch tenantType {
	case TypePlatform, TypeReseller, TypeDirectStore:
		return true
	default:
		return false
	}
}

func (status Status) Valid() bool {
	switch status {
	case StatusActive, StatusSuspended, StatusDisabled, StatusPendingSetup:
		return true
	default:
		return false
	}
}

type Tenant struct {
	ID               ID
	ParentID         ID
	Type             Type
	Name             string
	Slug             string
	Status           Status
	DefaultCurrency  string
	Timezone         string
	OwnerUserID      string
	BrandingSettings json.RawMessage
	BillingSettings  json.RawMessage
	RiskSettings     json.RawMessage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateTenantInput struct {
	ParentID         ID
	Type             Type
	Name             string
	Slug             string
	Status           Status
	DefaultCurrency  string
	Timezone         string
	OwnerUserID      string
	BrandingSettings json.RawMessage
	BillingSettings  json.RawMessage
	RiskSettings     json.RawMessage
}

func (input CreateTenantInput) Normalize() CreateTenantInput {
	output := input
	output.Name = strings.TrimSpace(output.Name)
	output.Slug = lowerTrim(output.Slug)
	output.DefaultCurrency = strings.ToUpper(strings.TrimSpace(output.DefaultCurrency))
	output.Timezone = strings.TrimSpace(output.Timezone)
	if output.Timezone == "" {
		output.Timezone = "Asia/Ho_Chi_Minh"
	}
	if output.Status == "" {
		output.Status = StatusPendingSetup
	}
	output.BrandingSettings = defaultJSON(output.BrandingSettings)
	output.BillingSettings = defaultJSON(output.BillingSettings)
	output.RiskSettings = defaultJSON(output.RiskSettings)
	return output
}

func (input CreateTenantInput) Validate() error {
	if !input.Type.Valid() {
		return ErrTenantTypeInvalid
	}
	if input.Name == "" {
		return ErrTenantNameMissing
	}
	if input.Slug == "" {
		return ErrTenantSlugMissing
	}
	if input.Status != "" && !input.Status.Valid() {
		return ErrTenantStatusInvalid
	}
	if input.DefaultCurrency == "" {
		return ErrCurrencyMissing
	}
	return nil
}

type DomainType string

const (
	DomainTypeSystemSubdomain DomainType = "system_subdomain"
	DomainTypeCustomDomain    DomainType = "custom_domain"
)

func (domainType DomainType) Valid() bool {
	switch domainType {
	case DomainTypeSystemSubdomain, DomainTypeCustomDomain:
		return true
	default:
		return false
	}
}

type DomainVerificationStatus string

const (
	DomainVerificationPending  DomainVerificationStatus = "pending"
	DomainVerificationVerified DomainVerificationStatus = "verified"
	DomainVerificationFailed   DomainVerificationStatus = "failed"
	DomainVerificationDisabled DomainVerificationStatus = "disabled"
)

func (status DomainVerificationStatus) Valid() bool {
	switch status {
	case DomainVerificationPending, DomainVerificationVerified, DomainVerificationFailed, DomainVerificationDisabled:
		return true
	default:
		return false
	}
}

type TLSStatus string

const (
	TLSStatusPending TLSStatus = "pending"
	TLSStatusActive  TLSStatus = "active"
	TLSStatusFailed  TLSStatus = "failed"
	TLSStatusExpired TLSStatus = "expired"
)

func (status TLSStatus) Valid() bool {
	switch status {
	case TLSStatusPending, TLSStatusActive, TLSStatusFailed, TLSStatusExpired:
		return true
	default:
		return false
	}
}

type Domain struct {
	ID                    string
	TenantID              ID
	Domain                string
	Type                  DomainType
	VerificationStatus    DomainVerificationStatus
	VerificationTokenHash string
	TLSStatus             TLSStatus
	IsPrimary             bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CreateDomainInput struct {
	TenantID              ID
	Domain                string
	Type                  DomainType
	VerificationStatus    DomainVerificationStatus
	VerificationTokenHash string
	TLSStatus             TLSStatus
	IsPrimary             bool
}

func (input CreateDomainInput) Normalize() CreateDomainInput {
	output := input
	output.Domain = lowerTrim(output.Domain)
	output.VerificationTokenHash = strings.TrimSpace(output.VerificationTokenHash)
	if output.VerificationStatus == "" {
		output.VerificationStatus = DomainVerificationPending
	}
	if output.TLSStatus == "" {
		output.TLSStatus = TLSStatusPending
	}
	return output
}

func (input CreateDomainInput) Validate() error {
	if input.TenantID.Empty() {
		return ErrTenantIDMissing
	}
	if input.Domain == "" {
		return ErrDomainMissing
	}
	if !input.Type.Valid() {
		return ErrDomainTypeInvalid
	}
	if input.VerificationStatus != "" && !input.VerificationStatus.Valid() {
		return ErrDomainStatusInvalid
	}
	if input.TLSStatus != "" && !input.TLSStatus.Valid() {
		return ErrDomainTLSStatusInvalid
	}
	return nil
}

func lowerTrim(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func defaultJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return append(json.RawMessage(nil), value...)
}
