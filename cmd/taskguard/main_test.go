package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckTaskBoardPassesCompleteBoard(t *testing.T) {
	root := t.TempDir()
	writeTaskBoard(t, root, boardMarkdown("1", "1", "0", "0", "1", "1", claimableRows("T113"), inFlightRows("T114", "IN_PROGRESS")))
	writeTask(t, root, activeTaskDir, "T113_first.md", "T113", "TODO")
	writeTask(t, root, activeTaskDir, "T114_first.md", "T114", "IN_PROGRESS")
	writeTask(t, root, activeTaskDir, "T115_done.md", "T115", "DONE")
	writeTask(t, root, removedTaskDir, "T010_removed.md", "T010", "REMOVED")

	failures, err := checkTaskBoard(root, defaultIndexPath)
	if err != nil {
		t.Fatalf("checkTaskBoard returned error: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got %#v", failures)
	}
}

func TestCheckTaskBoardReportsSnapshotMismatch(t *testing.T) {
	root := t.TempDir()
	writeTaskBoard(t, root, boardMarkdown("2", "0", "0", "0", "0", "0", claimableRows("T113"), "No task files are currently `IN_PROGRESS` or `REVIEW`."))
	writeTask(t, root, activeTaskDir, "T113_first.md", "T113", "TODO")

	failures, err := checkTaskBoard(root, defaultIndexPath)
	if err != nil {
		t.Fatalf("checkTaskBoard returned error: %v", err)
	}
	assertFailureContains(t, failures, "snapshot `TODO` is 2, expected 1")
}

func TestCheckTaskBoardReportsClaimableNonTodo(t *testing.T) {
	root := t.TempDir()
	writeTaskBoard(t, root, boardMarkdown("0", "1", "0", "0", "0", "0", claimableRows("T113"), inFlightRows("T113", "IN_PROGRESS")))
	writeTask(t, root, activeTaskDir, "T113_first.md", "T113", "IN_PROGRESS")

	failures, err := checkTaskBoard(root, defaultIndexPath)
	if err != nil {
		t.Fatalf("checkTaskBoard returned error: %v", err)
	}
	assertFailureContains(t, failures, "claimable row `T113` points to status `IN_PROGRESS`")
}

func TestCheckTaskBoardReportsMissingRequiredField(t *testing.T) {
	root := t.TempDir()
	writeTaskBoard(t, root, boardMarkdown("1", "0", "0", "0", "0", "0", claimableRows("T113"), "No task files are currently `IN_PROGRESS` or `REVIEW`."))
	writeFile(t, root, filepath.Join(activeTaskDir, "T113_first.md"), strings.ReplaceAll(taskMarkdown("T113", "TODO"), "Owner: -\n", ""))

	failures, err := checkTaskBoard(root, defaultIndexPath)
	if err != nil {
		t.Fatalf("checkTaskBoard returned error: %v", err)
	}
	assertFailureContains(t, failures, "missing required field `Owner`")
}

func TestCheckTaskBoardReportsInFlightMismatch(t *testing.T) {
	root := t.TempDir()
	writeTaskBoard(t, root, boardMarkdown("0", "0", "1", "0", "0", "0", "| - | - | - | - | No TODO task is currently claimable. |", inFlightRows("T113", "IN_PROGRESS")))
	writeTask(t, root, activeTaskDir, "T113_first.md", "T113", "REVIEW")

	failures, err := checkTaskBoard(root, defaultIndexPath)
	if err != nil {
		t.Fatalf("checkTaskBoard returned error: %v", err)
	}
	assertFailureContains(t, failures, "in-flight row `T113` status is `IN_PROGRESS`, expected `REVIEW`")
}

func assertFailureContains(t *testing.T, failures []string, want string) {
	t.Helper()
	for _, failure := range failures {
		if strings.Contains(failure, want) {
			return
		}
	}
	t.Fatalf("expected failure containing %q, got %#v", want, failures)
}

func boardMarkdown(todo, inProgress, review, blocked, done, removed, claimable, inFlight string) string {
	return strings.Join([]string{
		"# Shared Task Index",
		"",
		"## Board Snapshot",
		"",
		"- `TODO`: " + todo,
		"- `IN_PROGRESS`: " + inProgress,
		"- `REVIEW`: " + review,
		"- `BLOCKED`: " + blocked,
		"- `DONE` task files in `tasks/active/`: " + done,
		"- `REMOVED` task files in `tasks/removed/`: " + removed,
		"",
		"## Claimable Tasks",
		"",
		"| ID | Task File | Suggested Branch | Area | Summary |",
		"| --- | --- | --- | --- | --- |",
		claimable,
		"",
		"## In-Flight Task Files",
		"",
		inFlight,
		"",
		"## Done Task Files",
	}, "\n")
}

func claimableRows(ids ...string) string {
	rows := make([]string, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, "| "+id+" | [tasks/active/"+id+"_first.md](tasks/active/"+id+"_first.md) | codex/"+strings.ToLower(id)+" | workflow | Example task. |")
	}
	return strings.Join(rows, "\n")
}

func inFlightRows(id string, status string) string {
	return strings.Join([]string{
		"| ID | Task File | Status | Owner | Branch | Summary |",
		"| --- | --- | --- | --- | --- | --- |",
		"| " + id + " | [tasks/active/" + id + "_first.md](tasks/active/" + id + "_first.md) | " + status + " | Codex | codex/" + strings.ToLower(id) + " | Example task. |",
	}, "\n")
}

func writeTask(t *testing.T, root string, dir string, name string, id string, status string) {
	t.Helper()
	writeFile(t, root, filepath.Join(dir, name), taskMarkdown(id, status))
}

func writeTaskBoard(t *testing.T, root string, body string) {
	t.Helper()
	writeFile(t, root, defaultIndexPath, body)
}

func writeFile(t *testing.T, root string, path string, body string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("create test dir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(body), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}

func taskMarkdown(id string, status string) string {
	return strings.Join([]string{
		"# " + id + " - Example task",
		"",
		"Status: " + status,
		"Owner: -",
		"Branch: codex/" + strings.ToLower(id),
		"PR: -",
		"Risk: workflow",
		"Created: 2026-04-25",
		"Updated: 2026-04-25",
		"",
		"## Summary",
		"",
		"Example summary.",
		"",
		"## Scope",
		"",
		"- Example scope.",
		"",
		"## Acceptance Criteria",
		"",
		"- Example acceptance.",
		"",
		"## Notes",
		"",
		"- Example notes.",
		"",
		"## Agent Log",
		"",
		"- 2026-04-25: Task created.",
	}, "\n")
}
