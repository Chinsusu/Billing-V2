package rbac

import (
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestValidateTenantUserLookupRequiresTenant(t *testing.T) {
	err := validateTenantUserLookup("", "user_1")

	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant id error, got %v", err)
	}
}

func TestValidateTenantUserLookupRequiresUser(t *testing.T) {
	err := validateTenantUserLookup("tenant_a", "")

	if !errors.Is(err, identity.ErrUserIDMissing) {
		t.Fatalf("expected user id error, got %v", err)
	}
}

func TestRoleIDStringsDropsEmptyAndCopies(t *testing.T) {
	roleIDs := []identity.RoleID{"owner", "", "support"}
	values := roleIDStrings(roleIDs)
	roleIDs[0] = "changed"

	if len(values) != 2 {
		t.Fatalf("expected two role ids, got %d", len(values))
	}
	if values[0] != "owner" || values[1] != "support" {
		t.Fatalf("unexpected role ids: %#v", values)
	}
}
