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

func TestAdminReadHTTPHandlerTenantsRequiresTenant(t *testing.T) {
	service := &fakeAdminReadService{}
	handler := registerAdminReadTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/tenants", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTenantCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listTenantCalls)
	}
}

func TestAdminReadHTTPHandlerListTenantsUsesTenantScopeAndFilters(t *testing.T) {
	now := time.Date(2026, 4, 24, 1, 2, 3, 0, time.UTC)
	service := &fakeAdminReadService{
		tenants: []tenant.TenantSummary{{
			Tenant: tenant.Tenant{
				ID:              "tenant_1",
				DisplayID:       10010,
				Type:            tenant.TypeReseller,
				Name:            "Demo Reseller",
				Slug:            "demo-reseller",
				Status:          tenant.StatusActive,
				DefaultCurrency: "USD",
				Timezone:        "Asia/Ho_Chi_Minh",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			PrimaryDomain: "demo.local",
			UserCount:     3,
		}},
	}
	handler := registerAdminReadTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/tenants?type=reseller&status=active&display_id=10010&limit=10", nil)
	request.Header.Set(tenant.HeaderTenantID, " tenant_scope ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTenantCalls != 1 {
		t.Fatalf("expected list tenants once, got %d", service.listTenantCalls)
	}
	if service.tenantFilter.ScopeTenantID != "tenant_scope" || service.tenantFilter.Type != tenant.TypeReseller ||
		service.tenantFilter.Status != tenant.StatusActive || service.tenantFilter.DisplayID != 10010 || service.tenantFilter.Limit != 10 {
		t.Fatalf("unexpected tenant filter: %+v", service.tenantFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{`"display_id":10010`, `"primary_domain":"demo.local"`, `"user_count":3`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in response, got %s", expected, body)
		}
	}
}

func TestAdminReadHTTPHandlerCustomersDefaultsToClientType(t *testing.T) {
	service := &fakeAdminReadService{
		users: []UserSummary{{
			User: User{
				ID:              "user_1",
				DisplayID:       20001,
				TenantID:        "tenant_1",
				Email:           "client@example.com",
				FullName:        "Client One",
				Type:            UserTypeClient,
				Status:          UserStatusActive,
				TwoFactorStatus: TwoFactorStatusDisabled,
			},
			TenantName: "Demo Reseller",
			TenantSlug: "demo-reseller",
		}},
	}
	handler := registerAdminReadTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/customers?status=active", nil)
	request.Header.Set(tenant.HeaderTenantID, "tenant_scope")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listUserCalls != 1 {
		t.Fatalf("expected list users once, got %d", service.listUserCalls)
	}
	if service.userFilter.TenantID != "tenant_scope" || service.userFilter.Type != UserTypeClient || service.userFilter.Status != UserStatusActive {
		t.Fatalf("unexpected user filter: %+v", service.userFilter)
	}
	body := response.Body.String()
	if strings.Contains(body, "password_hash") {
		t.Fatalf("response leaked password hash: %s", body)
	}
	if !strings.Contains(body, `"display_id":20001`) || !strings.Contains(body, `"tenant_name":"Demo Reseller"`) {
		t.Fatalf("expected account fields in response, got %s", body)
	}
}

func TestAdminReadHTTPHandlerResellerCustomersUsesTenantScopeAndFilters(t *testing.T) {
	service := &fakeAdminReadService{
		users: []UserSummary{{
			User: User{
				ID:              "user_2",
				DisplayID:       20002,
				TenantID:        "reseller_tenant",
				Email:           "buyer@example.com",
				FullName:        "Buyer Two",
				Type:            UserTypeClient,
				Status:          UserStatusActive,
				TwoFactorStatus: TwoFactorStatusDisabled,
			},
			TenantName: "Demo Reseller",
			TenantSlug: "demo-reseller",
		}},
	}
	handler := registerAdminReadTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/customers?status=active&email=buyer@example.com&display_id=20002&limit=10", nil)
	request.Header.Set(tenant.HeaderTenantID, " reseller_tenant ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listUserCalls != 1 {
		t.Fatalf("expected list users once, got %d", service.listUserCalls)
	}
	if service.userFilter.TenantID != "reseller_tenant" || service.userFilter.Type != UserTypeClient ||
		service.userFilter.Status != UserStatusActive || service.userFilter.Email != "buyer@example.com" ||
		service.userFilter.DisplayID != 20002 || service.userFilter.Limit != 10 {
		t.Fatalf("unexpected user filter: %+v", service.userFilter)
	}
	body := response.Body.String()
	if !strings.Contains(body, `"display_id":20002`) || !strings.Contains(body, `"email":"buyer@example.com"`) {
		t.Fatalf("expected reseller customer fields in response, got %s", body)
	}
}

func TestAdminReadHTTPHandlerRejectsBadUserStatus(t *testing.T) {
	service := &fakeAdminReadService{}
	handler := registerAdminReadTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/accounts?status=bad", nil)
	request.Header.Set(tenant.HeaderTenantID, "tenant_scope")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listUserCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listUserCalls)
	}
	if !strings.Contains(response.Body.String(), "identity.user_status_invalid") {
		t.Fatalf("expected user status validation, got %s", response.Body.String())
	}
}

func TestAdminReadHTTPHandlerAdminMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeAdminReadService{}
	mux := http.NewServeMux()
	NewAdminReadHTTPHandlerWithOptions(service, AdminReadHTTPHandlerOptions{
		AdminMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/admin/accounts", nil)
	request.Header.Set(tenant.HeaderTenantID, "tenant_scope")
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.listUserCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listUserCalls)
	}
}

func TestAdminReadHTTPHandlerResellerMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeAdminReadService{}
	mux := http.NewServeMux()
	NewAdminReadHTTPHandlerWithOptions(service, AdminReadHTTPHandlerOptions{
		ResellerMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/reseller/customers", nil)
	request.Header.Set(tenant.HeaderTenantID, "tenant_scope")
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.listUserCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listUserCalls)
	}
}

func registerAdminReadTestHandler(service AdminReadHTTPService) http.Handler {
	mux := http.NewServeMux()
	NewAdminReadHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeAdminReadService struct {
	tenants          []tenant.TenantSummary
	users            []UserSummary
	tenantFilter     tenant.ListTenantsFilter
	userFilter       UserListFilter
	listTenantCalls  int
	listUserCalls    int
	listTenantsError error
	listUsersError   error
}

func (service *fakeAdminReadService) ListAdminTenants(ctx context.Context, filter tenant.ListTenantsFilter) ([]tenant.TenantSummary, error) {
	service.listTenantCalls++
	service.tenantFilter = filter
	if service.listTenantsError != nil {
		return nil, service.listTenantsError
	}
	return service.tenants, nil
}

func (service *fakeAdminReadService) ListAdminUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error) {
	service.listUserCalls++
	service.userFilter = filter
	if service.listUsersError != nil {
		return nil, service.listUsersError
	}
	return service.users, nil
}
