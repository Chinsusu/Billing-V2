package identity

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrPasswordMissing     = errors.New("password missing")
	ErrPasswordHashInvalid = errors.New("password hash invalid")
)

type Argon2idConfig struct {
	MemoryKiB   uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DefaultArgon2idConfig() Argon2idConfig {
	return Argon2idConfig{
		MemoryKiB:   64 * 1024,
		Iterations:  3,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func HashPasswordArgon2id(password string) (string, error) {
	return HashPasswordArgon2idWithConfig(password, DefaultArgon2idConfig())
}

func HashPasswordArgon2idWithConfig(password string, cfg Argon2idConfig) (string, error) {
	if password == "" {
		return "", ErrPasswordMissing
	}
	if err := cfg.Validate(); err != nil {
		return "", err
	}
	salt := make([]byte, cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.MemoryKiB, cfg.Parallelism, cfg.KeyLength)
	return formatArgon2idHash(cfg, salt, hash), nil
}

func VerifyPasswordArgon2id(password string, encodedHash string) (bool, error) {
	if password == "" {
		return false, ErrPasswordMissing
	}
	cfg, salt, expectedHash, err := parseArgon2idHash(encodedHash)
	if err != nil {
		return false, err
	}
	actualHash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.MemoryKiB, cfg.Parallelism, cfg.KeyLength)
	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1, nil
}

func (cfg Argon2idConfig) Validate() error {
	if cfg.MemoryKiB == 0 || cfg.Iterations == 0 || cfg.Parallelism == 0 || cfg.SaltLength == 0 || cfg.KeyLength == 0 {
		return ErrPasswordHashInvalid
	}
	return nil
}

func formatArgon2idHash(cfg Argon2idConfig, salt []byte, hash []byte) string {
	encoder := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		cfg.MemoryKiB,
		cfg.Iterations,
		cfg.Parallelism,
		encoder.EncodeToString(salt),
		encoder.EncodeToString(hash),
	)
}

func parseArgon2idHash(encodedHash string) (Argon2idConfig, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return Argon2idConfig{}, nil, nil, ErrPasswordHashInvalid
	}
	cfg, err := parseArgon2idParams(parts[3])
	if err != nil {
		return Argon2idConfig{}, nil, nil, err
	}
	encoder := base64.RawStdEncoding
	salt, err := encoder.DecodeString(parts[4])
	if err != nil {
		return Argon2idConfig{}, nil, nil, ErrPasswordHashInvalid
	}
	hash, err := encoder.DecodeString(parts[5])
	if err != nil {
		return Argon2idConfig{}, nil, nil, ErrPasswordHashInvalid
	}
	cfg.SaltLength = uint32(len(salt))
	cfg.KeyLength = uint32(len(hash))
	if err := cfg.Validate(); err != nil {
		return Argon2idConfig{}, nil, nil, err
	}
	return cfg, salt, hash, nil
}

func parseArgon2idParams(value string) (Argon2idConfig, error) {
	var cfg Argon2idConfig
	for _, part := range strings.Split(value, ",") {
		key, raw, ok := strings.Cut(part, "=")
		if !ok {
			return Argon2idConfig{}, ErrPasswordHashInvalid
		}
		parsed, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return Argon2idConfig{}, ErrPasswordHashInvalid
		}
		switch key {
		case "m":
			cfg.MemoryKiB = uint32(parsed)
		case "t":
			cfg.Iterations = uint32(parsed)
		case "p":
			if parsed > 255 {
				return Argon2idConfig{}, ErrPasswordHashInvalid
			}
			cfg.Parallelism = uint8(parsed)
		default:
			return Argon2idConfig{}, ErrPasswordHashInvalid
		}
	}
	if cfg.MemoryKiB == 0 || cfg.Iterations == 0 || cfg.Parallelism == 0 {
		return Argon2idConfig{}, ErrPasswordHashInvalid
	}
	return cfg, nil
}
