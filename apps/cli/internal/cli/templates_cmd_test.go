package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/template"
)

func resetTemplatesFlags() {
	templatesFormat = "table"
	taskDir = "."
}

func captureTemplatesListOutput(t *testing.T) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTemplatesList(templatesListCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func TestTemplatesList_BuiltinTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir

	// Change to tmpDir so resolveProjectRoot() doesn't find the real project root
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	output, err := captureTemplatesListOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should list built-in templates
	if !strings.Contains(output, "feature") {
		t.Error("expected 'feature' template in output")
	}
	if !strings.Contains(output, "bug") {
		t.Error("expected 'bug' template in output")
	}
	if !strings.Contains(output, "chore") {
		t.Error("expected 'chore' template in output")
	}
	if !strings.Contains(output, "built-in") {
		t.Error("expected 'built-in' source in output")
	}
}

func TestTemplatesList_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir
	templatesFormat = "json"

	output, err := captureTemplatesListOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var items []templateListItem
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		t.Fatalf("failed to parse JSON: %v\nOutput: %s", err, output)
	}

	if len(items) < 3 {
		t.Fatalf("expected at least 3 templates, got %d", len(items))
	}

	names := make(map[string]bool)
	for _, item := range items {
		names[item.Name] = true
		if item.Source == "" {
			t.Error("expected non-empty source")
		}
	}
	if !names["feature"] || !names["bug"] || !names["chore"] {
		t.Error("expected feature, bug, and chore templates")
	}
}

func TestTemplatesList_YAMLFormat(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir
	templatesFormat = "yaml"

	// Change to tmpDir so resolveProjectRoot() doesn't find the real project root
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	output, err := captureTemplatesListOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "name: feature") {
		t.Error("expected YAML output with feature template")
	}
	if !strings.Contains(output, "source: built-in") {
		t.Error("expected YAML output with built-in source")
	}
}

func TestTemplatesList_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir
	templatesFormat = "xml"

	_, err := captureTemplatesListOutput(t)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %v", err)
	}
}

func TestTemplatesList_ProjectTemplatesShown(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir

	// Create a project-level template in a known project root.
	// Use the template package Discover directly to avoid viper dependency.
	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	os.MkdirAll(projectDir, 0755)
	customTemplate := `---
_template:
  name: custom
  description: "Custom project template"
title: "{{title}}"
id: "{{id}}"
status: pending
---

# {{title}}
`
	os.WriteFile(filepath.Join(projectDir, "custom.md"), []byte(customTemplate), 0644)

	// Discover templates directly with known paths to verify project templates work.
	templates := template.Discover(tmpDir, "")
	foundCustom := false
	for _, tmpl := range templates {
		if tmpl.Name == "custom" && tmpl.Source == "project" {
			foundCustom = true
		}
	}
	if !foundCustom {
		t.Error("expected 'custom' template from project source")
	}
}

func TestTemplatesList_NoTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir

	// Change to tmpDir so resolveProjectRoot() doesn't find the real project root
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Clear built-in templates
	oldBuiltins := template.BuiltinTemplates
	template.BuiltinTemplates = nil
	defer func() { template.BuiltinTemplates = oldBuiltins }()

	output, err := captureTemplatesListOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "No templates found") {
		t.Errorf("expected 'No templates found', got: %s", output)
	}
}

func TestTemplatesList_TableHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	resetTemplatesFlags()
	taskDir = tmpDir

	output, err := captureTemplatesListOutput(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "NAME") {
		t.Error("expected NAME header")
	}
	if !strings.Contains(output, "DESCRIPTION") {
		t.Error("expected DESCRIPTION header")
	}
	if !strings.Contains(output, "SOURCE") {
		t.Error("expected SOURCE header")
	}
}
