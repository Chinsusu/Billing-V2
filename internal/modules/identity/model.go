package identity

import (
	"errors"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrActorIDMissing     = errors.New("actor id missing")
	ErrActorTypeMissing   = errors.New("actor type missing")
	ErrActorTenantMissing = errors.New("actor tenant missing")
)

type UserID string
type RoleID string

type ActorType string

const (
	ActorTypePlatformAdmin ActorType = "platform_admin"
	ActorTypePlatformStaff ActorType = "platform_staff"
	ActorTypeResellerOwner ActorType = "reseller_owner"
	ActorTypeResellerStaff ActorType = "reseller_staff"
	ActorTypeClient        ActorType = "client"
	ActorTypeSystem        ActorType = "system"
)

type Actor struct {
	ID              UserID
	TenantID        tenant.ID
	Type            ActorType
	RoleIDs         []RoleID
	IsPlatformAdmin bool
	IsSystem        bool
}

func NewActor(id UserID, tenantID tenant.ID, actorType ActorType, roleIDs ...RoleID) Actor {
	return Actor{
		ID:              id,
		TenantID:        tenantID,
		Type:            actorType,
		RoleIDs:         append([]RoleID(nil), roleIDs...),
		IsPlatformAdmin: actorType == ActorTypePlatformAdmin || actorType == ActorTypePlatformStaff,
		IsSystem:        actorType == ActorTypeSystem,
	}
}

func SystemActor(tenantID tenant.ID) Actor {
	return Actor{
		ID:       "system",
		TenantID: tenantID,
		Type:     ActorTypeSystem,
		IsSystem: true,
	}
}

func (actor Actor) Validate() error {
	if actor.ID == "" {
		return ErrActorIDMissing
	}
	if actor.Type == "" {
		return ErrActorTypeMissing
	}
	if actor.TenantID.Empty() {
		return ErrActorTenantMissing
	}
	return nil
}

func (actor Actor) HasRole(roleID RoleID) bool {
	for _, current := range actor.RoleIDs {
		if current == roleID {
			return true
		}
	}
	return false
}
