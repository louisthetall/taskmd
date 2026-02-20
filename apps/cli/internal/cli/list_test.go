package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/driangle/taskmd/apps/cli/internal/model"
)

func TestSortTasks(t *testing.T) {
	now := time.Now()
	tasks := []*model.Task{
		{ID: "003", Title: "C", Status: model.StatusPending, Priority: model.PriorityLow, Effort: model.EffortLarge, Created: now.Add(2 * time.Hour)},
		{ID: "001", Title: "A", Status: model.StatusCompleted, Priority: model.PriorityHigh, Effort: model.EffortSmall, Created: now},
		{ID: "002", Title: "B", Status: model.StatusInProgress, Priority: model.PriorityCritical, Effort: model.EffortMedium, Created: now.Add(1 * time.Hour)},
	}

	tests := []struct {
		name      string
		sortField string
		firstID   string
		wantErr   bool
	}{
		{"sort by id", "id", "001", false},
		{"sort by title", "title", "001", false},
		{"sort by priority", "priority", "002", false},
		{"sort by effort", "effort", "001", false},
		{"sort by created", "created", "001", false},
		{"invalid sort field", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasksCopy := make([]*model.Task, len(tasks))
			copy(tasksCopy, tasks)

			err := sortTasks(tasksCopy, tt.sortField)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tasksCopy[0].ID != tt.firstID {
				t.Errorf("sortTasks() first task ID = %s, want %s", tasksCopy[0].ID, tt.firstID)
			}
		})
	}
}

func TestGetColumnValue(t *testing.T) {
	created := time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC)
	task := &model.Task{
		ID:           "001",
		Title:        "Test Task",
		Status:       model.StatusPending,
		Priority:     model.PriorityHigh,
		Effort:       model.EffortSmall,
		Group:        "testing",
		Owner:        "alice",
		Created:      created,
		Dependencies: []string{"002", "003"},
		Tags:         []string{"cli", "test"},
	}

	tests := []struct {
		name     string
		column   string
		expected string
	}{
		{"id column", "id", "001"},
		{"title column", "title", "Test Task"},
		{"status column", "status", "pending"},
		{"priority column", "priority", "high"},
		{"effort column", "effort", "small"},
		{"group column", "group", "testing"},
		{"owner column", "owner", "alice"},
		{"created column", "created", "2026-02-08"},
		{"deps column", "deps", "002,003"},
		{"tags column", "tags", "cli,test"},
		{"unknown column", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getColumnValue(task, tt.column)
			if result != tt.expected {
				t.Errorf("getColumnValue(%s) = %s, want %s", tt.column, result, tt.expected)
			}
		})
	}
}

// resetListFlags resets list command flags to defaults before each test.
func resetListFlags() {
	listFilters = []string{}
	listSort = ""
	listColumns = "id,title,status,priority,file"
	noColor = true
}

