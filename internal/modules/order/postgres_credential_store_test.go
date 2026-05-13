package order

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateServiceCredentialArgsNormalizeAndValidate(t *testing.T) {
	args, err := createServiceCredentialArgs(CreateServiceCredentialInput{
		TenantID:         tenant.ID("tenant-1"),
		ServiceID:        ServiceID("service-1"),
		Type:             CredentialTypeVPSRoot,
		EncryptedPayload: " encrypted-fixture ",
		SecretVersion:    " secret-v1 ",
		MaskedHint:       " root / **** ",
	})
	if err != nil {
		t.Fatalf("expected service credential args: %v", err)
	}
	if len(args) != 9 ||
		args[3] != "encrypted-fixture" ||
		args[4] != DefaultCredentialEncryptionKeyVersion ||
		args[5] != DefaultCredentialEncryptionAlgorithm ||
		args[7] != "root / ****" ||
		args[8] != CredentialStatusActive {
		t.Fatalf("unexpected service credential args: %#v", args)
	}
	secretVersion, ok := args[6].(sql.NullString)
	if !ok || !secretVersion.Valid || secretVersion.String != "secret-v1" {
		t.Fatalf("expected nullable secret version, got %#v", args[6])
	}
}

func TestCreateServiceCredentialArgsRejectsInvalidInput(t *testing.T) {
	_, err := createServiceCredentialArgs(CreateServiceCredentialInput{})
	if err != tenant.ErrTenantIDMissing {
		t.Fatalf("expected tenant error, got %v", err)
	}

	_, err = createServiceCredentialArgs(CreateServiceCredentialInput{
		TenantID:         tenant.ID("tenant-1"),
		ServiceID:        ServiceID("service-1"),
		Type:             CredentialType("unknown"),
		EncryptedPayload: "encrypted-fixture",
	})
	if err != ErrCredentialTypeInvalid {
		t.Fatalf("expected credential type error, got %v", err)
	}
}

func TestCreateServiceCredentialSQLUpsertsActiveCredential(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO service_credentials",
		"ON CONFLICT (tenant_id, service_instance_id, credential_type) WHERE status = 'active'",
		"encrypted_payload = EXCLUDED.encrypted_payload",
		"RETURNING",
	} {
		if !strings.Contains(createServiceCredentialSQL, clause) {
			t.Fatalf("expected %q in service credential SQL: %s", clause, createServiceCredentialSQL)
		}
	}
}
