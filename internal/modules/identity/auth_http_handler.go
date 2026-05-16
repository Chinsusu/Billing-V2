package identity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/middleware"
)

const maxAuthJSONBodyBytes = 1 << 20

type AuthHTTPService interface {
	Login(ctx context.Context, input LoginInput) (LoginResult, error)
	Logout(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, input PasswordResetRequestInput) (PasswordResetRequestResult, error)
	ConfirmPasswordReset(ctx context.Context, input PasswordResetConfirmInput) (PasswordResetConfirmResult, error)
	SetupTwoFactor(ctx context.Context, token string) (SetupTwoFactorResult, error)
	VerifyTwoFactor(ctx context.Context, token string, code string) (VerifyTwoFactorResult, error)
}

type AuthHTTPHandlerOptions struct {
	CookieName             string
	CookieSecure           bool
	AllowLocalTenantHeader bool
}

type AuthHTTPHandler struct {
	service AuthHTTPService
	options AuthHTTPHandlerOptions
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	SessionID              string    `json:"session_id"`
	UserID                 UserID    `json:"user_id"`
	TenantID               tenant.ID `json:"tenant_id"`
	ActorType              ActorType `json:"actor_type"`
	ExpiresAt              time.Time `json:"expires_at"`
	TwoFactorRequired      bool      `json:"two_factor_required"`
	TwoFactorSatisfied     bool      `json:"two_factor_satisfied"`
	TwoFactorSetupRequired bool      `json:"two_factor_setup_required"`
}

type setupTwoFactorResponse struct {
	Method       string `json:"method"`
	Secret       string `json:"secret"`
	ProvisionURI string `json:"provision_uri"`
}

type verifyTwoFactorRequest struct {
	Code string `json:"code"`
}

type passwordResetRequest struct {
	Email string `json:"email"`
}

type passwordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type passwordResetResponse struct {
	Status string `json:"status"`
}

type verifyTwoFactorResponse struct {
	SessionID          string    `json:"session_id"`
	UserID             UserID    `json:"user_id"`
	TenantID           tenant.ID `json:"tenant_id"`
	TwoFactorSatisfied bool      `json:"two_factor_satisfied"`
}

func NewAuthHTTPHandlerWithOptions(service AuthHTTPService, options AuthHTTPHandlerOptions) *AuthHTTPHandler {
	if options.CookieName == "" {
		options.CookieName = "billing_session"
	}
	return &AuthHTTPHandler{service: service, options: options}
}

func (handler *AuthHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", middleware.RequireMethod(http.MethodPost, handler.handleLogin))
	mux.HandleFunc("/auth/logout", middleware.RequireMethod(http.MethodPost, handler.handleLogout))
	mux.HandleFunc("/auth/password-reset/request", middleware.RequireMethod(http.MethodPost, handler.handlePasswordResetRequest))
	mux.HandleFunc("/auth/password-reset/confirm", middleware.RequireMethod(http.MethodPost, handler.handlePasswordResetConfirm))
	mux.HandleFunc("/auth/2fa/setup", middleware.RequireMethod(http.MethodPost, handler.handleSetupTwoFactor))
	mux.HandleFunc("/auth/2fa/verify", middleware.RequireMethod(http.MethodPost, handler.handleVerifyTwoFactor))
}

func (handler *AuthHTTPHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request loginRequest
	if !decodeAuthJSON(w, r, &request) {
		return
	}
	if request.Email == "" || request.Password == "" {
		writeAuthValidationError(w, r, request)
		return
	}
	result, err := handler.service.Login(r.Context(), LoginInput{
		Email:                  request.Email,
		Password:               request.Password,
		Domain:                 requestDomain(r),
		LocalTenantID:          tenant.ID(r.Header.Get(tenant.HeaderTenantID)),
		AllowLocalTenantHeader: handler.options.AllowLocalTenantHeader,
		UserAgent:              r.UserAgent(),
		ClientIP:               requestClientIP(r),
	})
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	http.SetCookie(w, handler.sessionCookie(result.Token, result.ExpiresAt))
	httpserver.WriteSuccess(w, r, http.StatusOK, loginResponse{
		SessionID:              result.Session.ID,
		UserID:                 result.User.ID,
		TenantID:               result.User.TenantID,
		ActorType:              actorTypeForUser(result.User),
		ExpiresAt:              result.ExpiresAt,
		TwoFactorRequired:      result.TwoFactorRequired,
		TwoFactorSatisfied:     result.TwoFactorSatisfied,
		TwoFactorSetupRequired: result.TwoFactorSetupRequired,
	})
}

