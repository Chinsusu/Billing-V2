package wallet

type OwnerType string

const (
	OwnerTypeTenant             OwnerType = "tenant"
	OwnerTypeUser               OwnerType = "user"
	OwnerTypeResellerSettlement OwnerType = "reseller_settlement"
	OwnerTypePlatform           OwnerType = "platform"
)

func (ownerType OwnerType) Valid() bool {
	switch ownerType {
	case OwnerTypeTenant, OwnerTypeUser, OwnerTypeResellerSettlement, OwnerTypePlatform:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusActive Status = "active"
	StatusFrozen Status = "frozen"
	StatusClosed Status = "closed"
)

func (status Status) Valid() bool {
	switch status {
	case StatusActive, StatusFrozen, StatusClosed:
		return true
	default:
		return false
	}
}

type Direction string

const (
	DirectionCredit Direction = "credit"
	DirectionDebit  Direction = "debit"
)

func (direction Direction) Valid() bool {
	switch direction {
	case DirectionCredit, DirectionDebit:
		return true
	default:
		return false
	}
}

type EntryType string

const (
	EntryTypeTopup        EntryType = "topup"
	EntryTypePurchase     EntryType = "purchase"
	EntryTypeResellerCost EntryType = "reseller_cost"
	EntryTypeRefund       EntryType = "refund"
	EntryTypeAdjustment   EntryType = "adjustment"
	EntryTypeReversal     EntryType = "reversal"
	EntryTypeCommission   EntryType = "commission"
	EntryTypeLock         EntryType = "lock"
	EntryTypeUnlock       EntryType = "unlock"
)

func (entryType EntryType) Valid() bool {
	switch entryType {
	case EntryTypeTopup, EntryTypePurchase, EntryTypeResellerCost, EntryTypeRefund,
		EntryTypeAdjustment, EntryTypeReversal, EntryTypeCommission, EntryTypeLock, EntryTypeUnlock:
		return true
	default:
		return false
	}
}

type LedgerStatus string

const (
	LedgerStatusPosted           LedgerStatus = "posted"
	LedgerStatusVoidedByReversal LedgerStatus = "voided_by_reversal"
)

func (status LedgerStatus) Valid() bool {
	switch status {
	case LedgerStatusPosted, LedgerStatusVoidedByReversal:
		return true
	default:
		return false
	}
}
