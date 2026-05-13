package order

type OrderStatus string

const (
	OrderStatusDraft          OrderStatus = "draft"
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusFailed         OrderStatus = "failed"
	OrderStatusRefunded       OrderStatus = "refunded"
)

func (status OrderStatus) Valid() bool {
	switch status {
	case OrderStatusDraft, OrderStatusPendingPayment, OrderStatusPaid, OrderStatusCancelled, OrderStatusFailed, OrderStatusRefunded:
		return true
	default:
		return false
	}
}

type BillingStatus string

const (
	BillingStatusUnpaid            BillingStatus = "unpaid"
	BillingStatusPaid              BillingStatus = "paid"
	BillingStatusOverdue           BillingStatus = "overdue"
	BillingStatusRefunded          BillingStatus = "refunded"
	BillingStatusPartiallyRefunded BillingStatus = "partially_refunded"
)

func (status BillingStatus) Valid() bool {
	switch status {
	case BillingStatusUnpaid, BillingStatusPaid, BillingStatusOverdue, BillingStatusRefunded, BillingStatusPartiallyRefunded:
		return true
	default:
		return false
	}
}

type ReservationStatus string

const (
	ReservationStatusPendingReserve ReservationStatus = "pending_reserve"
	ReservationStatusReserved       ReservationStatus = "reserved"
	ReservationStatusExpired        ReservationStatus = "reservation_expired"
	ReservationStatusReleased       ReservationStatus = "reservation_released"
	ReservationStatusAllocated      ReservationStatus = "allocated"
)

func (status ReservationStatus) Valid() bool {
	switch status {
	case ReservationStatusPendingReserve, ReservationStatusReserved, ReservationStatusExpired, ReservationStatusReleased, ReservationStatusAllocated:
		return true
	default:
		return false
	}
}

type ProvisioningStatus string

const (
	ProvisioningStatusQueued       ProvisioningStatus = "queued"
	ProvisioningStatusProvisioning ProvisioningStatus = "provisioning"
	ProvisioningStatusProvisioned  ProvisioningStatus = "provisioned"
	ProvisioningStatusFailed       ProvisioningStatus = "failed"
	ProvisioningStatusManualReview ProvisioningStatus = "manual_review"
)

func (status ProvisioningStatus) Valid() bool {
	switch status {
	case ProvisioningStatusQueued, ProvisioningStatusProvisioning, ProvisioningStatusProvisioned, ProvisioningStatusFailed, ProvisioningStatusManualReview:
		return true
	default:
		return false
	}
}

type ServiceStatus string

const (
	ServiceStatusActive     ServiceStatus = "active"
	ServiceStatusSuspended  ServiceStatus = "suspended"
	ServiceStatusExpired    ServiceStatus = "expired"
	ServiceStatusCancelled  ServiceStatus = "cancelled"
	ServiceStatusTerminated ServiceStatus = "terminated"
)

func (status ServiceStatus) Valid() bool {
	switch status {
	case ServiceStatusActive, ServiceStatusSuspended, ServiceStatusExpired, ServiceStatusCancelled, ServiceStatusTerminated:
		return true
	default:
		return false
	}
}

type SuspensionReason string

const (
	SuspensionReasonExpiry         SuspensionReason = "expiry"
	SuspensionReasonManualAdmin    SuspensionReason = "manual_admin"
	SuspensionReasonManualReseller SuspensionReason = "manual_reseller"
	SuspensionReasonAbuse          SuspensionReason = "abuse"
	SuspensionReasonSystemIssue    SuspensionReason = "system_issue"
)

func (reason SuspensionReason) Valid() bool {
	switch reason {
	case SuspensionReasonExpiry, SuspensionReasonManualAdmin, SuspensionReasonManualReseller, SuspensionReasonAbuse, SuspensionReasonSystemIssue:
		return true
	default:
		return false
	}
}

func CanTransitionOrder(from OrderStatus, to OrderStatus) bool {
	if from == to {
		return from.Valid()
	}
	switch from {
	case OrderStatusDraft:
		return to == OrderStatusPendingPayment || to == OrderStatusCancelled
	case OrderStatusPendingPayment:
		return to == OrderStatusPaid || to == OrderStatusCancelled
	case OrderStatusPaid:
		return to == OrderStatusRefunded || to == OrderStatusFailed
	default:
		return false
	}
}

func CanTransitionProvisioning(from ProvisioningStatus, to ProvisioningStatus) bool {
	if from == to {
		return from.Valid()
	}
	switch from {
	case ProvisioningStatusQueued:
		return to == ProvisioningStatusProvisioning || to == ProvisioningStatusManualReview
	case ProvisioningStatusProvisioning:
		return to == ProvisioningStatusProvisioned || to == ProvisioningStatusFailed || to == ProvisioningStatusManualReview
	case ProvisioningStatusFailed, ProvisioningStatusManualReview:
		return to == ProvisioningStatusQueued
	default:
		return false
	}
}

func CanTransitionService(from ServiceStatus, to ServiceStatus) bool {
	if from == to {
		return from.Valid()
	}
	switch from {
	case ServiceStatusActive:
		return to == ServiceStatusSuspended || to == ServiceStatusExpired || to == ServiceStatusCancelled || to == ServiceStatusTerminated
	case ServiceStatusSuspended:
		return to == ServiceStatusActive || to == ServiceStatusTerminated
	case ServiceStatusExpired:
		return to == ServiceStatusActive || to == ServiceStatusSuspended || to == ServiceStatusTerminated
	case ServiceStatusCancelled:
		return to == ServiceStatusTerminated
	default:
		return false
	}
}
