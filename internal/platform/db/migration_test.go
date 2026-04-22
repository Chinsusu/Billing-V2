package db

import (
	"testing"
	"testing/fstest"
)

func TestLoadMigrationsSortsByVersion(t *testing.T) {
	migrations, err := LoadMigrations(fstest.MapFS{
		"0002_create_users.sql":   {Data: []byte("create table users (id text);")},
		"0001_create_tenants.sql": {Data: []byte("create table tenants (id text);")},
		"README.md":               {Data: []byte("notes")},
	})
	if err != nil {
		t.Fatalf("LoadMigrations returned error: %v", err)
	}
	if len(migrations) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migrations))
	}
	if migrations[0].Version != "0001" {
		t.Fatalf("expected first version 0001, got %q", migrations[0].Version)
	}
	if migrations[1].Version != "0002" {
		t.Fatalf("expected second version 0002, got %q", migrations[1].Version)
	}
}

func TestLoadMigrationsRejectsBadFileName(t *testing.T) {
	_, err := LoadMigrations(fstest.MapFS{
		"create_users.sql": {Data: []byte("create table users (id text);")},
	})
	if err == nil {
		t.Fatal("expected bad file name error")
	}
}

func TestLoadMigrationsRejectsDuplicateVersion(t *testing.T) {
	_, err := LoadMigrations(fstest.MapFS{
		"0001_create_users.sql":   {Data: []byte("create table users (id text);")},
		"0001_create_tenants.sql": {Data: []byte("create table tenants (id text);")},
	})
	if err == nil {
		t.Fatal("expected duplicate version error")
	}
}

func TestLoadMigrationsRejectsEmptyMigration(t *testing.T) {
	_, err := LoadMigrations(fstest.MapFS{
		"0001_create_users.sql": {Data: []byte("  \n\t")},
	})
	if err == nil {
		t.Fatal("expected empty migration error")
	}
}
