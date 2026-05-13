package identity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestSessionMiddlewareAddsActorAndTenantContext(t *testing.T) {
	resolver := &fakeSessionResolver{
		identity: SessionIdentity{
			Session: Session{ID: "session_1", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Type: UserTypeClient, Status: UserStatusActive},
			RoleIDs: []RoleID{"role_client"},
		},
	}
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName: "billing_session",
		Resolver:   resolver,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, actorOK := FromContext(r.Context())
		tenantContext, tenantOK := tenant.FromContext(r.Context())
		if !actorOK || !tenantOK {
			t.Fatalf("expected actor and tenant context")
		}
		if actor.ID != "user_1" || actor.Type != ActorTypeClient || !actor.HasRole("role_client") {
			t.Fatalf("unexpected actor: %+v", actor)
		}
		if tenantContext.ActorTenantID != "tenant_1" || tenantContext.EffectiveTenantID != "tenant_1" {
			t.Fatalf("unexpected tenant context: %+v", tenantContext)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/client/wallets", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "session-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", response.Code, response.Body.String())
	}
	if resolver.token != "session-token" {
		t.Fatalf("expected resolver token, got %q", resolver.token)
	}
}

func TestSessionMiddlewareRejectsInvalidCookie(t *testing.T) {
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName: "billing_session",
		Resolver:   &fakeSessionResolver{err: ErrSessionInvalid},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	}))

	request := httptest.NewRequest(http.MethodGet, "/client/wallets", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "bad-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "auth.session_invalid") {
		t.Fatalf("expected session error code, got %s", response.Body.String())
	}
}

func TestSessionMiddlewareAllowsAuthRoutesWithStaleCookie(t *testing.T) {
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName: "billing_session",
		Resolver:   &fakeSessionResolver{err: ErrSessionInvalid},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "bad-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected auth route to bypass session resolution, got %d: %s", response.Code, response.Body.String())
	}
}

func TestSessionMiddlewareAllowsPasswordResetWithStaleCookie(t *testing.T) {
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName: "billing_session",
		Resolver:   &fakeSessionResolver{err: ErrSessionInvalid},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, "/auth/password-reset/request", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "bad-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected password reset route to bypass session resolution, got %d: %s", response.Code, response.Body.String())
	}
}

func TestSessionMiddlewareRejectsUnsatisfiedAdminTwoFactor(t *testing.T) {
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName:            "billing_session",
		Resolver:              &fakeSessionResolver{identity: adminSessionIdentity(Session{})},
		RequireAdminTwoFactor: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin/wallets", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "session-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "auth.2fa_required") {
		t.Fatalf("expected 2FA error, got %s", response.Body.String())
	}
}

func TestSessionMiddlewareAllowsSatisfiedAdminTwoFactor(t *testing.T) {
	handler := SessionMiddleware(SessionMiddlewareOptions{
		CookieName: "billing_session",
		Resolver: &fakeSessionResolver{identity: adminSessionIdentity(Session{
			TwoFactorSatisfiedAt: time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC),
		})},
		RequireAdminTwoFactor: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin/wallets", nil)
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "session-token"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", response.Code, response.Body.String())
	}
}

func adminSessionIdentity(session Session) SessionIdentity {
	session.ID = "session_1"
	session.TenantID = "tenant_1"
	session.UserID = "user_1"
	return SessionIdentity{
		Session: session,
		User:    User{ID: "user_1", TenantID: "tenant_1", Type: UserTypePlatformStaff, Status: UserStatusActive},
		RoleIDs: []RoleID{"role_admin"},
	}
}

type fakeSessionResolver struct {
	identity SessionIdentity
	err      error
	token    string
}

func (resolver *fakeSessionResolver) ResolveSession(ctx context.Context, token string) (SessionIdentity, error) {
	resolver.token = token
	if err := resolver.err; err != nil {
		return SessionIdentity{}, err
	}
	return resolver.identity, nil
}
