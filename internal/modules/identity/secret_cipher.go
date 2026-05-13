package identity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrSecretCipherMissing  = errors.New("secret cipher missing")
	ErrEncryptionKeyInvalid = errors.New("encryption key invalid")
)

type SecretCipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type AESGCMSecretCipher struct {
	aead cipher.AEAD
}

func NewAESGCMSecretCipher(rawKey string) (*AESGCMSecretCipher, error) {
	key, err := parseEncryptionKey(rawKey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create secret cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create secret gcm: %w", err)
	}
	return &AESGCMSecretCipher{aead: aead}, nil
}

func (cipher *AESGCMSecretCipher) Encrypt(plaintext string) (string, error) {
	if cipher == nil || cipher.aead == nil {
		return "", ErrSecretCipherMissing
	}
	nonce := make([]byte, cipher.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate secret nonce: %w", err)
	}
	sealed := cipher.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (cipher *AESGCMSecretCipher) Decrypt(ciphertext string) (string, error) {
	if cipher == nil || cipher.aead == nil {
		return "", ErrSecretCipherMissing
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(ciphertext))
	if err != nil {
		return "", ErrEncryptionKeyInvalid
	}
	nonceSize := cipher.aead.NonceSize()
	if len(raw) <= nonceSize {
		return "", ErrEncryptionKeyInvalid
	}
	plaintext, err := cipher.aead.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return "", ErrEncryptionKeyInvalid
	}
	return string(plaintext), nil
}

func parseEncryptionKey(rawKey string) ([]byte, error) {
	rawKey = strings.TrimSpace(rawKey)
	if rawKey == "" {
		return nil, ErrSecretCipherMissing
	}
	if decoded, err := base64.StdEncoding.DecodeString(rawKey); err == nil && len(decoded) == 32 {
		return decoded, nil
	}
	if decoded, err := base64.RawStdEncoding.DecodeString(rawKey); err == nil && len(decoded) == 32 {
		return decoded, nil
	}
	if decoded, err := hex.DecodeString(rawKey); err == nil && len(decoded) == 32 {
		return decoded, nil
	}
	if len([]byte(rawKey)) == 32 {
		return []byte(rawKey), nil
	}
	return nil, ErrEncryptionKeyInvalid
}
