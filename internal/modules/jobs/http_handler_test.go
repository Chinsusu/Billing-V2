package jobs

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

func TestHTTPHandlerListAdminJobsUsesFilters(t *testing.T) {
	service := &fakeJobsHTTPService{jobs: []Job{testReadJob()}}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs?display_id=81001&job_type=provider.provision&status=failed_retryable&reference_type=order&reference_id=order_1&source_id=source_1&limit=20", nil)
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
		service.filter.DisplayID != 81001 ||
		service.filter.Type != Type("provider.provision") ||
		service.filter.Status != StatusFailedRetryable ||
		service.filter.ReferenceType != ReferenceType("order") ||
		service.filter.ReferenceID != ReferenceID("order_1") ||
		service.filter.SourceID != SourceID("source_1") ||
		service.filter.Limit != 20 {
		t.Fatalf("unexpected job filter: %+v", service.filter)
	}
	body := response.Body.String()
	if !strings.Contains(body, `"display_id":81001`) || !strings.Contains(body, `"job_type":"provider.provision"`) {
		t.Fatalf("expected job response, got %s", body)
	}
	if strings.Contains(body, "payload") || strings.Contains(body, "idempotency") {
		t.Fatalf("job response should not expose payload or idempotency key: %s", body)
	}
}

func TestHTTPHandlerGetResellerJobUsesTenantScope(t *testing.T) {
	service := &fakeJobsHTTPService{job: testReadJob()}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/jobs/job_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.getCalls != 1 {
		t.Fatalf("expected get once, got %d", service.getCalls)
	}
	if service.lookup.ID != ID("job_1") || service.lookup.TenantID != tenant.ID("tenant_1") {
		t.Fatalf("unexpected job lookup: %+v", service.lookup)
	}
}

func TestHTTPHandlerRejectsBadJobStatus(t *testing.T) {
	service := &fakeJobsHTTPService{}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs?status=lost", nil)
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
	if !strings.Contains(response.Body.String(), "job.status_invalid") {
		t.Fatalf("expected status validation response, got %s", response.Body.String())
	}
}

func registerJobsTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

func testReadJob() Job {
	return Job{
		ID:                       "job_1",
		DisplayID:                81001,
		TenantID:                 "tenant_1",
		Type:                     "provider.provision",
		ReferenceType:            "order",
		ReferenceID:              "order_1",
		SourceID:                 "source_1",
		Status:                   StatusFailedRetryable,
		Priority:                 50,
		AttemptCount:             2,
		MaxAttempts:              5,
		NextAttemptAt:            time.Date(2026, 4, 24, 1, 0, 0, 0, time.UTC),
		LastErrorCode:            "provider_timeout",
		LastErrorMessageRedacted: "provider timed out",
		ManualReviewReason:       "needs check",
		CorrelationID:            "correlation_1",
		CreatedAt:                time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC),
		UpdatedAt:                time.Date(2026, 4, 24, 0, 30, 0, 0, time.UTC),
	}
}

type fakeJobsHTTPService struct {
	job       Job
	jobs      []Job
	filter    Filter
	lookup    Lookup
	listCalls int
	getCalls  int
}

func (service *fakeJobsHTTPService) ListJobs(ctx context.Context, filter Filter) ([]Job, error) {
	service.listCalls++
	service.filter = filter
	return service.jobs, nil
}

func (service *fakeJobsHTTPService) GetJob(ctx context.Context, lookup Lookup) (Job, error) {
	service.getCalls++
	service.lookup = lookup
	return service.job, nil
}
