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

func TestAuthHTTPHandlerLoginSetsHttpOnlySessionCookie(t *testing.T) {
	expiresAt := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	service := &fakeAuthHTTPService{
		loginResult: LoginResult{
			Token:     "plain-session-token",
			Session:   Session{ID: "session_1", ExpiresAt: expiresAt},
			User:      User{ID: "user_1", TenantID: "tenant_1", Type: UserTypeClient},
			ExpiresAt: expiresAt,
		},
	}
	mux := http.NewServeMux()
	NewAuthHTTPHandlerWithOptions(service, AuthHTTPHandlerOptions{
		CookieName:             "billing_session",
		CookieSecure:           true,
		AllowLocalTenantHeader: true,
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"client@example.com","password":"admin123"}`))
	request.Header.Set(tenant.HeaderTenantID, "tenant_1")
	request.Header.Set("User-Agent", "test-agent")
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "billing_session" || cookie.Value != "plain-session-token" || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected session cookie: %+v", cookie)
	}
	if strings.Contains(response.Body.String(), "plain-session-token") {
		t.Fatalf("response body leaked session token: %s", response.Body.String())
	}
	if service.loginInput.LocalTenantID != "tenant_1" || !service.loginInput.AllowLocalTenantHeader {
		t.Fatalf("expected local tenant login input, got %+v", service.loginInput)
	}
}

func TestAuthHTTPHandlerLoginValidatesRequiredFields(t *testing.T) {
	mux := http.NewServeMux()
	NewAuthHTTPHandlerWithOptions(&fakeAuthHTTPService{}, AuthHTTPHandlerOptions{}).RegisterRoutes(mux)

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":""}`)))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "validation.failed") {
		t.Fatalf("expected validation response, got %s", response.Body.String())
	}
}

func TestAuthHTTPHandlerSetupTwoFactorRequiresSessionCookie(t *testing.T) {
	mux := http.NewServeMux()
	NewAuthHTTPHandlerWithOptions(&fakeAuthHTTPService{}, AuthHTTPHandlerOptions{}).RegisterRoutes(mux)

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/auth/2fa/setup", nil))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
}

func TestAuthHTTPHandlerVerifyTwoFactorCallsService(t *testing.T) {
	service := &fakeAuthHTTPService{
		verifyResult: VerifyTwoFactorResult{
			Session: Session{ID: "session_1", TwoFactorSatisfiedAt: time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)},
			User:    User{ID: "user_1", TenantID: "tenant_1"},
		},
	}
	mux := http.NewServeMux()
	NewAuthHTTPHandlerWithOptions(service, AuthHTTPHandlerOptions{CookieName: "billing_session"}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodPost, "/auth/2fa/verify", strings.NewReader(`{"code":"123456"}`))
	request.AddCookie(&http.Cookie{Name: "billing_session", Value: "session-token"})
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.verifyToken != "session-token" || service.verifyCode != "123456" {
		t.Fatalf("unexpected verify input token=%q code=%q", service.verifyToken, service.verifyCode)
	}
	if !strings.Contains(response.Body.String(), `"two_factor_satisfied":true`) {
		t.Fatalf("expected satisfied response, got %s", response.Body.String())
	}
}

type fakeAuthHTTPService struct {
	loginInput   LoginInput
	loginResult  LoginResult
	loginErr     error
	logoutToken  string
	logoutErr    error
	setupToken   string
	setupResult  SetupTwoFactorResult
	setupErr     error
	verifyToken  string
	verifyCode   string
	verifyResult VerifyTwoFactorResult
	verifyErr    error
}

func (service *fakeAuthHTTPService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	service.loginInput = input
	if service.loginErr != nil {
		return LoginResult{}, service.loginErr
	}
	return service.loginResult, nil
}

func (service *fakeAuthHTTPService) Logout(ctx context.Context, token string) error {
	service.logoutToken = token
	return service.logoutErr
}

func (service *fakeAuthHTTPService) SetupTwoFactor(ctx context.Context, token string) (SetupTwoFactorResult, error) {
	service.setupToken = token
	if service.setupErr != nil {
		return SetupTwoFactorResult{}, service.setupErr
	}
	return service.setupResult, nil
}

func (service *fakeAuthHTTPService) VerifyTwoFactor(ctx context.Context, token string, code string) (VerifyTwoFactorResult, error) {
	service.verifyToken = token
	service.verifyCode = code
	if service.verifyErr != nil {
		return VerifyTwoFactorResult{}, service.verifyErr
	}
	return service.verifyResult, nil
}
