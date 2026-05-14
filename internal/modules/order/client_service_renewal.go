package order

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const clientServiceRenewalDefaultReason = "Client service renewal"

type ClientServiceRenewalStore interface {
	RenewClientService(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error)
}

type ClientServiceRenewalInput struct {
	TenantID       tenant.ID
	BuyerUserID    identity.UserID
	ServiceID      ServiceID
	WalletID       wallet.WalletID
	ActorID        identity.UserID
	FromStatus     ServiceStatus
	Reason         string
	IdempotencyKey IdempotencyKey
}

type ClientServiceRenewal struct {
	Service                   ServiceInstance
	InvoiceID                 string
	InvoiceDisplayID          int64
	PaymentTransactionID      string
	PaymentTransactionDisplay int64
	WalletID                  wallet.WalletID
	LedgerEntryID             wallet.LedgerEntryID
	LedgerEntryDisplayID      int64
	AmountMinor               int64
	Currency                  string
	Renewed                   bool
	PreviousStatus            ServiceStatus
	PreviousTermEnd           time.Time
}

func (input ClientServiceRenewalInput) Normalize() ClientServiceRenewalInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.BuyerUserID = identity.UserID(trim(string(output.BuyerUserID)))
	output.ServiceID = ServiceID(trim(string(output.ServiceID)))
	output.WalletID = wallet.WalletID(trim(string(output.WalletID)))
	output.ActorID = identity.UserID(trim(string(output.ActorID)))
	output.FromStatus = ServiceStatus(trim(string(output.FromStatus)))
	output.Reason = trim(output.Reason)
	if output.Reason == "" {
		output.Reason = clientServiceRenewalDefaultReason
	}
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
}

func (input ClientServiceRenewalInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.BuyerUserID == "" {
		return ErrBuyerIDMissing
	}
	if input.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if input.WalletID.Empty() {
		return wallet.ErrWalletIDMissing
	}
	if input.ActorID == "" {
		return identity.ErrActorIDMissing
	}
	if !input.FromStatus.Valid() {
		return ErrServiceStatusInvalid
	}
	if input.Reason == "" {
		return ErrServiceLifecycleReasonMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}

func (service *Service) RenewClientService(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error) {
	if err := service.ready(); err != nil {
		return ClientServiceRenewal{}, err
	}
	store, ok := service.store.(ClientServiceRenewalStore)
	if !ok {
		return ClientServiceRenewal{}, ErrServiceStoreMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ClientServiceRenewal{}, err
	}
	result, err := store.RenewClientService(ctx, input)
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	if result.Renewed {
		if err := service.appendClientRenewalAudits(ctx, input, result); err != nil {
			return ClientServiceRenewal{}, err
		}
	}
	return result, nil
}

func (service *Service) appendClientRenewalAudits(ctx context.Context, input ClientServiceRenewalInput, result ClientServiceRenewal) error {
	if service.audit == nil {
		return nil
	}
	if err := service.appendServiceLifecycleAudit(ctx, TransitionServiceLifecycleInput{
		ID:          input.ServiceID,
		TenantID:    input.TenantID,
		BuyerUserID: input.BuyerUserID,
		ActorID:     audit.ActorID(input.ActorID),
		ActorType:   audit.ActorTypeUser,
		Action:      ServiceLifecycleActionRenew,
		FromStatus:  input.FromStatus,
		ToStatus:    ServiceStatusActive,
		Reason:      input.Reason,
	}, result.Service); err != nil {
		return err
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   input.TenantID,
		ActorID:    audit.ActorID(input.ActorID),
		ActorType:  audit.ActorTypeUser,
		Action:     "invoice.wallet_paid",
		TargetType: "invoice",
		TargetID:   audit.TargetID(result.InvoiceID),
		BeforeSnapshotRedacted: serviceLifecycleAuditJSON(map[string]string{
			"status": "issued",
		}),
		AfterSnapshotRedacted: serviceLifecycleAuditJSON(map[string]string{
			"status": "paid",
		}),
		MetadataRedacted: clientRenewalPaymentAuditJSON(result),
		CorrelationID:    audit.CorrelationID(result.InvoiceID),
	})
	return err
}

func clientRenewalPaymentAuditJSON(result ClientServiceRenewal) json.RawMessage {
	payload := struct {
		InvoiceDisplayID     int64                `json:"invoice_display_id"`
		PaymentTransactionID string               `json:"payment_transaction_id"`
		TransactionDisplayID int64                `json:"transaction_display_id"`
		WalletID             wallet.WalletID      `json:"wallet_id"`
		LedgerEntryID        wallet.LedgerEntryID `json:"ledger_entry_id"`
		AmountMinor          int64                `json:"amount_minor"`
		Currency             string               `json:"currency"`
	}{
		InvoiceDisplayID:     result.InvoiceDisplayID,
		PaymentTransactionID: result.PaymentTransactionID,
		TransactionDisplayID: result.PaymentTransactionDisplay,
		WalletID:             result.WalletID,
		LedgerEntryID:        result.LedgerEntryID,
		AmountMinor:          result.AmountMinor,
		Currency:             result.Currency,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
