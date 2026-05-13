package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Environment string

const (
	EnvironmentLocal      Environment = "local"
	EnvironmentDev        Environment = "dev"
	EnvironmentStaging    Environment = "staging"
	EnvironmentProduction Environment = "production"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Config struct {
	AppEnvironment      Environment
	AppName             string
	HTTPAddr            string
	LogLevel            LogLevel
	DatabaseDSN         string
	EncryptionKey       string
	SessionCookieName   string
	SessionCookieSecure bool
	SessionTokenTTL     time.Duration
	PasswordResetTTL    time.Duration
}

func LoadFromEnv() (Config, error) {
	environment := Environment(getenv("APP_ENV", string(EnvironmentLocal)))
	sessionCookieSecure, err := getenvBool("AUTH_SESSION_COOKIE_SECURE", defaultSessionCookieSecure(environment))
	if err != nil {
		return Config{}, fmt.Errorf("AUTH_SESSION_COOKIE_SECURE is invalid: %w", err)
	}
	cfg := Config{
		AppEnvironment:      environment,
		AppName:             getenv("APP_NAME", "billing-v2"),
		HTTPAddr:            getenv("APP_HTTP_ADDR", ":8080"),
		LogLevel:            LogLevel(getenv("LOG_LEVEL", string(LogLevelInfo))),
		DatabaseDSN:         os.Getenv("DB_DSN"),
		EncryptionKey:       os.Getenv("ENCRYPTION_KEY"),
		SessionCookieName:   getenv("AUTH_SESSION_COOKIE_NAME", "billing_session"),
		SessionCookieSecure: sessionCookieSecure,
	}
	sessionTTL, err := time.ParseDuration(getenv("AUTH_SESSION_TTL", "12h"))
	if err != nil {
		return Config{}, fmt.Errorf("AUTH_SESSION_TTL is invalid: %w", err)
	}
	cfg.SessionTokenTTL = sessionTTL
	passwordResetTTL, err := time.ParseDuration(getenv("AUTH_PASSWORD_RESET_TTL", "30m"))
	if err != nil {
		return Config{}, fmt.Errorf("AUTH_PASSWORD_RESET_TTL is invalid: %w", err)
	}
	cfg.PasswordResetTTL = passwordResetTTL
	return cfg, cfg.Validate()
}

func (cfg Config) Validate() error {
	if !validEnvironment(cfg.AppEnvironment) {
		return fmt.Errorf("APP_ENV must be one of local, dev, staging, production")
	}
	if cfg.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if cfg.HTTPAddr == "" {
		return fmt.Errorf("APP_HTTP_ADDR is required")
	}
	if err := validateHTTPAddr(cfg.HTTPAddr); err != nil {
		return fmt.Errorf("APP_HTTP_ADDR is invalid: %w", err)
	}
	if !validLogLevel(cfg.LogLevel) {
		return fmt.Errorf("LOG_LEVEL must be one of debug, info, warn, error")
	}
	if cfg.SessionCookieName == "" {
		return fmt.Errorf("AUTH_SESSION_COOKIE_NAME is required")
	}
	if cfg.SessionTokenTTL <= 0 {
		return fmt.Errorf("AUTH_SESSION_TTL must be positive")
	}
	if cfg.PasswordResetTTL <= 0 {
		return fmt.Errorf("AUTH_PASSWORD_RESET_TTL must be positive")
	}
	if cfg.AppEnvironment == EnvironmentProduction && !cfg.SessionCookieSecure {
		return fmt.Errorf("AUTH_SESSION_COOKIE_SECURE must be true in production")
	}
	if cfg.DatabaseDSN != "" && requiresProductionSecrets(cfg.AppEnvironment) && cfg.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required in staging and production")
	}
	return nil
}

func (cfg Config) AllowDevActorHeaders() bool {
	return cfg.AppEnvironment == EnvironmentLocal || cfg.AppEnvironment == EnvironmentDev
}

func validEnvironment(value Environment) bool {
	switch value {
	case EnvironmentLocal, EnvironmentDev, EnvironmentStaging, EnvironmentProduction:
		return true
	default:
		return false
	}
}

func validLogLevel(value LogLevel) bool {
	switch value {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return true
	default:
		return false
	}
}

func validateHTTPAddr(value string) error {
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}
	if port == "" {
		return fmt.Errorf("port is required")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port must be numeric")
	}
	if portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("port is out of range")
	}
	return nil
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return parsed, nil
}

func defaultSessionCookieSecure(environment Environment) bool {
	return environment == EnvironmentStaging || environment == EnvironmentProduction
}

func requiresProductionSecrets(environment Environment) bool {
	return environment == EnvironmentStaging || environment == EnvironmentProduction
}
