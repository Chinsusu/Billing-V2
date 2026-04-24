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

func TestHTTPHandlerListAdminJobAttemptsUsesTenantScope(t *testing.T) {
	service := &fakeJobsHTTPService{attempts: []Attempt{testReadAttempt()}}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs/job_1/attempts?limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listAttemptCalls != 1 {
		t.Fatalf("expected list attempts once, got %d", service.listAttemptCalls)
	}
	if service.attemptFilter.JobID != ID("job_1") ||
		service.attemptFilter.TenantID != tenant.ID("tenant_1") ||
		service.attemptFilter.Limit != 10 {
		t.Fatalf("unexpected attempt filter: %+v", service.attemptFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{`"display_id":82001`, `"worker_id":"worker_1"`, `"duration_ms":2500`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in attempt response, got %s", expected, body)
		}
	}
}

func TestHTTPHandlerAdminJobSummaryUsesTenantScope(t *testing.T) {
	service := &fakeJobsHTTPService{summary: testJobSummary()}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs/summary?job_type=provider.provision", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.summaryCalls != 1 {
		t.Fatalf("expected summary once, got %d", service.summaryCalls)
	}
	if service.summaryFilter.TenantID != tenant.ID("tenant_1") ||
		service.summaryFilter.Type != Type("provider.provision") {
		t.Fatalf("unexpected summary filter: %+v", service.summaryFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{
		`"job_type":"provider.provision"`,
		`"total":9`,
		`"attention_count":3`,
		`"failed_retryable":1`,
		`"manual_review":1`,
		`"failed_terminal":1`,
		`"oldest_queued_age_seconds":7200`,
		`"display_id":81009`,
		`"last_error_message_redacted":"provider timed out"`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in summary response, got %s", expected, body)
		}
	}
	if strings.Contains(body, "payload") || strings.Contains(body, "idempotency") || strings.Contains(body, "secret") {
		t.Fatalf("summary response should not expose internal or secret fields: %s", body)
	}
}

func TestHTTPHandlerAdminJobSummaryUsesSummaryMiddleware(t *testing.T) {
	service := &fakeJobsHTTPService{summary: testJobSummary()}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		AdminSummaryMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs/summary", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusTeapot {
		t.Fatalf("expected summary middleware status, got %d: %s", response.Code, response.Body.String())
	}
	if service.summaryCalls != 0 {
		t.Fatalf("expected middleware to stop summary service, got %d calls", service.summaryCalls)
	}
}

func TestHTTPHandlerListResellerJobAttemptsUsesTenantScope(t *testing.T) {
	service := &fakeJobsHTTPService{attempts: []Attempt{testReadAttempt()}}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/jobs/job_1/attempts?limit=7", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_2")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "tenant_2", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listAttemptCalls != 1 {
		t.Fatalf("expected list attempts once, got %d", service.listAttemptCalls)
	}
	if service.attemptFilter.JobID != ID("job_1") ||
		service.attemptFilter.TenantID != tenant.ID("tenant_2") ||
		service.attemptFilter.Limit != 7 {
		t.Fatalf("unexpected attempt filter: %+v", service.attemptFilter)
	}
}

func TestHTTPHandlerSummaryMissingServiceReturnsEnvelope(t *testing.T) {
	handler := registerJobsTestHandler(nil)

	request := httptest.NewRequest(http.MethodGet, "/admin/jobs/summary", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "job.service_unavailable") {
		t.Fatalf("expected service unavailable envelope, got %s", response.Body.String())
	}
}

func TestHTTPHandlerRetryAdminJobUsesTenantScope(t *testing.T) {
	service := &fakeJobsHTTPService{job: testReadJob()}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/jobs/job_1/retry", strings.NewReader(`{"next_attempt_at":"2026-04-24T02:00:00Z"}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.retryCalls != 1 {
		t.Fatalf("expected retry once, got %d", service.retryCalls)
	}
	if service.retryInput.ID != ID("job_1") ||
		service.retryInput.TenantID != tenant.ID("tenant_1") ||
		service.retryInput.ActorID != identity.UserID("admin_1") ||
		!service.retryInput.NextAttemptAt.Equal(time.Date(2026, 4, 24, 2, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected retry input: %+v", service.retryInput)
	}
}

