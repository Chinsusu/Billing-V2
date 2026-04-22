package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const DefaultDriverName = "postgres"

type Config struct {
	DriverName      string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (cfg Config) Validate() error {
	if cfg.DriverName == "" {
		return fmt.Errorf("database driver name is required")
	}
	if cfg.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}
	return nil
}

func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	if cfg.DriverName == "" {
		cfg.DriverName = DefaultDriverName
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	conn, err := sql.Open(cfg.DriverName, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if cfg.MaxOpenConns > 0 {
		conn.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		conn.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		conn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return conn, nil
}
