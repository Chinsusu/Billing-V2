package config

import "testing"

func TestLoadFromEnvUsesSafeDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_HTTP_ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DB_DSN", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	if cfg.AppEnvironment != EnvironmentLocal {
		t.Fatalf("expected local environment, got %q", cfg.AppEnvironment)
	}
	if cfg.AppName != "billing-v2" {
		t.Fatalf("expected default app name, got %q", cfg.AppName)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default HTTP addr, got %q", cfg.HTTPAddr)
	}
	if cfg.LogLevel != LogLevelInfo {
		t.Fatalf("expected info log level, got %q", cfg.LogLevel)
	}
}

func TestValidateRejectsInvalidEnvironment(t *testing.T) {
	cfg := Config{
		AppEnvironment: "prod",
		AppName:        "billing-v2",
		HTTPAddr:       ":8080",
		LogLevel:       LogLevelInfo,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid environment error")
	}
}

func TestValidateRejectsInvalidHTTPAddr(t *testing.T) {
	cfg := Config{
		AppEnvironment: EnvironmentLocal,
		AppName:        "billing-v2",
		HTTPAddr:       "8080",
		LogLevel:       LogLevelInfo,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid HTTP address error")
	}
}

func TestValidateRejectsInvalidLogLevel(t *testing.T) {
	cfg := Config{
		AppEnvironment: EnvironmentLocal,
		AppName:        "billing-v2",
		HTTPAddr:       ":8080",
		LogLevel:       "verbose",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid log level error")
	}
}
