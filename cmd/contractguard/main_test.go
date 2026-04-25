package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouteSectionStopsAtNextRoute(t *testing.T) {
	doc := strings.Join([]string{
		"- `GET /admin/jobs`",
		"  - auth: admin actor, `order.view`",
		"  - query: `display_id`, `limit`",
		"",
		"- `GET /admin/jobs/summary`",
		"  - auth: admin actor, `provisioning.job.view`",
	}, "\n")

	section := routeSection(doc, "GET /admin/jobs")
	if !strings.Contains(section, "`order.view`") {
		t.Fatalf("expected section to include first route: %q", section)
	}
	if strings.Contains(section, "provisioning.job.view") {
		t.Fatalf("expected section to stop before next route: %q", section)
	}
}

func TestCheckContractsReportsMissingRoute(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, defaultReferencePath, "- `GET /admin/jobs`\n  - auth: admin actor, `order.view`\n")
	writeTestFile(t, root, "internal/modules/jobs/http_handler.go", `mux.HandleFunc("/admin/jobs", handler.adminJobsRoute)`)

	failures, err := checkContracts(root, defaultReferencePath, []routeContract{
		contract("GET", "/admin/missing", "order.view", nil, nil,
			source("internal/modules/jobs/http_handler.go", `"/admin/jobs"`)),
	})
	if err != nil {
		t.Fatalf("checkContracts returned error: %v", err)
	}
	if len(failures) != 1 || !strings.Contains(failures[0], "missing docs route") {
		t.Fatalf("expected missing route failure, got %#v", failures)
	}
}

func TestCheckContractsReportsMissingQuery(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, defaultReferencePath, "- `GET /admin/jobs`\n  - auth: admin actor, `order.view`\n  - query: `display_id`\n")
	writeTestFile(t, root, "internal/modules/jobs/http_handler.go", `mux.HandleFunc("/admin/jobs", handler.adminJobsRoute)`)

	failures, err := checkContracts(root, defaultReferencePath, []routeContract{
		contract("GET", "/admin/jobs", "order.view", []string{"display_id", "source_id"}, nil,
			source("internal/modules/jobs/http_handler.go", `"/admin/jobs"`)),
	})
	if err != nil {
		t.Fatalf("checkContracts returned error: %v", err)
	}
	if len(failures) != 1 || !strings.Contains(failures[0], "source_id") {
		t.Fatalf("expected missing query failure, got %#v", failures)
	}
}

func TestCheckContractsPassesCompleteContract(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, defaultReferencePath, "- `GET /admin/jobs`\n  - auth: admin actor, `order.view`\n  - query: `display_id`, `source_id`\n  - note: does not expose `payload_json`\n")
	writeTestFile(t, root, "internal/modules/jobs/http_handler.go", `mux.HandleFunc("/admin/jobs", handler.adminJobsRoute)`)

	failures, err := checkContracts(root, defaultReferencePath, []routeContract{
		contract("GET", "/admin/jobs", "order.view", []string{"display_id", "source_id"}, []string{"`payload_json`"},
			source("internal/modules/jobs/http_handler.go", `"/admin/jobs"`)),
	})
	if err != nil {
		t.Fatalf("checkContracts returned error: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got %#v", failures)
	}
}

func writeTestFile(t *testing.T, root string, path string, body string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("create test dir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
