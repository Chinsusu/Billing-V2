package jobs

import (
	"context"
	"fmt"
)

func (store *PostgresStore) ListJobs(ctx context.Context, filter Filter) ([]Job, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListJobsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()
	return scanJobs(rows)
}

func (store *PostgresStore) GetJob(ctx context.Context, lookup Lookup) (Job, error) {
	if err := store.ready(); err != nil {
		return Job{}, err
	}
	query, args, err := buildGetJobQuery(lookup)
	if err != nil {
		return Job{}, err
	}
	return scanJob(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListJobsQuery(filter Filter) (string, []interface{}, error) {
	filter = normalizeFilter(filter)
	if err := validateFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + jobColumns + `
FROM jobs
WHERE tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND display_id = $%d", len(args))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf("\n  AND job_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND status = $%d", len(args))
	}
	if filter.ReferenceType != "" {
		args = append(args, filter.ReferenceType)
		query += fmt.Sprintf("\n  AND reference_type = $%d", len(args))
	}
	if filter.ReferenceID != "" {
		args = append(args, filter.ReferenceID)
		query += fmt.Sprintf("\n  AND reference_id = $%d", len(args))
	}
	if filter.SourceID != "" {
		args = append(args, filter.SourceID)
		query += fmt.Sprintf("\n  AND source_id = $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetJobQuery(lookup Lookup) (string, []interface{}, error) {
	lookup = normalizeLookup(lookup)
	if err := validateLookup(lookup); err != nil {
		return "", nil, err
	}
	return `SELECT ` + jobColumns + `
FROM jobs
WHERE job_id = $1
  AND tenant_id = $2`, []interface{}{lookup.ID, lookup.TenantID}, nil
}
