package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	AggregateTypeInvoice  = "invoice"
	EventInvoiceGenerated = "invoice.generated"
	EventInvoicePaid      = "invoice.paid"
)

type OrderReader interface {
	GetOrder(ctx context.Context, lookup order.OrderLookup) (order.Order, error)
}

type GenerateInvoiceInput struct {
	TenantID       tenant.ID
	OrderID        order.OrderID
	IdempotencyKey IdempotencyKey
}

type CreateInvoiceFromOrderInput struct {
	Invoice        CreateInvoiceInput
	Item           GeneratedInvoiceItemInput
	IdempotencyKey IdempotencyKey
	OrderDisplayID int64
}

type GeneratedInvoiceItemInput struct {
	OrderID        order.OrderID
	OrderItemID    OrderItemID
	ServiceID      order.ServiceID
	Description    string
	Quantity       int
	UnitPriceMinor int64
	TaxMinor       int64
	DiscountMinor  int64
	LineTotalMinor int64
	Metadata       json.RawMessage
}

func (service *Service) GenerateInvoiceForOrder(ctx context.Context, input GenerateInvoiceInput) (InvoiceDetail, error) {
	if err := service.ready(); err != nil {
		return InvoiceDetail{}, err
	}
	if service.orderReader == nil {
		return InvoiceDetail{}, ErrOrderReaderMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return InvoiceDetail{}, err
	}
	sourceOrder, err := service.orderReader.GetOrder(ctx, order.OrderLookup{ID: input.OrderID, TenantID: input.TenantID})
	if err != nil {
		return InvoiceDetail{}, err
	}
	if sourceOrder.TenantID != input.TenantID {
		return InvoiceDetail{}, tenant.ErrAccessDenied
	}
	if sourceOrder.OrderStatus != order.OrderStatusPaid || sourceOrder.BillingStatus != order.BillingStatusPaid {
		return InvoiceDetail{}, ErrOrderNotPaid
	}
	createInput, err := newCreateInvoiceFromOrderInput(sourceOrder, input.IdempotencyKey)
	if err != nil {
		return InvoiceDetail{}, err
	}
	return service.store.CreateInvoiceFromOrder(ctx, createInput)
}

func (input GenerateInvoiceInput) Normalize() GenerateInvoiceInput {
	output := input
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
}

