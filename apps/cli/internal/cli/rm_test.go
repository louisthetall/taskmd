package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const rmTaskPending = `---
id: "001"
title: "Setup project"
status: pending
priority: high
effort: small
dependencies: []
tags: ["infra"]
created: 2026-02-08
---

# Setup project
`

const rmTaskCompleted = `---
id: "002"
title: "Old feature"
status: completed
priority: low
effort: medium
dependencies: []
tags: ["backend"]
created: 2026-02-08
---

# Old feature
`

func createRmTestFiles(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	files := map[string]string{
		"001-setup.md": rmTaskPending,
		"002-old.md":   rmTaskCompleted,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", name, err)
		}
	}

	return tmpDir
}

func resetRmFlags() {
	rmForce = false
	rmDryRun = false
	taskDir = "."
}

func captureRmOutput(t *testing.T, args []string) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runRm(rmCmd, args)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func TestRm_WithForce(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir
	rmForce = true

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Deleted 1 task") {
		t.Errorf("expected delete confirmation, got: %s", output)
	}

	// File should be gone
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}

	// Other file should remain
	if _, err := os.Stat(filepath.Join(tmpDir, "002-old.md")); err != nil {
		t.Error("expected other file to remain")
	}
}

func TestRm_DryRun(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir
	rmDryRun = true

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Dry run") {
		t.Errorf("expected dry run message, got: %s", output)
	}

	if !strings.Contains(output, "Delete 1 task") {
		t.Errorf("expected preview of task, got: %s", output)
	}

	// File should NOT be deleted
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); err != nil {
		t.Error("expected file to remain after dry run")
	}
}

func TestRm_TaskNotFound(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir
	rmForce = true

	_, err := captureRmOutput(t, []string{"999"})
	if err == nil {
		t.Fatal("expected error for nonexistent task")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("expected 'task not found' error, got: %v", err)
	}
}

func TestRm_InteractiveConfirmYes(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir

	// Simulate user typing "y"
	oldStdin := rmStdinReader
	rmStdinReader = strings.NewReader("y\n")
	defer func() { rmStdinReader = oldStdin }()

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Deleted 1 task") {
		t.Errorf("expected delete confirmation, got: %s", output)
	}

	// File should be gone
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); !os.IsNotExist(err) {
		t.Error("expected file to be deleted after confirming")
	}
}

func TestRm_InteractiveConfirmNo(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir

	// Simulate user typing "n"
	oldStdin := rmStdinReader
	rmStdinReader = strings.NewReader("n\n")
	defer func() { rmStdinReader = oldStdin }()

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Cancelled") {
		t.Errorf("expected cancellation message, got: %s", output)
	}

	// File should remain
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); err != nil {
		t.Error("expected file to remain after declining")
	}
}

func TestRm_InteractiveConfirmEmpty(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir

	// Simulate user pressing Enter (empty input = default No)
	oldStdin := rmStdinReader
	rmStdinReader = strings.NewReader("\n")
	defer func() { rmStdinReader = oldStdin }()

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Cancelled") {
		t.Errorf("expected cancellation message, got: %s", output)
	}

	// File should remain
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); err != nil {
		t.Error("expected file to remain after empty input")
	}
}

func TestRm_BlockedByDependency(t *testing.T) {
	tmpDir := t.TempDir()
	// Task 001 exists; task 003 depends on 001
	os.WriteFile(filepath.Join(tmpDir, "001-setup.md"), []byte(rmTaskPending), 0644)
	depTask := `---
id: "003"
title: "Depends on setup"
status: pending
priority: medium
effort: small
dependencies: ["001"]
tags: []
created: 2026-02-08
---

# Depends on setup
`
	os.WriteFile(filepath.Join(tmpDir, "003-dep.md"), []byte(depTask), 0644)

	resetRmFlags()
	taskDir = tmpDir

	// Without --force, should fail
	_, err := captureRmOutput(t, []string{"001"})
	if err == nil {
		t.Fatal("expected error when task is referenced")
	}
	if !strings.Contains(err.Error(), "referenced by other tasks") {
		t.Errorf("expected reference error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "003") {
		t.Errorf("expected referencing task ID in error, got: %v", err)
	}

	// File should still exist
	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); err != nil {
		t.Error("expected file to remain when blocked by reference")
	}
}

