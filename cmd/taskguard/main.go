package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultIndexPath = "TASKS.md"
	activeTaskDir    = "tasks/active"
	removedTaskDir   = "tasks/removed"
)

var (
	taskFileNamePattern = regexp.MustCompile(`^(T[0-9]{3})_.+\.md$`)
	snapshotLinePattern = regexp.MustCompile("^- `([^`]+)`[^:]*:\\s*([0-9]+)\\s*$")
	linkPathPattern     = regexp.MustCompile(`\(([^)]+)\)`)
)

type taskRecord struct {
	ID     string
	Path   string
	Status string
	Fields map[string]string
}

type indexRow struct {
	ID     string
	Path   string
	Status string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("taskguard", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	index := flags.String("index", defaultIndexPath, "task index markdown")
	if err := flags.Parse(args); err != nil {
		return err
	}
	failures, err := checkTaskBoard(*root, *index)
	if err != nil {
		return err
	}
	if len(failures) > 0 {
		return fmt.Errorf("task board consistency guard failed:\n%s", strings.Join(failures, "\n"))
	}
	fmt.Println("Task board consistency guard passed")
	return nil
}

func checkTaskBoard(root string, indexPath string) ([]string, error) {
	indexBody, err := readText(filepath.Join(root, indexPath))
	if err != nil {
		return nil, err
	}
	active, removed, failures, err := loadTaskRecords(root)
	if err != nil {
		return nil, err
	}
	failures = append(failures, checkSnapshot(indexBody, active, removed)...)
	failures = append(failures, checkClaimableRows(indexBody, active)...)
	failures = append(failures, checkInFlightRows(indexBody, active)...)
	sort.Strings(failures)
	return failures, nil
}

func loadTaskRecords(root string) ([]taskRecord, []taskRecord, []string, error) {
	active, activeFailures, err := readTaskDir(root, activeTaskDir)
	if err != nil {
		return nil, nil, nil, err
	}
	removed, removedFailures, err := readTaskDir(root, removedTaskDir)
	if err != nil {
		return nil, nil, nil, err
	}
	failures := append(activeFailures, removedFailures...)
	failures = append(failures, checkTaskDirStatus(active, "active")...)
	failures = append(failures, checkTaskDirStatus(removed, "removed")...)
	failures = append(failures, checkDuplicateIDs(active, activeTaskDir)...)
	failures = append(failures, checkDuplicateIDs(removed, removedTaskDir)...)
	return active, removed, failures, nil
}

func readTaskDir(root string, relDir string) ([]taskRecord, []string, error) {
	entries, err := os.ReadDir(filepath.Join(root, relDir))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, []string{fmt.Sprintf("- missing task directory `%s`", relDir)}, nil
		}
		return nil, nil, fmt.Errorf("read %s: %w", relDir, err)
	}
	records := make([]taskRecord, 0)
	failures := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		relPath := filepath.ToSlash(filepath.Join(relDir, entry.Name()))
		body, err := readText(filepath.Join(root, relPath))
		if err != nil {
			return nil, nil, err
		}
		record, fileFailures := parseTaskFile(relPath, body)
		records = append(records, record)
		failures = append(failures, fileFailures...)
	}
	return records, failures, nil
}

func parseTaskFile(relPath string, body string) (taskRecord, []string) {
	base := filepath.Base(relPath)
	id := ""
	if matches := taskFileNamePattern.FindStringSubmatch(base); len(matches) == 2 {
		id = matches[1]
	}
	record := taskRecord{ID: id, Path: relPath, Fields: parseFields(body)}
	failures := make([]string, 0)
	if id == "" {
		failures = append(failures, fmt.Sprintf("- `%s` file name must start with `Txxx_` and end with `.md`", relPath))
	}
	if id != "" && !strings.HasPrefix(firstLine(body), "# "+id+" - ") {
		failures = append(failures, fmt.Sprintf("- `%s` heading must start with `# %s - `", relPath, id))
	}
	for _, field := range requiredFields() {
		value, ok := record.Fields[field]
		if !ok || strings.TrimSpace(value) == "" {
			failures = append(failures, fmt.Sprintf("- `%s` missing required field `%s`", relPath, field))
		}
	}
	record.Status = record.Fields["Status"]
	if record.Status != "" && !validStatuses()[record.Status] {
		failures = append(failures, fmt.Sprintf("- `%s` has invalid status `%s`", relPath, record.Status))
	}
	for _, section := range requiredSections() {
		if !hasSection(body, section) {
			failures = append(failures, fmt.Sprintf("- `%s` missing section `## %s`", relPath, section))
		}
	}
	return record, failures
}

