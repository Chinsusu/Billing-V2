package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestRunNotificationTelegramOnceUsesConfigAndPrintsSummary(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	setTelegramWorkerEnv(t, "test-token", "test-chat", "")
	var output bytes.Buffer
	factory := &fakeRunnerFactory{
		runner: fakeProvisionRunner{summary: jobs.RunSummary{
			Claimed:   1,
			Succeeded: 1,
		}},
	}

	err := runWithDependencies([]string{
		"notification-telegram-once",
		"-dsn", "postgres://localhost:5432/billing?sslmode=disable",
		"-worker-id", "telegram-test",
		"-batch-size", "2",
	}, workerDependencies{stdout: &output, newTelegramRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected telegram once success: %v", err)
	}
	if !factory.called || factory.cfg.WorkerID != "telegram-test" || factory.cfg.BatchSize != 2 {
		t.Fatalf("unexpected telegram runner config: called=%v cfg=%+v", factory.called, factory.cfg)
	}
	if !strings.Contains(output.String(), "notification-telegram-once claimed=1 succeeded=1") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestRunNotificationTelegramLoopPrintsPassSummary(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	setTelegramWorkerEnv(t, "test-token", "test-chat", "")
	var output bytes.Buffer
	calls := 0
	factory := &fakeRunnerFactory{
		runner: fakeProvisionRunner{
			summary: jobs.RunSummary{Claimed: 0},
			calls:   &calls,
		},
	}

	err := runWithDependencies([]string{
		"notification-telegram-loop",
		"-dsn", "postgres://localhost:5432/billing?sslmode=disable",
		"-worker-id", "telegram-loop-test",
		"-timeout", "25ms",
		"-interval", "100ms",
	}, workerDependencies{stdout: &output, newTelegramRunner: factory.newRunner})
	if err != nil {
		t.Fatalf("expected telegram loop success: %v", err)
	}
	if !factory.called || factory.cfg.WorkerID != "telegram-loop-test" {
		t.Fatalf("expected telegram runner factory call, got called=%v cfg=%+v", factory.called, factory.cfg)
	}
	if !strings.Contains(output.String(), "notification-telegram-loop pass=1 claimed=0") {
		t.Fatalf("unexpected output: %s", output.String())
	}
	if calls != 1 {
		t.Fatalf("expected one idle pass before timeout, got %d", calls)
	}
}

func TestRunNotificationTelegramRejectsProductionWithoutApproval(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	setTelegramWorkerEnv(t, "test-token", "test-chat", "")
	factory := &fakeRunnerFactory{runner: fakeProvisionRunner{}}

	err := runWithDependencies([]string{
		"notification-telegram-once",
		"-dsn", "postgres://localhost:5432/billing?sslmode=disable",
	}, workerDependencies{newTelegramRunner: factory.newRunner})
	if err == nil || !strings.Contains(err.Error(), "BILLING_TELEGRAM_DELIVERY_PRODUCTION_APPROVED=yes") {
		t.Fatalf("expected production approval guard error, got %v", err)
	}
	if factory.called {
		t.Fatal("telegram runner factory should not be called after production guard")
	}
}

func TestRunNotificationTelegramRequiresConfig(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	t.Setenv("TELEGRAM_CHAT_ID", "test-chat")

	err := runWithDependencies([]string{
		"notification-telegram-once",
		"-dsn", "postgres://localhost:5432/billing?sslmode=disable",
	}, workerDependencies{newTelegramRunner: (&fakeRunnerFactory{}).newRunner})
	if err == nil || !strings.Contains(err.Error(), "TELEGRAM_BOT_TOKEN is required") {
		t.Fatalf("expected missing telegram token error, got %v", err)
	}
}

func TestRunNotificationTelegramPreflightSendsSafeMessage(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	var requestBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requestBody = string(body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	setTelegramWorkerEnv(t, "test-token", "test-chat", server.URL)
	var output bytes.Buffer

	err := runWithDependencies([]string{
		"notification-telegram-preflight",
		"-timeout", "2s",
	}, workerDependencies{stdout: &output})
	if err != nil {
		t.Fatalf("expected telegram preflight success: %v", err)
	}
	got := output.String()
	for _, want := range []string{
		"telegram-preflight result=PASS",
		"telegram_api_called=yes",
		"message_payload_redacted=yes",
		"secrets_printed=no",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected output to contain %q, got %s", want, got)
		}
	}
	for _, forbidden := range []string{"test-token", "test-chat"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("preflight output leaked %q: %s", forbidden, got)
		}
	}
	if !strings.Contains(requestBody, "Billing Telegram delivery preflight") ||
		strings.Contains(requestBody, "test-token") {
		t.Fatalf("unexpected telegram preflight request body: %s", requestBody)
	}
}

func setTelegramWorkerEnv(t *testing.T, token string, chatID string, baseURL string) {
	t.Helper()
	t.Setenv("TELEGRAM_BOT_TOKEN", token)
	t.Setenv("TELEGRAM_CHAT_ID", chatID)
	t.Setenv("TELEGRAM_API_BASE_URL", baseURL)
}
