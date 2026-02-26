//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSet_DependsOn_Basic(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-setup.md", "001", "Setup", "pending", nil)
	writeTask(t, dir, "002-auth.md", "002", "Auth", "pending", []string{"001"})
	writeTask(t, dir, "003-ui.md", "003", "UI", "pending", []string{"002"})

	result := mustRun(t, dir, "set", "003", "--depends-on", "001,002")

	if !strings.Contains(result.Stdout, "dependencies:") {
		t.Errorf("Expected dependencies change in output, got: %s", result.Stdout)
	}

	content, err := os.ReadFile(filepath.Join(dir, "003-ui.md"))
	if err != nil {
		t.Fatalf("failed to read task file: %v", err)
	}
	fileStr := string(content)
	if !strings.Contains(fileStr, `"001"`) || !strings.Contains(fileStr, `"002"`) {
		t.Errorf("Expected file to contain both dep IDs, got:\n%s", fileStr)
	}
}

func TestSet_DependsOn_InvalidID(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-setup.md", "001", "Setup", "pending", nil)

	result := run(t, dir, "set", "001", "--depends-on", "999")

	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code for non-existent dependency")
	}
	combined := result.Stdout + result.Stderr
	if !strings.Contains(combined, "not found") {
		t.Errorf("Expected 'not found' in output, got: %s", combined)
	}
}

func TestSet_DependsOn_CircularDep(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-setup.md", "001", "Setup", "pending", nil)
	writeTask(t, dir, "002-auth.md", "002", "Auth", "pending", []string{"001"})
	writeTask(t, dir, "003-ui.md", "003", "UI", "pending", []string{"002"})

	// Setting 001 to depend on 003 creates a cycle: 001->003->002->001
	result := run(t, dir, "set", "001", "--depends-on", "003")

	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code for circular dependency")
	}
	combined := result.Stdout + result.Stderr
	if !strings.Contains(combined, "circular dependency") {
		t.Errorf("Expected 'circular dependency' in output, got: %s", combined)
	}
}

func TestSet_DependsOn_WithStatus(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-setup.md", "001", "Setup", "pending", nil)
	writeTask(t, dir, "002-auth.md", "002", "Auth", "pending", nil)

	result := mustRun(t, dir, "set", "002", "--depends-on", "001", "--status", "blocked")

	if !strings.Contains(result.Stdout, "status:") {
		t.Error("Expected status change in output")
	}
	if !strings.Contains(result.Stdout, "dependencies:") {
		t.Error("Expected dependencies change in output")
	}

	content, err := os.ReadFile(filepath.Join(dir, "002-auth.md"))
	if err != nil {
		t.Fatalf("failed to read task file: %v", err)
	}
	fileStr := string(content)
	if !strings.Contains(fileStr, "status: blocked") {
		t.Error("Expected file to contain updated status")
	}
	if !strings.Contains(fileStr, `"001"`) {
		t.Error("Expected file to contain dependency")
	}
}

func TestSet_DependsOn_Clear(t *testing.T) {
	dir := setupTaskDir(t)
	writeTask(t, dir, "001-setup.md", "001", "Setup", "pending", nil)
	writeTask(t, dir, "002-auth.md", "002", "Auth", "pending", []string{"001"})

	result := mustRun(t, dir, "set", "002", "--depends-on", "")

	if !strings.Contains(result.Stdout, "dependencies:") {
		t.Errorf("Expected dependencies change in output, got: %s", result.Stdout)
	}

	content, err := os.ReadFile(filepath.Join(dir, "002-auth.md"))
	if err != nil {
		t.Fatalf("failed to read task file: %v", err)
	}
	fileStr := string(content)
	if strings.Contains(fileStr, "dependencies:") {
		t.Errorf("Expected dependencies line to be removed, got:\n%s", fileStr)
	}
}
