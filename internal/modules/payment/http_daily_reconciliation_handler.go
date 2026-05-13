package payment

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const adminDailyReconciliationPath = "/admin/daily-reconciliation"

type dailyReconciliationHTTPService interface {
	BuildDailyReconciliationReport(ctx context.Context, input DailyReconciliationInput) (DailyReconciliationReport, error)
}

func (handler *HTTPHandler) adminDailyReconciliationRoute(w http.ResponseWriter, r *http.Request) {
	dispatchPaymentMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleGetAdminDailyReconciliation, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) handleGetAdminDailyReconciliation(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := actorFromContext(w, r); !ok {
		return
	}
	date, ok := dailyReconciliationDateFromRequest(w, r)
	if !ok {
		return
	}
	service, ok := handler.service.(dailyReconciliationHTTPService)
	if !ok {
		writePaymentError(w, r, ErrBillingDependencyMissing)
		return
	}
	report, err := service.BuildDailyReconciliationReport(r.Context(), DailyReconciliationInput{
		TenantID: tenantID,
		Date:     date,
	})
	if err != nil {
		writePaymentError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newDailyReconciliationReportResponse(report))
}

func dailyReconciliationDateFromRequest(w http.ResponseWriter, r *http.Request) (time.Time, bool) {
	value := strings.TrimSpace(r.URL.Query().Get("date"))
	if value == "" {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("date", "payment.created_time_invalid", "Date is required in YYYY-MM-DD format."),
		})
		return time.Time{}, false
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("date", "payment.created_time_invalid", "Date must use YYYY-MM-DD format."),
		})
		return time.Time{}, false
	}
	return parsed, true
}
