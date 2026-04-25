package audit

import (
	"context"
	"fmt"
)

func (store *PostgresStore) ListLogs(ctx context.Context, filter Filter) ([]Log, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListLogsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()
	logs := make([]Log, 0)
	for rows.Next() {
		record, err := scanLogRead(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read audit logs: %w", err)
	}
	return logs, nil
}

func (store *PostgresStore) GetLog(ctx context.Context, lookup Lookup) (Log, error) {
	if err := store.ready(); err != nil {
		return Log{}, err
	}
	query, args, err := buildGetLogQuery(lookup)
	if err != nil {
		return Log{}, err
	}
	return scanLogRead(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListLogsQuery(filter Filter) (string, []interface{}, error) {
	filter = normalizeFilter(filter)
	if err := validateFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + auditReadColumns + `
FROM audit_logs
WHERE tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if !filter.ActorID.Empty() {
		args = append(args, filter.ActorID)
		query += fmt.Sprintf("\n  AND actor_id = $%d", len(args))
	}
	if filter.ActorType != "" {
		args = append(args, filter.ActorType)
		query += fmt.Sprintf("\n  AND actor_type = $%d", len(args))
	}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND display_id = $%d", len(args))
	}
	if filter.Action != "" {
		args = append(args, filter.Action)
		query += fmt.Sprintf("\n  AND action = $%d", len(args))
	}
	if filter.TargetType != "" {
		args = append(args, filter.TargetType)
		query += fmt.Sprintf("\n  AND target_type = $%d", len(args))
	}
	if !filter.TargetID.Empty() {
		args = append(args, filter.TargetID)
		query += fmt.Sprintf("\n  AND target_id = $%d", len(args))
	}
	if !filter.CreatedFrom.IsZero() {
		args = append(args, filter.CreatedFrom)
		query += fmt.Sprintf("\n  AND created_at >= $%d", len(args))
	}
	if !filter.CreatedTo.IsZero() {
		args = append(args, filter.CreatedTo)
		query += fmt.Sprintf("\n  AND created_at <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetLogQuery(lookup Lookup) (string, []interface{}, error) {
	if err := validateLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + auditReadColumns + `
FROM audit_logs
WHERE audit_id = $1
  AND tenant_id = $2`
	return query, []interface{}{lookup.ID, lookup.TenantID}, nil
}
