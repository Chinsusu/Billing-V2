package main

import (
	"os"
	"testing"
)

func TestSmokeCheckValidateExact(t *testing.T) {
	check := exactCheck("demo", "SELECT 1", 1)

	if err := check.validate(1); err != nil {
		t.Fatalf("expected exact count to pass: %v", err)
	}
	if err := check.validate(2); err == nil {
		t.Fatal("expected exact count mismatch to fail")
	}
}

func TestSmokeCheckValidateMinimum(t *testing.T) {
	check := minCheck("demo", "SELECT 1", 3)

	if err := check.validate(4); err != nil {
		t.Fatalf("expected minimum count to pass: %v", err)
	}
	if err := check.validate(2); err == nil {
		t.Fatal("expected count below minimum to fail")
	}
}

func TestGuardDevEnvironmentRejectsProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	if err := guardDevEnvironment(); err == nil {
		t.Fatal("expected production environment to be rejected")
	}
}

func TestGuardDevEnvironmentAllowsUnset(t *testing.T) {
	os.Unsetenv("APP_ENV")

	if err := guardDevEnvironment(); err != nil {
		t.Fatalf("expected unset APP_ENV to pass: %v", err)
	}
}
