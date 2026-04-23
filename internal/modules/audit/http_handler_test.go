package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListAdminAuditLogsUsesFilters(t *testing.T) {
	service := &fakeAuditHTTPService{logs: []Log{testAuditLog()}}
	handler := registerAuditTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?actor_id=actor_1&actor_type=user&action=invoice.paid&target_type=invoice&target_id=target_1&created_from=2026-04-23T00:00:00Z&created_to=2026-04-24T00:00:00Z&limit=20", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 1 {
		t.Fatalf("expected list once, got %d", service.listCalls)
	}
	if service.filter.TenantID != tenant.ID("tenant_1") ||
		service.filter.ActorID != ActorID("actor_1") ||
		service.filter.ActorType != ActorTypeUser ||
		service.filter.Action != "invoice.paid" ||
		service.filter.TargetType != "invoice" ||
		service.filter.TargetID != TargetID("target_1") ||
		service.filter.Limit != 20 {
		t.Fatalf("unexpected audit filter: %+v", service.filter)
	}
	if strings.Contains(response.Body.String(), "before_snapshot_redacted") {
		t.Fatalf("list response should not include payload snapshots: %s", response.Body.String())
	}
}

func TestHTTPHandlerGetAdminAuditLogUsesTenantScope(t *testing.T) {
	service := &fakeAuditHTTPService{log: testAuditLog()}
	handler := registerAuditTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/audit-logs/audit_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.lookup.ID != ID("audit_1") || service.lookup.TenantID != tenant.ID("tenant_1") {
		t.Fatalf("unexpected audit lookup: %+v", service.lookup)
	}
	if !strings.Contains(response.Body.String(), "before_snapshot_redacted") {
		t.Fatalf("detail response should include redacted snapshot fields: %s", response.Body.String())
	}
}

func TestHTTPHandlerRejectsBadAuditTime(t *testing.T) {
	service := &fakeAuditHTTPService{}
	handler := registerAuditTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?created_from=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no service call, got %d", service.listCalls)
	}
}

func registerAuditTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

func testAuditLog() Log {
	return Log{
		ID:                     ID("audit_1"),
		DisplayID:              70001,
		TenantID:               tenant.ID("tenant_1"),
		ActorID:                ActorID("actor_1"),
		ActorType:              ActorTypeUser,
		Action:                 "invoice.paid",
		TargetType:             "invoice",
		TargetID:               TargetID("target_1"),
		BeforeSnapshotRedacted: []byte(`{"status":"issued"}`),
		AfterSnapshotRedacted:  []byte(`{"status":"paid"}`),
		MetadataRedacted:       []byte(`{"source":"test"}`),
		CorrelationID:          CorrelationID("correlation_1"),
		CreatedAt:              time.Date(2026, 4, 23, 1, 0, 0, 0, time.UTC),
	}
}

type fakeAuditHTTPService struct {
	log       Log
	logs      []Log
	filter    Filter
	lookup    Lookup
	listCalls int
	getCalls  int
}

func (service *fakeAuditHTTPService) ListLogs(ctx context.Context, filter Filter) ([]Log, error) {
	service.listCalls++
	service.filter = filter
	return service.logs, nil
}

func (service *fakeAuditHTTPService) GetLog(ctx context.Context, lookup Lookup) (Log, error) {
	service.getCalls++
	service.lookup = lookup
	return service.log, nil
}
