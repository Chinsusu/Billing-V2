package db

import "testing"

func TestPendingMigrationsReturnsUnappliedMigrations(t *testing.T) {
	all := []Migration{
		{Version: "0001", Name: "create_tenants", Checksum: "a"},
		{Version: "0002", Name: "create_users", Checksum: "b"},
	}
	applied := map[string]AppliedMigration{
		"0001": {Version: "0001", Name: "create_tenants", Checksum: "a"},
	}

	pending, err := PendingMigrations(all, applied)
	if err != nil {
		t.Fatalf("PendingMigrations returned error: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending migration, got %d", len(pending))
	}
	if pending[0].Version != "0002" {
		t.Fatalf("expected version 0002, got %q", pending[0].Version)
	}
}

func TestPendingMigrationsRejectsChecksumChange(t *testing.T) {
	all := []Migration{
		{Version: "0001", Name: "create_tenants", Checksum: "new"},
	}
	applied := map[string]AppliedMigration{
		"0001": {Version: "0001", Name: "create_tenants", Checksum: "old"},
	}

	if _, err := PendingMigrations(all, applied); err == nil {
		t.Fatal("expected checksum change error")
	}
}
