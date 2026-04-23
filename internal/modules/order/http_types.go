package order

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type createOrderRequest struct {
	TenantPlanID    catalog.TenantPlanID `json:"tenant_plan_id"`
	Quantity        int                  `json:"quantity"`
	Currency        string               `json:"currency"`
	UnitPriceMinor  int64                `json:"unit_price_minor"`
	DiscountMinor   int64                `json:"discount_minor"`
	TotalMinor      int64                `json:"total_minor"`
	ProductSnapshot json.RawMessage      `json:"product_snapshot"`
	PlanSnapshot    json.RawMessage      `json:"plan_snapshot"`
	PriceSnapshot   json.RawMessage      `json:"price_snapshot"`
}

func (request createOrderRequest) toInput(tenantID tenant.ID, buyerUserID identity.UserID, idempotencyKey IdempotencyKey) CreateOrderInput {
	return CreateOrderInput{
		TenantID:        tenantID,
		BuyerUserID:     buyerUserID,
		TenantPlanID:    request.TenantPlanID,
		Quantity:        request.Quantity,
		Currency:        request.Currency,
		UnitPriceMinor:  request.UnitPriceMinor,
		DiscountMinor:   request.DiscountMinor,
		TotalMinor:      request.TotalMinor,
		IdempotencyKey:  idempotencyKey,
		ProductSnapshot: request.ProductSnapshot,
		PlanSnapshot:    request.PlanSnapshot,
		PriceSnapshot:   request.PriceSnapshot,
	}
}

type orderResponse struct {
	ID              OrderID              `json:"id"`
	DisplayID       int64                `json:"display_id"`
	TenantID        tenant.ID            `json:"tenant_id"`
	BuyerUserID     identity.UserID      `json:"buyer_user_id"`
	TenantPlanID    catalog.TenantPlanID `json:"tenant_plan_id"`
	Quantity        int                  `json:"quantity"`
	Currency        string               `json:"currency"`
	UnitPriceMinor  int64                `json:"unit_price_minor"`
	DiscountMinor   int64                `json:"discount_minor"`
	TotalMinor      int64                `json:"total_minor"`
	OrderStatus     OrderStatus          `json:"order_status"`
	BillingStatus   BillingStatus        `json:"billing_status"`
	ProductSnapshot json.RawMessage      `json:"product_snapshot"`
	PlanSnapshot    json.RawMessage      `json:"plan_snapshot"`
	PriceSnapshot   json.RawMessage      `json:"price_snapshot"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

func newOrderResponse(order Order) orderResponse {
	return orderResponse{
		ID:              order.ID,
		DisplayID:       order.DisplayID,
		TenantID:        order.TenantID,
		BuyerUserID:     order.BuyerUserID,
		TenantPlanID:    order.TenantPlanID,
		Quantity:        order.Quantity,
		Currency:        order.Currency,
		UnitPriceMinor:  order.UnitPriceMinor,
		DiscountMinor:   order.DiscountMinor,
		TotalMinor:      order.TotalMinor,
		OrderStatus:     order.OrderStatus,
		BillingStatus:   order.BillingStatus,
		ProductSnapshot: order.ProductSnapshot,
		PlanSnapshot:    order.PlanSnapshot,
		PriceSnapshot:   order.PriceSnapshot,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func newOrderResponses(orders []Order) []orderResponse {
	responses := make([]orderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, newOrderResponse(order))
	}
	return responses
}
