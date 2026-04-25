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
	var ownerDisplayID sql.NullInt64
	var metadata []byte
	if err := row.Scan(
		&id, &record.DisplayID, &tenantID, &ownerType, &ownerID, &ownerDisplayID, &record.Currency, &status,
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
	if ownerDisplayID.Valid {
		record.OwnerDisplayID = ownerDisplayID.Int64
	}
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

func scanTopupRequest(row ledgerEntryScanner) (TopupRequest, error) {
	return scanTopupRequestFields(row, false)
}

func scanTopupRequestRead(row ledgerEntryScanner) (TopupRequest, error) {
	return scanTopupRequestFields(row, true)
}

func scanTopupRequestFields(row ledgerEntryScanner, includeRelatedDisplayIDs bool) (TopupRequest, error) {
	var record TopupRequest
	var id, tenantID, walletID, requestedBy, currency, method, status, idempotencyKey string
	var walletDisplayID, requestedByDisplayID, reviewedByDisplayID sql.NullInt64
	var paymentReference, reviewedBy, reviewNote, ledgerEntryID sql.NullString
	var reviewedAt sql.NullTime
	destinations := []interface{}{
		&id, &record.DisplayID, &tenantID, &walletID, &requestedBy, &record.AmountMinor,
		&currency, &method, &paymentReference, &status, &reviewedBy, &reviewedAt,
		&reviewNote, &ledgerEntryID, &idempotencyKey, &record.CreatedAt, &record.UpdatedAt,
	}
	if includeRelatedDisplayIDs {
		destinations = append(destinations, &walletDisplayID, &requestedByDisplayID, &reviewedByDisplayID)
	}
	if err := row.Scan(destinations...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TopupRequest{}, ErrTopupRequestNotFound
		}
		return TopupRequest{}, fmt.Errorf("scan wallet top-up request: %w", err)
	}
	record.ID = TopupRequestID(id)
	record.TenantID = tenant.ID(tenantID)
	record.WalletID = WalletID(walletID)
	if walletDisplayID.Valid {
		record.WalletDisplayID = walletDisplayID.Int64
	}
	record.RequestedBy = identity.UserID(requestedBy)
	if requestedByDisplayID.Valid {
		record.RequestedByDisplayID = requestedByDisplayID.Int64
	}
	record.Currency = currency
	record.PaymentMethod = PaymentMethod(method)
	record.PaymentReference = paymentReference.String
	record.Status = TopupStatus(status)
	record.ReviewedBy = identity.UserID(reviewedBy.String)
	if reviewedByDisplayID.Valid {
		record.ReviewedByDisplayID = reviewedByDisplayID.Int64
	}
	if reviewedAt.Valid {
		record.ReviewedAt = &reviewedAt.Time
	}
	record.ReviewNote = reviewNote.String
	record.LedgerEntryID = LedgerEntryID(ledgerEntryID.String)
	record.IdempotencyKey = IdempotencyKey(idempotencyKey)
	return record, nil
}
