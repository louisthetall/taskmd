package template

import (
	"os"
	"path/filepath"
	"testing"
)

const featureTemplate = `---
_template:
  name: feature
  description: "Feature request template"
title: "{{title}}"
id: "{{id}}"
status: pending
priority: medium
type: feature
tags: []
created: "{{date}}"
---

# {{title}}

## Objective

## Tasks

- [ ] TODO
`

const bugTemplate = `---
_template:
  name: bug
  description: "Bug report template"
title: "{{title}}"
id: "{{id}}"
status: pending
priority: high
type: bug
tags: []
created: "{{date}}"
---

# {{title}}

## Steps to Reproduce
1. ...
`

func writeTemplateFile(t *testing.T, dir, filename, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create dir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", filename, err)
	}
}

func TestDiscover_BuiltinOnly(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"feature": featureTemplate,
		"bug":     bugTemplate,
	}

	templates := Discover(tmpDir, filepath.Join(tmpDir, "home"))

	if len(templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(templates))
	}

	names := make(map[string]bool)
	for _, tmpl := range templates {
		names[tmpl.Name] = true
		if tmpl.Source != "built-in" {
			t.Errorf("expected source 'built-in', got %q", tmpl.Source)
		}
	}
	if !names["feature"] || !names["bug"] {
		t.Error("expected both feature and bug templates")
	}
}

func TestDiscover_ProjectOverridesBuiltin(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"feature": featureTemplate,
		"bug":     bugTemplate,
	}

	// Create a project-level feature template
	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	customFeature := `---
_template:
  name: feature
  description: "Custom project feature"
title: "{{title}}"
id: "{{id}}"
status: pending
priority: low
---

# {{title}}
`
	writeTemplateFile(t, projectDir, "feature.md", customFeature)

	templates := Discover(tmpDir, filepath.Join(tmpDir, "home"))

	// Should have 2 templates: project feature + built-in bug
	if len(templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(templates))
	}

	for _, tmpl := range templates {
		if tmpl.Name == "feature" {
			if tmpl.Source != "project" {
				t.Errorf("feature should come from project, got %q", tmpl.Source)
			}
			if tmpl.Description != "Custom project feature" {
				t.Errorf("feature should use project description, got %q", tmpl.Description)
			}
		}
		if tmpl.Name == "bug" && tmpl.Source != "built-in" {
			t.Errorf("bug should come from built-in, got %q", tmpl.Source)
		}
	}
}

func TestDiscover_UserOverridesBuiltin(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"bug": bugTemplate,
	}

	userDir := filepath.Join(homeDir, ".taskmd", "templates")
	customBug := `---
_template:
  name: bug
  description: "User bug template"
title: "{{title}}"
id: "{{id}}"
status: pending
priority: critical
---

# {{title}}
`
	writeTemplateFile(t, userDir, "bug.md", customBug)

	templates := Discover(tmpDir, homeDir)

	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
	if templates[0].Source != "user" {
		t.Errorf("bug should come from user, got %q", templates[0].Source)
	}
}

func TestDiscover_ProjectOverridesUser(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()
	BuiltinTemplates = nil

	// User-level
	userDir := filepath.Join(homeDir, ".taskmd", "templates")
	userContent := `---
_template:
  name: feature
  description: "User feature"
title: "{{title}}"
id: "{{id}}"
status: pending
---

# {{title}}
`
	writeTemplateFile(t, userDir, "feature.md", userContent)

	// Project-level
	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	projectContent := `---
_template:
  name: feature
  description: "Project feature"
title: "{{title}}"
id: "{{id}}"
status: pending
---

# {{title}}
`
	writeTemplateFile(t, projectDir, "feature.md", projectContent)

	templates := Discover(tmpDir, homeDir)

	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
	if templates[0].Source != "project" {
		t.Errorf("expected project source, got %q", templates[0].Source)
	}
	if templates[0].Description != "Project feature" {
		t.Errorf("expected project description, got %q", templates[0].Description)
	}
}

func TestResolve_FindsBuiltin(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"bug": bugTemplate,
	}

	tmpl, ok := Resolve("bug", tmpDir, filepath.Join(tmpDir, "home"))
	if !ok {
		t.Fatal("expected to find 'bug' template")
	}
	if tmpl.Name != "bug" {
		t.Errorf("expected name 'bug', got %q", tmpl.Name)
	}
	if tmpl.Source != "built-in" {
		t.Errorf("expected source 'built-in', got %q", tmpl.Source)
	}
}

func TestResolve_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()
	BuiltinTemplates = nil

	_, ok := Resolve("nonexistent", tmpDir, filepath.Join(tmpDir, "home"))
	if ok {
		t.Error("expected not found for nonexistent template")
	}
}

func TestResolve_ProjectFirst(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"bug": bugTemplate,
	}

	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	customBug := `---
_template:
  name: bug
  description: "Project bug"
title: "{{title}}"
id: "{{id}}"
status: pending
---

# {{title}}
`
	writeTemplateFile(t, projectDir, "bug.md", customBug)

	tmpl, ok := Resolve("bug", tmpDir, filepath.Join(tmpDir, "home"))
	if !ok {
		t.Fatal("expected to find 'bug' template")
	}
	if tmpl.Source != "project" {
		t.Errorf("expected project source, got %q", tmpl.Source)
	}
}

func TestDiscover_IgnoresNonMdFiles(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()
	BuiltinTemplates = nil

	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	writeTemplateFile(t, projectDir, "bug.md", bugTemplate)
	// Write a non-md file
	os.WriteFile(filepath.Join(projectDir, "README.txt"), []byte("ignore me"), 0644)

	templates := Discover(tmpDir, filepath.Join(tmpDir, "home"))

	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
}

func TestDiscover_IgnoresInvalidTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()
	BuiltinTemplates = nil

	projectDir := filepath.Join(tmpDir, ".taskmd", "templates")
	writeTemplateFile(t, projectDir, "valid.md", bugTemplate)
	// Write an invalid template (no _template block)
	writeTemplateFile(t, projectDir, "invalid.md", `---
title: no template block
---
`)

	templates := Discover(tmpDir, filepath.Join(tmpDir, "home"))

	if len(templates) != 1 {
		t.Fatalf("expected 1 template (skipping invalid), got %d", len(templates))
	}
}

func TestDiscover_EmptyHome(t *testing.T) {
	tmpDir := t.TempDir()

	oldBuiltins := BuiltinTemplates
	defer func() { BuiltinTemplates = oldBuiltins }()

	BuiltinTemplates = map[string]string{
		"feature": featureTemplate,
	}

	// Empty home dir string should skip user templates without error
	templates := Discover(tmpDir, "")
	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
}
