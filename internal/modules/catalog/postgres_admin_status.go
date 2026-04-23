package catalog

import (
	"context"
	"fmt"
)

func (store *PostgresStore) UpdateProductStatus(ctx context.Context, input UpdateProductStatusInput) (Product, error) {
	if err := store.ready(); err != nil {
		return Product{}, err
	}
	query, args, err := buildUpdateProductStatusQuery(input)
	if err != nil {
		return Product{}, err
	}
	return scanProduct(store.executor.QueryRowContext(ctx, query, args...))
}

func (store *PostgresStore) UpdatePlanStatus(ctx context.Context, input UpdatePlanStatusInput) (Plan, error) {
	if err := store.ready(); err != nil {
		return Plan{}, err
	}
	query, args, err := buildUpdatePlanStatusQuery(input)
	if err != nil {
		return Plan{}, err
	}
	return scanPlan(store.executor.QueryRowContext(ctx, query, args...))
}

func (store *PostgresStore) UpdateProviderSourceStatus(ctx context.Context, input UpdateProviderSourceStatusInput) (ProviderSource, error) {
	if err := store.ready(); err != nil {
		return ProviderSource{}, err
	}
	query, args, err := buildUpdateProviderSourceStatusQuery(input)
	if err != nil {
		return ProviderSource{}, err
	}
	return scanProviderSource(store.executor.QueryRowContext(ctx, query, args...))
}

func buildUpdateProductStatusQuery(input UpdateProductStatusInput) (string, []interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf(`UPDATE master_products
SET status = $2, updated_at = NOW()
WHERE product_id = $1
RETURNING %s`, productColumns)
	return query, []interface{}{input.ID, input.Status}, nil
}

func buildUpdatePlanStatusQuery(input UpdatePlanStatusInput) (string, []interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf(`UPDATE master_plans
SET status = $2, updated_at = NOW()
WHERE plan_id = $1
RETURNING %s`, planColumns)
	return query, []interface{}{input.ID, input.Status}, nil
}

func buildUpdateProviderSourceStatusQuery(input UpdateProviderSourceStatusInput) (string, []interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf(`UPDATE provider_sources
SET status = $2, updated_at = NOW()
WHERE source_id = $1
RETURNING %s`, providerSourceColumns)
	return query, []interface{}{input.ID, input.Status}, nil
}