func parseFields(body string) map[string]string {
	fields := map[string]string{}
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 2 {
			fields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return fields
}

func checkTaskDirStatus(records []taskRecord, dirType string) []string {
	failures := make([]string, 0)
	for _, record := range records {
		if dirType == "active" && record.Status == "REMOVED" {
			failures = append(failures, fmt.Sprintf("- `%s` is under active tasks but has status REMOVED", record.Path))
		}
		if dirType == "removed" && record.Status != "REMOVED" {
			failures = append(failures, fmt.Sprintf("- `%s` is under removed tasks but status is `%s`", record.Path, record.Status))
		}
	}
	return failures
}

func checkDuplicateIDs(records []taskRecord, relDir string) []string {
	seen := map[string]string{}
	failures := make([]string, 0)
	for _, record := range records {
		if record.ID == "" {
			continue
		}
		if previous, ok := seen[record.ID]; ok {
			failures = append(failures, fmt.Sprintf("- duplicate task ID `%s` in `%s`: `%s` and `%s`", record.ID, relDir, previous, record.Path))
			continue
		}
		seen[record.ID] = record.Path
	}
	return failures
}

func checkSnapshot(indexBody string, active []taskRecord, removed []taskRecord) []string {
	expected := map[string]int{
		"TODO":        countStatus(active, "TODO"),
		"IN_PROGRESS": countStatus(active, "IN_PROGRESS"),
		"REVIEW":      countStatus(active, "REVIEW"),
		"BLOCKED":     countStatus(active, "BLOCKED"),
		"DONE":        countStatus(active, "DONE"),
		"REMOVED":     len(removed),
	}
	actual := parseSnapshot(indexBody)
	failures := make([]string, 0)
	for _, status := range []string{"TODO", "IN_PROGRESS", "REVIEW", "BLOCKED", "DONE", "REMOVED"} {
		value, ok := actual[status]
		if !ok {
			failures = append(failures, fmt.Sprintf("- `TASKS.md` board snapshot missing `%s` count", status))
			continue
		}
		if value != expected[status] {
			failures = append(failures, fmt.Sprintf("- `TASKS.md` board snapshot `%s` is %d, expected %d", status, value, expected[status]))
		}
	}
	return failures
}

func parseSnapshot(indexBody string) map[string]int {
	counts := map[string]int{}
	for _, line := range sectionLines(indexBody, "Board Snapshot") {
		matches := snapshotLinePattern.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) != 3 {
			continue
		}
		value, err := strconv.Atoi(matches[2])
		if err == nil {
			counts[matches[1]] = value
		}
	}
	return counts
}

func checkClaimableRows(indexBody string, active []taskRecord) []string {
	rows := parseTableRows(indexBody, "Claimable Tasks")
	byID := recordsByID(active)
	todoIDs := idsByStatus(active, "TODO")
	failures := make([]string, 0)
	seen := map[string]bool{}
	for _, row := range rows {
		if row.ID == "-" {
			if len(todoIDs) > 0 {
				failures = append(failures, "- `TASKS.md` claimable placeholder is present while TODO tasks exist")
			}
			continue
		}
		record, ok := byID[row.ID]
		if !ok {
			failures = append(failures, fmt.Sprintf("- claimable row `%s` points to a missing active task", row.ID))
			continue
		}
		if normalizePath(row.Path) != record.Path {
			failures = append(failures, fmt.Sprintf("- claimable row `%s` path is `%s`, expected `%s`", row.ID, row.Path, record.Path))
		}
		if record.Status != "TODO" {
			failures = append(failures, fmt.Sprintf("- claimable row `%s` points to status `%s`, expected TODO", row.ID, record.Status))
		}
		seen[row.ID] = true
	}
	for _, id := range todoIDs {
		if !seen[id] {
			failures = append(failures, fmt.Sprintf("- TODO task `%s` is missing from `TASKS.md` claimable rows", id))
		}
	}
	return failures
}

