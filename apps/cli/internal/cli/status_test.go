package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createStatusTestFiles(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	tasks := map[string]string{
		"001-setup.md": `---
id: "001"
title: "Setup project"
status: completed
priority: high
effort: small
dependencies: []
tags: ["infra", "setup"]
created: 2026-02-08
---

# Setup project

Initial project setup with build tooling.
`,
		"002-auth.md": `---
id: "002"
title: "Implement authentication"
status: in-progress
priority: critical
effort: large
dependencies: ["001"]
tags: ["backend", "security"]
owner: "alice"
created: 2026-02-08
---

# Implement authentication

Add JWT-based auth with refresh tokens.
`,
		"003-ui.md": `---
id: "003"
title: "Build UI components"
status: pending
priority: medium
effort: medium
dependencies: ["002"]
tags: ["frontend"]
parent: "001"
created: 2026-02-08
---

# Build UI components

Create reusable component library.
`,
	}

	for filename, content := range tasks {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

func resetStatusFlags() {
	statusFormat = "text"
	statusExact = false
	statusThreshold = 0.6
	statusMinimal = false
	statusStatusline = false
	statusScope = ""
	taskDir = "."
}

func captureStatusOutput(t *testing.T, query string) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runStatus(statusCmd, []string{query})
	if err != nil {
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("runStatus failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestStatus_ExactMatchByID(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	output := captureStatusOutput(t, "001")

	if !strings.Contains(output, "Task: 001") {
		t.Error("Expected output to contain 'Task: 001'")
	}
	if !strings.Contains(output, "Title: Setup project") {
		t.Error("Expected output to contain task title")
	}
}

func TestStatus_TextFormat(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	output := captureStatusOutput(t, "002")

	expected := []string{
		"Task: 002",
		"Title: Implement authentication",
		"Status: in-progress",
		"Priority: critical",
		"Effort: large",
		"Tags: backend, security",
		"Owner: alice",
		"Created: 2026-02-08",
		"Dependencies: 001",
		"File:",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected output to contain %q", exp)
		}
	}

	// Verify no body content is present
	if strings.Contains(output, "Description:") {
		t.Error("Status output should not contain Description section")
	}
	if strings.Contains(output, "Add JWT-based auth") {
		t.Error("Status output should not contain body content")
	}
}

func TestStatus_TextFormat_ParentField(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	output := captureStatusOutput(t, "003")

	if !strings.Contains(output, "Parent: 001") {
		t.Error("Expected output to contain 'Parent: 001'")
	}
}

func TestStatus_JSONFormat(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "json"

	output := captureStatusOutput(t, "002")

	var result statusOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	if result.ID != "002" {
		t.Errorf("Expected ID '002', got %q", result.ID)
	}
	if result.Title != "Implement authentication" {
		t.Errorf("Expected title 'Implement authentication', got %q", result.Title)
	}
	if result.Status != "in-progress" {
		t.Errorf("Expected status 'in-progress', got %q", result.Status)
	}
	if result.Priority != "critical" {
		t.Errorf("Expected priority 'critical', got %q", result.Priority)
	}
	if result.Effort != "large" {
		t.Errorf("Expected effort 'large', got %q", result.Effort)
	}
	if result.Owner != "alice" {
		t.Errorf("Expected owner 'alice', got %q", result.Owner)
	}
	if len(result.Dependencies) != 1 || result.Dependencies[0] != "001" {
		t.Errorf("Expected dependencies [001], got %v", result.Dependencies)
	}

	// Verify no content/body field in JSON
	var raw map[string]any
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		t.Fatalf("Failed to parse raw JSON: %v", err)
	}
	if _, ok := raw["content"]; ok {
		t.Error("JSON output should not contain 'content' key")
	}
	if _, ok := raw["body"]; ok {
		t.Error("JSON output should not contain 'body' key")
	}
}

func TestStatus_YAMLFormat(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "yaml"

	output := captureStatusOutput(t, "001")

	expected := []string{"id: \"001\"", "title: Setup project", "status: completed"}
	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected YAML output to contain %q", exp)
		}
	}

	// Verify no content field
	if strings.Contains(output, "content:") {
		t.Error("YAML output should not contain 'content' field")
	}
}

