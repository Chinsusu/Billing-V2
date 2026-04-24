package jobs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const (
	jobActionRetry        = "retry"
	jobActionManualReview = "manual-review"
	jobActionCancel       = "cancel"
)

type retryJobRequest struct {
	NextAttemptAt time.Time `json:"next_attempt_at,omitempty"`
}

type manualReviewJobRequest struct {
	Reason string `json:"reason"`
}

type cancelJobRequest struct {
	Reason string `json:"reason,omitempty"`
}

func (handler *HTTPHandler) adminJobRecoveryRoute(w http.ResponseWriter, r *http.Request) {
	_, action, ok := jobIDAndRecoveryActionFromPath(w, r, adminJobPrefix)
	if !ok {
		return
	}
	switch action {
	case jobActionRetry:
		dispatchJobMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleRetryAdminJob, handler.options.AdminRetryMiddleware),
		})
	case jobActionManualReview:
		dispatchJobMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleMarkAdminJobManualReview, handler.options.AdminManualReviewMiddleware),
		})
	case jobActionCancel:
		dispatchJobMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(handler.handleCancelAdminJob, handler.options.AdminCancelMiddleware),
		})
	default:
		writeJobError(w, r, ErrJobNotFound)
	}
}

func (handler *HTTPHandler) handleRetryAdminJob(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := jobTenantAndActor(w, r)
	if !ok {
		return
	}
	jobID, _, ok := jobIDAndRecoveryActionFromPath(w, r, adminJobPrefix)
	if !ok {
		return
	}
	var request retryJobRequest
	if !decodeJobJSON(w, r, &request) {
		return
	}
	job, err := handler.service.RetryJob(r.Context(), request.toInput(jobID, tenantID, actor.ID))
	if err != nil {
		writeJobError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newJobResponse(job))
}

func (handler *HTTPHandler) handleMarkAdminJobManualReview(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := jobTenantAndActor(w, r)
	if !ok {
		return
	}
	jobID, _, ok := jobIDAndRecoveryActionFromPath(w, r, adminJobPrefix)
	if !ok {
		return
	}
	var request manualReviewJobRequest
	if !decodeJobJSON(w, r, &request) {
		return
	}
	job, err := handler.service.MarkManualReview(r.Context(), request.toInput(jobID, tenantID, actor.ID))
	if err != nil {
		writeJobError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newJobResponse(job))
}

func (handler *HTTPHandler) handleCancelAdminJob(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := jobTenantAndActor(w, r)
	if !ok {
		return
	}
	jobID, _, ok := jobIDAndRecoveryActionFromPath(w, r, adminJobPrefix)
	if !ok {
		return
	}
	var request cancelJobRequest
	if !decodeJobJSON(w, r, &request) {
		return
	}
	job, err := handler.service.CancelJob(r.Context(), request.toInput(jobID, tenantID, actor.ID))
	if err != nil {
		writeJobError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newJobResponse(job))
}

func (request retryJobRequest) toInput(jobID ID, tenantID tenant.ID, actorID identity.UserID) RetryJobInput {
	return RetryJobInput{ID: jobID, TenantID: tenantID, ActorID: actorID, NextAttemptAt: request.NextAttemptAt}
}

func (request manualReviewJobRequest) toInput(jobID ID, tenantID tenant.ID, actorID identity.UserID) ManualReviewJobInput {
	return ManualReviewJobInput{ID: jobID, TenantID: tenantID, ActorID: actorID, Reason: request.Reason}
}

func (request cancelJobRequest) toInput(jobID ID, tenantID tenant.ID, actorID identity.UserID) CancelJobInput {
	return CancelJobInput{ID: jobID, TenantID: tenantID, ActorID: actorID, Reason: request.Reason}
}

func jobTenantAndActor(w http.ResponseWriter, r *http.Request) (tenant.ID, identity.Actor, bool) {
	tenantID, ok := jobTenantIDFromContext(w, r)
	if !ok {
		return "", identity.Actor{}, false
	}
	actor, ok := jobActorFromContext(w, r)
	if !ok {
		return "", identity.Actor{}, false
	}
	return tenantID, actor, true
}

func jobRecoveryActionPath(path string, prefix string) bool {
	_, action, ok := splitJobRecoveryPath(path, prefix)
	return ok && knownJobRecoveryAction(action)
}

func jobIDAndRecoveryActionFromPath(w http.ResponseWriter, r *http.Request, prefix string) (ID, string, bool) {
	jobID, action, ok := splitJobRecoveryPath(r.URL.Path, prefix)
	if !ok || !knownJobRecoveryAction(action) {
		writeJobError(w, r, ErrJobIDMissing)
		return "", "", false
	}
	return jobID, action, true
}

func splitJobRecoveryPath(path string, prefix string) (ID, string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", "", false
	}
	value := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.Split(value, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return ID(strings.TrimSpace(parts[0])), strings.TrimSpace(parts[1]), true
}

func knownJobRecoveryAction(action string) bool {
	switch action {
	case jobActionRetry, jobActionManualReview, jobActionCancel:
		return true
	default:
		return false
	}
}

func decodeJobJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request.invalid_json", "Request body must be valid JSON.")
		return false
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return true
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
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