func (handler *AuthHTTPHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	if cookie, err := r.Cookie(handler.options.CookieName); err == nil {
		if err := handler.service.Logout(r.Context(), cookie.Value); err != nil {
			writeAuthError(w, r, err)
			return
		}
	}
	http.SetCookie(w, handler.expiredSessionCookie())
	httpserver.WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "logged_out"})
}

func (handler *AuthHTTPHandler) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request passwordResetRequest
	if !decodeAuthJSON(w, r, &request) {
		return
	}
	if strings.TrimSpace(request.Email) == "" {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			{Field: "email", Code: "required", Message: "Email is required."},
		})
		return
	}
	if _, err := handler.service.RequestPasswordReset(r.Context(), PasswordResetRequestInput{
		Email:                  request.Email,
		Domain:                 requestDomain(r),
		LocalTenantID:          tenant.ID(r.Header.Get(tenant.HeaderTenantID)),
		AllowLocalTenantHeader: handler.options.AllowLocalTenantHeader,
		ClientIP:               requestClientIP(r),
	}); err != nil {
		writeAuthError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusAccepted, passwordResetResponse{Status: "accepted"})
}

func (handler *AuthHTTPHandler) handlePasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	var request passwordResetConfirmRequest
	if !decodeAuthJSON(w, r, &request) {
		return
	}
	fields := []httpserver.ValidationField{}
	if strings.TrimSpace(request.Token) == "" {
		fields = append(fields, httpserver.ValidationField{Field: "token", Code: "required", Message: "Reset token is required."})
	}
	if strings.TrimSpace(request.NewPassword) == "" {
		fields = append(fields, httpserver.ValidationField{Field: "new_password", Code: "required", Message: "New password is required."})
	}
	if len(fields) > 0 {
		httpserver.WriteValidationError(w, r, fields)
		return
	}
	if _, err := handler.service.ConfirmPasswordReset(r.Context(), PasswordResetConfirmInput{
		Token:       request.Token,
		NewPassword: request.NewPassword,
		ClientIP:    requestClientIP(r),
	}); err != nil {
		writeAuthError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, passwordResetResponse{Status: "password_updated"})
}

func (handler *AuthHTTPHandler) handleSetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	token, ok := handler.sessionToken(w, r)
	if !ok {
		return
	}
	result, err := handler.service.SetupTwoFactor(r.Context(), token)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, setupTwoFactorResponse{
		Method:       result.Method,
		Secret:       result.Secret,
		ProvisionURI: result.ProvisionURI,
	})
}

func (handler *AuthHTTPHandler) handleVerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	token, ok := handler.sessionToken(w, r)
	if !ok {
		return
	}
	var request verifyTwoFactorRequest
	if !decodeAuthJSON(w, r, &request) {
		return
	}
	if request.Code == "" {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			{Field: "code", Code: "required", Message: "Two-factor code is required."},
		})
		return
	}
	result, err := handler.service.VerifyTwoFactor(r.Context(), token, request.Code)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, verifyTwoFactorResponse{
		SessionID:          result.Session.ID,
		UserID:             result.User.ID,
		TenantID:           result.User.TenantID,
		TwoFactorSatisfied: result.Session.TwoFactorSatisfied(),
	})
}

func (handler *AuthHTTPHandler) ready(w http.ResponseWriter, r *http.Request) bool {
	if handler == nil || handler.service == nil {
		writeAuthError(w, r, ErrAuthStoreMissing)
		return false
	}
	return true
}

func (handler *AuthHTTPHandler) sessionCookie(token string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     handler.options.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   handler.options.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (handler *AuthHTTPHandler) expiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     handler.options.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   handler.options.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (handler *AuthHTTPHandler) sessionToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(handler.options.CookieName)
	if err != nil || cookie.Value == "" {
		writeAuthError(w, r, ErrSessionTokenMissing)
		return "", false
	}
	return cookie.Value, true
}

func requestClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func requestDomain(r *http.Request) string {
	if r == nil {
		return ""
	}
	if host := firstForwardedHost(r.Header.Get("X-Forwarded-Host")); host != "" {
		return host
	}
	if host := forwardedHeaderHost(r.Header.Get("Forwarded")); host != "" {
		return host
	}
	return r.Host
}

func firstForwardedHost(value string) string {
	first, _, _ := strings.Cut(value, ",")
	return cleanForwardedHost(first)
}

func forwardedHeaderHost(value string) string {
	first, _, _ := strings.Cut(value, ",")
	for _, part := range strings.Split(first, ";") {
		key, raw, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok || !strings.EqualFold(strings.TrimSpace(key), "host") {
			continue
		}
		return cleanForwardedHost(raw)
	}
	return ""
}

func cleanForwardedHost(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"`)
	return strings.TrimSpace(value)
}

func decodeAuthJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxAuthJSONBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must be valid JSON.")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must contain a single JSON object.")
		return false
	}
	return true
}

func writeAuthValidationError(w http.ResponseWriter, r *http.Request, request loginRequest) {
	fields := []httpserver.ValidationField{}
	if request.Email == "" {
		fields = append(fields, httpserver.ValidationField{Field: "email", Code: "required", Message: "Email is required."})
	}
	if request.Password == "" {
		fields = append(fields, httpserver.ValidationField{Field: "password", Code: "required", Message: "Password is required."})
	}
	httpserver.WriteValidationError(w, r, fields)
}

func writeAuthError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrAuthRateLimited):
		httpserver.WriteError(w, r, http.StatusTooManyRequests, "auth.rate_limited", "Too many authentication attempts.")
	case errors.Is(err, ErrInvalidCredentials):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.invalid_credentials", "Invalid email or password.")
	case errors.Is(err, ErrLoginTenantInvalid), errors.Is(err, tenant.ErrDomainNotFound), errors.Is(err, tenant.ErrTenantIDMissing):
		httpserver.WriteError(w, r, http.StatusBadRequest, "tenant.context_missing", "Tenant context is required.")
	case errors.Is(err, tenant.ErrTenantStatusInvalid), errors.Is(err, ErrUserInactive):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.user_inactive", "User is not active.")
	case errors.Is(err, ErrTwoFactorRequired):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.2fa_required", "Two-factor authentication is required.")
	case errors.Is(err, ErrTwoFactorSetupRequired):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.2fa_setup_required", "Two-factor setup is required.")
	case errors.Is(err, ErrTwoFactorAlreadyEnabled):
		httpserver.WriteError(w, r, http.StatusConflict, "auth.2fa_already_enabled", "Two-factor authentication is already enabled.")
	case errors.Is(err, ErrTwoFactorNotAllowed):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.2fa_not_allowed", "Two-factor setup is not allowed for this actor.")
	case errors.Is(err, ErrTwoFactorCodeMissing), errors.Is(err, ErrTwoFactorCodeInvalid):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.2fa_invalid", "Two-factor code is invalid.")
	case errors.Is(err, ErrPasswordResetTokenMissing), errors.Is(err, ErrPasswordResetTokenInvalid), errors.Is(err, ErrPasswordResetTokenExpired), errors.Is(err, ErrPasswordResetTokenUsed):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.password_reset_invalid", "Password reset token is invalid or expired.")
	case errors.Is(err, ErrSessionTokenMissing), errors.Is(err, ErrSessionInvalid), errors.Is(err, ErrSessionExpired):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.session_invalid", "Session is invalid or expired.")
	case errors.Is(err, ErrAuthStoreMissing), errors.Is(err, ErrUserStoreExecutorMissing), errors.Is(err, ErrSessionStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.service_unavailable", "Authentication service is unavailable.")
	case errors.Is(err, ErrTwoFactorStoreExecutorMissing), errors.Is(err, ErrSecretCipherMissing), errors.Is(err, ErrEncryptionKeyInvalid):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.service_unavailable", "Authentication service is unavailable.")
	case errors.Is(err, ErrAuthRateLimitStoreExecutorMissing), errors.Is(err, ErrPasswordResetStoreExecutorMissing), errors.Is(err, ErrPasswordResetTTLMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.service_unavailable", "Authentication service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "auth.failed", "Authentication failed.")
	}
}
