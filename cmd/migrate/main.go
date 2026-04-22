package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "migrate failed: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	dir := flags.String("dir", "migrations", "migration directory")
	dsn := flags.String("dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	timeout := flags.Duration("timeout", 30*time.Second, "migration command timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}

	command := "validate"
	if flags.NArg() > 0 {
		command = flags.Arg(0)
	}

	migrations, err := loadMigrationsFromDir(*dir)
	if err != nil {
		return err
	}

	switch command {
	case "validate":
		fmt.Printf("validated %d migration(s)\n", len(migrations))
		return nil
	case "plan":
		return printPlan(*dsn, *timeout, migrations)
	case "up":
		return applyMigrations(*dsn, *timeout, migrations)
	default:
		return fmt.Errorf("unknown command %q; use validate, plan, or up", command)
	}
}

func loadMigrationsFromDir(dir string) ([]db.Migration, error) {
	fsys := os.DirFS(dir)
	migrations, err := db.LoadMigrations(fsys)
	if err != nil {
		return nil, err
	}
	return migrations, nil
}

func printPlan(dsn string, timeout time.Duration, migrations []db.Migration) error {
	if dsn == "" {
		fmt.Printf("available migration(s): %d\n", len(migrations))
		for _, migration := range migrations {
			fmt.Printf("- %s %s\n", migration.Version, migration.Name)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := db.Open(ctx, db.Config{DriverName: db.DefaultDriverName, DSN: dsn})
	if err != nil {
		return err
	}
	defer conn.Close()

	migrator, err := db.NewMigrator(conn, migrations)
	if err != nil {
		return err
	}
	pending, err := migrator.Pending(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("pending migration(s): %d\n", len(pending))
	for _, migration := range pending {
		fmt.Printf("- %s %s\n", migration.Version, migration.Name)
	}
	return nil
}

func applyMigrations(dsn string, timeout time.Duration, migrations []db.Migration) error {
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for up")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := db.Open(ctx, db.Config{DriverName: db.DefaultDriverName, DSN: dsn})
	if err != nil {
		return err
	}
	defer conn.Close()

	migrator, err := db.NewMigrator(conn, migrations)
	if err != nil {
		return err
	}
	applied, err := migrator.ApplyAll(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("applied migration(s): %d\n", len(applied))
	for _, migration := range applied {
		fmt.Printf("- %s %s\n", migration.Version, migration.Name)
	}
	return nil
}