func (input GenerateInvoiceInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.OrderID.Empty() {
		return order.ErrOrderIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}

func (input CreateInvoiceFromOrderInput) Normalize() CreateInvoiceFromOrderInput {
	output := input
	output.Invoice = output.Invoice.Normalize()
	output.Item = output.Item.Normalize()
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
}

func (input CreateInvoiceFromOrderInput) Validate() error {
	if input.Invoice.OrderID.Empty() {
		return order.ErrOrderIDMissing
	}
	if err := input.Invoice.Validate(); err != nil {
		return err
	}
	if err := input.Item.Validate(); err != nil {
		return err
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}

func (input GeneratedInvoiceItemInput) Normalize() GeneratedInvoiceItemInput {
	output := input
	output.Description = trim(output.Description)
	output.Metadata = defaultJSON(output.Metadata)
	return output
}

func (input GeneratedInvoiceItemInput) Validate() error {
	if input.OrderID.Empty() {
		return order.ErrOrderIDMissing
	}
	if input.Description == "" {
		return ErrDescriptionMissing
	}
	if input.Quantity <= 0 {
		return ErrQuantityInvalid
	}
	if err := validateMinorAmount(input.UnitPriceMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.TaxMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.DiscountMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.LineTotalMinor); err != nil {
		return err
	}
	expectedTotal := int64(input.Quantity)*input.UnitPriceMinor + input.TaxMinor - input.DiscountMinor
	if expectedTotal != input.LineTotalMinor {
		return ErrTotalInvalid
	}
	return nil
}

func newCreateInvoiceFromOrderInput(sourceOrder order.Order, idempotencyKey IdempotencyKey) (CreateInvoiceFromOrderInput, error) {
	if sourceOrder.Quantity <= 0 {
		return CreateInvoiceFromOrderInput{}, ErrQuantityInvalid
	}
	subtotalMinor := int64(sourceOrder.Quantity) * sourceOrder.UnitPriceMinor
	if subtotalMinor-sourceOrder.DiscountMinor != sourceOrder.TotalMinor {
		return CreateInvoiceFromOrderInput{}, ErrTotalInvalid
	}
	metadata, err := invoiceGenerationMetadata(sourceOrder, idempotencyKey)
	if err != nil {
		return CreateInvoiceFromOrderInput{}, err
	}
	return CreateInvoiceFromOrderInput{
		Invoice: CreateInvoiceInput{
			TenantID:      sourceOrder.TenantID,
			BuyerUserID:   sourceOrder.BuyerUserID,
			OrderID:       sourceOrder.ID,
			Status:        StatusIssued,
			Currency:      sourceOrder.Currency,
			SubtotalMinor: subtotalMinor,
			TaxMinor:      0,
			DiscountMinor: sourceOrder.DiscountMinor,
			TotalMinor:    sourceOrder.TotalMinor,
			IssuedAt:      time.Now().UTC(),
			Metadata:      metadata,
		},
		Item: GeneratedInvoiceItemInput{
			OrderID:        sourceOrder.ID,
			Description:    invoiceItemDescription(sourceOrder),
			Quantity:       sourceOrder.Quantity,
			UnitPriceMinor: sourceOrder.UnitPriceMinor,
			TaxMinor:       0,
			DiscountMinor:  sourceOrder.DiscountMinor,
			LineTotalMinor: sourceOrder.TotalMinor,
			Metadata:       metadata,
		},
		IdempotencyKey: idempotencyKey,
		OrderDisplayID: sourceOrder.DisplayID,
	}, nil
}

func invoiceGenerationMetadata(sourceOrder order.Order, idempotencyKey IdempotencyKey) (json.RawMessage, error) {
	payload := struct {
		Source          string          `json:"source"`
		IdempotencyKey  IdempotencyKey  `json:"idempotency_key"`
		OrderID         order.OrderID   `json:"order_id"`
		OrderDisplayID  int64           `json:"order_display_id"`
		ProductSnapshot json.RawMessage `json:"product_snapshot"`
		PlanSnapshot    json.RawMessage `json:"plan_snapshot"`
		PriceSnapshot   json.RawMessage `json:"price_snapshot"`
	}{
		Source:          "order",
		IdempotencyKey:  idempotencyKey,
		OrderID:         sourceOrder.ID,
		OrderDisplayID:  sourceOrder.DisplayID,
		ProductSnapshot: defaultJSON(sourceOrder.ProductSnapshot),
		PlanSnapshot:    defaultJSON(sourceOrder.PlanSnapshot),
		PriceSnapshot:   defaultJSON(sourceOrder.PriceSnapshot),
	}
	value, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("build invoice metadata: %w", err)
	}
	return value, nil
}

func invoiceItemDescription(sourceOrder order.Order) string {
	productName := snapshotName(sourceOrder.ProductSnapshot)
	planName := snapshotName(sourceOrder.PlanSnapshot)
	switch {
	case productName != "" && planName != "":
		return productName + " - " + planName
	case planName != "":
		return planName
	case productName != "":
		return productName
	default:
		return fmt.Sprintf("Order %d", sourceOrder.DisplayID)
	}
}

func snapshotName(value json.RawMessage) string {
	var payload struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	}
	if err := json.Unmarshal(defaultJSON(value), &payload); err != nil {
		return ""
	}
	if name := strings.TrimSpace(payload.Name); name != "" {
		return name
	}
	return strings.TrimSpace(payload.DisplayName)
}
