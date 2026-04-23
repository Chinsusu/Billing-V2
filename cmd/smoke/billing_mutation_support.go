package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type successEnvelope[T any] struct {
	Data T `json:"data"`
}

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type topupRequestBody struct {
	WalletID         string `json:"wallet_id"`
	AmountMinor      int64  `json:"amount_minor"`
	Currency         string `json:"currency"`
	PaymentMethod    string `json:"payment_method"`
	PaymentReference string `json:"payment_reference"`
}

type reviewTopupBody struct {
	ReviewNote string `json:"review_note"`
}

type transitionOrderStatusBody struct {
	FromStatus    string `json:"from_status"`
	ToStatus      string `json:"to_status"`
	BillingStatus string `json:"billing_status"`
}

type createOrderBody struct {
	TenantPlanID    string          `json:"tenant_plan_id"`
	Quantity        int             `json:"quantity"`
	Currency        string          `json:"currency"`
	UnitPriceMinor  int64           `json:"unit_price_minor"`
	DiscountMinor   int64           `json:"discount_minor"`
	TotalMinor      int64           `json:"total_minor"`
	ProductSnapshot json.RawMessage `json:"product_snapshot"`
	PlanSnapshot    json.RawMessage `json:"plan_snapshot"`
	PriceSnapshot   json.RawMessage `json:"price_snapshot"`
}

type invoiceWalletPaymentBody struct {
	InvoiceID string `json:"invoice_id"`
	WalletID  string `json:"wallet_id"`
}

type topupResponse struct {
	ID            string `json:"id"`
	DisplayID     int64  `json:"display_id"`
	Status        string `json:"status"`
	LedgerEntryID string `json:"ledger_entry_id,omitempty"`
}

type orderResponse struct {
	ID            string `json:"id"`
	DisplayID     int64  `json:"display_id"`
	OrderStatus   string `json:"order_status"`
	BillingStatus string `json:"billing_status"`
}

type invoiceSummaryResponse struct {
	ID        string `json:"id"`
	DisplayID int64  `json:"display_id"`
	Status    string `json:"status"`
	OrderID   string `json:"order_id"`
}

type invoiceWalletPaymentResponse struct {
	Invoice struct {
		ID         string `json:"id"`
		DisplayID  int64  `json:"display_id"`
		Status     string `json:"status"`
		TotalMinor int64  `json:"total_minor"`
		Currency   string `json:"currency"`
	} `json:"invoice"`
	Transaction struct {
		ID        string `json:"id"`
		DisplayID int64  `json:"display_id"`
		Status    string `json:"status"`
	} `json:"transaction"`
	Ledger *struct {
		ID          string `json:"id"`
		DisplayID   int64  `json:"display_id"`
		AmountMinor int64  `json:"amount_minor"`
	} `json:"ledger,omitempty"`
}

type auditSummaryResponse struct {
	ID         string `json:"id"`
	DisplayID  int64  `json:"display_id"`
	TenantID   string `json:"tenant_id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
}

type auditDetailResponse struct {
	ID                    string          `json:"id"`
	DisplayID             int64           `json:"display_id"`
	TenantID              string          `json:"tenant_id"`
	Action                string          `json:"action"`
	TargetType            string          `json:"target_type"`
	TargetID              string          `json:"target_id"`
	AfterSnapshotRedacted json.RawMessage `json:"after_snapshot_redacted"`
	MetadataRedacted      json.RawMessage `json:"metadata_redacted"`
	CorrelationID         string          `json:"correlation_id"`
}

type auditMutationCheck struct {
	Action           string
	TargetType       string
	TargetID         string
	MetadataContains string
	AfterContains    string
}

func newBillingMutationScenario() billingMutationScenario {
	return billingMutationScenario{RunID: fmt.Sprintf("%d", time.Now().UTC().UnixNano())}
}

func (scenario billingMutationScenario) topupIdempotencyKey() string {
	return "smoke-topup-" + scenario.RunID
}

func (scenario billingMutationScenario) topupPaymentReference() string {
	return "SMOKE-TOPUP-" + scenario.RunID
}

func (scenario billingMutationScenario) orderIdempotencyKey() string {
	return "smoke-order-" + scenario.RunID
}

func (scenario billingMutationScenario) invoiceIdempotencyKey() string {
	return "smoke-invoice-" + scenario.RunID
}

func (scenario billingMutationScenario) paymentIdempotencyKey() string {
	return "smoke-payment-" + scenario.RunID
}

func paymentLedgerDisplayID(record invoiceWalletPaymentResponse) int64 {
	if record.Ledger == nil {
		return 0
	}
	return record.Ledger.DisplayID
}

func cloneHeaders(source map[string]string) map[string]string {
	headers := make(map[string]string, len(source))
	for key, value := range source {
		headers[key] = value
	}
	return headers
}

func doJSON[T any](
	ctx context.Context,
	client *http.Client,
	method string,
	baseURL string,
	path string,
	headers map[string]string,
	requestBody interface{},
	wantStatus int,
) (T, error) {
	var zero T
	fullURL, err := normalizedAPIURL(baseURL, path)
	if err != nil {
		return zero, err
	}

	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return zero, fmt.Errorf("marshal request %s %s: %w", method, path, err)
		}
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return zero, fmt.Errorf("build request %s %s: %w", method, path, err)
	}
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return zero, fmt.Errorf("request %s %s: %w", method, path, err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return zero, fmt.Errorf("read response %s %s: %w", method, path, err)
	}
	if response.StatusCode != wantStatus {
		var apiError errorEnvelope
		if err := json.Unmarshal(payload, &apiError); err == nil && apiError.Error.Code != "" {
			return zero, fmt.Errorf("%s %s expected HTTP %d, got %d (%s: %s)", method, path, wantStatus, response.StatusCode, apiError.Error.Code, apiError.Error.Message)
		}
		return zero, fmt.Errorf("%s %s expected HTTP %d, got %d: %s", method, path, wantStatus, response.StatusCode, strings.TrimSpace(string(payload)))
	}

	var envelope successEnvelope[T]
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return zero, fmt.Errorf("decode response %s %s: %w", method, path, err)
	}
	return envelope.Data, nil
}
