package wallet

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Wallet struct {
	ID                    WalletID
	DisplayID             int64
	TenantID              tenant.ID
	OwnerType             OwnerType
	OwnerID               OwnerID
	Currency              string
	Status                Status
	AvailableBalanceMinor int64
	LockedBalanceMinor    int64
	Metadata              json.RawMessage
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CreateWalletInput struct {
	TenantID              tenant.ID
	OwnerType             OwnerType
	OwnerID               OwnerID
	Currency              string
	Status                Status
	AvailableBalanceMinor int64
	LockedBalanceMinor    int64
	Metadata              json.RawMessage
}

func (input CreateWalletInput) Normalize() CreateWalletInput {
	output := input
	output.OwnerID = OwnerID(trim(string(output.OwnerID)))
	output.Currency = upperTrim(output.Currency)
	output.Metadata = defaultJSON(output.Metadata)
	if output.Status == "" {
		output.Status = StatusActive
	}
	return output
}

func (input CreateWalletInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.OwnerType.Valid() {
		return ErrOwnerTypeInvalid
	}
	if input.OwnerID.Empty() {
		return ErrOwnerIDMissing
	}
	if !input.Status.Valid() {
		return ErrStatusInvalid
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if err := validateBalance(input.AvailableBalanceMinor); err != nil {
		return err
	}
	if err := validateBalance(input.LockedBalanceMinor); err != nil {
		return err
	}
	return nil
}
