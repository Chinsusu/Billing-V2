package wallet

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ledgerEntryScanner interface {
	Scan(dest ...interface{}) error
}

func scanWallet(row ledgerEntryScanner) (Wallet, error) {
	var record Wallet
	var id, tenantID, ownerType, ownerID, status string
	var metadata []byte
	if err := row.Scan(
		&id, &record.DisplayID, &tenantID, &ownerType, &ownerID, &record.Currency, &status,
		&record.AvailableBalanceMinor, &record.LockedBalanceMinor, &metadata, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Wallet{}, ErrWalletNotFound
		}
		return Wallet{}, fmt.Errorf("scan wallet: %w", err)
	}
	record.ID = WalletID(id)
	record.TenantID = tenant.ID(tenantID)
	record.OwnerType = OwnerType(ownerType)
	record.OwnerID = OwnerID(ownerID)
	record.Status = Status(status)
	record.Metadata = append(record.Metadata, metadata...)
	return record, nil
}

func scanLedgerEntry(row ledgerEntryScanner) (LedgerEntry, error) {
	var record LedgerEntry
	var id, walletID, tenantID, direction, currency, entryType, status, referenceType, referenceID, idempotencyKey, correlationID string
	var createdBy, reason sql.NullString
	if err := row.Scan(
		&id, &record.DisplayID, &walletID, &tenantID, &direction, &record.AmountMinor, &currency,
		&entryType, &status, &record.BalanceAfterMinor, &referenceType, &referenceID, &idempotencyKey,
		&createdBy, &reason, &correlationID, &record.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LedgerEntry{}, ErrLedgerEntryNotFound
		}
		return LedgerEntry{}, fmt.Errorf("scan wallet ledger entry: %w", err)
	}
	record.ID = LedgerEntryID(id)
	record.WalletID = WalletID(walletID)
	record.TenantID = tenant.ID(tenantID)
	record.Direction = Direction(direction)
	record.Currency = currency
	record.EntryType = EntryType(entryType)
	record.Status = LedgerStatus(status)
	record.ReferenceType = ReferenceType(referenceType)
	record.ReferenceID = ReferenceID(referenceID)
	record.IdempotencyKey = IdempotencyKey(idempotencyKey)
	record.CreatedBy = identity.UserID(createdBy.String)
	record.Reason = reason.String
	record.CorrelationID = CorrelationID(correlationID)
	return record, nil
}
