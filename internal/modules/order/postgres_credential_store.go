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

const listServiceCredentialsSQL = `
SELECT ` + serviceCredentialColumns + `
FROM service_credentials
WHERE tenant_id = $1
  AND service_instance_id = $2
ORDER BY created_at DESC`

const listServiceCredentialsByStatusSQL = `
SELECT ` + serviceCredentialColumns + `
FROM service_credentials
WHERE tenant_id = $1
  AND service_instance_id = $2
  AND status = $3
ORDER BY created_at DESC`

const getServiceCredentialSQL = `
SELECT ` + serviceCredentialColumns + `
FROM service_credentials
WHERE credential_id = $1
  AND tenant_id = $2
  AND service_instance_id = $3`

const markServiceCredentialRevealedSQL = `
UPDATE service_credentials
SET last_revealed_at = $4,
    last_revealed_by = $5,
    updated_at = NOW()
WHERE credential_id = $1
  AND tenant_id = $2
  AND service_instance_id = $3
  AND status = 'active'
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

func (store *PostgresStore) ListServiceCredentials(ctx context.Context, filter ServiceCredentialFilter) ([]ServiceCredential, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	filter = filter.Normalize()
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	query := listServiceCredentialsSQL
	args := []interface{}{filter.TenantID, filter.ServiceID}
	if filter.Status != "" {
		query = listServiceCredentialsByStatusSQL
		args = append(args, filter.Status)
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	credentials := make([]ServiceCredential, 0)
	for rows.Next() {
		record, err := scanServiceCredential(rows)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return credentials, nil
}

func (store *PostgresStore) GetServiceCredential(ctx context.Context, lookup ServiceCredentialLookup) (ServiceCredential, error) {
	if err := store.ready(); err != nil {
		return ServiceCredential{}, err
	}
	lookup = lookup.Normalize()
	if err := lookup.Validate(); err != nil {
		return ServiceCredential{}, err
	}
	return scanServiceCredential(store.executor.QueryRowContext(ctx, getServiceCredentialSQL, lookup.ID, lookup.TenantID, lookup.ServiceID))
}

func (store *PostgresStore) MarkServiceCredentialRevealed(ctx context.Context, input MarkServiceCredentialRevealedInput) (ServiceCredential, error) {
	if err := store.ready(); err != nil {
		return ServiceCredential{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ServiceCredential{}, err
	}
	return scanServiceCredential(store.executor.QueryRowContext(
		ctx,
		markServiceCredentialRevealedSQL,
		input.ID,
		input.TenantID,
		input.ServiceID,
		input.RevealedAt,
		input.ActorID,
	))
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
