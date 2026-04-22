package tenant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrStoreExecutorMissing = errors.New("store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const tenantColumns = `tenant_id, parent_tenant_id, tenant_type, name, slug, status, default_currency, timezone, owner_user_id, branding_settings, billing_settings, risk_settings, created_at, updated_at`
const domainColumns = `domain_id, tenant_id, domain, domain_type, verification_status, verification_token_hash, tls_status, is_primary, created_at, updated_at`

func (store *PostgresStore) Create(ctx context.Context, input CreateTenantInput) (Tenant, error) {
	if err := store.ready(); err != nil {
		return Tenant{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Tenant{}, err
	}

	row := store.executor.QueryRowContext(ctx, `
INSERT INTO tenants (parent_tenant_id, tenant_type, name, slug, status, default_currency, timezone, owner_user_id, branding_settings, billing_settings, risk_settings)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10::jsonb, $11::jsonb)
RETURNING `+tenantColumns,
		nullableTenantID(input.ParentID), input.Type, input.Name, input.Slug, input.Status, input.DefaultCurrency,
		input.Timezone, nullableString(input.OwnerUserID), string(input.BrandingSettings), string(input.BillingSettings), string(input.RiskSettings))
	return scanTenant(row)
}

func (store *PostgresStore) GetByID(ctx context.Context, tenantID ID) (Tenant, error) {
	if err := store.ready(); err != nil {
		return Tenant{}, err
	}
	if tenantID.Empty() {
		return Tenant{}, ErrTenantIDMissing
	}
	row := store.executor.QueryRowContext(ctx, `SELECT `+tenantColumns+` FROM tenants WHERE tenant_id = $1`, tenantID)
	return scanTenant(row)
}

func (store *PostgresStore) FindBySlug(ctx context.Context, slug string) (Tenant, error) {
	if err := store.ready(); err != nil {
		return Tenant{}, err
	}
	slug = lowerTrim(slug)
	if slug == "" {
		return Tenant{}, ErrTenantSlugMissing
	}
	row := store.executor.QueryRowContext(ctx, `SELECT `+tenantColumns+` FROM tenants WHERE slug = $1`, slug)
	return scanTenant(row)
}

func (store *PostgresStore) CreateDomain(ctx context.Context, input CreateDomainInput) (Domain, error) {
	if err := store.ready(); err != nil {
		return Domain{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Domain{}, err
	}

	row := store.executor.QueryRowContext(ctx, `
INSERT INTO tenant_domains (tenant_id, domain, domain_type, verification_status, verification_token_hash, tls_status, is_primary)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING `+domainColumns,
		input.TenantID, input.Domain, input.Type, input.VerificationStatus, nullableString(input.VerificationTokenHash), input.TLSStatus, input.IsPrimary)
	return scanDomain(row)
}

func (store *PostgresStore) FindActiveDomain(ctx context.Context, domain string) (Domain, error) {
	if err := store.ready(); err != nil {
		return Domain{}, err
	}
	domain = lowerTrim(domain)
	if domain == "" {
		return Domain{}, ErrDomainMissing
	}
	row := store.executor.QueryRowContext(ctx, `
SELECT d.domain_id, d.tenant_id, d.domain, d.domain_type, d.verification_status, d.verification_token_hash, d.tls_status, d.is_primary, d.created_at, d.updated_at
FROM tenant_domains d
JOIN tenants t ON t.tenant_id = d.tenant_id
WHERE d.domain = $1
  AND d.verification_status = 'verified'
  AND d.tls_status = 'active'
  AND t.status = 'active'`, domain)
	return scanDomain(row)
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanTenant(row scanner) (Tenant, error) {
	var record Tenant
	var id, tenantType, status string
	var parentID, ownerUserID sql.NullString
	var brandingSettings, billingSettings, riskSettings []byte

	if err := row.Scan(
		&id, &parentID, &tenantType, &record.Name, &record.Slug, &status, &record.DefaultCurrency, &record.Timezone,
		&ownerUserID, &brandingSettings, &billingSettings, &riskSettings, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return Tenant{}, mapTenantNotFound(err)
	}
	record.ID = ID(id)
	record.ParentID = ID(parentID.String)
	record.Type = Type(tenantType)
	record.Status = Status(status)
	record.OwnerUserID = ownerUserID.String
	record.BrandingSettings = append(record.BrandingSettings, brandingSettings...)
	record.BillingSettings = append(record.BillingSettings, billingSettings...)
	record.RiskSettings = append(record.RiskSettings, riskSettings...)
	return record, nil
}

func scanDomain(row scanner) (Domain, error) {
	var record Domain
	var tenantID, domainType, verificationStatus, tlsStatus string
	var verificationTokenHash sql.NullString
	if err := row.Scan(
		&record.ID, &tenantID, &record.Domain, &domainType, &verificationStatus, &verificationTokenHash,
		&tlsStatus, &record.IsPrimary, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Domain{}, ErrDomainNotFound
		}
		return Domain{}, fmt.Errorf("scan tenant domain: %w", err)
	}
	record.TenantID = ID(tenantID)
	record.Type = DomainType(domainType)
	record.VerificationStatus = DomainVerificationStatus(verificationStatus)
	record.VerificationTokenHash = verificationTokenHash.String
	record.TLSStatus = TLSStatus(tlsStatus)
	return record, nil
}

func mapTenantNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTenantNotFound
	}
	return fmt.Errorf("scan tenant: %w", err)
}

func nullableTenantID(id ID) sql.NullString {
	return nullableString(string(id))
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
