package catalog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresStore struct {
	executor platformdb.Executor
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const productColumns = `product_id, display_id, product_type, name, description, status, display_order, created_by, created_at, updated_at`
const planColumns = `plan_id, display_id, product_id, plan_code, name, specs, billing_cycle_type, billing_cycle_value, base_cost_minor, suggested_price_minor, reseller_min_price_minor, currency, status, version, created_at, updated_at`
const providerSourceColumns = `source_id, display_id, source_type, name, provider_account_id, location, status, capability_profile, inventory_mode, risk_level, created_at, updated_at`
const planSourceColumns = `plan_source_id, display_id, plan_id, source_id, priority, cost_override_minor, capacity_policy, capability_override, status, created_at, updated_at`
const tenantProductColumns = `tenant_product_id, display_id, tenant_id, master_product_id, name_override, description_override, status, clone_version, created_at, updated_at`
const tenantPlanColumns = `tenant_plan_id, display_id, tenant_id, tenant_product_id, master_plan_id, selling_price_minor, reseller_cost_minor, currency, margin_policy, visibility, status, clone_version, product_snapshot, plan_snapshot, price_snapshot, capability_snapshot, created_at, updated_at`

func (store *PostgresStore) CreateProduct(ctx context.Context, input CreateProductInput) (Product, error) {
	if err := store.ready(); err != nil {
		return Product{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Product{}, err
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO master_products (product_type, name, description, status, display_order, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING `+productColumns,
		input.Type, input.Name, nullableString(input.Description), input.Status, input.DisplayOrder, input.CreatedBy)
	return scanProduct(row)
}

func (store *PostgresStore) CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error) {
	if err := store.ready(); err != nil {
		return Plan{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Plan{}, err
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO master_plans (product_id, plan_code, name, specs, billing_cycle_type, billing_cycle_value, base_cost_minor, suggested_price_minor, reseller_min_price_minor, currency, status, version)
VALUES ($1, $2, $3, $4::jsonb, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING `+planColumns,
		input.ProductID, input.Code, input.Name, string(input.Specs), input.BillingCycle.Type, input.BillingCycle.Value,
		input.BaseCostMinor, input.SuggestedPriceMinor, input.ResellerMinPriceMinor, input.Currency, input.Status, input.Version)
	return scanPlan(row)
}

func (store *PostgresStore) CreateProviderSource(ctx context.Context, input CreateProviderSourceInput) (ProviderSource, error) {
	if err := store.ready(); err != nil {
		return ProviderSource{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ProviderSource{}, err
	}
	profileJSON, err := json.Marshal(input.CapabilityProfile)
	if err != nil {
		return ProviderSource{}, fmt.Errorf("marshal capability profile: %w", err)
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO provider_sources (source_type, name, provider_account_id, location, status, capability_profile, inventory_mode, risk_level)
VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8)
RETURNING `+providerSourceColumns,
		input.Type, input.Name, nullableString(string(input.ProviderAccountID)), nullableString(input.Location), input.Status,
		string(profileJSON), input.InventoryMode, input.RiskLevel)
	return scanProviderSource(row)
}

func (store *PostgresStore) CreatePlanSource(ctx context.Context, input CreatePlanSourceInput) (PlanSource, error) {
	if err := store.ready(); err != nil {
		return PlanSource{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return PlanSource{}, err
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO plan_sources (plan_id, source_id, priority, cost_override_minor, capacity_policy, capability_override, status)
VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7)
RETURNING `+planSourceColumns,
		input.PlanID, input.SourceID, input.Priority, input.CostOverrideMinor, string(input.CapacityPolicy), string(input.CapabilityOverride), input.Status)
	return scanPlanSource(row)
}

func (store *PostgresStore) CreateTenantProduct(ctx context.Context, input CreateTenantProductInput) (TenantProduct, error) {
	if err := store.ready(); err != nil {
		return TenantProduct{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TenantProduct{}, err
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO tenant_products (tenant_id, master_product_id, name_override, description_override, status, clone_version)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING `+tenantProductColumns,
		input.TenantID, input.MasterProductID, nullableString(input.NameOverride), nullableString(input.DescriptionOverride), input.Status, input.CloneVersion)
	return scanTenantProduct(row)
}

func (store *PostgresStore) CreateTenantPlan(ctx context.Context, input CreateTenantPlanInput) (TenantPlan, error) {
	if err := store.ready(); err != nil {
		return TenantPlan{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TenantPlan{}, err
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO tenant_plans (tenant_id, tenant_product_id, master_plan_id, selling_price_minor, reseller_cost_minor, currency, margin_policy, visibility, status, clone_version, product_snapshot, plan_snapshot, price_snapshot, capability_snapshot)
VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9, $10, $11::jsonb, $12::jsonb, $13::jsonb, $14::jsonb)
RETURNING `+tenantPlanColumns,
		input.TenantID, input.TenantProductID, input.MasterPlanID, input.SellingPriceMinor, input.ResellerCostMinor,
		input.Currency, string(input.MarginPolicy), input.Visibility, input.Status, input.CloneVersion, string(input.ProductSnapshot),
		string(input.PlanSnapshot), string(input.PriceSnapshot), string(input.CapabilitySnapshot))
	return scanTenantPlan(row)
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrCatalogStoreExecutorMissing
	}
	return nil
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
