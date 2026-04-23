package order

import "testing"

func TestOrderLifecycleEventNamesAreStable(t *testing.T) {
	if OrderAggregateType != "order" {
		t.Fatalf("unexpected aggregate type %q", OrderAggregateType)
	}
	if OrderEventCreated != "order.created" {
		t.Fatalf("unexpected created event %q", OrderEventCreated)
	}
	if OrderEventStatusChanged != "order.status_changed" {
		t.Fatalf("unexpected status event %q", OrderEventStatusChanged)
	}
	if OrderProvisioningTrigger != OrderEventStatusChanged {
		t.Fatalf("expected provisioning to trigger from status change, got %q", OrderProvisioningTrigger)
	}
}
