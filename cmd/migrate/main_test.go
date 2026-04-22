package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunValidateAcceptsEmptyMigrationDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := run([]string{"-dir", dir, "validate"}); err != nil {
		t.Fatalf("run validate returned error: %v", err)
	}
}

func TestRunValidateRejectsInvalidMigrationName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "create_users.sql")
	if err := os.WriteFile(path, []byte("create table users (id text);"), 0600); err != nil {
		t.Fatalf("write migration: %v", err)
	}

	if err := run([]string{"-dir", dir, "validate"}); err == nil {
		t.Fatal("expected invalid migration name error")
	}
}

func TestRunUpRequiresDSN(t *testing.T) {
	dir := t.TempDir()

	if err := run([]string{"-dir", dir, "up"}); err == nil {
		t.Fatal("expected missing DSN error")
	}
}
