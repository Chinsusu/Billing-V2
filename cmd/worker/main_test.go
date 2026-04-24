package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestRunProvisionOnceUsesConfigAndPrintsSummary(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var output bytes.Buffer
	factory := &fakeRunnerFactory{
		runner: fakeProvisionRunner{summary: jobs.RunSummary{
			Claimed:   2,
			Succeeded: 1,
			Retried:   1,
		}},
	}

	err := runWithDependencies([]string{
		"provision-once",
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"-worker-id", "worker-test",
		"-batch-size", "2",
	}, workerDependencies{stdout: &output, newRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected provision once success: %v", err)
	}
	if factory.cfg.WorkerID != "worker-test" || factory.cfg.BatchSize != 2 {
		t.Fatalf("unexpected worker config: %+v", factory.cfg)
	}
	if !strings.Contains(output.String(), "claimed=2 succeeded=1 retried=1") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestRunProvisionOnceAcceptsCommandAfterFlags(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	factory := &fakeRunnerFactory{runner: fakeProvisionRunner{}}

	err := runWithDependencies([]string{
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"provision-once",
	}, workerDependencies{newRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected provision once success: %v", err)
	}
	if !factory.called {
		t.Fatal("expected runner factory call")
	}
}

func TestRunProvisionOnceRejectsProductionEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	factory := &fakeRunnerFactory{runner: fakeProvisionRunner{}}

	err := runWithDependencies([]string{
		"provision-once",
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
	}, workerDependencies{newRunner: factory.newRunner})
	if err == nil || !strings.Contains(err.Error(), "refusing to run worker") {
		t.Fatalf("expected production guard error, got %v", err)
	}
	if factory.called {
		t.Fatal("runner factory should not be called after production guard")
	}
}

func TestRunProvisionOnceRequiresDSN(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("DB_DSN", "")

	err := runWithDependencies([]string{"provision-once"}, workerDependencies{newRunner: (&fakeRunnerFactory{}).newRunner})
	if err == nil || !strings.Contains(err.Error(), "DB_DSN or -dsn is required") {
		t.Fatalf("expected dsn error, got %v", err)
	}
}

func TestRunProvisionLoopAcceptsCommandAfterFlagsAndPrintsPassSummary(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var output bytes.Buffer
	calls := 0
	factory := &fakeRunnerFactory{
		runner: fakeProvisionRunner{
			summary: jobs.RunSummary{Claimed: 0},
			calls:   &calls,
		},
	}

	err := runWithDependencies([]string{
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"-worker-id", "loop-test",
		"-timeout", "25ms",
		"-interval", "100ms",
		"provision-loop",
	}, workerDependencies{stdout: &output, newRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected provision loop success: %v", err)
	}
	if !factory.called {
		t.Fatal("expected runner factory call")
	}
	if factory.cfg.WorkerID != "loop-test" || factory.cfg.Interval != 100*time.Millisecond {
		t.Fatalf("unexpected worker config: %+v", factory.cfg)
	}
	if !strings.Contains(output.String(), "provision-loop pass=1 claimed=0") {
		t.Fatalf("unexpected output: %s", output.String())
	}
	if calls != 1 {
		t.Fatalf("expected one idle pass before timeout, got %d", calls)
	}
}

func TestRunProvisionLoopRejectsInvalidInterval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	factory := &fakeRunnerFactory{runner: fakeProvisionRunner{}}

	err := runWithDependencies([]string{
		"provision-loop",
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"-interval", "0",
	}, workerDependencies{newRunner: factory.newRunner})
	if err == nil || !strings.Contains(err.Error(), "worker interval must be positive") {
		t.Fatalf("expected interval error, got %v", err)
	}
	if factory.called {
		t.Fatal("runner factory should not be called with invalid interval")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := runWithDependencies([]string{"unknown"}, workerDependencies{newRunner: (&fakeRunnerFactory{}).newRunner})
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

type fakeRunnerFactory struct {
	called bool
	cfg    workerConfig
	runner fakeProvisionRunner
	err    error
}

func (factory *fakeRunnerFactory) newRunner(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error) {
	factory.called = true
	factory.cfg = cfg
	if factory.err != nil {
		return nil, nil, factory.err
	}
	return factory.runner, func() error { return nil }, nil
}

type fakeProvisionRunner struct {
	summary jobs.RunSummary
	err     error
	calls   *int
}

func (runner fakeProvisionRunner) RunOnce(ctx context.Context) (jobs.RunSummary, error) {
	if runner.calls != nil {
		*runner.calls++
	}
	if runner.err != nil {
		return jobs.RunSummary{}, runner.err
	}
	return runner.summary, nil
}
