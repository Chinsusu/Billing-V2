package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/platform/config"
)

func TestLoggerWritesStructuredJSON(t *testing.T) {
	var output bytes.Buffer
	log := New(&output, config.LogLevelDebug)

	log.Info("api started", String("module", "api"), String("request_id", "req_123"))

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v", err)
	}
	if entry["level"] != "info" {
		t.Fatalf("expected info level, got %v", entry["level"])
	}
	if entry["message"] != "api started" {
		t.Fatalf("expected message, got %v", entry["message"])
	}
	if entry["request_id"] != "req_123" {
		t.Fatalf("expected request id, got %v", entry["request_id"])
	}
}

func TestLoggerRedactsSecretFields(t *testing.T) {
	var output bytes.Buffer
	log := New(&output, config.LogLevelDebug)

	log.Info("provider configured", String("provider_api_key", "secret-value"))

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v", err)
	}
	if entry["provider_api_key"] != "[REDACTED]" {
		t.Fatalf("expected redacted value, got %v", entry["provider_api_key"])
	}
}

func TestLoggerFiltersBelowMinimumLevel(t *testing.T) {
	var output bytes.Buffer
	log := New(&output, config.LogLevelWarn)

	log.Info("hidden message")
	if output.Len() != 0 {
		t.Fatalf("expected no output, got %q", output.String())
	}
}
