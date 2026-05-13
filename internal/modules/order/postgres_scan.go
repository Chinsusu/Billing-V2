package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type orderScanner interface {
	Scan(dest ...interface{}) error
}

func scanOrder(row orderScanner) (Order, error) {
	return scanOrderFields(row, false)
}

func scanOrderRead(row orderScanner) (Order, error) {
	return scanOrderFields(row, true)
}

func scanOrderFields(row orderScanner, includeBuyerDisplayID bool) (Order, error) {
	var record Order
	var id, tenantID, buyerUserID, tenantPlanID, orderStatus, billingStatus, idempotencyKey string
	var buyerDisplayID sql.NullInt64
	var productSnapshot, planSnapshot, priceSnapshot []byte
	destinations := []interface{}{
		&id, &record.DisplayID, &tenantID, &buyerUserID, &tenantPlanID, &record.Quantity, &record.Currency,
		&record.UnitPriceMinor, &record.DiscountMinor, &record.TotalMinor, &orderStatus, &billingStatus, &idempotencyKey,
		&productSnapshot, &planSnapshot, &priceSnapshot, &record.CreatedAt, &record.UpdatedAt,
	}
	if includeBuyerDisplayID {
		destinations = append(destinations, &buyerDisplayID)
	}
	if err := row.Scan(destinations...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, ErrOrderNotFound
		}
		return Order{}, fmt.Errorf("scan order: %w", err)
	}
	record.ID = OrderID(id)
	record.TenantID = tenant.ID(tenantID)
	record.BuyerUserID = identity.UserID(buyerUserID)
	if buyerDisplayID.Valid {
		record.BuyerDisplayID = buyerDisplayID.Int64
	}
	record.TenantPlanID = catalog.TenantPlanID(tenantPlanID)
	record.OrderStatus = OrderStatus(orderStatus)
	record.BillingStatus = BillingStatus(billingStatus)
	record.IdempotencyKey = IdempotencyKey(idempotencyKey)
	record.ProductSnapshot = append(json.RawMessage(nil), productSnapshot...)
	record.PlanSnapshot = append(json.RawMessage(nil), planSnapshot...)
	record.PriceSnapshot = append(json.RawMessage(nil), priceSnapshot...)
	return record, nil
}

func scanReservation(row orderScanner) (Reservation, error) {
	var record Reservation
	var id, orderID, tenantID, providerSourceID, status string
	if err := row.Scan(
		&id, &record.DisplayID, &orderID, &tenantID, &providerSourceID,
		&record.Quantity, &status, &record.ExpiresAt, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return Reservation{}, fmt.Errorf("scan order reservation: %w", err)
	}
	record.ID = ReservationID(id)
	record.OrderID = OrderID(orderID)
	record.TenantID = tenant.ID(tenantID)
	record.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	record.Status = ReservationStatus(status)
	return record, nil
}

func scanProvisioningJob(row orderScanner) (ProvisioningJob, error) {
	var record ProvisioningJob
	var id, orderID, tenantID, providerSourceID, operationID, status, idempotencyKey string
	var lastErrorCode, lastErrorMessage sql.NullString
	if err := row.Scan(
		&id, &record.DisplayID, &orderID, &tenantID, &providerSourceID, &operationID, &status,
		&idempotencyKey, &record.AttemptNumber, &lastErrorCode, &lastErrorMessage, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return ProvisioningJob{}, fmt.Errorf("scan order provisioning job: %w", err)
	}
	record.ID = ProvisioningJobID(id)
	record.OrderID = OrderID(orderID)
	record.TenantID = tenant.ID(tenantID)
	record.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	record.ProviderOperationID = ProviderOperationID(operationID)
	record.Status = ProvisioningStatus(status)
	record.IdempotencyKey = IdempotencyKey(idempotencyKey)
	record.LastErrorCode = lastErrorCode.String
	record.LastErrorMessage = lastErrorMessage.String
	return record, nil
}

func scanServiceInstance(row orderScanner) (ServiceInstance, error) {
	return scanServiceInstanceFields(row, false)
}

func scanServiceInstanceRead(row orderScanner) (ServiceInstance, error) {
	return scanServiceInstanceFields(row, true)
}

func scanServiceInstanceFields(row orderScanner, includeRelatedDisplayIDs bool) (ServiceInstance, error) {
	var record ServiceInstance
	var id, tenantID, orderID, tenantPlanID, providerSourceID, externalResourceID, status, billingStatus string
	var orderDisplayID, buyerDisplayID, providerSourceDisplayID sql.NullInt64
	var suspensionReason sql.NullString
	destinations := []interface{}{
		&id, &record.DisplayID, &tenantID, &orderID, &tenantPlanID, &providerSourceID, &externalResourceID,
		&status, &billingStatus, &suspensionReason, &record.TermStart, &record.TermEnd, &record.CreatedAt, &record.UpdatedAt,
	}
	if includeRelatedDisplayIDs {
		destinations = append(destinations, &orderDisplayID, &buyerDisplayID, &providerSourceDisplayID)
	}
	if err := row.Scan(destinations...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceInstance{}, ErrServiceNotFound
		}
		return ServiceInstance{}, fmt.Errorf("scan service instance: %w", err)
	}
	record.ID = ServiceID(id)
	record.TenantID = tenant.ID(tenantID)
	record.OrderID = OrderID(orderID)
	if orderDisplayID.Valid {
		record.OrderDisplayID = orderDisplayID.Int64
	}
	if buyerDisplayID.Valid {
		record.BuyerDisplayID = buyerDisplayID.Int64
	}
	record.TenantPlanID = catalog.TenantPlanID(tenantPlanID)
	record.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	if providerSourceDisplayID.Valid {
		record.ProviderSourceDisplayID = providerSourceDisplayID.Int64
	}
	record.ExternalResourceID = provider.ExternalResourceID(externalResourceID)
	record.Status = ServiceStatus(status)
	record.BillingStatus = BillingStatus(billingStatus)
	record.SuspensionReason = SuspensionReason(suspensionReason.String)
	return record, nil
}
