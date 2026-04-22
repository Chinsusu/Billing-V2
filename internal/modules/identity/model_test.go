package identity

import (
	"errors"
	"testing"
)

func TestNewActorCopiesRoleIDs(t *testing.T) {
	roles := []RoleID{"owner"}
	actor := NewActor("user_1", "tenant_a", ActorTypeResellerOwner, roles...)
	roles[0] = "changed"

	if !actor.HasRole("owner") {
		t.Fatalf("expected copied role id")
	}
}

func TestActorValidateRequiresTenant(t *testing.T) {
	actor := NewActor("user_1", "", ActorTypeClient)

	if err := actor.Validate(); !errors.Is(err, ErrActorTenantMissing) {
		t.Fatalf("expected missing tenant, got %v", err)
	}
}

func TestSystemActorIsValidForTenantJob(t *testing.T) {
	actor := SystemActor("tenant_a")

	if err := actor.Validate(); err != nil {
		t.Fatalf("expected system actor to be valid, got %v", err)
	}
	if !actor.IsSystem {
		t.Fatalf("expected system actor flag")
	}
}
