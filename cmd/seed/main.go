package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/seed"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("seed", flag.ContinueOnError)
	dsn := flags.String("dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	timeout := flags.Duration("timeout", 30*time.Second, "seed command timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}

	command := "dev"
	if flags.NArg() > 0 {
		command = flags.Arg(0)
	}

	switch command {
	case "plan":
		return printPlan(seed.DevStatements())
	case "dev":
		return applyDevSeeds(*dsn, *timeout)
	default:
		return fmt.Errorf("unknown command %q; use dev or plan", command)
	}
}

func printPlan(statements []seed.Statement) error {
	fmt.Printf("available seed statement(s): %d\n", len(statements))
	for _, statement := range statements {
		fmt.Printf("- %s\n", statement.Name)
	}
	return nil
}

func applyDevSeeds(dsn string, timeout time.Duration) error {
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := db.Open(ctx, db.Config{DriverName: db.DefaultDriverName, DSN: dsn})
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := seed.ApplyDev(ctx, conn); err != nil {
		return err
	}
	fmt.Printf("applied seed statement(s): %d\n", len(seed.DevStatements()))
	return nil
}
