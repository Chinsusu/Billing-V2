package identity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/middleware"
)

const maxAuthJSONBodyBytes = 1 << 20

type AuthHTTPService interface {
	Login(ctx context.Context, input LoginInput) (LoginResult, error)
	Logout(ctx context.Context, token string) error
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
	SessionID string    `json:"session_id"`
	UserID    UserID    `json:"user_id"`
	TenantID  tenant.ID `json:"tenant_id"`
	ActorType ActorType `json:"actor_type"`
	ExpiresAt time.Time `json:"expires_at"`
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
		Domain:                 r.Host,
		LocalTenantID:          tenant.ID(r.Header.Get(tenant.HeaderTenantID)),
		AllowLocalTenantHeader: handler.options.AllowLocalTenantHeader,
		UserAgent:              r.UserAgent(),
	})
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	http.SetCookie(w, handler.sessionCookie(result.Token, result.ExpiresAt))
	httpserver.WriteSuccess(w, r, http.StatusOK, loginResponse{
		SessionID: result.Session.ID,
		UserID:    result.User.ID,
		TenantID:  result.User.TenantID,
		ActorType: actorTypeForUser(result.User),
		ExpiresAt: result.ExpiresAt,
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
	case errors.Is(err, ErrInvalidCredentials):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.invalid_credentials", "Invalid email or password.")
	case errors.Is(err, ErrLoginTenantInvalid), errors.Is(err, tenant.ErrDomainNotFound), errors.Is(err, tenant.ErrTenantIDMissing):
		httpserver.WriteError(w, r, http.StatusBadRequest, "tenant.context_missing", "Tenant context is required.")
	case errors.Is(err, tenant.ErrTenantStatusInvalid), errors.Is(err, ErrUserInactive):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.user_inactive", "User is not active.")
	case errors.Is(err, ErrTwoFactorRequired):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.2fa_required", "Two-factor authentication is required.")
	case errors.Is(err, ErrAuthStoreMissing), errors.Is(err, ErrUserStoreExecutorMissing), errors.Is(err, ErrSessionStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.service_unavailable", "Authentication service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "auth.failed", "Authentication failed.")
	}
}
