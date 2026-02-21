package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/metrics"
	"github.com/driangle/taskmd/apps/cli/internal/model"
)

// resetStatsFlags resets all stats command flags to defaults.
func resetStatsFlags() {
	statsFormat = "table"
}

// createStatsTestFiles creates test task files with varying status/priority/effort.
func createStatsTestFiles(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	tasks := map[string]string{
		"001-pending-high.md": `---
id: "001"
title: "High Priority Pending"
status: pending
priority: high
effort: small
dependencies: []
tags: ["api"]
created: 2026-02-08
---

Pending high-priority task.
`,
		"002-completed-medium.md": `---
id: "002"
title: "Completed Medium"
status: completed
priority: medium
effort: medium
dependencies: []
tags: ["api"]
created: 2026-02-08
---

Completed medium-priority task.
`,
		"003-pending-low.md": `---
id: "003"
title: "Pending Low with Dep"
status: pending
priority: low
effort: large
dependencies: ["002"]
tags: ["frontend"]
created: 2026-02-08
---

Pending task with dependency.
`,
	}

	for filename, content := range tasks {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

// captureStatsOutput runs the stats command and captures stdout.
func captureStatsOutput(t *testing.T, args []string) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runStats(statsCmd, args)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func TestRunStats_JSONOutput(t *testing.T) {
	tmpDir := createStatsTestFiles(t)
	resetStatsFlags()
	statsFormat = "json"

	output, err := captureStatsOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runStats failed: %v", err)
	}

	var m metrics.Metrics
	if err := json.Unmarshal([]byte(output), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, output)
	}

	if m.TotalTasks != 3 {
		t.Errorf("total_tasks = %d, want 3", m.TotalTasks)
	}
	if m.TasksByStatus[model.StatusPending] != 2 {
		t.Errorf("tasks_by_status[pending] = %d, want 2", m.TasksByStatus[model.StatusPending])
	}
	if m.TasksByStatus[model.StatusCompleted] != 1 {
		t.Errorf("tasks_by_status[completed] = %d, want 1", m.TasksByStatus[model.StatusCompleted])
	}
}

func TestRunStats_YAMLOutput(t *testing.T) {
	tmpDir := createStatsTestFiles(t)
	resetStatsFlags()
	statsFormat = "yaml"

	output, err := captureStatsOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runStats failed: %v", err)
	}

	if !strings.Contains(output, "totaltasks:") {
		t.Errorf("YAML output missing 'totaltasks:':\n%s", output)
	}
}

func TestRunStats_TableOutput(t *testing.T) {
	tmpDir := createStatsTestFiles(t)
	resetStatsFlags()
	statsFormat = "table"

	output, err := captureStatsOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runStats failed: %v", err)
	}

	for _, expected := range []string{"TASK STATISTICS", "BY STATUS:", "BY PRIORITY:", "BY EFFORT:"} {
		if !strings.Contains(output, expected) {
			t.Errorf("table output missing %q:\n%s", expected, output)
		}
	}
}

func TestRunStats_InvalidFormat(t *testing.T) {
	tmpDir := createStatsTestFiles(t)
	resetStatsFlags()
	statsFormat = "invalid"

	_, err := captureStatsOutput(t, []string{tmpDir})
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %v", err)
	}
}

func TestRunStats_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	resetStatsFlags()
	statsFormat = "json"

	output, err := captureStatsOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runStats failed: %v", err)
	}

	var m metrics.Metrics
	if err := json.Unmarshal([]byte(output), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, output)
	}

	if m.TotalTasks != 0 {
		t.Errorf("total_tasks = %d, want 0", m.TotalTasks)
	}
}

func TestOutputStatsTable_EmptyBreakdowns(t *testing.T) {
	m := &metrics.Metrics{
		TotalTasks:      0,
		TasksByStatus:   map[model.Status]int{},
		TasksByPriority: map[model.Priority]int{},
		TasksByEffort:   map[model.Effort]int{},
		TasksByType:     map[model.TaskType]int{},
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputStatsTable(m)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("outputStatsTable failed: %v", err)
	}

	// Each empty breakdown should print "(none)"
	count := strings.Count(output, "(none)")
	if count < 3 {
		t.Errorf("expected at least 3 '(none)' strings (status, priority, effort), got %d\noutput:\n%s", count, output)
	}
}
