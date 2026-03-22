package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadGlobalRegistry_FileNotExist(t *testing.T) {
	t.Setenv("TASKMD_HOME_CONFIG", filepath.Join(t.TempDir(), "nonexistent.yaml"))

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(entries))
	}
}

func TestLoadGlobalRegistry_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte("dir: ./tasks\n"), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(entries))
	}
}

func TestLoadGlobalRegistry_ValidEntries(t *testing.T) {
	dir := t.TempDir()
	projA := filepath.Join(dir, "project-a")
	projB := filepath.Join(dir, "project-b")
	os.MkdirAll(projA, 0755)
	os.MkdirAll(projB, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: alpha
    name: "Project Alpha"
    path: ` + projA + `
  - id: beta
    name: "Project Beta"
    path: ` + projB + `
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].ID != "alpha" || entries[0].Name != "Project Alpha" || entries[0].Path != projA {
		t.Errorf("entry 0 mismatch: %+v", entries[0])
	}
	if entries[1].ID != "beta" || entries[1].Name != "Project Beta" || entries[1].Path != projB {
		t.Errorf("entry 1 mismatch: %+v", entries[1])
	}
}

func TestLoadGlobalRegistry_DeriveIDFromPath(t *testing.T) {
	dir := t.TempDir()
	projDir := filepath.Join(dir, "my-project")
	os.MkdirAll(projDir, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - path: ` + projDir + `
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "my-project" {
		t.Errorf("expected id 'my-project', got %q", entries[0].ID)
	}
	if entries[0].Name != "my-project" {
		t.Errorf("expected name 'my-project', got %q", entries[0].Name)
	}
}

func TestLoadGlobalRegistry_NameFallsBackToID(t *testing.T) {
	dir := t.TempDir()
	projDir := filepath.Join(dir, "proj")
	os.MkdirAll(projDir, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: custom-id
    path: ` + projDir + `
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Name != "custom-id" {
		t.Errorf("expected name to fall back to id 'custom-id', got %q", entries[0].Name)
	}
}

func TestLoadGlobalRegistry_TildeExpansion(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	projDir := filepath.Join(dir, "workplace", "my-app")
	os.MkdirAll(projDir, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: myapp
    path: ~/workplace/my-app
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path != projDir {
		t.Errorf("expected path %q, got %q", projDir, entries[0].Path)
	}
}

func TestLoadGlobalRegistry_RelativePath(t *testing.T) {
	dir := t.TempDir()
	projDir := filepath.Join(dir, "projects", "foo")
	os.MkdirAll(projDir, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: foo
    path: projects/foo
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path != projDir {
		t.Errorf("expected path %q, got %q", projDir, entries[0].Path)
	}
}

func TestLoadGlobalRegistry_MissingPath(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: broken
    name: "No Path"
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	_, err := LoadGlobalRegistry()
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "path is required") {
		t.Errorf("expected 'path is required' in error, got: %v", err)
	}
}

func TestLoadGlobalRegistry_DuplicateID(t *testing.T) {
	dir := t.TempDir()
	projA := filepath.Join(dir, "a")
	projB := filepath.Join(dir, "b")
	os.MkdirAll(projA, 0755)
	os.MkdirAll(projB, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - id: same
    path: ` + projA + `
  - id: same
    path: ` + projB + `
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	_, err := LoadGlobalRegistry()
	if err == nil {
		t.Fatal("expected error for duplicate id")
	}
	if !strings.Contains(err.Error(), "duplicate id") {
		t.Errorf("expected 'duplicate id' in error, got: %v", err)
	}
}

func TestLoadGlobalRegistry_DuplicateDerivedID(t *testing.T) {
	dir := t.TempDir()
	// Two different paths that both have basename "proj"
	projA := filepath.Join(dir, "a", "proj")
	projB := filepath.Join(dir, "b", "proj")
	os.MkdirAll(projA, 0755)
	os.MkdirAll(projB, 0755)

	cfgPath := filepath.Join(dir, "config.yaml")
	yaml := `projects:
  - path: ` + projA + `
  - path: ` + projB + `
`
	os.WriteFile(cfgPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	_, err := LoadGlobalRegistry()
	if err == nil {
		t.Fatal("expected error for duplicate derived id")
	}
	if !strings.Contains(err.Error(), "duplicate id") {
		t.Errorf("expected 'duplicate id' in error, got: %v", err)
	}
}

func TestLoadGlobalRegistry_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	projDir := filepath.Join(dir, "proj")
	os.MkdirAll(projDir, 0755)

	customPath := filepath.Join(dir, "custom-config.yaml")
	yaml := `projects:
  - id: from-env
    path: ` + projDir + `
`
	os.WriteFile(customPath, []byte(yaml), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", customPath)

	entries, err := LoadGlobalRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 || entries[0].ID != "from-env" {
		t.Errorf("expected entry from env config, got: %+v", entries)
	}
}

func TestLoadGlobalRegistry_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte("projects: [[[invalid"), 0644)
	t.Setenv("TASKMD_HOME_CONFIG", cfgPath)

	_, err := LoadGlobalRegistry()
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home dir")
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/foo/bar", filepath.Join(home, "foo/bar")},
		{"~", home},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		got, err := expandTilde(tt.input)
		if err != nil {
			t.Errorf("expandTilde(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("expandTilde(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
