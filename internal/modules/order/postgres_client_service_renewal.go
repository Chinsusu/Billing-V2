package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var errClientServiceRenewalNotFound = errors.New("client service renewal not found")

type clientServiceRenewalContext struct {
	Service           ServiceInstance
	BuyerUserID       identity.UserID
	SellingPriceMinor int64
	Currency          string
	TenantPlanStatus  catalog.TenantPlanStatus
	Cycle             ServiceRenewalCycle
}

type renewalInvoiceRef struct {
	ID        string
	DisplayID int64
}

type renewalPaymentRef struct {
	ID        string
	DisplayID int64
}

const clientServiceRenewalContextSQL = `
SELECT ` + serviceInstanceRelatedReadColumns + `,
       ord.buyer_user_id,
       tpl.selling_price_minor,
       tpl.currency,
       tpl.status,
       mp.billing_cycle_type,
       mp.billing_cycle_value
FROM service_instances svc
JOIN orders ord
  ON ord.order_id = svc.order_id
 AND ord.tenant_id = svc.tenant_id
JOIN tenant_plans tpl
  ON tpl.tenant_plan_id = svc.tenant_plan_id
 AND tpl.tenant_id = svc.tenant_id
JOIN master_plans mp
  ON mp.plan_id = tpl.master_plan_id
WHERE svc.service_instance_id = $1
  AND svc.tenant_id = $2
  AND ord.buyer_user_id = $3
FOR UPDATE OF svc`

const existingClientServiceRenewalSQL = `
SELECT ` + serviceInstanceRelatedReadColumns + `,
       inv.invoice_id::text,
       inv.display_id,
       txn.payment_transaction_id::text,
       txn.display_id,
       ledger.ledger_entry_id::text,
       ledger.display_id,
       inv.total_minor,
       inv.currency
FROM payment_transactions txn
JOIN invoices inv
  ON inv.invoice_id = txn.invoice_id
 AND inv.tenant_id = txn.tenant_id
JOIN wallet_ledger_entries ledger
  ON ledger.tenant_id = txn.tenant_id
 AND ledger.wallet_id = $5::uuid
 AND ledger.idempotency_key = $6
 AND ledger.reference_type = 'invoice'
 AND ledger.reference_id = inv.invoice_id
JOIN service_instances svc
  ON svc.service_instance_id::text = txn.metadata->>'service_id'
 AND svc.tenant_id = txn.tenant_id
JOIN orders ord
  ON ord.order_id = svc.order_id
 AND ord.tenant_id = svc.tenant_id
WHERE txn.tenant_id = $1
  AND txn.account_user_id = $2
  AND txn.idempotency_key = $3
  AND txn.metadata->>'source' = 'service_renewal'
  AND txn.metadata->>'service_id' = $4
  AND txn.metadata->>'wallet_id' = $5::text`

const paymentTransactionExistsSQL = `
SELECT EXISTS (
    SELECT 1
    FROM payment_transactions
    WHERE tenant_id = $1
      AND idempotency_key = $2
)`

const createRenewalInvoiceSQL = `
INSERT INTO invoices (tenant_id, buyer_user_id, status, currency, subtotal_minor, tax_minor, discount_minor, total_minor, issued_at, metadata)
VALUES ($1, $2, 'issued', $3, $4, 0, 0, $4, NOW(), $5::jsonb)
RETURNING invoice_id::text, display_id`

const createRenewalInvoiceItemSQL = `
INSERT INTO invoice_items (invoice_id, tenant_id, service_instance_id, description, quantity, unit_price_minor, tax_minor, discount_minor, line_total_minor, metadata)
VALUES ($1, $2, $3, $4, 1, $5, 0, 0, $5, $6::jsonb)`

const createRenewalPaymentTransactionSQL = `
INSERT INTO payment_transactions (tenant_id, account_user_id, invoice_id, transaction_type, status, currency, amount_minor, description, idempotency_key, metadata)
VALUES ($1, $2, $3, 'charge', 'posted', $4, $5, $6, $7, $8::jsonb)
ON CONFLICT (tenant_id, idempotency_key) DO NOTHING
RETURNING payment_transaction_id::text, display_id`

