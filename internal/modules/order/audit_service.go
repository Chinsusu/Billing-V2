package order

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
)

type serviceLifecycleAuditStatus struct {
	Status           ServiceStatus    `json:"status"`
	BillingStatus    BillingStatus    `json:"billing_status,omitempty"`
	SuspensionReason SuspensionReason `json:"suspension_reason,omitempty"`
	TermEnd          time.Time        `json:"term_end,omitempty"`
}

type serviceLifecycleAuditMetadata struct {
	Action         ServiceLifecycleAction `json:"action"`
	DisplayID      int64                  `json:"display_id"`
	OrderID        OrderID                `json:"order_id"`
	OrderDisplayID int64                  `json:"order_display_id,omitempty"`
	Reason         string                 `json:"reason,omitempty"`
}

func (service *Service) appendServiceLifecycleAudit(ctx context.Context, input TransitionServiceLifecycleInput, record ServiceInstance) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   record.TenantID,
		ActorID:    input.ActorID,
		ActorType:  input.ActorType,
		Action:     serviceLifecycleEventType(input.Action),
		TargetType: "service",
		TargetID:   audit.TargetID(record.ID),
		BeforeSnapshotRedacted: serviceLifecycleAuditJSON(serviceLifecycleAuditStatus{
			Status: input.FromStatus,
		}),
		AfterSnapshotRedacted: serviceLifecycleAuditJSON(serviceLifecycleAuditStatus{
			Status:           record.Status,
			BillingStatus:    record.BillingStatus,
			SuspensionReason: record.SuspensionReason,
			TermEnd:          record.TermEnd,
		}),
		MetadataRedacted: serviceLifecycleAuditJSON(serviceLifecycleAuditMetadata{
			Action:         input.Action,
			DisplayID:      record.DisplayID,
			OrderID:        record.OrderID,
			OrderDisplayID: record.OrderDisplayID,
			Reason:         input.Reason,
		}),
		CorrelationID: audit.CorrelationID(record.ID),
	})
	return err
}

func serviceLifecycleAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