func TestHTTPHandlerDoesNotExposeResellerRecoveryActions(t *testing.T) {
	service := &fakeJobsHTTPService{job: testReadJob()}
	handler := registerJobsTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/reseller/jobs/job_1/retry", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d: %s", response.Code, response.Body.String())
	}
	if service.retryCalls != 0 {
		t.Fatalf("expected no retry call, got %d", service.retryCalls)
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

func testJobSummary() JobSummary {
	return JobSummary{
		TenantID:       "tenant_1",
		Type:           "provider.provision",
		Total:          9,
		AttentionCount: 3,
		Counts: JobStatusCounts{
			Queued:          2,
			Claimed:         1,
			Running:         1,
			Succeeded:       2,
			FailedRetryable: 1,
			FailedTerminal:  1,
			ManualReview:    1,
		},
		OldestQueuedAt: time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC),
		GeneratedAt:    time.Date(2026, 4, 24, 2, 0, 0, 0, time.UTC),
		LatestFailure: &JobFailureContext{
			ID:                       "job_9",
			DisplayID:                81009,
			Status:                   StatusFailedRetryable,
			LastErrorCode:            "provider_timeout",
			LastErrorMessageRedacted: "provider timed out",
			CreatedAt:                time.Date(2026, 4, 24, 1, 0, 0, 0, time.UTC),
			UpdatedAt:                time.Date(2026, 4, 24, 1, 30, 0, 0, time.UTC),
		},
	}
}

func testReadAttempt() Attempt {
	return Attempt{
		ID:                   "attempt_1",
		DisplayID:            82001,
		JobID:                "job_1",
		WorkerID:             "worker_1",
		AttemptNumber:        2,
		StartedAt:            time.Date(2026, 4, 24, 1, 0, 0, 0, time.UTC),
		FinishedAt:           time.Date(2026, 4, 24, 1, 0, 2, 500000000, time.UTC),
		Result:               AttemptResultFailedRetryable,
		ErrorCode:            "provider_timeout",
		ErrorMessageRedacted: "provider timed out",
		Duration:             2500 * time.Millisecond,
		CorrelationID:        "correlation_1",
	}
}

type fakeJobsHTTPService struct {
	job               Job
	jobs              []Job
	attempts          []Attempt
	summary           JobSummary
	retryInput        RetryJobInput
	manualReviewInput ManualReviewJobInput
	cancelInput       CancelJobInput
	filter            Filter
	lookup            Lookup
	attemptFilter     AttemptFilter
	summaryFilter     SummaryFilter
	listCalls         int
	getCalls          int
	listAttemptCalls  int
	summaryCalls      int
	retryCalls        int
	manualReviewCalls int
	cancelCalls       int
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

func (service *fakeJobsHTTPService) ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error) {
	service.listAttemptCalls++
	service.attemptFilter = filter
	return service.attempts, nil
}

func (service *fakeJobsHTTPService) SummarizeJobs(ctx context.Context, filter SummaryFilter) (JobSummary, error) {
	service.summaryCalls++
	service.summaryFilter = filter
	return service.summary, nil
}

func (service *fakeJobsHTTPService) RetryJob(ctx context.Context, input RetryJobInput) (Job, error) {
	service.retryCalls++
	service.retryInput = input
	return service.job, nil
}

func (service *fakeJobsHTTPService) MarkManualReview(ctx context.Context, input ManualReviewJobInput) (Job, error) {
	service.manualReviewCalls++
	service.manualReviewInput = input
	return service.job, nil
}

func (service *fakeJobsHTTPService) CancelJob(ctx context.Context, input CancelJobInput) (Job, error) {
	service.cancelCalls++
	service.cancelInput = input
	return service.job, nil
}