const markRenewalInvoicePaidSQL = `
UPDATE invoices
SET status = 'paid',
    paid_at = NOW(),
    updated_at = NOW()
WHERE invoice_id = $1
  AND tenant_id = $2
  AND status = 'issued'
RETURNING invoice_id::text`

func (store *PostgresStore) RenewClientService(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error) {
	if err := store.ready(); err != nil {
		return ClientServiceRenewal{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ClientServiceRenewal{}, err
	}
	if conn, ok := store.executor.(*sql.DB); ok {
		var result ClientServiceRenewal
		err := platformdb.WithTx(ctx, conn, func(ctx context.Context, tx *sql.Tx) error {
			var runErr error
			result, runErr = NewPostgresStore(tx).renewClientService(ctx, input)
			return runErr
		})
		if err != nil {
			return ClientServiceRenewal{}, err
		}
		return result, nil
	}
	return store.renewClientService(ctx, input)
}

func (store *PostgresStore) renewClientService(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error) {
	existing, found, err := store.loadExistingClientServiceRenewal(ctx, input)
	if err != nil || found {
		return existing, err
	}
	renewalContext, err := store.loadClientServiceRenewalContext(ctx, input)
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	if renewalContext.TenantPlanStatus != catalog.TenantPlanStatusActive {
		return ClientServiceRenewal{}, ErrServiceRenewalUnavailable
	}
	if renewalContext.Service.Status != input.FromStatus {
		return ClientServiceRenewal{}, ErrServiceStatusConflict
	}
	if renewalContext.SellingPriceMinor <= 0 {
		return ClientServiceRenewal{}, ErrAmountInvalid
	}
	newTermEnd, err := CalculateRenewedTermEnd(renewalContext.Service, renewalContext.Cycle)
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	accountWallet, err := wallet.NewService(wallet.NewPostgresStore(store.executor)).GetWallet(ctx, wallet.WalletLookup{
		ID:        input.WalletID,
		TenantID:  input.TenantID,
		OwnerType: wallet.OwnerTypeUser,
		OwnerID:   wallet.UserOwnerID(input.BuyerUserID),
	})
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	if accountWallet.Status != wallet.StatusActive {
		return ClientServiceRenewal{}, wallet.ErrWalletStatusConflict
	}
	if accountWallet.Currency != renewalContext.Currency {
		return ClientServiceRenewal{}, wallet.ErrWalletCurrencyMismatch
	}
	invoiceRef, err := store.createRenewalInvoice(ctx, input, renewalContext, newTermEnd)
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	ledger, err := wallet.NewPostgresStore(store.executor).PostLedgerEntryResult(ctx, wallet.PostLedgerEntryInput{
		WalletID:       input.WalletID,
		TenantID:       input.TenantID,
		Direction:      wallet.DirectionDebit,
		AmountMinor:    renewalContext.SellingPriceMinor,
		Currency:       renewalContext.Currency,
		EntryType:      wallet.EntryTypePurchase,
		ReferenceType:  wallet.ReferenceType("invoice"),
		ReferenceID:    wallet.ReferenceID(invoiceRef.ID),
		IdempotencyKey: renewalLedgerIdempotency(input),
		CreatedBy:      input.ActorID,
		Reason:         fmt.Sprintf("Service %d renewal invoice %d", renewalContext.Service.DisplayID, invoiceRef.DisplayID),
		CorrelationID:  wallet.CorrelationID(invoiceRef.ID),
	})
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	paymentRef, err := store.createRenewalPaymentTransaction(ctx, input, renewalContext, invoiceRef, ledger.Entry.ID, newTermEnd)
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	if err := store.markRenewalInvoicePaid(ctx, input.TenantID, invoiceRef.ID); err != nil {
		return ClientServiceRenewal{}, err
	}
	updatedService, err := store.TransitionServiceLifecycle(ctx, TransitionServiceLifecycleInput{
		ID:                       input.ServiceID,
		TenantID:                 input.TenantID,
		BuyerUserID:              input.BuyerUserID,
		ActorID:                  audit.ActorID(input.ActorID),
		ActorType:                audit.ActorTypeUser,
		Action:                   ServiceLifecycleActionRenew,
		FromStatus:               input.FromStatus,
		ToStatus:                 ServiceStatusActive,
		BillingStatus:            BillingStatusPaid,
		Reason:                   input.Reason,
		TermEnd:                  newTermEnd,
		ExpectedTermEnd:          renewalContext.Service.TermEnd,
		ExpectedBillingStatus:    renewalContext.Service.BillingStatus,
		ExpectedSuspensionReason: renewalContext.Service.SuspensionReason,
	})
	if err != nil {
		return ClientServiceRenewal{}, err
	}
	return ClientServiceRenewal{
		Service:                   updatedService,
		InvoiceID:                 invoiceRef.ID,
		InvoiceDisplayID:          invoiceRef.DisplayID,
		PaymentTransactionID:      paymentRef.ID,
		PaymentTransactionDisplay: paymentRef.DisplayID,
		WalletID:                  input.WalletID,
		LedgerEntryID:             ledger.Entry.ID,
		LedgerEntryDisplayID:      ledger.Entry.DisplayID,
		AmountMinor:               renewalContext.SellingPriceMinor,
		Currency:                  renewalContext.Currency,
		Renewed:                   true,
		PreviousStatus:            renewalContext.Service.Status,
		PreviousTermEnd:           renewalContext.Service.TermEnd,
	}, nil
}

func (store *PostgresStore) loadClientServiceRenewalContext(ctx context.Context, input ClientServiceRenewalInput) (clientServiceRenewalContext, error) {
	return scanClientServiceRenewalContext(store.executor.QueryRowContext(
		ctx,
		clientServiceRenewalContextSQL,
		input.ServiceID,
		input.TenantID,
		input.BuyerUserID,
	))
}

func (store *PostgresStore) loadExistingClientServiceRenewal(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, bool, error) {
	result, err := scanExistingClientServiceRenewal(store.executor.QueryRowContext(
		ctx,
		existingClientServiceRenewalSQL,
		input.TenantID,
		input.BuyerUserID,
		input.IdempotencyKey,
		string(input.ServiceID),
		string(input.WalletID),
		string(renewalLedgerIdempotency(input)),
	))
	if err == nil {
		result.WalletID = input.WalletID
		return result, true, nil
	}
	if !errors.Is(err, errClientServiceRenewalNotFound) {
		return ClientServiceRenewal{}, false, err
	}
	exists, err := store.paymentTransactionExists(ctx, input.TenantID, input.IdempotencyKey)
	if err != nil {
		return ClientServiceRenewal{}, false, err
	}
	if exists {
		return ClientServiceRenewal{}, false, ErrIdempotencyConflict
	}
	return ClientServiceRenewal{}, false, nil
}

func (store *PostgresStore) paymentTransactionExists(ctx context.Context, tenantID tenant.ID, idempotencyKey IdempotencyKey) (bool, error) {
	var exists bool
	if err := store.executor.QueryRowContext(ctx, paymentTransactionExistsSQL, tenantID, idempotencyKey).Scan(&exists); err != nil {
		return false, fmt.Errorf("check renewal payment idempotency: %w", err)
	}
	return exists, nil
}

func (store *PostgresStore) createRenewalInvoice(ctx context.Context, input ClientServiceRenewalInput, renewalContext clientServiceRenewalContext, termEnd time.Time) (renewalInvoiceRef, error) {
	metadata := renewalMetadata(input, renewalContext.Service, termEnd)
	var invoice renewalInvoiceRef
	if err := store.executor.QueryRowContext(
		ctx,
		createRenewalInvoiceSQL,
		input.TenantID,
		input.BuyerUserID,
		renewalContext.Currency,
		renewalContext.SellingPriceMinor,
		string(metadata),
	).Scan(&invoice.ID, &invoice.DisplayID); err != nil {
		return renewalInvoiceRef{}, fmt.Errorf("create renewal invoice: %w", err)
	}
	if _, err := store.executor.ExecContext(
		ctx,
		createRenewalInvoiceItemSQL,
		invoice.ID,
		input.TenantID,
		input.ServiceID,
		fmt.Sprintf("Service %d renewal", renewalContext.Service.DisplayID),
		renewalContext.SellingPriceMinor,
		string(metadata),
	); err != nil {
		return renewalInvoiceRef{}, fmt.Errorf("create renewal invoice item: %w", err)
	}
	return invoice, nil
}

func (store *PostgresStore) createRenewalPaymentTransaction(
	ctx context.Context,
	input ClientServiceRenewalInput,
	renewalContext clientServiceRenewalContext,
	invoice renewalInvoiceRef,
	ledgerEntryID wallet.LedgerEntryID,
	termEnd time.Time,
) (renewalPaymentRef, error) {
	metadata := renewalPaymentMetadata(input, ledgerEntryID, termEnd)
	var payment renewalPaymentRef
	err := store.executor.QueryRowContext(
		ctx,
		createRenewalPaymentTransactionSQL,
		input.TenantID,
		input.BuyerUserID,
		invoice.ID,
		renewalContext.Currency,
		renewalContext.SellingPriceMinor,
		fmt.Sprintf("Service %d renewal invoice %d", renewalContext.Service.DisplayID, invoice.DisplayID),
		input.IdempotencyKey,
		string(metadata),
	).Scan(&payment.ID, &payment.DisplayID)
	if err == nil {
		return payment, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return renewalPaymentRef{}, fmt.Errorf("create renewal payment transaction: %w", err)
	}
	existing, found, loadErr := store.loadExistingClientServiceRenewal(ctx, input)
	if loadErr != nil {
		return renewalPaymentRef{}, loadErr
	}
	if !found || existing.InvoiceID != invoice.ID {
		return renewalPaymentRef{}, ErrIdempotencyConflict
	}
	return renewalPaymentRef{ID: existing.PaymentTransactionID, DisplayID: existing.PaymentTransactionDisplay}, nil
}

func (store *PostgresStore) markRenewalInvoicePaid(ctx context.Context, tenantID tenant.ID, invoiceID string) error {
	var paidInvoiceID string
	if err := store.executor.QueryRowContext(ctx, markRenewalInvoicePaidSQL, invoiceID, tenantID).Scan(&paidInvoiceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderStatusConflict
		}
		return fmt.Errorf("mark renewal invoice paid: %w", err)
	}
	return nil
}