// captureListTableOutput runs outputTable and captures stdout.
func captureListTableOutput(t *testing.T, tasks []*model.Task, columns string) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputTable(tasks, columns)
	if err != nil {
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("outputTable failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestListCommand_TableColorEnabled(t *testing.T) {
	resetListFlags()
	noColor = false
	forceColor = true
	defer func() { forceColor = false }()
	os.Unsetenv("NO_COLOR")

	tasks := []*model.Task{
		{ID: "001", Title: "Test Task", Status: model.StatusPending, Priority: model.PriorityHigh, FilePath: "test.md"},
	}

	output := captureListTableOutput(t, tasks, "id,title,status,priority")

	// With colors enabled, output should contain ANSI escape codes
	if !strings.Contains(output, "\x1b[") {
		t.Error("Expected colored table output to contain ANSI escape codes")
	}

	// Task data should still be present
	if !strings.Contains(output, "001") {
		t.Error("Expected task ID in output")
	}
	if !strings.Contains(output, "pending") {
		t.Error("Expected status in output")
	}
}

func TestListCommand_TableNoColorFlag(t *testing.T) {
	resetListFlags()
	// noColor is already true from resetListFlags
	os.Unsetenv("NO_COLOR")

	tasks := []*model.Task{
		{ID: "001", Title: "Test Task", Status: model.StatusPending, Priority: model.PriorityHigh, FilePath: "test.md"},
	}

	output := captureListTableOutput(t, tasks, "id,title,status,priority")

	// With no-color, output should NOT contain ANSI escape codes
	if strings.Contains(output, "\x1b[") {
		t.Error("Expected no ANSI codes in no-color table output")
	}

	// Task data should still be present
	if !strings.Contains(output, "001") {
		t.Error("Expected task ID in output")
	}
	if !strings.Contains(output, "pending") {
		t.Error("Expected status in output")
	}
}

func TestListCommand_TableNoColorEnvVar(t *testing.T) {
	resetListFlags()
	noColor = false // enable via flag, but env var should override
	forceColor = true
	defer func() { forceColor = false }()

	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	tasks := []*model.Task{
		{ID: "001", Title: "Test Task", Status: model.StatusPending, Priority: model.PriorityHigh, FilePath: "test.md"},
	}

	output := captureListTableOutput(t, tasks, "id,title,status,priority")

	// NO_COLOR env var should disable colors
	if strings.Contains(output, "\x1b[") {
		t.Error("Expected no ANSI codes when NO_COLOR env var is set")
	}
}

func TestListCommand_TableColorColumns(t *testing.T) {
	resetListFlags()
	noColor = false
	forceColor = true
	defer func() { forceColor = false }()
	os.Unsetenv("NO_COLOR")

	tasks := []*model.Task{
		{ID: "001", Title: "Test Task", Status: model.StatusCompleted, Priority: model.PriorityCritical, FilePath: "test.md"},
		{ID: "002", Title: "Another", Status: model.StatusInProgress, Priority: model.PriorityLow, FilePath: "test2.md"},
	}

	output := captureListTableOutput(t, tasks, "id,title,status,priority")

	// Verify colored output contains task data
	if !strings.Contains(output, "001") {
		t.Error("Expected task 001 in output")
	}
	if !strings.Contains(output, "002") {
		t.Error("Expected task 002 in output")
	}
	if !strings.Contains(output, "Test Task") {
		t.Error("Expected title in output")
	}
}

func TestGetColumnValue_Parent(t *testing.T) {
	task := &model.Task{
		ID:     "001",
		Title:  "Test Task",
		Parent: "010",
	}

	result := getColumnValue(task, "parent")
	if result != "010" {
		t.Errorf("getColumnValue(parent) = %s, want 010", result)
	}

	// Empty parent
	task2 := &model.Task{ID: "002", Title: "No Parent"}
	result2 := getColumnValue(task2, "parent")
	if result2 != "" {
		t.Errorf("getColumnValue(parent) for task without parent = %q, want empty", result2)
	}
}

func TestListCommand_ParentColumn(t *testing.T) {
	resetListFlags()

	tasks := []*model.Task{
		{ID: "001", Title: "Parent", Status: model.StatusPending},
		{ID: "002", Title: "Child", Status: model.StatusPending, Parent: "001"},
	}

	output := captureListTableOutput(t, tasks, "id,title,parent")

	if !strings.Contains(output, "parent") {
		t.Error("Expected 'parent' column header in output")
	}
	if !strings.Contains(output, "001") {
		t.Error("Expected parent value '001' in output")
	}
}

func TestListCommand_TypeColumn(t *testing.T) {
	resetListFlags()

	tasks := []*model.Task{
		{ID: "001", Title: "Feature", Status: model.StatusPending, Type: model.TypeFeature},
		{ID: "002", Title: "Bug", Status: model.StatusPending, Type: model.TypeBug},
		{ID: "003", Title: "No type", Status: model.StatusPending},
	}

	output := captureListTableOutput(t, tasks, "id,title,type")

	if !strings.Contains(output, "type") {
		t.Error("Expected 'type' column header in output")
	}
	if !strings.Contains(output, "feature") {
		t.Error("Expected 'feature' in output")
	}
	if !strings.Contains(output, "bug") {
		t.Error("Expected 'bug' in output")
	}
}

func TestListCommand_SeparatorAlignment(t *testing.T) {
	resetListFlags()

	tasks := []*model.Task{
		{ID: "001", Title: "A very long task title here", Status: model.StatusInProgress, Priority: model.PriorityHigh, FilePath: "tasks/cli/001.md"},
		{ID: "002", Title: "Short", Status: model.StatusPending, Priority: model.PriorityCritical, FilePath: "tasks/002.md"},
	}

	output := captureListTableOutput(t, tasks, "id,title,status,priority,file")
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines (header, separator, data), got %d", len(lines))
	}

	// The separator is the second line. With tabwriter, columns are tab-aligned,
	// so we check that each separator segment (dashes) is at least as wide as
	// its corresponding header.
	headerLine := lines[0]
	separatorLine := lines[1]

	// tabwriter replaces tabs with spaces, so split on 2+ spaces
	headerCols := splitTableColumns(headerLine)
	sepCols := splitTableColumns(separatorLine)

	if len(headerCols) != len(sepCols) {
		t.Fatalf("Header has %d columns, separator has %d columns", len(headerCols), len(sepCols))
	}

	for i, header := range headerCols {
		sep := sepCols[i]
		if len(sep) < len(header) {
			t.Errorf("Separator column %d (%q, len=%d) is shorter than header (%q, len=%d)",
				i, sep, len(sep), header, len(header))
		}
		// Separator should be all dashes
		if strings.Trim(sep, "-") != "" {
			t.Errorf("Separator column %d should be all dashes, got %q", i, sep)
		}
	}

	// Also verify that "title" separator matches the longest title, not just the header
	titleIdx := -1
	for i, h := range headerCols {
		if h == "title" {
			titleIdx = i
			break
		}
	}
	if titleIdx >= 0 {
		titleSep := sepCols[titleIdx]
		longestTitle := "A very long task title here"
		if len(titleSep) < len(longestTitle) {
			t.Errorf("Title separator (%d dashes) should be at least as wide as longest title (%d chars)",
				len(titleSep), len(longestTitle))
		}
	}
}

