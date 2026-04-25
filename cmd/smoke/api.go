package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	demoTenantID   = "00000000-0000-0000-0000-000000000010"
	demoResellerID = "00000000-0000-0000-0000-000000000102"
	demoCustomerID = "00000000-0000-0000-0000-000000000103"
)

type apiSmokeCheck struct {
	Name                string
	Path                string
	Headers             map[string]string
	Contains            []string
	NotContains         []string
	RedactBodyOnFailure bool
	SummaryFields       []string
}

func runDevAPISmoke(baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-api smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{Timeout: timeout}
	checks := apiSmokeChecks()
	for _, check := range checks {
		summary, err := runAPICheck(ctx, client, baseURL, check)
		if err != nil {
			return err
		}
		if summary != "" {
			fmt.Printf("api check passed: %s %s %s\n", check.Name, check.Path, summary)
			continue
		}
		fmt.Printf("api check passed: %s %s\n", check.Name, check.Path)
	}
	rbacChecks := apiRBACNegativeChecks()
	for _, check := range rbacChecks {
		if err := runAPIRBACNegativeCheck(ctx, client, baseURL, check); err != nil {
			return err
		}
		fmt.Printf("api RBAC negative check passed: %s %s %s\n", check.Name, check.Method, check.Path)
	}
	fmt.Printf("dev API smoke passed: %d check(s)\n", len(checks)+len(rbacChecks))
	return nil
}

func runAPICheck(ctx context.Context, client *http.Client, baseURL string, check apiSmokeCheck) (string, error) {
	fullURL, err := normalizedAPIURL(baseURL, check.Path)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request %q: %w", check.Name, err)
	}
	for key, value := range check.Headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("request %q: %w", check.Name, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read response %q: %w", check.Name, err)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("check %q expected HTTP 200, got %d: %s", check.Name, response.StatusCode, responseFailureBody(check, string(body)))
	}
	bodyText := string(body)
	for _, expected := range check.Contains {
		if !strings.Contains(bodyText, expected) {
			return "", fmt.Errorf("check %q response missing %q: %s", check.Name, expected, responseFailureBody(check, bodyText))
		}
	}
	if err := assertResponseOmitsTokens(check.Name, bodyText, check.NotContains); err != nil {
		return "", err
	}
	summary, err := apiResponseSummary(check, bodyText)
	if err != nil {
		return "", err
	}
	return summary, nil
}

func responseFailureBody(check apiSmokeCheck, bodyText string) string {
	if check.RedactBodyOnFailure {
		return "response body omitted for sensitive smoke check"
	}
	return strings.TrimSpace(bodyText)
}

