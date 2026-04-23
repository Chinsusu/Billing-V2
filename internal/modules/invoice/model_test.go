package invoice

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateInvoiceInputNormalizeValidate(t *testing.T) {
	input := CreateInvoiceInput{
		TenantID:      tenant.ID("tenant-1"),
		BuyerUserID:   identity.UserID("buyer-1"),
		Currency:      " usd ",
		SubtotalMinor: 1000,
		TaxMinor:      100,
		DiscountMinor: 50,
		TotalMinor:    1050,
		Metadata:      json.RawMessage(`{"source":"order"}`),
	}.Normalize()

	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.Status != StatusDraft {
		t.Fatalf("expected draft status, got %q", input.Status)
	}
	if string(input.Metadata) != `{"source":"order"}` {
		t.Fatalf("expected metadata to be preserved, got %s", input.Metadata)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid invoice input: %v", err)
	}

	input.TotalMinor = 999
	if err := input.Validate(); !errors.Is(err, ErrTotalInvalid) {
		t.Fatalf("expected total error, got %v", err)
	}
}

func TestCreateInvoiceInputRequiresBuyer(t *testing.T) {
	err := CreateInvoiceInput{
		TenantID:      tenant.ID("tenant-1"),
		Currency:      "USD",
		SubtotalMinor: 1000,
		TotalMinor:    1000,
	}.Normalize().Validate()
	if !errors.Is(err, ErrBuyerIDMissing) {
		t.Fatalf("expected buyer error, got %v", err)
	}
}

func TestCreateItemInputNormalizeValidate(t *testing.T) {
	input := CreateItemInput{
		InvoiceID:      InvoiceID("invoice-1"),
		TenantID:       tenant.ID("tenant-1"),
		Description:    "  VPS monthly plan  ",
		UnitPriceMinor: 2000,
		TaxMinor:       200,
		DiscountMinor:  100,
		LineTotalMinor: 2100,
	}.Normalize()

	if input.Quantity != 1 {
		t.Fatalf("expected default quantity 1, got %d", input.Quantity)
	}
	if input.Description != "VPS monthly plan" {
		t.Fatalf("expected trimmed description, got %q", input.Description)
	}
	if string(input.Metadata) != "{}" {
		t.Fatalf("expected default metadata, got %s", input.Metadata)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid invoice item input: %v", err)
	}

	input.Quantity = -1
	if err := input.Validate(); !errors.Is(err, ErrQuantityInvalid) {
		t.Fatalf("expected quantity error, got %v", err)
	}
}

func TestInvoiceTransitionGuards(t *testing.T) {
	if !CanTransition(StatusDraft, StatusIssued) {
		t.Fatal("expected draft invoice to become issued")
	}
	if !CanTransition(StatusIssued, StatusPaid) {
		t.Fatal("expected issued invoice to become paid")
	}
	if CanTransition(StatusVoided, StatusPaid) {
		t.Fatal("voided invoice should not become paid")
	}
}
