package wallet

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type walletResponse struct {
	ID                    WalletID        `json:"id"`
	DisplayID             int64           `json:"display_id"`
	TenantID              tenant.ID       `json:"tenant_id"`
	OwnerType             OwnerType       `json:"owner_type"`
	OwnerID               OwnerID         `json:"owner_id"`
	Currency              string          `json:"currency"`
	Status                Status          `json:"status"`
	AvailableBalanceMinor int64           `json:"available_balance_minor"`
	LockedBalanceMinor    int64           `json:"locked_balance_minor"`
	Metadata              json.RawMessage `json:"metadata"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type ledgerEntryResponse struct {
	ID                LedgerEntryID   `json:"id"`
	DisplayID         int64           `json:"display_id"`
	WalletID          WalletID        `json:"wallet_id"`
	TenantID          tenant.ID       `json:"tenant_id"`
	Direction         Direction       `json:"direction"`
	AmountMinor       int64           `json:"amount_minor"`
	Currency          string          `json:"currency"`
	EntryType         EntryType       `json:"entry_type"`
	Status            LedgerStatus    `json:"status"`
	BalanceAfterMinor int64           `json:"balance_after_minor"`
	ReferenceType     ReferenceType   `json:"reference_type"`
	ReferenceID       ReferenceID     `json:"reference_id"`
	CreatedBy         identity.UserID `json:"created_by,omitempty"`
	Reason            string          `json:"reason,omitempty"`
	CorrelationID     CorrelationID   `json:"correlation_id"`
	CreatedAt         time.Time       `json:"created_at"`
}

func newWalletResponse(wallet Wallet) walletResponse {
	return walletResponse{
		ID:                    wallet.ID,
		DisplayID:             wallet.DisplayID,
		TenantID:              wallet.TenantID,
		OwnerType:             wallet.OwnerType,
		OwnerID:               wallet.OwnerID,
		Currency:              wallet.Currency,
		Status:                wallet.Status,
		AvailableBalanceMinor: wallet.AvailableBalanceMinor,
		LockedBalanceMinor:    wallet.LockedBalanceMinor,
		Metadata:              wallet.Metadata,
		CreatedAt:             wallet.CreatedAt,
		UpdatedAt:             wallet.UpdatedAt,
	}
}

func newWalletResponses(wallets []Wallet) []walletResponse {
	responses := make([]walletResponse, 0, len(wallets))
	for _, wallet := range wallets {
		responses = append(responses, newWalletResponse(wallet))
	}
	return responses
}

func newLedgerEntryResponse(entry LedgerEntry) ledgerEntryResponse {
	return ledgerEntryResponse{
		ID:                entry.ID,
		DisplayID:         entry.DisplayID,
		WalletID:          entry.WalletID,
		TenantID:          entry.TenantID,
		Direction:         entry.Direction,
		AmountMinor:       entry.AmountMinor,
		Currency:          entry.Currency,
		EntryType:         entry.EntryType,
		Status:            entry.Status,
		BalanceAfterMinor: entry.BalanceAfterMinor,
		ReferenceType:     entry.ReferenceType,
		ReferenceID:       entry.ReferenceID,
		CreatedBy:         entry.CreatedBy,
		Reason:            entry.Reason,
		CorrelationID:     entry.CorrelationID,
		CreatedAt:         entry.CreatedAt,
	}
}

func newLedgerEntryResponses(entries []LedgerEntry) []ledgerEntryResponse {
	responses := make([]ledgerEntryResponse, 0, len(entries))
	for _, entry := range entries {
		responses = append(responses, newLedgerEntryResponse(entry))
	}
	return responses
}
