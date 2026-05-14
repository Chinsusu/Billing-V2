package order

import (
	"errors"
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func writeOrderError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if field, ok := orderValidationField(err); ok {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{field})
		return
	}
	switch {
	case errors.Is(err, ErrOrderNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "order.not_found", "Order was not found.")
	case errors.Is(err, ErrServiceNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "service.not_found", "Service instance was not found.")
	case errors.Is(err, ErrCredentialNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "credential.not_found", "Credential was not found.")
	case errors.Is(err, wallet.ErrWalletNotFound):
		httpserver.WriteError(w, r, http.StatusNotFound, "wallet.not_found", "Wallet was not found.")
	case errors.Is(err, ErrCredentialRevealRateLimited):
		httpserver.WriteError(w, r, http.StatusTooManyRequests, "credential.reveal_rate_limited", "Credential reveal limit was reached.")
	case errors.Is(err, ErrCredentialStatusInvalid):
		httpserver.WriteError(w, r, http.StatusForbidden, "credential.reveal_denied", "Credential cannot be revealed.")
	case errors.Is(err, ErrOrderStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "order.status_conflict", "Order status no longer matches the expected value.")
	case errors.Is(err, ErrServiceStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "service.status_conflict", "Service status no longer matches the expected value.")
	case errors.Is(err, ErrServiceRenewalUnavailable):
		httpserver.WriteError(w, r, http.StatusConflict, "service.renewal_unavailable", "Service renewal is unavailable.")
	case errors.Is(err, ErrIdempotencyConflict), errors.Is(err, wallet.ErrIdempotencyConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "order.idempotency_conflict", "Idempotency key conflicts with a different request.")
	case errors.Is(err, wallet.ErrInsufficientBalance):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.insufficient_balance", "Wallet balance is insufficient.")
	case errors.Is(err, wallet.ErrWalletStatusConflict):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.status_conflict", "Wallet status does not allow this operation.")
	case errors.Is(err, wallet.ErrWalletCurrencyMismatch):
		httpserver.WriteError(w, r, http.StatusConflict, "wallet.currency_mismatch", "Wallet currency does not match the charge currency.")
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing),
		errors.Is(err, audit.ErrActorIDMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, ErrServiceStoreMissing),
		errors.Is(err, ErrStoreExecutorMissing),
		errors.Is(err, ErrCredentialStoreMissing),
		errors.Is(err, ErrCredentialCipherMissing),
		errors.Is(err, ErrCredentialRevealLimiterMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "order.service_unavailable", "Order service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "order.operation_failed", "Order operation failed.")
	}
}

func orderValidationField(err error) (httpserver.ValidationField, bool) {
	switch {
	case errors.Is(err, tenant.ErrTenantIDMissing), errors.Is(err, tenant.ErrContextMissing):
		return validationField("tenant_id", "tenant.context_missing", "Tenant context is required."), true
	case errors.Is(err, tenant.ErrTenantMismatch), errors.Is(err, tenant.ErrAccessDenied):
		return validationField("tenant_id", "tenant.context_invalid", "Tenant context is invalid."), true
	case errors.Is(err, ErrBuyerIDMissing):
		return validationField("actor_id", "order.buyer_missing", "Buyer actor is required."), true
	case errors.Is(err, ErrOrderIDMissing):
		return validationField("order_id", "order.order_id_missing", "Order id is required."), true
	case errors.Is(err, ErrServiceIDMissing):
		return validationField("service_id", "service.service_id_missing", "Service id is required."), true
	case errors.Is(err, wallet.ErrWalletIDMissing):
		return validationField("wallet_id", "wallet.wallet_id_missing", "Wallet id is required."), true
	case errors.Is(err, ErrCredentialIDMissing):
		return validationField("credential_id", "credential.credential_id_missing", "Credential id is required."), true
	case errors.Is(err, ErrTenantPlanIDMissing):
		return validationField("tenant_plan_id", "order.tenant_plan_id_missing", "Tenant plan id is required."), true
	case errors.Is(err, ErrIdempotencyKeyMissing):
		return validationField("idempotency_key", "order.idempotency_key_missing", "Idempotency key is required."), true
	case errors.Is(err, ErrCurrencyMissing):
		return validationField("currency", "order.currency_missing", "Currency is required."), true
	case errors.Is(err, ErrCurrencyInvalid):
		return validationField("currency", "order.currency_invalid", "Currency is invalid."), true
	case errors.Is(err, ErrAmountInvalid):
		return validationField("amount_minor", "order.amount_invalid", "Money amount must not be negative."), true
	case errors.Is(err, ErrQuantityInvalid):
		return validationField("quantity", "order.quantity_invalid", "Quantity must be greater than zero."), true
	case errors.Is(err, ErrOrderStatusInvalid):
		return validationField("order_status", "order.status_invalid", "Order status is invalid."), true
	case errors.Is(err, ErrBillingStatusInvalid):
		return validationField("billing_status", "order.billing_status_invalid", "Billing status is invalid."), true
	case errors.Is(err, ErrServiceStatusInvalid):
		return validationField("status", "service.status_invalid", "Service status is invalid."), true
	case errors.Is(err, ErrServiceStatusTransitionInvalid):
		return validationField("to_status", "service.status_transition_invalid", "Service status change is not allowed."), true
	case errors.Is(err, ErrServiceLifecycleActionInvalid):
		return validationField("action", "service.lifecycle_action_invalid", "Service lifecycle action is invalid."), true
	case errors.Is(err, ErrServiceLifecycleReasonMissing):
		return validationField("reason", "service.reason_missing", "A reason is required for this service action."), true
	case errors.Is(err, ErrStatusTransitionInvalid):
		return validationField("to_status", "order.status_transition_invalid", "Order status change is not allowed."), true
	case errors.Is(err, ErrSuspensionReasonInvalid):
		return validationField("suspension_reason", "service.suspension_reason_invalid", "Suspension reason is invalid."), true
	case errors.Is(err, catalog.ErrBillingCycleInvalid):
		return validationField("billing_cycle_type", "service.billing_cycle_invalid", "Billing cycle is invalid."), true
	case errors.Is(err, catalog.ErrBillingCycleValue):
		return validationField("billing_cycle_value", "service.billing_cycle_value_invalid", "Billing cycle value is invalid."), true
	default:
		return httpserver.ValidationField{}, false
	}
}

func validationField(field string, code string, message string) httpserver.ValidationField {
	return httpserver.ValidationField{
		Field:   field,
		Code:    code,
		Message: message,
	}
}
