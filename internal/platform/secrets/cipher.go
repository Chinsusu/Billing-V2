package secrets

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
	ErrCipherMissing  = errors.New("secret cipher missing")
	ErrKeyInvalid     = errors.New("encryption key invalid")
	ErrPayloadInvalid = ErrKeyInvalid
)

type Cipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type AESGCMCipher struct {
	aead cipher.AEAD
}

func NewAESGCMCipher(rawKey string) (*AESGCMCipher, error) {
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
	return &AESGCMCipher{aead: aead}, nil
}

func (cipher *AESGCMCipher) Encrypt(plaintext string) (string, error) {
	if cipher == nil || cipher.aead == nil {
		return "", ErrCipherMissing
	}
	nonce := make([]byte, cipher.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate secret nonce: %w", err)
	}
	sealed := cipher.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (cipher *AESGCMCipher) Decrypt(ciphertext string) (string, error) {
	if cipher == nil || cipher.aead == nil {
		return "", ErrCipherMissing
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(ciphertext))
	if err != nil {
		return "", ErrPayloadInvalid
	}
	nonceSize := cipher.aead.NonceSize()
	if len(raw) <= nonceSize {
		return "", ErrPayloadInvalid
	}
	plaintext, err := cipher.aead.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return "", ErrPayloadInvalid
	}
	return string(plaintext), nil
}

func parseEncryptionKey(rawKey string) ([]byte, error) {
	rawKey = strings.TrimSpace(rawKey)
	if rawKey == "" {
		return nil, ErrCipherMissing
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
	return nil, ErrKeyInvalid
}
