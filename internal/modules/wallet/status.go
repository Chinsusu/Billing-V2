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
