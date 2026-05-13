package order

import (
	"context"
)

const serviceCredentialColumns = `credential_id, tenant_id, service_instance_id, credential_type, encrypted_payload, encryption_key_version, encryption_algorithm, secret_version, masked_hint, status, last_revealed_at, last_revealed_by, rotated_at, rotated_by, created_at, updated_at`

const createServiceCredentialSQL = `
INSERT INTO service_credentials (tenant_id, service_instance_id, credential_type, encrypted_payload, encryption_key_version, encryption_algorithm, secret_version, masked_hint, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (tenant_id, service_instance_id, credential_type) WHERE status = 'active'
DO UPDATE SET
    encrypted_payload = EXCLUDED.encrypted_payload,
    encryption_key_version = EXCLUDED.encryption_key_version,
    encryption_algorithm = EXCLUDED.encryption_algorithm,
    secret_version = EXCLUDED.secret_version,
    masked_hint = EXCLUDED.masked_hint,
    updated_at = NOW()
RETURNING ` + serviceCredentialColumns

func (store *PostgresStore) CreateServiceCredential(ctx context.Context, input CreateServiceCredentialInput) (ServiceCredential, error) {
	if err := store.ready(); err != nil {
		return ServiceCredential{}, err
	}
	args, err := createServiceCredentialArgs(input)
	if err != nil {
		return ServiceCredential{}, err
	}
	return scanServiceCredential(store.executor.QueryRowContext(ctx, createServiceCredentialSQL, args...))
}

func createServiceCredentialArgs(input CreateServiceCredentialInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.TenantID,
		input.ServiceID,
		input.Type,
		input.EncryptedPayload,
		input.EncryptionKeyVersion,
		input.EncryptionAlgorithm,
		nullableString(input.SecretVersion),
		input.MaskedHint,
		input.Status,
	}, nil
}
