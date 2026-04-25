package order

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListServiceInstancesQueryAddsClientScope(t *testing.T) {
	query, args, err := buildListServiceInstancesQuery(ServiceInstanceFilter{
		TenantID:                tenant.ID("tenant-1"),
		BuyerUserID:             identity.UserID("buyer-1"),
		BuyerDisplayID:          10002,
		DisplayID:               50001,
		OrderID:                 OrderID("order-1"),
		OrderDisplayID:          30001,
		ProviderSourceDisplayID: 10003,
		Status:                  ServiceStatusActive,
		Limit:                   25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"JOIN orders ord",
		"ord.tenant_id = svc.tenant_id",
		"svc.tenant_id = $1",
		"ord.buyer_user_id = $2",
		"buyer.user_id = ord.buyer_user_id",
		"buyer.display_id = $3",
		"svc.display_id = $4",
		"svc.order_id = $5",
		"ord.display_id = $6",
		"source.source_id = svc.provider_source_id",
		"source.display_id = $7",
		"svc.status = $8",
		"LIMIT $9",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 9 || args[8] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListServiceInstancesQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListServiceInstancesQuery(ServiceInstanceFilter{TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "LIMIT $2") {
		t.Fatalf("expected default limit placeholder: %s", query)
	}
	if len(args) != 2 || args[1] != defaultServiceInstanceListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListServiceInstancesQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListServiceInstancesQuery(ServiceInstanceFilter{
		TenantID: tenant.ID("tenant-1"),
		Status:   ServiceStatus("bad"),
	})
	if !errors.Is(err, ErrServiceStatusInvalid) {
		t.Fatalf("expected service status error, got %v", err)
	}
}

func TestBuildGetServiceInstanceQueryAddsBuyerScope(t *testing.T) {
	query, args, err := buildGetServiceInstanceQuery(ServiceInstanceLookup{
		ID:          ServiceID("service-1"),
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: identity.UserID("buyer-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"svc.service_instance_id = $1", "svc.tenant_id = $2", "ord.buyer_user_id = $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}
