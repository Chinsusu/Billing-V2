package checkout

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type invoiceResponse struct {
	ID            invoice.InvoiceID   `json:"id"`
	DisplayID     int64               `json:"display_id"`
	TenantID      tenant.ID           `json:"tenant_id"`
	BuyerUserID   identity.UserID     `json:"buyer_user_id"`
	OrderID       order.OrderID       `json:"order_id,omitempty"`
	Status        invoice.Status      `json:"status"`
	Currency      string              `json:"currency"`
	SubtotalMinor int64               `json:"subtotal_minor"`
	TaxMinor      int64               `json:"tax_minor"`
	DiscountMinor int64               `json:"discount_minor"`
	TotalMinor    int64               `json:"total_minor"`
	IssuedAt      *time.Time          `json:"issued_at,omitempty"`
	DueAt         *time.Time          `json:"due_at,omitempty"`
	PaidAt        *time.Time          `json:"paid_at,omitempty"`
	Metadata      json.RawMessage     `json:"metadata"`
	Items         []invoiceItemRecord `json:"items"`
}

type invoiceItemRecord struct {
	ID             invoice.InvoiceItemID `json:"id"`
	Description    string                `json:"description"`
	Quantity       int                   `json:"quantity"`
	UnitPriceMinor int64                 `json:"unit_price_minor"`
	LineTotalMinor int64                 `json:"line_total_minor"`
}

func newInvoiceDetailResponse(detail invoice.InvoiceDetail) invoiceResponse {
	items := make([]invoiceItemRecord, 0, len(detail.Items))
	for _, item := range detail.Items {
		items = append(items, invoiceItemRecord{
			ID:             item.ID,
			Description:    item.Description,
			Quantity:       item.Quantity,
			UnitPriceMinor: item.UnitPriceMinor,
			LineTotalMinor: item.LineTotalMinor,
		})
	}
	return invoiceResponse{
		ID:            detail.Invoice.ID,
		DisplayID:     detail.Invoice.DisplayID,
		TenantID:      detail.Invoice.TenantID,
		BuyerUserID:   detail.Invoice.BuyerUserID,
		OrderID:       detail.Invoice.OrderID,
		Status:        detail.Invoice.Status,
		Currency:      detail.Invoice.Currency,
		SubtotalMinor: detail.Invoice.SubtotalMinor,
		TaxMinor:      detail.Invoice.TaxMinor,
		DiscountMinor: detail.Invoice.DiscountMinor,
		TotalMinor:    detail.Invoice.TotalMinor,
		IssuedAt:      timeIfSet(detail.Invoice.IssuedAt),
		DueAt:         timeIfSet(detail.Invoice.DueAt),
		PaidAt:        timeIfSet(detail.Invoice.PaidAt),
		Metadata:      detail.Invoice.Metadata,
		Items:         items,
	}
}

func timeIfSet(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
