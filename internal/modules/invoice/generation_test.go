package invoice

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestGenerateInvoiceForPaidOrder(t *testing.T) {
	store := &fakeInvoiceStore{}
	reader := &fakeOrderReader{record: paidOrder()}
	service := NewServiceWithOrderReader(store, reader)

	detail, err := service.GenerateInvoiceForOrder(context.Background(), GenerateInvoiceInput{
		TenantID:       tenant.ID("tenant-1"),
		OrderID:        order.OrderID("order-1"),
		IdempotencyKey: IdempotencyKey(" key-1 "),
	})
	if err != nil {
		t.Fatalf("expected invoice detail: %v", err)
	}
	if reader.lookup.ID != order.OrderID("order-1") || reader.lookup.TenantID != tenant.ID("tenant-1") {
		t.Fatalf("unexpected order lookup: %+v", reader.lookup)
	}
	if store.createCalls != 1 {
		t.Fatalf("expected create once, got %d", store.createCalls)
	}
	if store.createInput.Invoice.Status != StatusIssued ||
		store.createInput.Invoice.SubtotalMinor != 2000 ||
		store.createInput.Invoice.DiscountMinor != 200 ||
		store.createInput.Invoice.TotalMinor != 1800 {
		t.Fatalf("unexpected invoice input: %+v", store.createInput.Invoice)
	}
	if store.createInput.Item.Description != "Cloud VPS - Starter" ||
		store.createInput.Item.LineTotalMinor != 1800 ||
		store.createInput.IdempotencyKey != IdempotencyKey("key-1") {
		t.Fatalf("unexpected item/idempotency input: %+v", store.createInput)
	}
	if detail.Invoice.ID != InvoiceID("invoice-1") {
		t.Fatalf("expected returned invoice detail, got %+v", detail)
	}
}

func TestGenerateInvoiceRejectsUnpaidOrder(t *testing.T) {
	sourceOrder := paidOrder()
	sourceOrder.BillingStatus = order.BillingStatusUnpaid
	store := &fakeInvoiceStore{}
	service := NewServiceWithOrderReader(store, &fakeOrderReader{record: sourceOrder})

	_, err := service.GenerateInvoiceForOrder(context.Background(), GenerateInvoiceInput{
		TenantID:       tenant.ID("tenant-1"),
		OrderID:        order.OrderID("order-1"),
		IdempotencyKey: IdempotencyKey("key-1"),
	})
	if !errors.Is(err, ErrOrderNotPaid) {
		t.Fatalf("expected not paid error, got %v", err)
	}
	if store.createCalls != 0 {
		t.Fatalf("expected no invoice create, got %d", store.createCalls)
	}
}

func TestGenerateInvoiceReturnsExistingDuplicate(t *testing.T) {
	existing := InvoiceDetail{Invoice: Invoice{ID: "invoice-existing", TenantID: "tenant-1", OrderID: "order-1"}}
	store := &fakeInvoiceStore{detail: existing}
	service := NewServiceWithOrderReader(store, &fakeOrderReader{record: paidOrder()})

	detail, err := service.GenerateInvoiceForOrder(context.Background(), GenerateInvoiceInput{
		TenantID:       tenant.ID("tenant-1"),
		OrderID:        order.OrderID("order-1"),
		IdempotencyKey: IdempotencyKey("key-1"),
	})
	if err != nil {
		t.Fatalf("expected existing invoice detail: %v", err)
	}
	if detail.Invoice.ID != InvoiceID("invoice-existing") {
		t.Fatalf("expected existing invoice, got %+v", detail.Invoice)
	}
}

func TestGenerateInvoiceRejectsCrossTenantOrder(t *testing.T) {
	sourceOrder := paidOrder()
	sourceOrder.TenantID = tenant.ID("tenant-2")
	service := NewServiceWithOrderReader(&fakeInvoiceStore{}, &fakeOrderReader{record: sourceOrder})

	_, err := service.GenerateInvoiceForOrder(context.Background(), GenerateInvoiceInput{
		TenantID:       tenant.ID("tenant-1"),
		OrderID:        order.OrderID("order-1"),
		IdempotencyKey: IdempotencyKey("key-1"),
	})
	if !errors.Is(err, tenant.ErrAccessDenied) {
		t.Fatalf("expected cross-tenant error, got %v", err)
	}
}

func TestGenerateInvoiceRequiresIdempotencyKey(t *testing.T) {
	service := NewServiceWithOrderReader(&fakeInvoiceStore{}, &fakeOrderReader{record: paidOrder()})

	_, err := service.GenerateInvoiceForOrder(context.Background(), GenerateInvoiceInput{
		TenantID: tenant.ID("tenant-1"),
		OrderID:  order.OrderID("order-1"),
	})
	if !errors.Is(err, ErrIdempotencyKeyMissing) {
		t.Fatalf("expected idempotency error, got %v", err)
	}
}

func paidOrder() order.Order {
	return order.Order{
		ID:              "order-1",
		DisplayID:       60001,
		TenantID:        "tenant-1",
		BuyerUserID:     identity.UserID("buyer-1"),
		Quantity:        2,
		Currency:        "USD",
		UnitPriceMinor:  1000,
		DiscountMinor:   200,
		TotalMinor:      1800,
		OrderStatus:     order.OrderStatusPaid,
		BillingStatus:   order.BillingStatusPaid,
		ProductSnapshot: json.RawMessage(`{"name":"Cloud VPS"}`),
		PlanSnapshot:    json.RawMessage(`{"name":"Starter"}`),
		PriceSnapshot:   json.RawMessage(`{"cycle":"monthly"}`),
	}
}

type fakeOrderReader struct {
	record order.Order
	lookup order.OrderLookup
	err    error
}

func (reader *fakeOrderReader) GetOrder(ctx context.Context, lookup order.OrderLookup) (order.Order, error) {
	reader.lookup = lookup
	return reader.record, reader.err
}

type fakeInvoiceStore struct {
	detail      InvoiceDetail
	createInput CreateInvoiceFromOrderInput
	createCalls int
}

func (store *fakeInvoiceStore) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error) {
	return nil, nil
}

func (store *fakeInvoiceStore) GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error) {
	return InvoiceDetail{}, nil
}

func (store *fakeInvoiceStore) CreateInvoiceFromOrder(ctx context.Context, input CreateInvoiceFromOrderInput) (InvoiceDetail, error) {
	store.createCalls++
	store.createInput = input.Normalize()
	if store.detail.Invoice.ID != "" {
		return store.detail, nil
	}
	return InvoiceDetail{Invoice: Invoice{ID: "invoice-1", TenantID: input.Invoice.TenantID, OrderID: input.Invoice.OrderID}}, nil
}

func (store *fakeInvoiceStore) MarkInvoicePaid(ctx context.Context, input MarkInvoicePaidInput) (InvoiceDetail, error) {
	return store.detail, nil
}
