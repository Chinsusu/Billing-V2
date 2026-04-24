package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const (
	smokeProvisioningWorkerID = "smoke-provisioner"
	smokeProvisioningPasses   = 3
	smokeProvisioningBatch    = 100
)

type provisioningJobSmokeRecord struct {
	ID                       string
	DisplayID                int64
	Status                   string
	AttemptCount             int
	LastErrorCode            string
	LastErrorMessageRedacted string
}

type serviceInstanceSmokeResponse struct {
	ID                 string `json:"id"`
	DisplayID          int64  `json:"display_id"`
	OrderID            string `json:"order_id"`
	TenantPlanID       string `json:"tenant_plan_id"`
	ProviderSourceID   string `json:"provider_source_id"`
	ExternalResourceID string `json:"external_resource_id"`
	Status             string `json:"status"`
	BillingStatus      string `json:"billing_status"`
}

func runProvisioningFulfillmentSmoke(
	ctx context.Context,
	conn *sql.DB,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
	orderID string,
	orderDisplayID int64,
) (serviceInstanceSmokeResponse, error) {
	job, err := readProvisioningJobForOrder(ctx, conn, orderID)
	if err != nil {
		return serviceInstanceSmokeResponse{}, err
	}
	if !provisioningJobSmokeStatusOK(job.Status) {
		return serviceInstanceSmokeResponse{}, provisioningJobFailure(orderDisplayID, job, jobs.RunSummary{})
	}

	summary, err := runProvisioningWorkerUntilDone(ctx, conn, scenario, orderID)
	if err != nil {
		return serviceInstanceSmokeResponse{}, err
	}
	job, err = readProvisioningJobForOrder(ctx, conn, orderID)
	if err != nil {
		return serviceInstanceSmokeResponse{}, err
	}
	if job.Status != "succeeded" {
		return serviceInstanceSmokeResponse{}, provisioningJobFailure(orderDisplayID, job, summary)
	}

	service, err := verifyProvisionedServiceVisibleViaAPI(ctx, client, baseURL, orderID)
	if err != nil {
		return serviceInstanceSmokeResponse{}, err
	}
	fmt.Printf(
		"billing mutation passed: worker fulfilled order %d job=%d service=%d\n",
		orderDisplayID,
		job.DisplayID,
		service.DisplayID,
	)
	return service, nil
}

func runProvisioningWorkerUntilDone(
	ctx context.Context,
	conn *sql.DB,
	scenario billingMutationScenario,
	orderID string,
) (jobs.RunSummary, error) {
	total := jobs.RunSummary{}
	workerID := jobs.WorkerID(smokeProvisioningWorkerID + "-" + scenario.RunID)
	for pass := 0; pass < smokeProvisioningPasses; pass++ {
		summary, err := runProvisioningWorkerOnce(ctx, conn, workerID)
		if err != nil {
			return total, err
		}
		addRunSummary(&total, summary)
		job, err := readProvisioningJobForOrder(ctx, conn, orderID)
		if err != nil {
			return total, err
		}
		if job.Status == "succeeded" || !provisioningJobSmokeStatusOK(job.Status) || summary.Claimed == 0 {
			break
		}
	}
	return total, nil
}

func runProvisioningWorkerOnce(ctx context.Context, conn *sql.DB, workerID jobs.WorkerID) (jobs.RunSummary, error) {
	registry, err := provider.NewFakeRegistry()
	if err != nil {
		return jobs.RunSummary{}, err
	}
	runner := order.NewProviderProvisioningRunner(
		jobs.NewPostgresStore(conn),
		registry,
		order.NewPostgresStore(conn),
		workerID,
	)
	runner.BatchSize = smokeProvisioningBatch
	runner.LockFor = time.Minute
	return runner.RunOnce(ctx)
}

func readProvisioningJobForOrder(ctx context.Context, conn *sql.DB, orderID string) (provisioningJobSmokeRecord, error) {
	record := provisioningJobSmokeRecord{}
	err := conn.QueryRowContext(ctx, `
SELECT job_id, display_id, status::text, attempt_count, COALESCE(last_error_code, ''), COALESCE(last_error_message_redacted, '')
FROM jobs
WHERE tenant_id = $1
  AND job_type = 'provider.provision'
  AND reference_type = 'order'
  AND reference_id = $2
ORDER BY created_at DESC
LIMIT 1`, demoTenantID, orderID).Scan(
		&record.ID,
		&record.DisplayID,
		&record.Status,
		&record.AttemptCount,
		&record.LastErrorCode,
		&record.LastErrorMessageRedacted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return provisioningJobSmokeRecord{}, fmt.Errorf("expected provider.provision job for paid order %s, got none", orderID)
		}
		return provisioningJobSmokeRecord{}, fmt.Errorf("read provider.provision job for order %s: %w", orderID, err)
	}
	return record, nil
}

func verifyProvisionedServiceVisibleViaAPI(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	orderID string,
) (serviceInstanceSmokeResponse, error) {
	query := url.Values{}
	query.Set("order_id", orderID)
	query.Set("status", "active")
	query.Set("limit", "20")
	services, err := doJSON[[]serviceInstanceSmokeResponse](ctx, client, http.MethodGet, baseURL, "/client/services?"+query.Encode(), clientHeaders(), nil, http.StatusOK)
	if err != nil {
		return serviceInstanceSmokeResponse{}, err
	}
	for _, service := range services {
		if service.OrderID != orderID {
			continue
		}
		if service.DisplayID <= 0 || service.Status != "active" || service.BillingStatus != "paid" || service.ExternalResourceID == "" {
			return serviceInstanceSmokeResponse{}, fmt.Errorf("service for order %s has wrong state: %+v", orderID, service)
		}
		return service, nil
	}
	return serviceInstanceSmokeResponse{}, fmt.Errorf("expected active service for order %s after worker fulfillment, got %+v", orderID, services)
}

func provisioningJobFailure(orderDisplayID int64, job provisioningJobSmokeRecord, summary jobs.RunSummary) error {
	return fmt.Errorf(
		"expected provisioning job %d for order %d to succeed after worker, got status=%q attempts=%d error=%q message=%q summary=%+v",
		job.DisplayID,
		orderDisplayID,
		job.Status,
		job.AttemptCount,
		job.LastErrorCode,
		job.LastErrorMessageRedacted,
		summary,
	)
}

func addRunSummary(total *jobs.RunSummary, summary jobs.RunSummary) {
	total.Claimed += summary.Claimed
	total.Succeeded += summary.Succeeded
	total.Retried += summary.Retried
	total.ManualReview += summary.ManualReview
	total.TerminalFailed += summary.TerminalFailed
	total.Cancelled += summary.Cancelled
}
