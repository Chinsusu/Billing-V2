package order

import (
	"context"
	"encoding/json"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
)

const orderAuditActionStatusChanged = "order.status_changed"

type AuditAppender interface {
	Append(ctx context.Context, input audit.AppendInput) (audit.Log, error)
}

func (service *Service) appendOrderStatusAudit(ctx context.Context, input TransitionOrderStatusInput, record Order) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:               record.TenantID,
		ActorID:                audit.ActorID(input.ActorID),
		ActorType:              audit.ActorTypeUser,
		Action:                 orderAuditActionStatusChanged,
		TargetType:             "order",
		TargetID:               audit.TargetID(record.ID),
		BeforeSnapshotRedacted: orderAuditJSON(orderAuditStatus{OrderStatus: input.FromStatus}),
		AfterSnapshotRedacted: orderAuditJSON(orderAuditStatus{
			OrderStatus:   record.OrderStatus,
			BillingStatus: record.BillingStatus,
		}),
		MetadataRedacted: orderAuditJSON(orderAuditMetadata{
			DisplayID:   record.DisplayID,
			BuyerUserID: record.BuyerUserID,
			TotalMinor:  record.TotalMinor,
			Currency:    record.Currency,
		}),
		CorrelationID: audit.CorrelationID(record.ID),
	})
	return err
}

type orderAuditStatus struct {
	OrderStatus   OrderStatus   `json:"order_status"`
	BillingStatus BillingStatus `json:"billing_status,omitempty"`
}

type orderAuditMetadata struct {
	DisplayID   int64           `json:"display_id"`
	BuyerUserID identity.UserID `json:"buyer_user_id"`
	TotalMinor  int64           `json:"total_minor"`
	Currency    string          `json:"currency"`
}

func orderAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
