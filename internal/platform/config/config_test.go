package config

import "testing"

func TestLoadFromEnvUsesSafeDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_HTTP_ADDR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DB_DSN", "")
	t.Setenv("AUTH_SESSION_COOKIE_NAME", "")
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "")
	t.Setenv("AUTH_SESSION_TTL", "")

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
	if cfg.SessionCookieName != "billing_session" {
		t.Fatalf("expected default session cookie name, got %q", cfg.SessionCookieName)
	}
	if cfg.SessionCookieSecure {
		t.Fatal("expected local session cookie to default to insecure transport")
	}
	if cfg.SessionTokenTTL.String() != "12h0m0s" {
		t.Fatalf("expected default session TTL, got %s", cfg.SessionTokenTTL)
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

func TestValidateRejectsInvalidSessionTTL(t *testing.T) {
	cfg := Config{
		AppEnvironment:    EnvironmentLocal,
		AppName:           "billing-v2",
		HTTPAddr:          ":8080",
		LogLevel:          LogLevelInfo,
		SessionCookieName: "billing_session",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid session TTL error")
	}
}

func TestValidateRequiresSecureSessionCookieInProduction(t *testing.T) {
	cfg := Config{
		AppEnvironment:    EnvironmentProduction,
		AppName:           "billing-v2",
		HTTPAddr:          ":8080",
		LogLevel:          LogLevelInfo,
		SessionCookieName: "billing_session",
		SessionTokenTTL:   12,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected production cookie security error")
	}
}

func TestLoadFromEnvRejectsInvalidSessionCookieSecure(t *testing.T) {
	t.Setenv("AUTH_SESSION_COOKIE_SECURE", "maybe")

	if _, err := LoadFromEnv(); err == nil {
		t.Fatal("expected invalid session cookie secure error")
	}
}
