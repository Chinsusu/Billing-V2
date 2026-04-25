package catalog

import (
	"context"
	"fmt"
)

const productListColumns = `product_id, display_id, product_type, name, description, status, display_order, created_by, created_at, updated_at`
const providerSourceListColumns = `source_id, display_id, source_type, name, provider_account_id, location, status, capability_profile, inventory_mode, risk_level, created_at, updated_at`

func (store *PostgresStore) ListProducts(ctx context.Context, filter ProductFilter) ([]Product, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListProductsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalog products: %w", err)
	}
	defer rows.Close()
	products := make([]Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read catalog products: %w", err)
	}
	return products, nil
}

func (store *PostgresStore) ListProviderSources(ctx context.Context, filter ProviderSourceFilter) ([]ProviderSource, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListProviderSourcesQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalog provider sources: %w", err)
	}
	defer rows.Close()
	sources := make([]ProviderSource, 0)
	for rows.Next() {
		source, err := scanProviderSource(rows)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read catalog provider sources: %w", err)
	}
	return sources, nil
}

func buildListProductsQuery(filter ProductFilter) (string, []interface{}, error) {
	if err := validateProductFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + productListColumns + `
FROM master_products
WHERE TRUE`
	args := make([]interface{}, 0, 3)
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf("\n  AND product_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND status = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY display_order ASC, created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildListProviderSourcesQuery(filter ProviderSourceFilter) (string, []interface{}, error) {
	if err := validateProviderSourceFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + providerSourceListColumns + `
FROM provider_sources
WHERE TRUE`
	args := make([]interface{}, 0, 3)
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND display_id = $%d", len(args))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf("\n  AND source_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND status = $%d", len(args))
	}
	args = append(args, normalizeCatalogListLimit(filter.Limit))
	query += fmt.Sprintf("\nORDER BY created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func validateProductFilter(filter ProductFilter) error {
	if filter.Type != "" && !filter.Type.Valid() {
		return ErrProductTypeInvalid
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrProductStatusInvalid
	}
	return nil
}

func validateProviderSourceFilter(filter ProviderSourceFilter) error {
	if filter.Type != "" && !providerTypeValid(filter.Type) {
		return ErrSourceTypeInvalid
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrSourceStatusInvalid
	}
	return nil
}