func TestListCommand_ColorAlignmentMatchesPlain(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Task One", Status: model.StatusInProgress, Priority: model.PriorityHigh, FilePath: "a.md"},
		{ID: "002", Title: "Task Two", Status: model.StatusPending, Priority: model.PriorityCritical, FilePath: "b.md"},
	}

	// Capture with no color
	resetListFlags()
	noColor = true
	plainOutput := captureListTableOutput(t, tasks, "id,title,status,priority,file")

	// Capture with color
	resetListFlags()
	noColor = false
	forceColor = true
	defer func() { forceColor = false }()
	os.Unsetenv("NO_COLOR")
	colorOutput := captureListTableOutput(t, tasks, "id,title,status,priority,file")

	plainLines := strings.Split(strings.TrimRight(plainOutput, "\n"), "\n")
	colorLines := strings.Split(strings.TrimRight(colorOutput, "\n"), "\n")

	if len(plainLines) != len(colorLines) {
		t.Fatalf("Line count mismatch: plain=%d, color=%d", len(plainLines), len(colorLines))
	}

	// Strip ANSI codes and compare visible widths per line
	for i := range plainLines {
		plainLen := len(strings.TrimRight(plainLines[i], " "))
		strippedColor := StripANSI(colorLines[i])
		colorLen := len(strings.TrimRight(strippedColor, " "))
		if plainLen != colorLen {
			t.Errorf("Line %d visible width mismatch: plain=%d, color(stripped)=%d\n  plain: %q\n  color: %q",
				i, plainLen, colorLen, plainLines[i], strippedColor)
		}
	}
}

// splitTableColumns splits a table-formatted line into columns.
// Columns are separated by 2+ spaces.
func splitTableColumns(line string) []string {
	parts := strings.Fields(line)
	return parts
}
