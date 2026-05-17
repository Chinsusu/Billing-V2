package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
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

func TestRunLifecycleScheduleOnceUsesGracePeriodAndPrintsSummary(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var output bytes.Buffer
	schedulerFactory := &fakeLifecycleSchedulerFactory{
		scheduler: fakeLifecycleScheduler{
			summary: order.ServiceLifecycleScheduleSummary{Due: 3, Scheduled: 2},
		},
	}

	err := runWithDependencies([]string{
		"lifecycle-schedule-once",
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"-batch-size", "7",
		"-grace-period", "48h",
	}, workerDependencies{stdout: &output, newLifecycleScheduler: schedulerFactory.newScheduler})
	if err != nil {
		t.Fatalf("expected lifecycle schedule success: %v", err)
	}
	if schedulerFactory.cfg.BatchSize != 7 || schedulerFactory.cfg.GracePeriod != 48*time.Hour {
		t.Fatalf("unexpected lifecycle schedule config: %+v", schedulerFactory.cfg)
	}
	if schedulerFactory.scheduler.input.Limit != 7 || schedulerFactory.scheduler.input.GracePeriod != 48*time.Hour {
		t.Fatalf("unexpected lifecycle schedule input: %+v", schedulerFactory.scheduler.input)
	}
	if !strings.Contains(output.String(), "lifecycle-schedule-once due=3 scheduled=2") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestRunLifecycleOnceUsesLifecycleRunner(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var output bytes.Buffer
	factory := &fakeRunnerFactory{
		runner: fakeProvisionRunner{summary: jobs.RunSummary{
			Claimed:      1,
			ManualReview: 1,
		}},
	}

	err := runWithDependencies([]string{
		"lifecycle-once",
		"-dsn", "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
		"-worker-id", "lifecycle-test",
	}, workerDependencies{stdout: &output, newLifecycleRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected lifecycle once success: %v", err)
	}
	if !factory.called || factory.cfg.WorkerID != "lifecycle-test" {
		t.Fatalf("expected lifecycle runner factory call, got called=%v cfg=%+v", factory.called, factory.cfg)
	}
	if !strings.Contains(output.String(), "lifecycle-once claimed=1 succeeded=0 retried=0 manual_review=1") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestRunProviderRegistryCheckDoesNotRequireDSNAndRedactsConfig(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("DB_DSN", "")
	setCloudminiWorkerProviderEnv(t, validCloudminiWorkerProviderEnv())
	var output bytes.Buffer

	err := runWithDependencies([]string{commandProviderRegistryCheck}, workerDependencies{stdout: &output})
	if err != nil {
		t.Fatalf("expected provider registry check success: %v", err)
	}
	got := output.String()
	for _, want := range []string{
		"provider-registry-check result=PASS",
		"mode=cloudmini_v3",
		"cloudmini_v3_adapter=real",
		"cloudmini_source_mappings=1",
		"cloudmini_account_mappings=0",
		"provider_api_called=no",
		"mutating_routes_called=no",
		"jobs_claimed=0",
		"secrets_printed=no",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected output to contain %q, got %s", want, got)
		}
	}
	for _, forbidden := range []string{
		"sandbox-token",
		"source-cloudmini-1",
		"group-cloudmini-1",
		"http://cloudmini.example",
	} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("provider registry check output leaked %q: %s", forbidden, got)
		}
	}
}

func TestRunProviderRegistryCheckRejectsProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	setCloudminiWorkerProviderEnv(t, validCloudminiWorkerProviderEnv())

	err := runWithDependencies([]string{commandProviderRegistryCheck}, workerDependencies{})
	if err == nil || !strings.Contains(err.Error(), "refusing to run worker") {
		t.Fatalf("expected production guard error, got %v", err)
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := runWithDependencies([]string{"unknown"}, workerDependencies{newRunner: (&fakeRunnerFactory{}).newRunner})
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func setCloudminiWorkerProviderEnv(t *testing.T, env workerProviderEnv) {
	t.Helper()
	t.Setenv("PROVIDER_DEFAULT_MODE", env.Mode)
	t.Setenv("ENCRYPTION_KEY", env.EncryptionKey)
	t.Setenv("CLOUDMINI_V3_BASE_URL", env.CloudminiV3BaseURL)
	t.Setenv("CLOUDMINI_V3_API_TOKEN", env.CloudminiV3APIToken)
	t.Setenv("CLOUDMINI_V3_SOURCE_ID", env.CloudminiV3SourceID)
	t.Setenv("CLOUDMINI_V3_KIND", env.CloudminiV3Kind)
	t.Setenv("CLOUDMINI_V3_GROUP_ID", env.CloudminiV3GroupID)
	t.Setenv("CLOUDMINI_V3_PROTOCOL", env.CloudminiV3Protocol)
	t.Setenv("CLOUDMINI_V3_BANDWIDTH_LIMIT_MB", env.CloudminiV3BandwidthMB)
	t.Setenv("CLOUDMINI_V3_SPEED_LIMIT_MBPS", env.CloudminiV3SpeedMBps)
	t.Setenv("CLOUDMINI_V3_POLL_INTERVAL", env.CloudminiV3PollInterval)
	t.Setenv("CLOUDMINI_V3_POLL_TIMEOUT", env.CloudminiV3PollTimeout)
}

type fakeLifecycleSchedulerFactory struct {
	called    bool
	cfg       workerConfig
	scheduler fakeLifecycleScheduler
	err       error
}

func (factory *fakeLifecycleSchedulerFactory) newScheduler(ctx context.Context, cfg workerConfig) (lifecycleScheduler, func() error, error) {
	factory.called = true
	factory.cfg = cfg
	if factory.err != nil {
		return nil, nil, factory.err
	}
	return &factory.scheduler, func() error { return nil }, nil
}

type fakeLifecycleScheduler struct {
	input   order.ListDueServiceLifecycleActionsInput
	summary order.ServiceLifecycleScheduleSummary
	err     error
}

func (scheduler *fakeLifecycleScheduler) ScheduleDue(ctx context.Context, input order.ListDueServiceLifecycleActionsInput) (order.ServiceLifecycleScheduleSummary, error) {
	scheduler.input = input
	if scheduler.err != nil {
		return order.ServiceLifecycleScheduleSummary{}, scheduler.err
	}
	return scheduler.summary, nil
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
