//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRm_DeletesUnreferencedTask(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-hello.md", "001", "Hello Task", "pending", nil)

	result := mustRun(t, dir, "rm", "001", "--force")

	if !strings.Contains(result.Stdout, "Deleted 1 task") {
		t.Errorf("expected delete confirmation, got:\n%s", result.Stdout)
	}

	if _, err := os.Stat(filepath.Join(dir, "001-hello.md")); !os.IsNotExist(err) {
		t.Error("expected task file to be deleted")
	}
}

func TestRm_BlocksWhenReferenced(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-base.md", "001", "Base Task", "pending", nil)
	writeTask(t, dir, "002-dep.md", "002", "Dependent Task", "pending", []string{"001"})

	result := run(t, dir, "rm", "001")

	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code when task is referenced")
	}
	if !strings.Contains(result.Stderr, "referenced by other tasks") {
		t.Errorf("expected reference error in stderr, got:\n%s", result.Stderr)
	}
	if !strings.Contains(result.Stderr, "002") {
		t.Errorf("expected referencing task ID in stderr, got:\n%s", result.Stderr)
	}

	// File should still exist
	if _, err := os.Stat(filepath.Join(dir, "001-base.md")); err != nil {
		t.Error("expected task file to remain when blocked")
	}
}

func TestRm_ForceDeletesReferenced(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-base.md", "001", "Base Task", "pending", nil)
	writeTask(t, dir, "002-dep.md", "002", "Dependent Task", "pending", []string{"001"})

	result := mustRun(t, dir, "rm", "001", "--force")

	if !strings.Contains(result.Stdout, "Deleted 1 task") {
		t.Errorf("expected delete confirmation, got:\n%s", result.Stdout)
	}

	if _, err := os.Stat(filepath.Join(dir, "001-base.md")); !os.IsNotExist(err) {
		t.Error("expected task file to be deleted with --force")
	}
}

func TestRm_DeletesWorklog(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-hello.md", "001", "Hello Task", "pending", nil)

	// Create worklog file
	wlDir := filepath.Join(dir, ".worklogs")
	if err := os.MkdirAll(wlDir, 0o755); err != nil {
		t.Fatalf("failed to create .worklogs dir: %v", err)
	}
	wlFile := filepath.Join(wlDir, "001.md")
	if err := os.WriteFile(wlFile, []byte("## 2026-02-08T10:00:00Z\n\nStarted.\n"), 0o644); err != nil {
		t.Fatalf("failed to write worklog: %v", err)
	}

	result := mustRun(t, dir, "rm", "001", "--force")

	if !strings.Contains(result.Stdout, "Deleted worklog") {
		t.Errorf("expected worklog deletion message, got:\n%s", result.Stdout)
	}

	if _, err := os.Stat(wlFile); !os.IsNotExist(err) {
		t.Error("expected worklog file to be deleted")
	}
}

func TestRm_DryRun(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-hello.md", "001", "Hello Task", "pending", nil)

	result := mustRun(t, dir, "rm", "001", "--dry-run")

	if !strings.Contains(result.Stdout, "Dry run") {
		t.Errorf("expected dry run message, got:\n%s", result.Stdout)
	}

	// File should still exist
	if _, err := os.Stat(filepath.Join(dir, "001-hello.md")); err != nil {
		t.Error("expected task file to remain after dry run")
	}
}

func TestRm_TaskNotFound(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-hello.md", "001", "Hello Task", "pending", nil)

	result := run(t, dir, "rm", "999")

	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code for missing task")
	}
	if !strings.Contains(result.Stderr, "task not found") {
		t.Errorf("expected 'task not found' in stderr, got:\n%s", result.Stderr)
	}
}
