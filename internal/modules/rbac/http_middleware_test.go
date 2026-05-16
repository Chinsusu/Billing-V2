package rbac

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestRequirePermissionRejectsMissingActor(t *testing.T) {
	handler := RequirePermission(&fakeAuthorizer{}, PermissionCatalogView, RiskLow)(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not run")
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	handler(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "auth.actor_required") {
		t.Fatalf("expected actor error, got %s", response.Body.String())
	}
}

func TestRequirePermissionCallsAuthorizer(t *testing.T) {
	authorizer := &fakeAuthorizer{}
	called := false
	handler := RequirePermission(authorizer, PermissionCatalogView, RiskLow)(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))
	ctx = identity.WithActor(ctx, identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil).WithContext(ctx)
	handler(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", response.Code, response.Body.String())
	}
	if !called {
		t.Fatal("expected protected handler to run")
	}
	if authorizer.request.Permission != PermissionCatalogView || authorizer.request.ResourceTenantID != "tenant_a" {
		t.Fatalf("unexpected auth request: %+v", authorizer.request)
	}
}

func TestRequirePermissionRejectsDisallowedActorType(t *testing.T) {
	authorizer := &fakeAuthorizer{}
	handler := RequirePermissionWithOptions(PermissionMiddlewareOptions{
		Authorizer:        authorizer,
		Permission:        PermissionCatalogView,
		Risk:              RiskLow,
		AllowedActorTypes: []identity.ActorType{identity.ActorTypeResellerStaff},
	})(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not run")
	})
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))
	ctx = identity.WithActor(ctx, identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil).WithContext(ctx)
	handler(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "auth.permission_denied") {
		t.Fatalf("expected permission denied response, got %s", response.Body.String())
	}
	if authorizer.called {
		t.Fatal("authorizer should not run for disallowed actor type")
	}
}

func TestRequirePermissionMapsDenied(t *testing.T) {
	handler := RequirePermission(&fakeAuthorizer{err: ErrPermissionDenied}, PermissionCatalogManage, RiskHigh)(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not run")
	})
	ctx := identity.WithActor(context.Background(), identity.NewActor("admin_1", "platform", identity.ActorTypePlatformAdmin))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/protected", nil).WithContext(ctx)
	handler(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", response.Code, response.Body.String())
	}
}

type fakeAuthorizer struct {
	called  bool
	request CheckRequest
	err     error
}

func (authorizer *fakeAuthorizer) Check(ctx context.Context, request CheckRequest) error {
	authorizer.called = true
	authorizer.request = request
	return authorizer.err
}