func scanClientServiceRenewalContext(row orderScanner) (clientServiceRenewalContext, error) {
	var renewalContext clientServiceRenewalContext
	var id, tenantID, orderID, tenantPlanID, providerSourceID, externalResourceID, status, billingStatus string
	var suspensionReason sql.NullString
	var orderDisplayID, buyerDisplayID, providerSourceDisplayID sql.NullInt64
	var buyerUserID, planStatus, cycleType string
	if err := row.Scan(
		&id, &renewalContext.Service.DisplayID, &tenantID, &orderID, &tenantPlanID, &providerSourceID, &externalResourceID,
		&status, &billingStatus, &suspensionReason, &renewalContext.Service.TermStart, &renewalContext.Service.TermEnd,
		&renewalContext.Service.CreatedAt, &renewalContext.Service.UpdatedAt, &orderDisplayID, &buyerDisplayID, &providerSourceDisplayID,
		&buyerUserID, &renewalContext.SellingPriceMinor, &renewalContext.Currency, &planStatus, &cycleType, &renewalContext.Cycle.Value,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return clientServiceRenewalContext{}, ErrServiceNotFound
		}
		return clientServiceRenewalContext{}, fmt.Errorf("scan client service renewal context: %w", err)
	}
	renewalContext.Service.ID = ServiceID(id)
	renewalContext.Service.TenantID = tenant.ID(tenantID)
	renewalContext.Service.OrderID = OrderID(orderID)
	renewalContext.Service.OrderDisplayID = orderDisplayID.Int64
	renewalContext.Service.BuyerDisplayID = buyerDisplayID.Int64
	renewalContext.Service.TenantPlanID = catalog.TenantPlanID(tenantPlanID)
	renewalContext.Service.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	renewalContext.Service.ProviderSourceDisplayID = providerSourceDisplayID.Int64
	renewalContext.Service.ExternalResourceID = provider.ExternalResourceID(externalResourceID)
	renewalContext.Service.Status = ServiceStatus(status)
	renewalContext.Service.BillingStatus = BillingStatus(billingStatus)
	renewalContext.Service.SuspensionReason = SuspensionReason(suspensionReason.String)
	renewalContext.BuyerUserID = identity.UserID(buyerUserID)
	renewalContext.TenantPlanStatus = catalog.TenantPlanStatus(planStatus)
	renewalContext.Cycle.Type = catalog.BillingCycleType(cycleType)
	return renewalContext, nil
}

