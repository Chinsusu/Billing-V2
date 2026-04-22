package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
	AppEnvironment Environment
	AppName        string
	HTTPAddr       string
	LogLevel       LogLevel
	DatabaseDSN    string
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		AppEnvironment: Environment(getenv("APP_ENV", string(EnvironmentLocal))),
		AppName:        getenv("APP_NAME", "billing-v2"),
		HTTPAddr:       getenv("APP_HTTP_ADDR", ":8080"),
		LogLevel:       LogLevel(getenv("LOG_LEVEL", string(LogLevelInfo))),
		DatabaseDSN:    os.Getenv("DB_DSN"),
	}
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
	return nil
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