func checkInFlightRows(indexBody string, active []taskRecord) []string {
	rows := parseTableRows(indexBody, "In-Flight Task Files")
	byID := recordsByID(active)
	inFlightIDs := append(idsByStatus(active, "IN_PROGRESS"), idsByStatus(active, "REVIEW")...)
	sort.Strings(inFlightIDs)
	failures := make([]string, 0)
	seen := map[string]bool{}
	for _, row := range rows {
		if row.ID == "-" {
			continue
		}
		record, ok := byID[row.ID]
		if !ok {
			failures = append(failures, fmt.Sprintf("- in-flight row `%s` points to a missing active task", row.ID))
			continue
		}
		if normalizePath(row.Path) != record.Path {
			failures = append(failures, fmt.Sprintf("- in-flight row `%s` path is `%s`, expected `%s`", row.ID, row.Path, record.Path))
		}
		if record.Status != "IN_PROGRESS" && record.Status != "REVIEW" {
			failures = append(failures, fmt.Sprintf("- in-flight row `%s` points to status `%s`", row.ID, record.Status))
		}
		if row.Status != "" && row.Status != record.Status {
			failures = append(failures, fmt.Sprintf("- in-flight row `%s` status is `%s`, expected `%s`", row.ID, row.Status, record.Status))
		}
		seen[row.ID] = true
	}
	for _, id := range inFlightIDs {
		if !seen[id] {
			failures = append(failures, fmt.Sprintf("- in-flight task `%s` is missing from `TASKS.md` in-flight rows", id))
		}
	}
	return failures
}

func parseTableRows(indexBody string, heading string) []indexRow {
	rows := make([]indexRow, 0)
	for _, line := range sectionLines(indexBody, heading) {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") || strings.Contains(trimmed, "---") {
			continue
		}
		cells := splitMarkdownRow(trimmed)
		if len(cells) < 2 || cells[0] == "ID" {
			continue
		}
		row := indexRow{ID: cells[0], Path: extractPath(cells[1])}
		if heading == "In-Flight Task Files" && len(cells) >= 3 {
			row.Status = cells[2]
		}
		rows = append(rows, row)
	}
	return rows
}

func splitMarkdownRow(line string) []string {
	parts := strings.Split(strings.Trim(line, "|"), "|")
	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells
}

func extractPath(cell string) string {
	matches := linkPathPattern.FindStringSubmatch(cell)
	if len(matches) == 2 {
		return normalizePath(matches[1])
	}
	return normalizePath(cell)
}

func sectionLines(body string, heading string) []string {
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## "+heading {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return nil
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			end = i
			break
		}
	}
	return lines[start:end]
}

func recordsByID(records []taskRecord) map[string]taskRecord {
	byID := map[string]taskRecord{}
	for _, record := range records {
		if record.ID != "" {
			byID[record.ID] = record
		}
	}
	return byID
}

func idsByStatus(records []taskRecord, status string) []string {
	ids := make([]string, 0)
	for _, record := range records {
		if record.Status == status && record.ID != "" {
			ids = append(ids, record.ID)
		}
	}
	sort.Strings(ids)
	return ids
}

func countStatus(records []taskRecord, status string) int {
	count := 0
	for _, record := range records {
		if record.Status == status {
			count++
		}
	}
	return count
}

func validStatuses() map[string]bool {
	return map[string]bool{
		"TODO":        true,
		"IN_PROGRESS": true,
		"REVIEW":      true,
		"BLOCKED":     true,
		"DONE":        true,
		"REMOVED":     true,
	}
}

func requiredFields() []string {
	return []string{"Status", "Owner", "Branch", "PR", "Risk", "Created", "Updated"}
}

func requiredSections() []string {
	return []string{"Summary", "Scope", "Acceptance Criteria", "Notes", "Agent Log"}
}

func hasSection(body string, section string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) == "## "+section {
			return true
		}
	}
	return false
}

func firstLine(body string) string {
	if index := strings.Index(body, "\n"); index >= 0 {
		return strings.TrimSpace(body[:index])
	}
	return strings.TrimSpace(body)
}

func normalizePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
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
