package order

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func decodeOrderJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must be valid JSON.")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must contain one JSON object.")
		return false
	}
	return true
}

func decodeOptionalOrderJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if r.Body == nil || r.ContentLength == 0 {
		return true
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return true
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must be valid JSON.")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must contain one JSON object.")
		return false
	}
	return true
}

func tenantIDFromContext(w http.ResponseWriter, r *http.Request) (tenant.ID, bool) {
	tenantContext, err := tenant.RequireContext(r.Context())
	if err != nil {
		writeOrderError(w, r, err)
		return "", false
	}
	return tenantContext.EffectiveTenantID, true
}

func actorFromContext(w http.ResponseWriter, r *http.Request) (identity.Actor, bool) {
	actor, err := identity.RequireActor(r.Context())
	if err != nil {
		writeOrderError(w, r, err)
		return identity.Actor{}, false
	}
	return actor, true
}

func idempotencyKeyFromHeader(r *http.Request) IdempotencyKey {
	return IdempotencyKey(strings.TrimSpace(r.Header.Get(IdempotencyKeyHeader)))
}

func adminOrderIDFromPath(w http.ResponseWriter, r *http.Request) (OrderID, bool) {
	return orderIDFromPrefix(w, r, adminOrderPrefix)
}

func orderIDFromPath(w http.ResponseWriter, r *http.Request) (OrderID, bool) {
	return orderIDFromPrefix(w, r, clientOrderPrefix)
}

func orderIDFromPrefix(w http.ResponseWriter, r *http.Request, prefix string) (OrderID, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
	if value == "" || strings.Contains(value, "/") {
		writeOrderError(w, r, ErrOrderIDMissing)
		return "", false
	}
	return OrderID(value), true
}