func scanExistingClientServiceRenewal(row orderScanner) (ClientServiceRenewal, error) {
	var result ClientServiceRenewal
	var id, tenantID, orderID, tenantPlanID, providerSourceID, externalResourceID, status, billingStatus string
	var suspensionReason sql.NullString
	var orderDisplayID, buyerDisplayID, providerSourceDisplayID sql.NullInt64
	if err := row.Scan(
		&id, &result.Service.DisplayID, &tenantID, &orderID, &tenantPlanID, &providerSourceID, &externalResourceID,
		&status, &billingStatus, &suspensionReason, &result.Service.TermStart, &result.Service.TermEnd,
		&result.Service.CreatedAt, &result.Service.UpdatedAt, &orderDisplayID, &buyerDisplayID, &providerSourceDisplayID,
		&result.InvoiceID, &result.InvoiceDisplayID, &result.PaymentTransactionID, &result.PaymentTransactionDisplay,
		&result.LedgerEntryID, &result.LedgerEntryDisplayID, &result.AmountMinor, &result.Currency,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ClientServiceRenewal{}, errClientServiceRenewalNotFound
		}
		return ClientServiceRenewal{}, fmt.Errorf("scan existing client service renewal: %w", err)
	}
	result.Service.ID = ServiceID(id)
	result.Service.TenantID = tenant.ID(tenantID)
	result.Service.OrderID = OrderID(orderID)
	result.Service.OrderDisplayID = orderDisplayID.Int64
	result.Service.BuyerDisplayID = buyerDisplayID.Int64
	result.Service.TenantPlanID = catalog.TenantPlanID(tenantPlanID)
	result.Service.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	result.Service.ProviderSourceDisplayID = providerSourceDisplayID.Int64
	result.Service.ExternalResourceID = provider.ExternalResourceID(externalResourceID)
	result.Service.Status = ServiceStatus(status)
	result.Service.BillingStatus = BillingStatus(billingStatus)
	result.Service.SuspensionReason = SuspensionReason(suspensionReason.String)
	result.Renewed = false
	return result, nil
}

