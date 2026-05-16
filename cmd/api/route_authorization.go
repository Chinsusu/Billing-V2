package main

import (
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
)

var (
	adminRouteActorTypes = []identity.ActorType{
		identity.ActorTypePlatformAdmin,
		identity.ActorTypePlatformStaff,
	}
	resellerRouteActorTypes = []identity.ActorType{
		identity.ActorTypeResellerOwner,
		identity.ActorTypeResellerStaff,
	}
	clientRouteActorTypes = []identity.ActorType{
		identity.ActorTypeClient,
	}
)

func permissionMiddleware(
	authorizer rbac.Authorizer,
	permission rbac.Permission,
	risk rbac.RiskLevel,
	allowedActorTypes []identity.ActorType,
) rbac.RouteMiddleware {
	return rbac.RequirePermissionWithOptions(rbac.PermissionMiddlewareOptions{
		Authorizer:        authorizer,
		Permission:        permission,
		Risk:              risk,
		AllowedActorTypes: allowedActorTypes,
	})
}
