package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckErrorContractsPassesCompleteContract(t *testing.T) {
	root := t.TempDir()
	writeGuardFile(t, root, defaultResponseStandardPath, responseDoc("validation.failed", "wallet.insufficient_balance"))
	writeGuardFile(t, root, defaultAPIReferencePath, apiDoc("validation.failed", "wallet.insufficient_balance"))
	writeGuardFile(t, root, responseSourcePath, responseSource())
	writeGuardFile(t, root, "internal/modules/payment/http_handler.go", `httpserver.WriteError(w, r, 409, "wallet.insufficient_balance", "Wallet balance is insufficient.")`)

	failures, err := checkErrorContracts(root, defaultResponseStandardPath, defaultAPIReferencePath, []errorCodeContract{
		code("wallet.insufficient_balance", source("internal/modules/payment/http_handler.go")),
	})
	if err != nil {
		t.Fatalf("checkErrorContracts returned error: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got %#v", failures)
	}
}

func TestCheckErrorContractsReportsMissingAPIDocCode(t *testing.T) {
	root := t.TempDir()
	writeGuardFile(t, root, defaultResponseStandardPath, responseDoc("validation.failed"))
	writeGuardFile(t, root, defaultAPIReferencePath, apiDoc("validation.failed"))
	writeGuardFile(t, root, responseSourcePath, responseSource())
	writeGuardFile(t, root, "internal/modules/payment/http_handler.go", `httpserver.WriteError(w, r, 409, "wallet.insufficient_balance", "Wallet balance is insufficient.")`)

	failures, err := checkErrorContracts(root, defaultResponseStandardPath, defaultAPIReferencePath, []errorCodeContract{
		code("wallet.insufficient_balance", source("internal/modules/payment/http_handler.go")),
	})
	if err != nil {
		t.Fatalf("checkErrorContracts returned error: %v", err)
	}
	assertGuardFailure(t, failures, "API reference missing stable error code `wallet.insufficient_balance`")
}

func TestCheckErrorContractsReportsMissingSourceCode(t *testing.T) {
	root := t.TempDir()
	writeGuardFile(t, root, defaultResponseStandardPath, responseDoc("validation.failed", "wallet.insufficient_balance"))
	writeGuardFile(t, root, defaultAPIReferencePath, apiDoc("validation.failed", "wallet.insufficient_balance"))
	writeGuardFile(t, root, responseSourcePath, responseSource())
	writeGuardFile(t, root, "internal/modules/payment/http_handler.go", `httpserver.WriteError(w, r, 409, "wallet.other", "Other.")`)

	failures, err := checkErrorContracts(root, defaultResponseStandardPath, defaultAPIReferencePath, []errorCodeContract{
		code("wallet.insufficient_balance", source("internal/modules/payment/http_handler.go")),
	})
	if err != nil {
		t.Fatalf("checkErrorContracts returned error: %v", err)
	}
	assertGuardFailure(t, failures, "source internal/modules/payment/http_handler.go missing")
}

func TestCheckErrorContractsReportsMissingEnvelopeToken(t *testing.T) {
	root := t.TempDir()
	writeGuardFile(t, root, defaultResponseStandardPath, strings.ReplaceAll(responseDoc("validation.failed"), `"fields"`, ""))
	writeGuardFile(t, root, defaultAPIReferencePath, apiDoc("validation.failed"))
	writeGuardFile(t, root, responseSourcePath, responseSource())

	failures, err := checkErrorContracts(root, defaultResponseStandardPath, defaultAPIReferencePath, nil)
	if err != nil {
		t.Fatalf("checkErrorContracts returned error: %v", err)
	}
	assertGuardFailure(t, failures, `response standard missing envelope token "\"fields\""`)
}

func assertGuardFailure(t *testing.T, failures []string, want string) {
	t.Helper()
	for _, failure := range failures {
		if strings.Contains(failure, want) {
			return
		}
	}
	t.Fatalf("expected failure containing %q, got %#v", want, failures)
}

func responseDoc(codes ...string) string {
	return strings.Join([]string{
		`{ "error": { "code": "validation.failed", "message": "Request validation failed.", "details": {}, "fields": [] }, "request_id": "req_1" }`,
		strings.Join(codes, "\n"),
	}, "\n")
}

func apiDoc(codes ...string) string {
	lines := []string{
		"## Tracked Stable API Error Codes",
		"### Shared codes",
		"### Route-specific codes",
	}
	for _, code := range codes {
		lines = append(lines, "- `"+code+"`")
	}
	return strings.Join(lines, "\n")
}

func responseSource() string {
	return strings.Join([]string{
		"`json:\"error\"`",
		"`json:\"code\"`",
		"`json:\"message\"`",
		"`json:\"details,omitempty\"`",
		"`json:\"fields,omitempty\"`",
		"`json:\"request_id\"`",
		"WriteValidationError",
	}, "\n")
}

func writeGuardFile(t *testing.T, root string, path string, body string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("create test dir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
