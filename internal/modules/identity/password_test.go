package identity

import (
	"errors"
	"testing"
)

func TestArgon2idPasswordHashVerifies(t *testing.T) {
	hash, err := HashPasswordArgon2idWithConfig("correct horse", testArgon2idConfig())
	if err != nil {
		t.Fatalf("HashPasswordArgon2id returned error: %v", err)
	}
	if hash == "correct horse" {
		t.Fatal("password hash must not equal plaintext")
	}
	ok, err := VerifyPasswordArgon2id("correct horse", hash)
	if err != nil {
		t.Fatalf("VerifyPasswordArgon2id returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected password to verify")
	}
}

func TestArgon2idPasswordRejectsWrongPassword(t *testing.T) {
	hash, err := HashPasswordArgon2idWithConfig("correct horse", testArgon2idConfig())
	if err != nil {
		t.Fatalf("HashPasswordArgon2id returned error: %v", err)
	}
	ok, err := VerifyPasswordArgon2id("wrong", hash)
	if err != nil {
		t.Fatalf("VerifyPasswordArgon2id returned error: %v", err)
	}
	if ok {
		t.Fatal("expected wrong password to be rejected")
	}
}

func TestArgon2idPasswordRejectsInvalidHash(t *testing.T) {
	_, err := VerifyPasswordArgon2id("password", "dev-only-placeholder-hash")
	if !errors.Is(err, ErrPasswordHashInvalid) {
		t.Fatalf("expected invalid hash error, got %v", err)
	}
}

func testArgon2idConfig() Argon2idConfig {
	return Argon2idConfig{
		MemoryKiB:   32,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  8,
		KeyLength:   16,
	}
}