func TestRm_BlockedByParent(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "001-setup.md"), []byte(rmTaskPending), 0644)
	childTask := `---
id: "004"
title: "Child task"
status: pending
priority: medium
effort: small
dependencies: []
parent: "001"
tags: []
created: 2026-02-08
---

# Child task
`
	os.WriteFile(filepath.Join(tmpDir, "004-child.md"), []byte(childTask), 0644)

	resetRmFlags()
	taskDir = tmpDir

	_, err := captureRmOutput(t, []string{"001"})
	if err == nil {
		t.Fatal("expected error when task has child referencing it as parent")
	}
	if !strings.Contains(err.Error(), "referenced by other tasks") {
		t.Errorf("expected reference error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "004") {
		t.Errorf("expected child task ID in error, got: %v", err)
	}
}

func TestRm_ForceDeletesReferencedTask(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "001-setup.md"), []byte(rmTaskPending), 0644)
	depTask := `---
id: "003"
title: "Depends on setup"
status: pending
priority: medium
effort: small
dependencies: ["001"]
tags: []
created: 2026-02-08
---

# Depends on setup
`
	os.WriteFile(filepath.Join(tmpDir, "003-dep.md"), []byte(depTask), 0644)

	resetRmFlags()
	taskDir = tmpDir
	rmForce = true

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("expected --force to override reference check, got: %v", err)
	}
	if !strings.Contains(output, "Deleted 1 task") {
		t.Errorf("expected delete confirmation, got: %s", output)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "001-setup.md")); !os.IsNotExist(err) {
		t.Error("expected file to be deleted with --force")
	}
}

func TestRm_DeletesWorklog(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "001-setup.md"), []byte(rmTaskPending), 0644)

	// Create worklog
	wlDir := filepath.Join(tmpDir, ".worklogs")
	os.MkdirAll(wlDir, 0755)
	os.WriteFile(filepath.Join(wlDir, "001.md"), []byte("## 2026-02-08T10:00:00Z\n\nStarted work.\n"), 0644)

	resetRmFlags()
	taskDir = tmpDir
	rmForce = true

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Deleted worklog") {
		t.Errorf("expected worklog deletion message, got: %s", output)
	}

	// Worklog file should be gone
	if _, err := os.Stat(filepath.Join(wlDir, "001.md")); !os.IsNotExist(err) {
		t.Error("expected worklog file to be deleted")
	}

	// Empty .worklogs dir should be removed
	if _, err := os.Stat(wlDir); !os.IsNotExist(err) {
		t.Error("expected empty .worklogs directory to be removed")
	}
}

func TestRm_WorklogDirKeptIfNotEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "001-setup.md"), []byte(rmTaskPending), 0644)
	os.WriteFile(filepath.Join(tmpDir, "002-old.md"), []byte(rmTaskCompleted), 0644)

	// Create worklogs for both tasks
	wlDir := filepath.Join(tmpDir, ".worklogs")
	os.MkdirAll(wlDir, 0755)
	os.WriteFile(filepath.Join(wlDir, "001.md"), []byte("## 2026-02-08T10:00:00Z\n\nWork.\n"), 0644)
	os.WriteFile(filepath.Join(wlDir, "002.md"), []byte("## 2026-02-08T10:00:00Z\n\nWork.\n"), 0644)

	resetRmFlags()
	taskDir = tmpDir
	rmForce = true

	_, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 001 worklog gone, but .worklogs dir should remain (002.md still there)
	if _, err := os.Stat(filepath.Join(wlDir, "001.md")); !os.IsNotExist(err) {
		t.Error("expected 001 worklog to be deleted")
	}
	if _, err := os.Stat(wlDir); err != nil {
		t.Error("expected .worklogs directory to remain (still has 002.md)")
	}
}

func TestRm_ShowsTaskDetails(t *testing.T) {
	tmpDir := createRmTestFiles(t)
	resetRmFlags()
	taskDir = tmpDir
	rmDryRun = true

	output, err := captureRmOutput(t, []string{"001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "001") {
		t.Errorf("expected task ID in output, got: %s", output)
	}
	if !strings.Contains(output, "Setup project") {
		t.Errorf("expected task title in output, got: %s", output)
	}
}
