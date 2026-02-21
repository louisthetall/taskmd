package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/validator"
)

func TestParseScopeEntries_WithDescription(t *testing.T) {
	scopeMap := map[string]any{
		"cli/graph": map[string]any{
			"description": "Graph visualization",
			"paths":       []any{"apps/cli/internal/graph/"},
		},
	}

	scopes := parseScopeEntries(scopeMap)

	sc, ok := scopes["cli/graph"]
	if !ok {
		t.Fatal("expected scope cli/graph to exist")
	}
	if sc.Description != "Graph visualization" {
		t.Errorf("Description = %q, want %q", sc.Description, "Graph visualization")
	}
	if len(sc.Paths) != 1 || sc.Paths[0] != "apps/cli/internal/graph/" {
		t.Errorf("Paths = %v, want [apps/cli/internal/graph/]", sc.Paths)
	}
}

func TestParseScopeEntries_WithoutDescription(t *testing.T) {
	scopeMap := map[string]any{
		"cli/output": map[string]any{
			"paths": []any{"apps/cli/internal/cli/format.go"},
		},
	}

	scopes := parseScopeEntries(scopeMap)

	sc, ok := scopes["cli/output"]
	if !ok {
		t.Fatal("expected scope cli/output to exist")
	}
	if sc.Description != "" {
		t.Errorf("Description = %q, want empty string", sc.Description)
	}
	if len(sc.Paths) != 1 {
		t.Errorf("Paths = %v, want 1 element", sc.Paths)
	}
}

// --- Helpers ---

func resetValidateFlags() {
	validateFormat = "text"
	validateStrict = false
}

func createValidateTestFiles(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	tasks := map[string]string{
		"001-valid.md": `---
id: "001"
title: "Valid Task"
status: pending
priority: high
effort: small
dependencies: []
tags: ["test"]
created: 2026-02-08
---

A valid task.
`,
		"002-valid.md": `---
id: "002"
title: "Another Valid Task"
status: completed
priority: medium
effort: medium
dependencies: ["001"]
tags: ["test"]
created: 2026-02-08
---

Another valid task.
`,
	}

	for filename, content := range tasks {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

func captureValidateOutput(t *testing.T, args []string) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runValidate(validateCmd, args)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

// --- Unit tests (no file I/O) ---

func TestMergeValidationResults(t *testing.T) {
	target := &validator.ValidationResult{
		Issues:   []validator.ValidationIssue{{Level: validator.LevelError, Message: "err1"}},
		Errors:   1,
		Warnings: 0,
	}
	source := &validator.ValidationResult{
		Issues:   []validator.ValidationIssue{{Level: validator.LevelWarning, Message: "warn1"}},
		Errors:   0,
		Warnings: 1,
	}

	mergeValidationResults(target, source)

	if len(target.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(target.Issues))
	}
	if target.Errors != 1 {
		t.Errorf("expected 1 error, got %d", target.Errors)
	}
	if target.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", target.Warnings)
	}
}

func TestOutputValidationText_NoIssues(t *testing.T) {
	result := &validator.ValidationResult{
		Issues:    []validator.ValidationIssue{},
		TaskCount: 5,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputValidationText(result, false)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "5 task(s) are valid") {
		t.Errorf("expected success message, got:\n%s", output)
	}
}

func TestOutputValidationText_WithErrors(t *testing.T) {
	result := &validator.ValidationResult{
		Issues: []validator.ValidationIssue{
			{Level: validator.LevelError, TaskID: "001", Message: "missing title"},
			{Level: validator.LevelError, TaskID: "002", Message: "bad status"},
		},
		Errors:    2,
		TaskCount: 3,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputValidationText(result, false)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "2 error(s)") {
		t.Errorf("expected error count, got:\n%s", output)
	}
}

