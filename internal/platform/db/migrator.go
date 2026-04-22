package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const schemaMigrationsTable = "schema_migrations"

type AppliedMigration struct {
	Version   string
	Name      string
	Checksum  string
	AppliedAt time.Time
}

type Migrator struct {
	conn       *sql.DB
	migrations []Migration
}

func NewMigrator(conn *sql.DB, migrations []Migration) (*Migrator, error) {
	if conn == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	return &Migrator{conn: conn, migrations: migrations}, nil
}

func (migrator *Migrator) Pending(ctx context.Context) ([]Migration, error) {
	if err := migrator.ensureSchemaMigrations(ctx); err != nil {
		return nil, err
	}
	applied, err := migrator.appliedVersions(ctx)
	if err != nil {
		return nil, err
	}
	return PendingMigrations(migrator.migrations, applied)
}

func (migrator *Migrator) ApplyAll(ctx context.Context) ([]Migration, error) {
	pending, err := migrator.Pending(ctx)
	if err != nil {
		return nil, err
	}

	applied := make([]Migration, 0, len(pending))
	for _, migration := range pending {
		if err := migrator.apply(ctx, migration); err != nil {
			return applied, err
		}
		applied = append(applied, migration)
	}
	return applied, nil
}

func PendingMigrations(all []Migration, applied map[string]AppliedMigration) ([]Migration, error) {
	pending := make([]Migration, 0, len(all))
	for _, migration := range all {
		appliedMigration, exists := applied[migration.Version]
		if !exists {
			pending = append(pending, migration)
			continue
		}
		if appliedMigration.Checksum != migration.Checksum {
			return nil, fmt.Errorf("migration %s checksum changed after apply", migration.Version)
		}
	}
	return pending, nil
}

func (migrator *Migrator) ensureSchemaMigrations(ctx context.Context) error {
	_, err := migrator.conn.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version text PRIMARY KEY,
	name text NOT NULL,
	checksum text NOT NULL,
	applied_at timestamptz NOT NULL DEFAULT now()
)`)
	if err != nil {
		return fmt.Errorf("ensure %s table: %w", schemaMigrationsTable, err)
	}
	return nil
}

func (migrator *Migrator) appliedVersions(ctx context.Context) (map[string]AppliedMigration, error) {
	rows, err := migrator.conn.QueryContext(ctx, `
SELECT version, name, checksum, applied_at
FROM schema_migrations
ORDER BY version`)
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]AppliedMigration)
	for rows.Next() {
		var migration AppliedMigration
		if err := rows.Scan(&migration.Version, &migration.Name, &migration.Checksum, &migration.AppliedAt); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied[migration.Version] = migration
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	return applied, nil
}

func (migrator *Migrator) apply(ctx context.Context, migration Migration) error {
	tx, err := migrator.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", migration.Version, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("apply migration %s: %w", migration.Version, err)
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO schema_migrations (version, name, checksum)
VALUES ($1, $2, $3)`, migration.Version, migration.Name, migration.Checksum); err != nil {
		return fmt.Errorf("record migration %s: %w", migration.Version, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", migration.Version, err)
	}
	return nil
}
