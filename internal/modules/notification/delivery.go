package notification

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

type LocalDeliveryHandler struct {
	Store DeliveryStore
	Now   func() time.Time
}

func NewLocalDeliveryHandler(store DeliveryStore) *LocalDeliveryHandler {
	return &LocalDeliveryHandler{Store: store}
}

func NewLocalDeliveryRunner(store jobs.Store, deliveryStore DeliveryStore, workerID jobs.WorkerID) jobs.Runner {
	return jobs.Runner{
		Store:     store,
		Handler:   NewLocalDeliveryHandler(deliveryStore),
		WorkerID:  workerID,
		BatchSize: 10,
		Types:     []jobs.Type{DeliveryJobType},
	}
}

func (handler *LocalDeliveryHandler) Handle(ctx context.Context, job jobs.Job) (jobs.Completion, error) {
	if handler == nil || handler.Store == nil {
		return jobs.Completion{}, ErrStoreMissing
	}
	if job.ReferenceType != DeliveryReferenceType || job.ReferenceID == "" {
		return jobs.Completion{
			Status:                   jobs.StatusFailedTerminal,
			RetrySafety:              jobs.RetrySafetyDoNotRetry,
			LastErrorCode:            "notification_job_invalid",
			LastErrorMessageRedacted: "notification delivery job is invalid",
			FinishedAt:               handler.now(),
		}, nil
	}
	if _, err := handler.Store.MarkNotificationSent(ctx, ID(job.ReferenceID), handler.now()); err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			return jobs.Completion{
				Status:                   jobs.StatusFailedTerminal,
				RetrySafety:              jobs.RetrySafetyDoNotRetry,
				LastErrorCode:            "notification_not_found",
				LastErrorMessageRedacted: "notification record was not found",
				FinishedAt:               handler.now(),
			}, nil
		}
		return jobs.Completion{}, err
	}
	return jobs.Completion{Status: jobs.StatusSucceeded, FinishedAt: handler.now()}, nil
}

func (handler *LocalDeliveryHandler) now() time.Time {
	if handler.Now == nil {
		return time.Now().UTC()
	}
	return handler.Now()
}
