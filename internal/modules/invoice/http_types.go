package invoice

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type invoiceResponse struct {
	ID            InvoiceID       `json:"id"`
	DisplayID     int64           `json:"display_id"`
	TenantID      tenant.ID       `json:"tenant_id"`
	BuyerUserID   identity.UserID `json:"buyer_user_id"`
	OrderID       order.OrderID   `json:"order_id,omitempty"`
	Status        Status          `json:"status"`
	Currency      string          `json:"currency"`
	SubtotalMinor int64           `json:"subtotal_minor"`
	TaxMinor      int64           `json:"tax_minor"`
	DiscountMinor int64           `json:"discount_minor"`
	TotalMinor    int64           `json:"total_minor"`
	IssuedAt      *time.Time      `json:"issued_at,omitempty"`
	DueAt         *time.Time      `json:"due_at,omitempty"`
	PaidAt        *time.Time      `json:"paid_at,omitempty"`
	VoidedAt      *time.Time      `json:"voided_at,omitempty"`
	Metadata      json.RawMessage `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type invoiceDetailResponse struct {
	invoiceResponse
	Items []invoiceItemResponse `json:"items"`
}

type invoiceItemResponse struct {
	ID             InvoiceItemID   `json:"id"`
	InvoiceID      InvoiceID       `json:"invoice_id"`
	TenantID       tenant.ID       `json:"tenant_id"`
	OrderID        order.OrderID   `json:"order_id,omitempty"`
	OrderItemID    OrderItemID     `json:"order_item_id,omitempty"`
	ServiceID      order.ServiceID `json:"service_id,omitempty"`
	Description    string          `json:"description"`
	Quantity       int             `json:"quantity"`
	UnitPriceMinor int64           `json:"unit_price_minor"`
	TaxMinor       int64           `json:"tax_minor"`
	DiscountMinor  int64           `json:"discount_minor"`
	LineTotalMinor int64           `json:"line_total_minor"`
	Metadata       json.RawMessage `json:"metadata"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func newInvoiceResponse(invoice Invoice) invoiceResponse {
	return invoiceResponse{
		ID:            invoice.ID,
		DisplayID:     invoice.DisplayID,
		TenantID:      invoice.TenantID,
		BuyerUserID:   invoice.BuyerUserID,
		OrderID:       invoice.OrderID,
		Status:        invoice.Status,
		Currency:      invoice.Currency,
		SubtotalMinor: invoice.SubtotalMinor,
		TaxMinor:      invoice.TaxMinor,
		DiscountMinor: invoice.DiscountMinor,
		TotalMinor:    invoice.TotalMinor,
		IssuedAt:      timeIfSet(invoice.IssuedAt),
		DueAt:         timeIfSet(invoice.DueAt),
		PaidAt:        timeIfSet(invoice.PaidAt),
		VoidedAt:      timeIfSet(invoice.VoidedAt),
		Metadata:      invoice.Metadata,
		CreatedAt:     invoice.CreatedAt,
		UpdatedAt:     invoice.UpdatedAt,
	}
}

func newInvoiceResponses(invoices []Invoice) []invoiceResponse {
	responses := make([]invoiceResponse, 0, len(invoices))
	for _, invoice := range invoices {
		responses = append(responses, newInvoiceResponse(invoice))
	}
	return responses
}

func newInvoiceDetailResponse(detail InvoiceDetail) invoiceDetailResponse {
	return invoiceDetailResponse{
		invoiceResponse: newInvoiceResponse(detail.Invoice),
		Items:           newInvoiceItemResponses(detail.Items),
	}
}

func newInvoiceItemResponses(items []Item) []invoiceItemResponse {
	responses := make([]invoiceItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, newInvoiceItemResponse(item))
	}
	return responses
}

func newInvoiceItemResponse(item Item) invoiceItemResponse {
	return invoiceItemResponse{
		ID:             item.ID,
		InvoiceID:      item.InvoiceID,
		TenantID:       item.TenantID,
		OrderID:        item.OrderID,
		OrderItemID:    item.OrderItemID,
		ServiceID:      item.ServiceID,
		Description:    item.Description,
		Quantity:       item.Quantity,
		UnitPriceMinor: item.UnitPriceMinor,
		TaxMinor:       item.TaxMinor,
		DiscountMinor:  item.DiscountMinor,
		LineTotalMinor: item.LineTotalMinor,
		Metadata:       item.Metadata,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func timeIfSet(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