func TestStatus_UnsupportedFormat(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "csv"

	err := runStatus(statusCmd, []string{"001"})
	if err == nil {
		t.Fatal("Expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("Expected 'unsupported format' error, got: %v", err)
	}
}

func TestStatus_TaskNotFound_ExactMode(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusExact = true

	err := runStatus(statusCmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("Expected error for non-matching query in exact mode")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestStatus_FuzzyMatch(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	// "auth" is a substring of "Implement authentication"
	statusStdinReader = strings.NewReader("1\n")
	defer func() { statusStdinReader = os.Stdin }()

	output := captureStatusOutput(t, "auth")

	if !strings.Contains(output, "Task: 002") {
		t.Error("Expected fuzzy match to find task 002")
	}
}

func TestStatus_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	resetStatusFlags()
	taskDir = tmpDir

	err := runStatus(statusCmd, []string{"anything"})
	if err == nil {
		t.Fatal("Expected error for empty directory")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func createStatusTestFilesWithChildren(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	tasks := map[string]string{
		"010-parent.md": `---
id: "010"
title: "Parent task"
status: in-progress
tags: []
dependencies: []
---

# Parent task
`,
		"011-child-a.md": `---
id: "011"
title: "Child A"
status: pending
parent: "010"
tags: []
dependencies: []
---

# Child A
`,
		"012-child-b.md": `---
id: "012"
title: "Child B"
status: completed
parent: "010"
tags: []
dependencies: []
---

# Child B
`,
		"013-grandchild.md": `---
id: "013"
title: "Grandchild"
status: pending
parent: "011"
tags: []
dependencies: []
---

# Grandchild
`,
	}

	for filename, content := range tasks {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

func TestStatus_ChildrenTree(t *testing.T) {
	tmpDir := createStatusTestFilesWithChildren(t)
	resetStatusFlags()
	taskDir = tmpDir

	output := captureStatusOutput(t, "010")

	if !strings.Contains(output, "Children:") {
		t.Error("Expected output to contain 'Children:' section")
	}
	if !strings.Contains(output, "011") {
		t.Error("Expected output to contain child ID '011'")
	}
	if !strings.Contains(output, "Child A") {
		t.Error("Expected output to contain child title 'Child A'")
	}
	if !strings.Contains(output, "012") {
		t.Error("Expected output to contain child ID '012'")
	}
	if !strings.Contains(output, "Child B") {
		t.Error("Expected output to contain child title 'Child B'")
	}
	if !strings.Contains(output, "013") {
		t.Error("Expected output to contain grandchild ID '013'")
	}
	if !strings.Contains(output, "Grandchild") {
		t.Error("Expected output to contain grandchild title 'Grandchild'")
	}
}

func TestStatus_ChildrenTree_JSON(t *testing.T) {
	tmpDir := createStatusTestFilesWithChildren(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "json"

	output := captureStatusOutput(t, "010")

	var result statusOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, output)
	}

	if len(result.Children) == 0 {
		t.Fatal("Expected children in JSON output")
	}

	// Find child with grandchild
	var foundGrandchild bool
	for _, child := range result.Children {
		if child.ID == "011" {
			if len(child.Children) == 0 {
				t.Error("Expected child 011 to have grandchild")
			} else if child.Children[0].ID == "013" {
				foundGrandchild = true
			}
		}
	}
	if !foundGrandchild {
		t.Error("Expected to find grandchild 013 under child 011")
	}
}

func TestStatus_ChildrenTree_Circular(t *testing.T) {
	tmpDir := t.TempDir()

	tasks := map[string]string{
		"020-a.md": `---
id: "020"
title: "Task A"
status: pending
parent: "021"
tags: []
dependencies: []
---
`,
		"021-b.md": `---
id: "021"
title: "Task B"
status: pending
parent: "020"
tags: []
dependencies: []
---
`,
	}

	for filename, content := range tasks {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	resetStatusFlags()
	taskDir = tmpDir

	// Should not hang or panic
	output := captureStatusOutput(t, "020")

	if !strings.Contains(output, "Task: 020") {
		t.Error("Expected output to contain 'Task: 020'")
	}
}

func TestStatus_MinimalFlag(t *testing.T) {
	tmpDir := createStatusTestFilesWithChildren(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusMinimal = true

	output := captureStatusOutput(t, "010")

	if strings.Contains(output, "Children:") {
		t.Error("--minimal should suppress children section")
	}
}

func TestStatus_MinimalFlag_JSON(t *testing.T) {
	tmpDir := createStatusTestFilesWithChildren(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusMinimal = true
	statusFormat = "json"

	output := captureStatusOutput(t, "010")

	var raw map[string]any
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if _, ok := raw["children"]; ok {
		t.Error("--minimal JSON output should not contain 'children' key")
	}
}

func TestStatus_NoChildren(t *testing.T) {
	tmpDir := createStatusTestFilesWithChildren(t)
	resetStatusFlags()
	taskDir = tmpDir

	// Task 012 has no children
	output := captureStatusOutput(t, "012")

	if strings.Contains(output, "Children:") {
		t.Error("Task with no children should not show 'Children:' section")
	}
}

func TestStatus_NoBodyInOutput(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	// Text format
	output := captureStatusOutput(t, "001")
	if strings.Contains(output, "Initial project setup") {
		t.Error("Text output should not contain task body")
	}

	// JSON format
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "json"
	output = captureStatusOutput(t, "001")

	var raw map[string]any
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if _, ok := raw["content"]; ok {
		t.Error("JSON should not have 'content' field")
	}
}

// captureStatusNoArgsOutput runs status with no arguments and captures stdout.
func captureStatusNoArgsOutput(t *testing.T) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runStatus(statusCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func TestStatus_NoArgs_ZeroInProgress(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "001.md"), []byte(`---
id: "001"
title: "Done task"
status: completed
tags: []
dependencies: []
---
`), 0644); err != nil {
		t.Fatal(err)
	}

	resetStatusFlags()
	taskDir = tmpDir

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output, got: %q", output)
	}
}

func TestStatus_NoArgs_OneInProgress(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Task: 002") {
		t.Error("Expected output to contain 'Task: 002'")
	}
	if !strings.Contains(output, "Implement authentication") {
		t.Error("Expected output to contain task title")
	}
}

func createStatusTestFilesMultipleInProgress(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	tasks := map[string]string{
		"001.md": `---
id: "001"
title: "First task"
status: in-progress
tags: []
dependencies: []
---
`,
		"002.md": `---
id: "002"
title: "Second task"
status: in-progress
tags: []
dependencies: []
---
`,
		"003.md": `---
id: "003"
title: "Pending task"
status: pending
tags: []
dependencies: []
---
`,
	}

	for filename, content := range tasks {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

func TestStatus_NoArgs_MultipleInProgress(t *testing.T) {
	tmpDir := createStatusTestFilesMultipleInProgress(t)
	resetStatusFlags()
	taskDir = tmpDir

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "First task") {
		t.Error("Expected output to contain 'First task'")
	}
	if !strings.Contains(output, "Second task") {
		t.Error("Expected output to contain 'Second task'")
	}
	if strings.Contains(output, "Pending task") {
		t.Error("Should not contain pending task")
	}
}

func TestStatus_Statusline_OneTask(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusStatusline = true

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "#002 Implement authentication\n"
	if output != expected {
		t.Errorf("expected %q, got: %q", expected, output)
	}
}

func TestStatus_Statusline_MultipleTasks(t *testing.T) {
	tmpDir := createStatusTestFilesMultipleInProgress(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusStatusline = true

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "(+1 more)") {
		t.Errorf("expected '(+1 more)' in output, got: %q", output)
	}
}

func TestStatus_Statusline_NoTasks(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "001.md"), []byte(`---
id: "001"
title: "Done"
status: completed
tags: []
dependencies: []
---
`), 0644); err != nil {
		t.Fatal(err)
	}

	resetStatusFlags()
	taskDir = tmpDir
	statusStatusline = true

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output, got: %q", output)
	}
}

func TestStatus_Statusline_LongTitle(t *testing.T) {
	tmpDir := t.TempDir()
	longTitle := "This is a very long task title that exceeds thirty characters"
	if err := os.WriteFile(filepath.Join(tmpDir, "001.md"), []byte(`---
id: "001"
title: "`+longTitle+`"
status: in-progress
tags: []
dependencies: []
---
`), 0644); err != nil {
		t.Fatal(err)
	}

	resetStatusFlags()
	taskDir = tmpDir
	statusStatusline = true

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "#001 " + longTitle + "\n"
	if output != expected {
		t.Errorf("expected %q, got: %q", expected, output)
	}
}

func TestStatus_NoArgs_Scope(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectories with tasks
	cliDir := filepath.Join(tmpDir, "cli")
	webDir := filepath.Join(tmpDir, "web")
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(webDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(cliDir, "001.md"), []byte(`---
id: "001"
title: "CLI task"
status: in-progress
tags: []
dependencies: []
---
`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(webDir, "002.md"), []byte(`---
id: "002"
title: "Web task"
status: in-progress
tags: []
dependencies: []
---
`), 0644); err != nil {
		t.Fatal(err)
	}

	resetStatusFlags()
	taskDir = tmpDir
	statusScope = "cli"

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "CLI task") {
		t.Error("Expected output to contain 'CLI task'")
	}
	if strings.Contains(output, "Web task") {
		t.Error("Should not contain 'Web task' when scope is 'cli'")
	}
}

func TestStatus_NoArgs_JSON(t *testing.T) {
	tmpDir := createStatusTestFilesMultipleInProgress(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "json"

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []statusOutput
	if err := json.Unmarshal([]byte(output), &results); err != nil {
		t.Fatalf("Failed to parse JSON array: %v\nOutput: %s", err, output)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestStatus_NoArgs_YAML(t *testing.T) {
	tmpDir := createStatusTestFiles(t)
	resetStatusFlags()
	taskDir = tmpDir
	statusFormat = "yaml"

	output, err := captureStatusNoArgsOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "id: \"002\"") {
		t.Errorf("Expected YAML output to contain in-progress task, got: %s", output)
	}
}
