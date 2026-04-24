package tenant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrStoreExecutorMissing = errors.New("store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const tenantColumns = `tenant_id, display_id, parent_tenant_id, tenant_type, name, slug, status, default_currency, timezone, owner_user_id, branding_settings, billing_settings, risk_settings, created_at, updated_at`
const domainColumns = `domain_id, display_id, tenant_id, domain, domain_type, verification_status, verification_token_hash, tls_status, is_primary, created_at, updated_at`

var tenantColumnNames = []string{
	"tenant_id",
	"display_id",
	"parent_tenant_id",
	"tenant_type",
	"name",
	"slug",
	"status",
	"default_currency",
	"timezone",
	"owner_user_id",
	"branding_settings",
	"billing_settings",
	"risk_settings",
	"created_at",
	"updated_at",
}

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

func (store *PostgresStore) ListTenants(ctx context.Context, filter ListTenantsFilter) ([]TenantSummary, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args := buildListTenantsQuery(filter)
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	records := []TenantSummary{}
	for rows.Next() {
		record, err := scanTenantSummary(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tenants: %w", err)
	}
	return records, nil
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
SELECT d.domain_id, d.display_id, d.tenant_id, d.domain, d.domain_type, d.verification_status, d.verification_token_hash, d.tls_status, d.is_primary, d.created_at, d.updated_at
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

func buildListTenantsQuery(filter ListTenantsFilter) (string, []interface{}) {
	args := []interface{}{}
	conditions := []string{}
	if !filter.ScopeTenantID.Empty() {
		args = append(args, filter.ScopeTenantID)
		placeholder := fmt.Sprintf("$%d", len(args))
		conditions = append(conditions, fmt.Sprintf("(t.tenant_id = %s OR t.parent_tenant_id = %s)", placeholder, placeholder))
	}
	if !filter.ParentID.Empty() {
		args = append(args, filter.ParentID)
		conditions = append(conditions, fmt.Sprintf("t.parent_tenant_id = $%d", len(args)))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		conditions = append(conditions, fmt.Sprintf("t.tenant_type = $%d", len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		conditions = append(conditions, fmt.Sprintf("t.status = $%d", len(args)))
	}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		conditions = append(conditions, fmt.Sprintf("t.display_id = $%d", len(args)))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	args = append(args, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))

	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(tenantSelectColumns("t"))
	builder.WriteString(`,
       COALESCE(primary_domain.domain, '') AS primary_domain,
       COALESCE(user_counts.user_count, 0) AS user_count
FROM tenants t
LEFT JOIN LATERAL (
    SELECT domain
    FROM tenant_domains d
    WHERE d.tenant_id = t.tenant_id
      AND d.is_primary = TRUE
    ORDER BY d.created_at DESC
    LIMIT 1
) primary_domain ON TRUE
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS user_count
    FROM users u
    WHERE u.tenant_id = t.tenant_id
) user_counts ON TRUE`)
	if len(conditions) > 0 {
		builder.WriteString("\nWHERE ")
		builder.WriteString(strings.Join(conditions, "\n  AND "))
	}
	builder.WriteString("\nORDER BY t.created_at DESC, t.display_id DESC")
	builder.WriteString("\nLIMIT ")
	builder.WriteString(limitPlaceholder)
	return builder.String(), args
}

func tenantSelectColumns(alias string) string {
	columns := make([]string, 0, len(tenantColumnNames))
	for _, column := range tenantColumnNames {
		if alias == "" {
			columns = append(columns, column)
			continue
		}
		columns = append(columns, alias+"."+column)
	}
	return strings.Join(columns, ", ")
}

func scanTenant(row scanner) (Tenant, error) {
	var record Tenant
	var id, tenantType, status string
	var parentID, ownerUserID sql.NullString
	var brandingSettings, billingSettings, riskSettings []byte

	if err := row.Scan(
		&id, &record.DisplayID, &parentID, &tenantType, &record.Name, &record.Slug, &status, &record.DefaultCurrency, &record.Timezone,
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

func scanTenantSummary(row scanner) (TenantSummary, error) {
	var record Tenant
	var id, tenantType, status string
	var parentID, ownerUserID sql.NullString
	var brandingSettings, billingSettings, riskSettings []byte
	var primaryDomain sql.NullString
	var userCount int64

	if err := row.Scan(
		&id, &record.DisplayID, &parentID, &tenantType, &record.Name, &record.Slug, &status, &record.DefaultCurrency, &record.Timezone,
		&ownerUserID, &brandingSettings, &billingSettings, &riskSettings, &record.CreatedAt, &record.UpdatedAt,
		&primaryDomain, &userCount,
	); err != nil {
		return TenantSummary{}, mapTenantNotFound(err)
	}
	record.ID = ID(id)
	record.ParentID = ID(parentID.String)
	record.Type = Type(tenantType)
	record.Status = Status(status)
	record.OwnerUserID = ownerUserID.String
	record.BrandingSettings = append(record.BrandingSettings, brandingSettings...)
	record.BillingSettings = append(record.BillingSettings, billingSettings...)
	record.RiskSettings = append(record.RiskSettings, riskSettings...)
	return TenantSummary{
		Tenant:        record,
		PrimaryDomain: primaryDomain.String,
		UserCount:     userCount,
	}, nil
}

func scanDomain(row scanner) (Domain, error) {
	var record Domain
	var tenantID, domainType, verificationStatus, tlsStatus string
	var verificationTokenHash sql.NullString
	if err := row.Scan(
		&record.ID, &record.DisplayID, &tenantID, &record.Domain, &domainType, &verificationStatus, &verificationTokenHash,
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