func renewalLedgerIdempotency(input ClientServiceRenewalInput) wallet.IdempotencyKey {
	return wallet.IdempotencyKey(fmt.Sprintf("service-renewal:%s:%s", input.ServiceID, input.IdempotencyKey))
}

func renewalMetadata(input ClientServiceRenewalInput, service ServiceInstance, termEnd time.Time) json.RawMessage {
	payload := map[string]interface{}{
		"source":          "service_renewal",
		"service_id":      input.ServiceID,
		"service_display": service.DisplayID,
		"from_status":     input.FromStatus,
		"term_end":        termEnd.UTC().Format(time.RFC3339),
	}
	return mustRenewalJSON(payload)
}

func renewalPaymentMetadata(input ClientServiceRenewalInput, ledgerEntryID wallet.LedgerEntryID, termEnd time.Time) json.RawMessage {
	payload := map[string]interface{}{
		"source":          "service_renewal",
		"service_id":      input.ServiceID,
		"wallet_id":       input.WalletID,
		"ledger_entry_id": ledgerEntryID,
		"term_end":        termEnd.UTC().Format(time.RFC3339),
	}
	return mustRenewalJSON(payload)
}

func mustRenewalJSON(payload interface{}) json.RawMessage {
	data, err := json.Marshal(payload)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
