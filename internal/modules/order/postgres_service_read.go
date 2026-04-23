package order

import (
	"context"
	"fmt"
)

const serviceInstanceReadColumns = `svc.service_instance_id, svc.display_id, svc.tenant_id, svc.order_id, svc.tenant_plan_id, svc.provider_source_id, svc.external_resource_id, svc.status, svc.billing_status, svc.suspension_reason, svc.term_start, svc.term_end, svc.created_at, svc.updated_at`

func (store *PostgresStore) ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListServiceInstancesQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list service instances: %w", err)
	}
	defer rows.Close()
	services := make([]ServiceInstance, 0)
	for rows.Next() {
		service, err := scanServiceInstance(rows)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read service instances: %w", err)
	}
	return services, nil
}

func (store *PostgresStore) GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	if err := store.ready(); err != nil {
		return ServiceInstance{}, err
	}
	query, args, err := buildGetServiceInstanceQuery(lookup)
	if err != nil {
		return ServiceInstance{}, err
	}
	return scanServiceInstance(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListServiceInstancesQuery(filter ServiceInstanceFilter) (string, []interface{}, error) {
	filter = normalizeServiceInstanceFilter(filter)
	if err := validateServiceInstanceFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + serviceInstanceReadColumns + `
FROM service_instances svc
JOIN orders ord ON ord.order_id = svc.order_id
WHERE svc.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.BuyerUserID != "" {
		args = append(args, filter.BuyerUserID)
		query += fmt.Sprintf("\n  AND ord.buyer_user_id = $%d", len(args))
	}
	if filter.OrderID != "" {
		args = append(args, filter.OrderID)
		query += fmt.Sprintf("\n  AND svc.order_id = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND svc.status = $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY svc.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetServiceInstanceQuery(lookup ServiceInstanceLookup) (string, []interface{}, error) {
	if err := validateServiceInstanceLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + serviceInstanceReadColumns + `
FROM service_instances svc
JOIN orders ord ON ord.order_id = svc.order_id
WHERE svc.service_instance_id = $1
  AND svc.tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.BuyerUserID != "" {
		args = append(args, lookup.BuyerUserID)
		query += fmt.Sprintf("\n  AND ord.buyer_user_id = $%d", len(args))
	}
	return query, args, nil
}
