package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func scanServiceCredential(row orderScanner) (ServiceCredential, error) {
	var record ServiceCredential
	var id, tenantID, serviceID, credentialType, status string
	var secretVersion, lastRevealedBy, rotatedBy sql.NullString
	var lastRevealedAt, rotatedAt sql.NullTime
	if err := row.Scan(
		&id,
		&tenantID,
		&serviceID,
		&credentialType,
		&record.EncryptedPayload,
		&record.EncryptionKeyVersion,
		&record.EncryptionAlgorithm,
		&secretVersion,
		&record.MaskedHint,
		&status,
		&lastRevealedAt,
		&lastRevealedBy,
		&rotatedAt,
		&rotatedBy,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceCredential{}, ErrCredentialNotFound
		}
		return ServiceCredential{}, fmt.Errorf("scan service credential: %w", err)
	}
	record.ID = CredentialID(id)
	record.TenantID = tenant.ID(tenantID)
	record.ServiceID = ServiceID(serviceID)
	record.Type = CredentialType(credentialType)
	record.SecretVersion = secretVersion.String
	record.Status = CredentialStatus(status)
	if lastRevealedAt.Valid {
		record.LastRevealedAt = lastRevealedAt.Time
	}
	if lastRevealedBy.Valid {
		record.LastRevealedBy = identity.UserID(lastRevealedBy.String)
	}
	if rotatedAt.Valid {
		record.RotatedAt = rotatedAt.Time
	}
	if rotatedBy.Valid {
		record.RotatedBy = identity.UserID(rotatedBy.String)
	}
	return record, nil
}
