package catalog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func (store *PostgresStore) ListProviderSourceReadiness(ctx context.Context, filter ProviderSourceReadinessFilter) ([]ProviderSourceReadiness, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListProviderSourceReadinessQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list provider source readiness: %w", err)
	}
	defer rows.Close()
	records := make([]ProviderSourceReadiness, 0)
	for rows.Next() {
		record, err := scanProviderSourceReadiness(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read provider source readiness: %w", err)
	}
	return records, nil
}

func buildListProviderSourceReadinessQuery(filter ProviderSourceReadinessFilter) (string, []interface{}, error) {
	filter = normalizeProviderSourceReadinessFilter(filter)
	if err := validateProviderSourceReadinessFilter(filter); err != nil {
		return "", nil, err
	}
	query := `
WITH ranked_plan_sources AS (
  SELECT
    ps.plan_id,
    ps.display_id AS plan_source_display_id,
    ps.status AS plan_source_status,
    src.display_id AS source_display_id,
    src.source_type,
    src.name AS source_name,
    src.status AS source_status,
    src.inventory_mode,
    src.capability_profile,
    ROW_NUMBER() OVER (
      PARTITION BY ps.plan_id
      ORDER BY CASE WHEN ps.status = 'active' AND src.status = 'active' THEN 0 ELSE 1 END,
        ps.priority ASC,
        ps.created_at ASC
    ) AS row_number
  FROM plan_sources ps
  JOIN provider_sources src ON src.source_id = ps.source_id
)
SELECT
  mp.display_id,
  mp.plan_code,
  mp.name,
  product.product_type,
  mp.status,
  selected.plan_source_display_id,
  selected.plan_source_status,
  selected.source_display_id,
  selected.source_type,
  selected.source_name,
  selected.source_status,
  selected.inventory_mode,
  selected.capability_profile
FROM master_plans mp
JOIN master_products product ON product.product_id = mp.product_id
LEFT JOIN ranked_plan_sources selected ON selected.plan_id = mp.plan_id AND selected.row_number = 1
WHERE mp.status = $1`
	args := []interface{}{filter.PlanStatus}
	if filter.ProductType != "" {
		args = append(args, filter.ProductType)
		query += fmt.Sprintf("\n  AND product.product_type = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY mp.display_id ASC\nLIMIT $%d", len(args))
	return query, args, nil
}

func normalizeProviderSourceReadinessFilter(filter ProviderSourceReadinessFilter) ProviderSourceReadinessFilter {
	if filter.PlanStatus == "" {
		filter.PlanStatus = PlanStatusActive
	}
	return filter
}

func validateProviderSourceReadinessFilter(filter ProviderSourceReadinessFilter) error {
	if filter.ProductType != "" && !filter.ProductType.Valid() {
		return ErrProductTypeInvalid
	}
	if filter.PlanStatus != "" && !filter.PlanStatus.Valid() {
		return ErrPlanStatusInvalid
	}
	return nil
}

func scanProviderSourceReadiness(row catalogScanner) (ProviderSourceReadiness, error) {
	var record ProviderSourceReadiness
	var productType, planStatus string
	var planSourceDisplayID, sourceDisplayID sql.NullInt64
	var planSourceStatus, sourceType, sourceName, sourceStatus, inventoryMode, capabilityProfile sql.NullString
	if err := row.Scan(
		&record.PlanDisplayID,
		&record.PlanCode,
		&record.PlanName,
		&productType,
		&planStatus,
		&planSourceDisplayID,
		&planSourceStatus,
		&sourceDisplayID,
		&sourceType,
		&sourceName,
		&sourceStatus,
		&inventoryMode,
		&capabilityProfile,
	); err != nil {
		return ProviderSourceReadiness{}, fmt.Errorf("scan provider source readiness: %w", err)
	}
	record.ProductType = ProductType(productType)
	record.PlanStatus = PlanStatus(planStatus)
	record.PlanSourceDisplayID = planSourceDisplayID.Int64
	record.PlanSourceStatus = PlanSourceStatus(planSourceStatus.String)
	record.SourceDisplayID = sourceDisplayID.Int64
	record.SourceType = provider.Type(sourceType.String)
	record.SourceName = sourceName.String
	record.SourceStatus = ProviderSourceStatus(sourceStatus.String)
	record.InventoryMode = InventoryMode(inventoryMode.String)
	if capabilityProfile.Valid {
		if err := json.Unmarshal([]byte(capabilityProfile.String), &record.capabilityProfile); err != nil {
			return ProviderSourceReadiness{}, fmt.Errorf("decode readiness capability profile: %w", err)
		}
	}
	return record.withReadinessState(), nil
}
