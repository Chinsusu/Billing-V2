package identity

import "testing"

func TestAESGCMSecretCipherRoundTrips(t *testing.T) {
	cipher, err := NewAESGCMSecretCipher(testCipherKey())
	if err != nil {
		t.Fatalf("NewAESGCMSecretCipher returned error: %v", err)
	}
	ciphertext, err := cipher.Encrypt("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}
	if ciphertext == "JBSWY3DPEHPK3PXP" {
		t.Fatal("ciphertext must not equal plaintext")
	}
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}
	if plaintext != "JBSWY3DPEHPK3PXP" {
		t.Fatalf("unexpected plaintext %q", plaintext)
	}
}

func TestAESGCMSecretCipherRejectsInvalidKey(t *testing.T) {
	if _, err := NewAESGCMSecretCipher("short"); err != ErrEncryptionKeyInvalid {
		t.Fatalf("expected invalid key error, got %v", err)
	}
}

func testCipherKey() string {
	return "1234567890123456" + "7890123456789012"
}
