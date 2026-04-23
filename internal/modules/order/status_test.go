package order

import "testing"

func TestOrderTransitionGuards(t *testing.T) {
	if !CanTransitionOrder(OrderStatusDraft, OrderStatusPendingPayment) {
		t.Fatal("expected draft to pending payment transition")
	}
	if CanTransitionOrder(OrderStatusRefunded, OrderStatusPaid) {
		t.Fatal("refunded order should not transition back to paid")
	}
}

func TestProvisioningTransitionGuards(t *testing.T) {
	if !CanTransitionProvisioning(ProvisioningStatusQueued, ProvisioningStatusProvisioning) {
		t.Fatal("expected queued to provisioning transition")
	}
	if !CanTransitionProvisioning(ProvisioningStatusFailed, ProvisioningStatusQueued) {
		t.Fatal("expected failed job to retry as queued")
	}
	if CanTransitionProvisioning(ProvisioningStatusProvisioned, ProvisioningStatusQueued) {
		t.Fatal("provisioned job should not transition back to queued")
	}
}

func TestServiceTransitionGuards(t *testing.T) {
	if !CanTransitionService(ServiceStatusActive, ServiceStatusSuspended) {
		t.Fatal("expected active to suspended transition")
	}
	if !CanTransitionService(ServiceStatusSuspended, ServiceStatusActive) {
		t.Fatal("expected suspended to active transition")
	}
	if CanTransitionService(ServiceStatusTerminated, ServiceStatusActive) {
		t.Fatal("terminated service should not transition back to active")
	}
}