func TestOutputValidationText_WithWarnings(t *testing.T) {
	result := &validator.ValidationResult{
		Issues: []validator.ValidationIssue{
			{Level: validator.LevelWarning, TaskID: "001", Message: "no priority"},
		},
		Warnings:  1,
		TaskCount: 2,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputValidationText(result, false)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "1 warning(s)") {
		t.Errorf("expected warning count, got:\n%s", output)
	}
}

func TestOutputValidationText_Quiet(t *testing.T) {
	result := &validator.ValidationResult{
		Issues:    []validator.ValidationIssue{},
		TaskCount: 3,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputValidationText(result, true)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output != "" {
		t.Errorf("expected no output in quiet mode, got:\n%s", output)
	}
}

func TestPrintIssue_WithTaskID(t *testing.T) {
	issue := validator.ValidationIssue{
		Level:   validator.LevelError,
		TaskID:  "042",
		Message: "something wrong",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printIssue(issue, getRenderer())

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "042") {
		t.Errorf("expected task ID in output, got:\n%s", output)
	}
	if !strings.Contains(output, "something wrong") {
		t.Errorf("expected message in output, got:\n%s", output)
	}
}

func TestPrintIssue_WithoutTaskID(t *testing.T) {
	issue := validator.ValidationIssue{
		Level:   validator.LevelError,
		Message: "global error",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printIssue(issue, getRenderer())

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "global error") {
		t.Errorf("expected message in output, got:\n%s", output)
	}
	// Should not contain bracket prefix for task ID
	if strings.Contains(output, "[") && strings.Contains(output, "]") {
		t.Errorf("expected no [ID] prefix, got:\n%s", output)
	}
}

func TestPrintIssue_WithFilePath(t *testing.T) {
	issue := validator.ValidationIssue{
		Level:    validator.LevelError,
		TaskID:   "001",
		FilePath: "tasks/001-test.md",
		Message:  "some issue",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printIssue(issue, getRenderer())

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "tasks/001-test.md") {
		t.Errorf("expected file path in output, got:\n%s", output)
	}
}

func TestOutputValidationJSON(t *testing.T) {
	result := &validator.ValidationResult{
		Issues: []validator.ValidationIssue{
			{Level: validator.LevelError, TaskID: "001", Message: "missing title"},
		},
		Errors:    1,
		TaskCount: 2,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputValidationJSON(result)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("outputValidationJSON failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	var parsed validator.ValidationResult
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, output)
	}

	if parsed.Errors != 1 {
		t.Errorf("errors = %d, want 1", parsed.Errors)
	}
	if parsed.TaskCount != 2 {
		t.Errorf("task_count = %d, want 2", parsed.TaskCount)
	}
	if len(parsed.Issues) != 1 {
		t.Errorf("issues count = %d, want 1", len(parsed.Issues))
	}
}

// --- Command-level tests ---

func TestRunValidate_ValidTasks(t *testing.T) {
	tmpDir := createValidateTestFiles(t)
	resetValidateFlags()

	output, err := captureValidateOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runValidate failed: %v", err)
	}

	if !strings.Contains(output, "2 task(s) are valid") {
		t.Errorf("expected success message, got:\n%s", output)
	}
}

func TestRunValidate_JSONFormat(t *testing.T) {
	tmpDir := createValidateTestFiles(t)
	resetValidateFlags()
	validateFormat = "json"

	output, err := captureValidateOutput(t, []string{tmpDir})
	if err != nil {
		t.Fatalf("runValidate failed: %v", err)
	}

	var parsed validator.ValidationResult
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, output)
	}

	if parsed.TaskCount != 2 {
		t.Errorf("task_count = %d, want 2", parsed.TaskCount)
	}
}

func TestRunValidate_InvalidFormat(t *testing.T) {
	tmpDir := createValidateTestFiles(t)
	resetValidateFlags()
	validateFormat = "invalid"

	_, err := captureValidateOutput(t, []string{tmpDir})
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %v", err)
	}
}
