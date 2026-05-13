package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultResponseStandardPath = "docs/05_development_standards/50_API_Response_Error_Logging_Standard.md"
	defaultAPIReferencePath     = "docs/05_development_standards/62_API_Error_Code_Drift_Guard.md"
	responseSourcePath          = "internal/platform/httpserver/response.go"
)

type errorCodeContract struct {
	Code    string
	Sources []sourceExpectation
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
	flags := flag.NewFlagSet("errorcodeguard", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	responseStandard := flags.String("response-standard", defaultResponseStandardPath, "API response standard markdown")
	apiReference := flags.String("api-reference", defaultAPIReferencePath, "tracked API error code reference markdown")
	if err := flags.Parse(args); err != nil {
		return err
	}
	failures, err := checkErrorContracts(*root, *responseStandard, *apiReference, trackedErrorCodeContracts())
	if err != nil {
		return err
	}
	if len(failures) > 0 {
		return fmt.Errorf("API error code drift guard failed:\n%s", strings.Join(failures, "\n"))
	}
	fmt.Printf("API error code drift guard passed: %d error code contract(s)\n", len(trackedErrorCodeContracts()))
	return nil
}

func checkErrorContracts(root string, responseStandardPath string, apiReferencePath string, contracts []errorCodeContract) ([]string, error) {
	responseDoc, err := readText(filepath.Join(root, responseStandardPath))
	if err != nil {
		return nil, err
	}
	apiDoc, err := readText(filepath.Join(root, apiReferencePath))
	if err != nil {
		return nil, err
	}
	responseSource, err := readText(filepath.Join(root, responseSourcePath))
	if err != nil {
		return nil, err
	}
	failures := checkEnvelope(responseDoc, apiDoc, responseSource)
	for _, contract := range contracts {
		if !containsCode(apiDoc, contract.Code) {
			failures = append(failures, fmt.Sprintf("- API reference missing stable error code `%s`", contract.Code))
		}
		sourceFailures, err := checkSources(root, contract)
		if err != nil {
			return nil, err
		}
		failures = append(failures, sourceFailures...)
	}
	return failures, nil
}

func checkEnvelope(responseDoc string, apiDoc string, responseSource string) []string {
	failures := make([]string, 0)
	sourceTokens := []string{
		`json:"error"`,
		`json:"code"`,
		`json:"message"`,
		`json:"details,omitempty"`,
		`json:"fields,omitempty"`,
		`json:"request_id"`,
		"WriteValidationError",
	}
	for _, token := range sourceTokens {
		if !strings.Contains(responseSource, token) {
			failures = append(failures, fmt.Sprintf("- response source missing envelope token %q", token))
		}
	}
	docTokens := []string{`"error"`, `"code"`, `"message"`, `"details"`, `"fields"`, `"request_id"`, "validation.failed"}
	for _, token := range docTokens {
		if !strings.Contains(responseDoc, token) {
			failures = append(failures, fmt.Sprintf("- response standard missing envelope token %q", token))
		}
	}
	for _, token := range []string{"Tracked Stable API Error Codes", "Shared codes", "Route-specific codes"} {
		if !strings.Contains(apiDoc, token) {
			failures = append(failures, fmt.Sprintf("- API reference missing error documentation token %q", token))
		}
	}
	return failures
}

func checkSources(root string, contract errorCodeContract) ([]string, error) {
	failures := make([]string, 0)
	for _, source := range contract.Sources {
		body, err := readText(filepath.Join(root, source.Path))
		if err != nil {
			return nil, err
		}
		for _, token := range append([]string{contract.Code}, source.Tokens...) {
			if !strings.Contains(body, token) {
				failures = append(failures, fmt.Sprintf("- error code `%s` source %s missing %q", contract.Code, source.Path, token))
			}
		}
	}
	return failures, nil
}

func containsCode(doc string, code string) bool {
	return strings.Contains(doc, "`"+code+"`") || strings.Contains(doc, `"`+code+`"`)
}

func trackedErrorCodeContracts() []errorCodeContract {
	return []errorCodeContract{
		code("validation.failed", source(responseSourcePath, "WriteValidationError")),
		code("request.invalid_json", source("internal/modules/order/http_shared.go")),
		code("request.method_not_allowed", source("internal/platform/middleware/method.go")),
		code("request.limit_invalid", source("internal/modules/jobs/http_handler.go")),
		code("request.limit_too_large", source("internal/modules/jobs/http_handler.go")),
		code("request.display_id_invalid", source("internal/modules/jobs/http_filter_query.go")),
		code("request.amount_invalid", source("internal/modules/invoice/http_filter_query.go")),
		code("request.amount_range_invalid", source("internal/modules/invoice/http_filter_query.go")),

		code("tenant.context_missing", source("internal/modules/rbac/http_middleware.go")),
		code("tenant.context_invalid", source("internal/modules/identity/admin_read_http_handler.go")),
		code("tenant.context_mismatch", source("internal/modules/rbac/http_middleware.go")),
		code("auth.actor_required", source("internal/modules/rbac/http_middleware.go")),
		code("auth.permission_denied", source("internal/modules/rbac/http_middleware.go")),
		code("auth.reason_required", source("internal/modules/rbac/http_middleware.go")),

		code("catalog.not_found", source("internal/modules/catalog/http_handler.go")),
		code("order.not_found", source("internal/modules/order/http_handler.go")),
		code("order.status_conflict", source("internal/modules/order/http_handler.go")),
		code("order.status_transition_invalid", source("internal/modules/order/http_handler.go")),
		code("order.provisioning_source_not_found", source("internal/modules/payment/http_handler.go")),
		code("service.not_found", source("internal/modules/order/http_handler.go")),
		code("service.status_invalid", source("internal/modules/order/http_handler.go")),
		code("credential.not_found", source("internal/modules/order/http_handler.go")),
		code("credential.reveal_rate_limited", source("internal/modules/order/http_handler.go")),
		code("credential.reveal_denied", source("internal/modules/order/http_handler.go")),
		code("invoice.not_found", source("internal/modules/invoice/http_handler.go")),
		code("invoice.status_conflict", source("internal/modules/payment/http_handler.go")),

		code("wallet.not_found", source("internal/modules/wallet/http_handler.go")),
		code("wallet.ledger_not_found", source("internal/modules/wallet/http_handler.go")),
		code("wallet.topup_not_found", source("internal/modules/wallet/http_handler.go")),
		code("wallet.topup_status_conflict", source("internal/modules/wallet/http_handler.go")),
		code("wallet.payment_method_invalid", source("internal/modules/wallet/http_handler.go")),
		code("wallet.status_conflict", source("internal/modules/wallet/http_handler.go")),
		code("wallet.currency_mismatch", source("internal/modules/wallet/http_handler.go")),
		code("wallet.idempotency_conflict", source("internal/modules/wallet/http_handler.go")),
		code("wallet.insufficient_balance", source("internal/modules/payment/http_handler.go"), source("internal/modules/wallet/http_handler.go")),

		code("checkout.order_not_checkoutable", source("internal/modules/checkout/http_handler.go")),
		code("payment.transaction_not_found", source("internal/modules/payment/http_handler.go")),
		code("payment.invoice_not_payable", source("internal/modules/payment/http_handler.go")),
		code("payment.idempotency_conflict", source("internal/modules/payment/http_handler.go")),
		code("payment.wallet_currency_mismatch", source("internal/modules/payment/http_handler.go")),

		code("job.not_found", source("internal/modules/jobs/http_handler.go")),
		code("job.status_invalid", source("internal/modules/jobs/http_handler.go")),
		code("job.status_conflict", source("internal/modules/jobs/http_handler.go")),
		code("job.manual_review_reason_missing", source("internal/modules/jobs/http_handler.go")),
		code("audit.created_time_invalid", source("internal/modules/audit/http_handler.go")),
		code("audit.log_not_found", source("internal/modules/audit/http_handler.go")),
	}
}

func code(value string, sources ...sourceExpectation) errorCodeContract {
	return errorCodeContract{Code: value, Sources: sources}
}

func source(path string, tokens ...string) sourceExpectation {
	return sourceExpectation{Path: path, Tokens: tokens}
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