func apiResponseSummary(check apiSmokeCheck, bodyText string) (string, error) {
	if len(check.SummaryFields) == 0 {
		return "", nil
	}
	var body any
	if err := json.Unmarshal([]byte(bodyText), &body); err != nil {
		return "", fmt.Errorf("check %q response summary is not valid JSON", check.Name)
	}
	values := firstJSONFields(body, check.SummaryFields)
	if len(values) == 0 {
		return "", fmt.Errorf("check %q response summary missing display IDs", check.Name)
	}
	parts := make([]string, 0, len(check.SummaryFields))
	for _, field := range check.SummaryFields {
		value, ok := values[field]
		if !ok {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", field, formatJSONSummaryValue(value)))
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("check %q response summary missing display IDs", check.Name)
	}
	return "display_ids " + strings.Join(parts, " "), nil
}

func firstJSONFields(value any, fields []string) map[string]any {
	switch typed := value.(type) {
	case map[string]any:
		matches := make(map[string]any, len(fields))
		for _, field := range fields {
			if value, ok := typed[field]; ok {
				matches[field] = value
			}
		}
		if len(matches) > 0 {
			return matches
		}
		for _, nested := range typed {
			if matches := firstJSONFields(nested, fields); len(matches) > 0 {
				return matches
			}
		}
	case []any:
		for _, item := range typed {
			if matches := firstJSONFields(item, fields); len(matches) > 0 {
				return matches
			}
		}
	}
	return nil
}

func formatJSONSummaryValue(value any) string {
	switch typed := value.(type) {
	case float64:
		if typed == float64(int64(typed)) {
			return fmt.Sprintf("%.0f", typed)
		}
		return fmt.Sprintf("%g", typed)
	case string:
		return typed
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func normalizedAPIURL(baseURL string, path string) (string, error) {
	parsed, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "", fmt.Errorf("parse API base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("API base URL must include scheme and host")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return parsed.String() + path, nil
}

func apiSmokeChecks() []apiSmokeCheck {
	client := clientHeaders()
	admin := adminHeaders()
	return []apiSmokeCheck{
		{Name: "health", Path: "/healthz", Contains: []string{`"status":"ok"`}},
		{Name: "readiness", Path: "/readyz", Contains: []string{`"status":"ready"`}},
		{Name: "client wallet list", Path: "/client/wallets", Headers: client, Contains: []string{`"display_id":41001`}},
		{Name: "client wallet detail", Path: "/client/wallets/00000000-0000-0000-0000-000000000901", Headers: client, Contains: []string{`"available_balance_minor":3200`}},
		{Name: "client wallet ledger", Path: "/client/wallets/00000000-0000-0000-0000-000000000901/ledger", Headers: client, Contains: []string{`"display_id":50001`, `"display_id":50002`}},
		{Name: "client order list", Path: "/client/orders", Headers: client, Contains: []string{`"display_id":42001`}},
		{Name: "client order detail", Path: "/client/orders/00000000-0000-0000-0000-000000000903", Headers: client, Contains: []string{`"order_status":"paid"`}},
		{Name: "client service list", Path: "/client/services", Headers: client, Contains: []string{`"display_id":43001`}},
		{Name: "client service detail", Path: "/client/services/00000000-0000-0000-0000-000000000909", Headers: client, Contains: []string{`"status":"active"`}},
		{Name: "client invoice list", Path: "/client/invoices", Headers: client, Contains: []string{`"display_id":44001`}},
		{Name: "client invoice detail", Path: "/client/invoices/00000000-0000-0000-0000-000000000904", Headers: client, Contains: []string{`"status":"paid"`, `"items":[`}},
		{Name: "client transaction list", Path: "/client/transactions", Headers: client, Contains: []string{`"display_id":51001`}},
		{Name: "client transaction detail", Path: "/client/transactions/00000000-0000-0000-0000-000000000907", Headers: client, Contains: []string{`"status":"posted"`}},
		{Name: "admin order list", Path: "/admin/orders", Headers: admin, Contains: []string{`"display_id":42001`}},
		{Name: "admin order detail", Path: "/admin/orders/00000000-0000-0000-0000-000000000903", Headers: admin, Contains: []string{`"billing_status":"paid"`}},
		{Name: "admin service list", Path: "/admin/services", Headers: admin, Contains: []string{`"display_id":43001`, `"order_display_id":42001`, `"buyer_display_id":10002`, `"provider_source_display_id":10000`}},
		{Name: "admin service detail", Path: "/admin/services/00000000-0000-0000-0000-000000000909", Headers: admin, Contains: []string{`"billing_status":"paid"`, `"order_display_id":42001`, `"buyer_display_id":10002`, `"provider_source_display_id":10000`}},
		{Name: "admin service public id filter", Path: "/admin/services?display_id=43001&order_display_id=42001&provider_source_display_id=10000", Headers: admin, Contains: []string{`"display_id":43001`, `"order_display_id":42001`, `"provider_source_display_id":10000`}},
		{Name: "admin wallet list", Path: "/admin/wallets", Headers: admin, Contains: []string{`"display_id":41001`}},
		{Name: "admin wallet detail", Path: "/admin/wallets/00000000-0000-0000-0000-000000000901", Headers: admin, Contains: []string{`"currency":"USD"`}},
		{Name: "admin topup list", Path: "/admin/topup-requests", Headers: admin, Contains: []string{`"display_id":52001`, `"wallet_display_id":41001`, `"requested_by_display_id":10002`, `"reviewed_by_display_id":10001`}},
		{Name: "admin topup detail", Path: "/admin/topup-requests/00000000-0000-0000-0000-000000000908", Headers: admin, Contains: []string{`"status":"approved"`, `"wallet_display_id":41001`, `"requested_by_display_id":10002`, `"reviewed_by_display_id":10001`}},
		{Name: "admin transaction list", Path: "/admin/transactions", Headers: admin, Contains: []string{`"display_id":51001`, `"account_display_id":10002`, `"order_display_id":42001`, `"invoice_display_id":44001`}},
		{Name: "admin transaction detail", Path: "/admin/transactions/00000000-0000-0000-0000-000000000907", Headers: admin, Contains: []string{`"type":"charge"`, `"account_display_id":10002`, `"order_display_id":42001`, `"invoice_display_id":44001`}},
		{Name: "admin invoice public id filter", Path: "/admin/invoices?display_id=44001&buyer_display_id=10002&order_display_id=42001", Headers: admin, Contains: []string{`"display_id":44001`, `"buyer_display_id":10002`, `"order_display_id":42001`}},
		{Name: "admin invoice public id filter miss", Path: "/admin/invoices?display_id=999999", Headers: admin, NotContains: []string{`"display_id":44001`}},
		{Name: "admin reconciliation list", Path: "/admin/payment-reconciliation", Headers: admin, Contains: []string{`"provider":"wallet"`, `"display_id":51001`}},
		{Name: "admin reconciliation detail", Path: "/admin/payment-reconciliation/00000000-0000-0000-0000-000000000907", Headers: admin, Contains: []string{`"wallet_display_id":41001`}},
		{Name: "admin job public id filter", Path: "/admin/jobs?job_type=provider.provision&display_id=53001&source_display_id=10000", Headers: admin, Contains: []string{`"display_id":53001`, `"source_display_id":10000`, `"reference_display_id":42001`}},
		{
			Name:    "admin provider readiness",
			Path:    "/admin/catalog/provider-readiness?status=active&limit=20",
			Headers: admin,
			Contains: []string{
				`"plan_display_id":`,
				`"source_display_id":`,
				`"state":`,
				`"reason":`,
			},
			NotContains:         sensitiveAPIRedactionTokens(),
			RedactBodyOnFailure: true,
			SummaryFields:       []string{"plan_display_id", "source_display_id"},
		},
		{Name: "admin audit list", Path: "/admin/audit-logs", Headers: admin, Contains: []string{`"display_id":70001`, `"actor_display_id":10001`, `"target_display_id":53001`}},
	}
}

func clientHeaders() map[string]string {
	return actorHeaders(demoCustomerID, "client")
}

func adminHeaders() map[string]string {
	return actorHeaders(demoResellerID, "reseller_owner")
}

func actorHeaders(actorID string, actorType string) map[string]string {
	return map[string]string{
		"X-Tenant-Id":       demoTenantID,
		"X-Actor-Id":        actorID,
		"X-Actor-Type":      actorType,
		"X-Actor-Tenant-Id": demoTenantID,
	}
}
