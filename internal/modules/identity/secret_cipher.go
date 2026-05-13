package identity

import (
	"errors"

	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

var (
	ErrSecretCipherMissing  = secrets.ErrCipherMissing
	ErrEncryptionKeyInvalid = secrets.ErrKeyInvalid
)

type SecretCipher = secrets.Cipher

type AESGCMSecretCipher = secrets.AESGCMCipher

func NewAESGCMSecretCipher(rawKey string) (*AESGCMSecretCipher, error) {
	cipher, err := secrets.NewAESGCMCipher(rawKey)
	if err != nil {
		if errors.Is(err, secrets.ErrKeyInvalid) || errors.Is(err, secrets.ErrPayloadInvalid) {
			return nil, ErrEncryptionKeyInvalid
		}
		return nil, err
	}
	return cipher, nil
}
