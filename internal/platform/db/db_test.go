package db

import "testing"

func TestConfigValidateRequiresDriverName(t *testing.T) {
	cfg := Config{DSN: "postgres://billing:billing@localhost:5432/billing?sslmode=disable"}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected driver name error")
	}
}

func TestConfigValidateRequiresDSN(t *testing.T) {
	cfg := Config{DriverName: DefaultDriverName}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected DSN error")
	}
}

func TestConfigValidateAcceptsRequiredFields(t *testing.T) {
	cfg := Config{
		DriverName: DefaultDriverName,
		DSN:        "postgres://billing:billing@localhost:5432/billing?sslmode=disable",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}
