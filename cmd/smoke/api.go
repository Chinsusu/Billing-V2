package main

import (
	"context"
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
	Name     string
	Path     string
	Headers  map[string]string
	Contains []string
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
		if err := runAPICheck(ctx, client, baseURL, check); err != nil {
			return err
		}
		fmt.Printf("api check passed: %s %s\n", check.Name, check.Path)
	}
	fmt.Printf("dev API smoke passed: %d check(s)\n", len(checks))
	return nil
}

func runAPICheck(ctx context.Context, client *http.Client, baseURL string, check apiSmokeCheck) error {
	fullURL, err := normalizedAPIURL(baseURL, check.Path)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("build request %q: %w", check.Name, err)
	}
	for key, value := range check.Headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request %q: %w", check.Name, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response %q: %w", check.Name, err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("check %q expected HTTP 200, got %d: %s", check.Name, response.StatusCode, strings.TrimSpace(string(body)))
	}
	bodyText := string(body)
	for _, expected := range check.Contains {
		if !strings.Contains(bodyText, expected) {
			return fmt.Errorf("check %q response missing %q: %s", check.Name, expected, strings.TrimSpace(bodyText))
		}
	}
	return nil
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
		{Name: "admin service list", Path: "/admin/services", Headers: admin, Contains: []string{`"display_id":43001`}},
		{Name: "admin service detail", Path: "/admin/services/00000000-0000-0000-0000-000000000909", Headers: admin, Contains: []string{`"billing_status":"paid"`}},
		{Name: "admin wallet list", Path: "/admin/wallets", Headers: admin, Contains: []string{`"display_id":41001`}},
		{Name: "admin wallet detail", Path: "/admin/wallets/00000000-0000-0000-0000-000000000901", Headers: admin, Contains: []string{`"currency":"USD"`}},
		{Name: "admin topup list", Path: "/admin/topup-requests", Headers: admin, Contains: []string{`"display_id":52001`}},
		{Name: "admin topup detail", Path: "/admin/topup-requests/00000000-0000-0000-0000-000000000908", Headers: admin, Contains: []string{`"status":"approved"`}},
		{Name: "admin transaction list", Path: "/admin/transactions", Headers: admin, Contains: []string{`"display_id":51001`}},
		{Name: "admin transaction detail", Path: "/admin/transactions/00000000-0000-0000-0000-000000000907", Headers: admin, Contains: []string{`"type":"charge"`}},
		{Name: "admin reconciliation list", Path: "/admin/payment-reconciliation", Headers: admin, Contains: []string{`"provider":"wallet"`, `"display_id":51001`}},
		{Name: "admin reconciliation detail", Path: "/admin/payment-reconciliation/00000000-0000-0000-0000-000000000907", Headers: admin, Contains: []string{`"wallet_display_id":41001`}},
		{Name: "admin audit list", Path: "/admin/audit-logs", Headers: admin},
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
