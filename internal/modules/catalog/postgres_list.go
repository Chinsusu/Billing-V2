package catalog

import (
	"context"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const defaultCatalogListLimit = 100
const maxCatalogListLimit = 500

const planListColumns = `mp.plan_id, mp.display_id, mp.product_id, mp.plan_code, mp.name, mp.specs, mp.billing_cycle_type, mp.billing_cycle_value, mp.base_cost_minor, mp.suggested_price_minor, mp.reseller_min_price_minor, mp.currency, mp.status, mp.version, mp.created_at, mp.updated_at`
const tenantProductListColumns = `tp.tenant_product_id, tp.display_id, tp.tenant_id, tp.master_product_id, tp.name_override, tp.description_override, tp.status, tp.clone_version, tp.created_at, tp.updated_at`
const tenantPlanListColumns = `tpl.tenant_plan_id, tpl.display_id, tpl.tenant_id, tpl.tenant_product_id, tpl.master_plan_id, tpl.selling_price_minor, tpl.reseller_cost_minor, tpl.currency, tpl.margin_policy, tpl.visibility, tpl.status, tpl.clone_version, tpl.product_snapshot, tpl.plan_snapshot, tpl.price_snapshot, tpl.capability_snapshot, tpl.created_at, tpl.updated_at`

func (store *PostgresStore) ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListMasterPlansQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalog master plans: %w", err)
	}
	defer rows.Close()
	plans := make([]Plan, 0)
	for rows.Next() {
		plan, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read catalog master plans: %w", err)
	}
	return plans, nil
}

func (store *PostgresStore) ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error) {
	if err := store.ready(); err != nil {
		return TenantCatalog{}, err
	}
	if err := validateTenantCatalogFilter(filter); err != nil {
		return TenantCatalog{}, err
	}
	products, err := store.listTenantProducts(ctx, filter)
	if err != nil {
		return TenantCatalog{}, err
	}
	plans, err := store.listTenantPlans(ctx, filter)
	if err != nil {
		return TenantCatalog{}, err
	}
	return TenantCatalog{Products: products, Plans: plans}, nil
}

func (store *PostgresStore) listTenantProducts(ctx context.Context, filter TenantCatalogFilter) ([]TenantProduct, error) {
	query, args, err := buildListTenantProductsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalog tenant products: %w", err)
	}
	defer rows.Close()
	products := make([]TenantProduct, 0)
	for rows.Next() {
		product, err := scanTenantProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read catalog tenant products: %w", err)
	}
	return products, nil
}

func (store *PostgresStore) listTenantPlans(ctx context.Context, filter TenantCatalogFilter) ([]TenantPlan, error) {
	query, args, err := buildListTenantPlansQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalog tenant plans: %w", err)
	}
	defer rows.Close()
	plans := make([]TenantPlan, 0)
	for rows.Next() {
		plan, err := scanTenantPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read catalog tenant plans: %w", err)
	}
	return plans, nil
}

func buildListMasterPlansQuery(filter MasterPlanFilter) (string, []interface{}, error) {
	if err := validateMasterPlanFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + planListColumns + `
FROM master_plans mp
JOIN master_products p ON p.product_id = mp.product_id
WHERE TRUE`
	args := make([]interface{}, 0, 3)
	if filter.ProductType != "" {
		args = append(args, filter.ProductType)
		query += fmt.Sprintf("\n  AND p.product_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND mp.status = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY mp.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildListTenantProductsQuery(filter TenantCatalogFilter) (string, []interface{}, error) {
	if err := validateTenantCatalogFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + tenantProductListColumns + `
FROM tenant_products tp
JOIN master_products mp ON mp.product_id = tp.master_product_id
WHERE tp.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.ProductType != "" {
		args = append(args, filter.ProductType)
		query += fmt.Sprintf("\n  AND mp.product_type = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY tp.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildListTenantPlansQuery(filter TenantCatalogFilter) (string, []interface{}, error) {
	if err := validateTenantCatalogFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + tenantPlanListColumns + `
FROM tenant_plans tpl
JOIN master_plans mpl ON mpl.plan_id = tpl.master_plan_id
JOIN master_products mp ON mp.product_id = mpl.product_id
WHERE tpl.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.ProductType != "" {
		args = append(args, filter.ProductType)
		query += fmt.Sprintf("\n  AND mp.product_type = $%d", len(args))
	}
	if filter.Visibility != "" {
		args = append(args, filter.Visibility)
		query += fmt.Sprintf("\n  AND tpl.visibility = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND tpl.status = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY tpl.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func validateMasterPlanFilter(filter MasterPlanFilter) error {
	if filter.ProductType != "" && !filter.ProductType.Valid() {
		return ErrProductTypeInvalid
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrPlanStatusInvalid
	}
	return nil
}

func validateTenantCatalogFilter(filter TenantCatalogFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.ProductType != "" && !filter.ProductType.Valid() {
		return ErrProductTypeInvalid
	}
	if filter.Visibility != "" && !filter.Visibility.Valid() {
		return ErrTenantPlanVisibility
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrTenantPlanStatus
	}
	return nil
}

func normalizeCatalogListLimit(limit int) int {
	if limit <= 0 {
		return defaultCatalogListLimit
	}
	if limit > maxCatalogListLimit {
		return maxCatalogListLimit
	}
	return limit
}
