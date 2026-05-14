package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultReferencePath = "docs/05_development_standards/56_Billing_API_Operational_Reference.md"

type routeContract struct {
	Method    string
	Path      string
	Perm      string
	Queries   []string
	DocTokens []string
	Sources   []sourceExpectation
}

type sourceExpectation struct {
	Path   string
	Tokens []string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("contractguard", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	reference := flags.String("reference", defaultReferencePath, "billing API reference markdown")
	if err := flags.Parse(args); err != nil {
		return err
	}
	failures, err := checkContracts(*root, *reference, trackedBillingContracts())
	if err != nil {
		return err
	}
	if len(failures) > 0 {
		return fmt.Errorf("API contract drift guard failed:\n%s", strings.Join(failures, "\n"))
	}
	fmt.Printf("API contract drift guard passed: %d route contract(s)\n", len(trackedBillingContracts()))
	return nil
}

func checkContracts(root string, referencePath string, contracts []routeContract) ([]string, error) {
	doc, err := readText(filepath.Join(root, referencePath))
	if err != nil {
		return nil, err
	}
	failures := make([]string, 0)
	for _, contract := range contracts {
		failures = append(failures, checkDocContract(doc, contract)...)
		sourceFailures, err := checkSourceContract(root, contract)
		if err != nil {
			return nil, err
		}
		failures = append(failures, sourceFailures...)
	}
	return failures, nil
}

func checkDocContract(doc string, contract routeContract) []string {
	routeName := contract.Method + " " + contract.Path
	section := routeSection(doc, routeName)
	if section == "" {
		return []string{fmt.Sprintf("- missing docs route: `%s`", routeName)}
	}
	failures := make([]string, 0)
	if contract.Perm != "" && !strings.Contains(section, "`"+contract.Perm+"`") {
		failures = append(failures, fmt.Sprintf("- `%s` docs missing permission `%s`", routeName, contract.Perm))
	}
	for _, query := range contract.Queries {
		if !strings.Contains(section, "`"+query+"`") {
			failures = append(failures, fmt.Sprintf("- `%s` docs missing query `%s`", routeName, query))
		}
	}
	for _, token := range contract.DocTokens {
		if !strings.Contains(section, token) && !strings.Contains(doc, token) {
			failures = append(failures, fmt.Sprintf("- `%s` docs missing note token %q", routeName, token))
		}
	}
	return failures
}

func checkSourceContract(root string, contract routeContract) ([]string, error) {
	failures := make([]string, 0)
	for _, source := range contract.Sources {
		body, err := readText(filepath.Join(root, source.Path))
		if err != nil {
			return nil, err
		}
		for _, token := range source.Tokens {
			if !strings.Contains(body, token) {
				failures = append(failures, fmt.Sprintf("- `%s %s` source %s missing %q", contract.Method, contract.Path, source.Path, token))
			}
		}
	}
	return failures, nil
}

func routeSection(doc string, routeName string) string {
	target := "`" + routeName + "`"
	lines := strings.Split(doc, "\n")
	start := -1
	for index, line := range lines {
		if strings.Contains(line, target) && strings.HasPrefix(strings.TrimSpace(line), "- `") {
			start = index
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for index := start + 1; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "- `") || strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "## ") {
			end = index
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

func readText(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("required file is missing: %s", path)
		}
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(body), nil
}

func trackedBillingContracts() []routeContract {
	commonList := []string{"limit", "cursor"}
	return []routeContract{
		contract("GET", "/admin/catalog/provider-readiness", "catalog.view", []string{"plan_display_id", "source_display_id", "product_type", "status", "limit", "cursor"},
			[]string{"does not expose", "provider credentials", "raw provider payloads", "capability JSON"},
			source("internal/modules/catalog/http_handler.go", `"/admin/catalog/provider-readiness"`),
			source("cmd/api/main.go", "PermissionCatalogView")),
		contract("GET", "/admin/catalog/provider-sources", "catalog.view", []string{"display_id", "source_type", "status", "limit", "cursor"},
			[]string{"without credentials", "raw provider payloads"},
			source("internal/modules/catalog/http_handler.go", `"/admin/catalog/provider-sources"`),
			source("cmd/api/main.go", "PermissionCatalogView")),

		contract("GET", "/client/orders", "order.create", []string{"display_id", "status", "billing_status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/order/http_handler.go", `"/client/orders"`),
			source("cmd/api/main.go", "PermissionOrderCreate")),
		contract("POST", "/client/orders", "order.create", nil, []string{"Idempotency-Key", "`tenant_plan_id`"},
			source("internal/modules/order/http_handler.go", `"/client/orders"`),
			source("cmd/api/main.go", "PermissionOrderCreate")),
		contract("GET", "/admin/orders", "order.view", []string{"buyer_user_id", "buyer_display_id", "display_id", "status", "billing_status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/order/http_handler.go", `"/admin/orders"`),
			source("cmd/api/main.go", "PermissionOrderView")),
		contract("PATCH", "/admin/orders/{order_id}/status", "order.manage", nil, []string{"`from_status`", "`to_status`", "`billing_status`"},
			source("internal/modules/order/http_handler.go", `"/admin/orders/"`),
			source("cmd/api/main.go", "PermissionOrderManage")),

		contract("GET", "/client/services", "service.view", []string{"display_id", "order_id", "order_display_id", "status", "limit", "cursor"}, nil,
			source("internal/modules/order/http_handler.go", `"/client/services"`),
			source("cmd/api/main.go", "PermissionServiceView")),
		contract("GET", "/admin/services", "service.view", []string{"buyer_user_id", "buyer_display_id", "display_id", "order_id", "order_display_id", "provider_source_display_id", "status", "limit", "cursor"}, nil,
			source("internal/modules/order/http_handler.go", `"/admin/services"`),
			source("cmd/api/main.go", "PermissionServiceView")),
		contract("POST", "/admin/services/{service_id}/suspend", "service.suspend", nil, []string{"`from_status`", "`reason`", "service.suspended"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceSuspendSuffix"),
			source("cmd/api/main.go", "PermissionServiceSuspend")),
		contract("POST", "/admin/services/{service_id}/unsuspend", "service.unsuspend", nil, []string{"`from_status`", "`reason`", "service.unsuspended"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceUnsuspendSuffix"),
			source("cmd/api/main.go", "PermissionServiceUnsuspend")),
		contract("POST", "/admin/services/{service_id}/terminate", "service.terminate", nil, []string{"`from_status`", "`reason`", "service.terminated"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceTerminateSuffix"),
			source("cmd/api/main.go", "PermissionServiceTerminate")),
		contract("POST", "/reseller/services/{service_id}/suspend", "service.suspend", nil, []string{"`from_status`", "`reason`", "service.suspended"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceSuspendSuffix"),
			source("cmd/api/main.go", "PermissionServiceSuspend")),
		contract("POST", "/reseller/services/{service_id}/unsuspend", "service.unsuspend", nil, []string{"`from_status`", "`reason`", "service.unsuspended"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceUnsuspendSuffix"),
			source("cmd/api/main.go", "PermissionServiceUnsuspend")),
		contract("POST", "/reseller/services/{service_id}/terminate", "service.terminate", nil, []string{"`from_status`", "`reason`", "service.terminated"},
			source("internal/modules/order/http_service_lifecycle.go", "serviceTerminateSuffix"),
			source("cmd/api/main.go", "PermissionServiceTerminate")),
		contract("POST", "/client/services/{service_id}/credentials/{credential_id}/reveal", "service.view", nil, []string{"`reason`", "no-store", "credential.revealed"},
			source("internal/modules/order/http_service_handler.go", "clientServicePrefix", "credentials", "reveal"),
			source("cmd/api/main.go", "PermissionServiceView")),
		contract("POST", "/client/services/{service_id}/renew", "service.renew", nil, []string{"`wallet_id`", "`from_status`", "`Idempotency-Key`", "standalone renewal invoice", "service.renewed"},
			source("internal/modules/order/http_service_renewal.go", "serviceRenewSuffix"),
			source("cmd/api/main.go", "PermissionServiceRenew")),
		contract("POST", "/admin/services/{service_id}/credentials/{credential_id}/reveal", "service.credential.reveal", nil, []string{"`reason`", "rate-limited", "audited without plaintext"},
			source("internal/modules/order/http_service_handler.go", "adminServicePrefix", "credentials", "reveal"),
			source("cmd/api/main.go", "PermissionServiceReveal")),
		contract("POST", "/reseller/services/{service_id}/credentials/{credential_id}/reveal", "service.credential.reveal", nil, []string{"`reason`", "rate-limited", "audited without plaintext"},
			source("internal/modules/order/http_service_handler.go", "resellerServicePrefix", "credentials", "reveal"),
			source("cmd/api/main.go", "PermissionServiceReveal")),

		contract("GET", "/client/invoices", "wallet.view", []string{"display_id", "order_id", "order_display_id", "status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/invoice/http_handler.go", `"/client/invoices"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("GET", "/admin/invoices", "wallet.view", []string{"buyer_user_id", "buyer_display_id", "display_id", "order_id", "order_display_id", "status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/invoice/http_handler.go", `"/admin/invoices"`),
			source("cmd/api/main.go", "PermissionWalletView")),

		contract("GET", "/client/wallets", "wallet.view", append([]string{"display_id", "status"}, commonList...), nil,
			source("internal/modules/wallet/http_handler.go", `"/client/wallets"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("GET", "/admin/wallets", "wallet.view", []string{"display_id", "owner_type", "owner_id", "status", "limit", "cursor"}, nil,
			source("internal/modules/wallet/http_handler.go", `"/admin/wallets"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("POST", "/admin/wallet-refunds", "wallet.adjustment.create", nil, []string{"`Idempotency-Key`", "`X-Access-Reason`", "`reason`", "`wallet.refund.created`"},
			source("internal/modules/wallet/http_handler.go", "adminWalletRefundsPath"),
			source("internal/modules/wallet/http_manual_ledger_handler.go", "CreateWalletRefund", "walletIdempotencyKeyHeader"),
			source("cmd/api/main.go", "PermissionWalletAdjustment")),
		contract("POST", "/admin/wallet-adjustments", "wallet.adjustment.create", nil, []string{"`Idempotency-Key`", "`X-Access-Reason`", "`reason`", "`direction`", "`wallet.adjustment.created`"},
			source("internal/modules/wallet/http_handler.go", "adminWalletAdjustmentsPath"),
			source("internal/modules/wallet/http_manual_ledger_handler.go", "CreateWalletAdjustment", "walletIdempotencyKeyHeader"),
			source("cmd/api/main.go", "PermissionWalletAdjustment")),
		contract("GET", "/admin/topup-requests", "wallet.view", []string{"requested_by", "requested_by_display_id", "display_id", "wallet_id", "wallet_display_id", "payment_method", "status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/wallet/http_handler.go", `"/admin/topup-requests"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("POST", "/admin/topup-requests/{topup_request_id}/approve", "wallet.topup.approve", nil, []string{"`review_note`", "`ledger_entry_id`"},
			source("internal/modules/wallet/http_handler.go", `"/admin/topup-requests/"`),
			source("cmd/api/main.go", "PermissionWalletTopupApprove")),

		contract("GET", "/client/transactions", "wallet.view", []string{"display_id", "order_id", "order_display_id", "invoice_id", "invoice_display_id", "type", "status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/payment/http_handler.go", `"/client/transactions"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("GET", "/admin/transactions", "wallet.view", []string{"account_user_id", "account_display_id", "display_id", "order_id", "order_display_id", "invoice_id", "invoice_display_id", "type", "status", "amount_min", "amount_max", "limit", "cursor"}, nil,
			source("internal/modules/payment/http_handler.go", `"/admin/transactions"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("GET", "/admin/payment-reconciliation", "wallet.view", []string{"account_user_id", "display_id", "status", "provider", "invoice_id", "invoice_display_id", "wallet_id", "wallet_display_id", "amount_min", "amount_max", "created_from", "created_to", "limit", "cursor"}, nil,
			source("internal/modules/payment/http_handler.go", `"/admin/payment-reconciliation"`),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("GET", "/admin/daily-reconciliation", "wallet.view", []string{"date"}, []string{"wallet", "invoice", "duplicate", "UTC", "display IDs"},
			source("internal/modules/payment/http_handler.go", "adminDailyReconciliationPath"),
			source("internal/modules/payment/http_daily_reconciliation_handler.go", "BuildDailyReconciliationReport"),
			source("cmd/api/main.go", "PermissionWalletView")),
		contract("POST", "/client/invoice-wallet-payments", "wallet.view", nil, []string{"`invoice_id`", "`wallet_id`", "`Idempotency-Key`"},
			source("internal/modules/payment/http_handler.go", "clientInvoiceWalletPaymentsPath"),
			source("cmd/api/main.go", "PermissionWalletView")),

		contract("GET", "/admin/audit-logs", "audit.view", []string{"actor_id", "actor_type", "display_id", "action", "target_type", "target_id", "created_from", "created_to", "limit", "cursor"}, nil,
			source("internal/modules/audit/http_handler.go", `"/admin/audit-logs"`),
			source("cmd/api/main.go", "PermissionAuditView")),

		contract("GET", "/admin/jobs", "order.view", []string{"display_id", "job_type", "status", "reference_type", "reference_id", "source_id", "source_display_id", "limit", "cursor"}, []string{"`payload_json`", "`idempotency_key`"},
			source("internal/modules/jobs/http_handler.go", `"/admin/jobs"`),
			source("cmd/api/main.go", "PermissionOrderView")),
		contract("GET", "/admin/jobs/summary", "provisioning.job.view", []string{"job_type"}, []string{"redacted error fields"},
			source("internal/modules/jobs/http_handler.go", `"/admin/jobs/summary"`),
			source("cmd/api/main.go", "PermissionProvisioningJobView")),
		contract("POST", "/admin/jobs/{job_id}/retry", "provisioning.job.retry", nil, []string{"`next_attempt_at`", "`failed_retryable`", "`manual_review`"},
			source("internal/modules/jobs/http_handler.go", "adminJobPrefix"),
			source("internal/modules/jobs/http_recovery.go", `jobActionRetry        = "retry"`),
			source("cmd/api/main.go", "PermissionProvisioningJobRetry")),
		contract("POST", "/admin/jobs/{job_id}/manual-review", "provisioning.manual_review.resolve", nil, []string{"`reason`", "`job.status_conflict`"},
			source("internal/modules/jobs/http_handler.go", "adminJobPrefix"),
			source("internal/modules/jobs/http_recovery.go", `jobActionManualReview = "manual-review"`),
			source("cmd/api/main.go", "PermissionManualReviewResolve")),
		contract("POST", "/admin/jobs/{job_id}/cancel", "provisioning.manual_review.resolve", nil, []string{"`reason`", "`job.status_conflict`"},
			source("internal/modules/jobs/http_handler.go", "adminJobPrefix"),
			source("internal/modules/jobs/http_recovery.go", `jobActionCancel       = "cancel"`),
			source("cmd/api/main.go", "PermissionManualReviewResolve")),
		contract("GET", "/reseller/jobs", "order.view", []string{"display_id", "job_type", "status", "reference_type", "reference_id", "source_id", "source_display_id", "limit", "cursor"}, []string{"tenant scope"},
			source("internal/modules/jobs/http_handler.go", `"/reseller/jobs"`),
			source("cmd/api/main.go", "PermissionOrderView")),
	}
}

func contract(method string, path string, perm string, queries []string, docTokens []string, sources ...sourceExpectation) routeContract {
	return routeContract{Method: method, Path: path, Perm: perm, Queries: queries, DocTokens: docTokens, Sources: sources}
}

func source(path string, tokens ...string) sourceExpectation {
	return sourceExpectation{Path: path, Tokens: tokens}
}
