package order

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Order struct {
	ID              OrderID
	DisplayID       int64
	TenantID        tenant.ID
	BuyerUserID     identity.UserID
	BuyerDisplayID  int64
	TenantPlanID    catalog.TenantPlanID
	Quantity        int
	Currency        string
	UnitPriceMinor  int64
	DiscountMinor   int64
	TotalMinor      int64
	OrderStatus     OrderStatus
	BillingStatus   BillingStatus
	IdempotencyKey  IdempotencyKey
	ProductSnapshot json.RawMessage
	PlanSnapshot    json.RawMessage
	PriceSnapshot   json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateOrderInput struct {
	TenantID        tenant.ID
	BuyerUserID     identity.UserID
	TenantPlanID    catalog.TenantPlanID
	Quantity        int
	Currency        string
	UnitPriceMinor  int64
	DiscountMinor   int64
	TotalMinor      int64
	OrderStatus     OrderStatus
	BillingStatus   BillingStatus
	IdempotencyKey  IdempotencyKey
	ProductSnapshot json.RawMessage
	PlanSnapshot    json.RawMessage
	PriceSnapshot   json.RawMessage
}

func (input CreateOrderInput) Normalize() CreateOrderInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.ProductSnapshot = defaultJSON(output.ProductSnapshot)
	output.PlanSnapshot = defaultJSON(output.PlanSnapshot)
	output.PriceSnapshot = defaultJSON(output.PriceSnapshot)
	if output.Quantity == 0 {
		output.Quantity = 1
	}
	if output.OrderStatus == "" {
		output.OrderStatus = OrderStatusPendingPayment
	}
	if output.BillingStatus == "" {
		output.BillingStatus = BillingStatusUnpaid
	}
	return output
}

func (input CreateOrderInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.BuyerUserID == "" {
		return ErrBuyerIDMissing
	}
	if input.TenantPlanID.Empty() {
		return ErrTenantPlanIDMissing
	}
	if input.Quantity <= 0 {
		return ErrQuantityInvalid
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if err := validateMinorAmount(input.UnitPriceMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.DiscountMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.TotalMinor); err != nil {
		return err
	}
	if !input.OrderStatus.Valid() {
		return ErrOrderStatusInvalid
	}
	if !input.BillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}
