package mcpb_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// manifestTemplate mirrors the subset of manifest.json fields we validate.
type manifestTemplate struct {
	ManifestVersion string `json:"manifest_version"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Description     string `json:"description"`
	Server          struct {
		Type       string `json:"type"`
		EntryPoint string `json:"entry_point"`
		MCPConfig  struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcp_config"`
	} `json:"server"`
	Tools []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"tools"`
}

func templatePath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file path")
	}
	return filepath.Join(filepath.Dir(filename), "manifest.template.json")
}

func loadTemplate(t *testing.T) manifestTemplate {
	t.Helper()
	data, err := os.ReadFile(templatePath(t))
	if err != nil {
		t.Fatalf("reading manifest template: %v", err)
	}

	var m manifestTemplate
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("manifest template is not valid JSON: %v", err)
	}
	return m
}

func TestManifestTemplate_ValidJSON(t *testing.T) {
	data, err := os.ReadFile(templatePath(t))
	if err != nil {
		t.Fatalf("reading manifest template: %v", err)
	}
	if !json.Valid(data) {
		t.Fatal("manifest.template.json is not valid JSON")
	}
}

func TestManifestTemplate_RequiredFields(t *testing.T) {
	m := loadTemplate(t)

	if m.ManifestVersion == "" {
		t.Error("manifest_version is missing")
	}
	if m.Name == "" {
		t.Error("name is missing")
	}
	if m.Version == "" {
		t.Error("version is missing (should be VERSION_PLACEHOLDER)")
	}
	if m.Description == "" {
		t.Error("description is missing")
	}
	if m.Server.Type == "" {
		t.Error("server.type is missing")
	}
	if m.Server.Type != "binary" {
		t.Errorf("server.type = %q, want %q", m.Server.Type, "binary")
	}
	if m.Server.EntryPoint == "" {
		t.Error("server.entry_point is missing")
	}
	if m.Server.MCPConfig.Command == "" {
		t.Error("server.mcp_config.command is missing")
	}
	if len(m.Server.MCPConfig.Args) == 0 || m.Server.MCPConfig.Args[0] != "mcp" {
		t.Errorf("server.mcp_config.args = %v, want [\"mcp\"]", m.Server.MCPConfig.Args)
	}
}

func TestManifestTemplate_ToolList(t *testing.T) {
	m := loadTemplate(t)

	// These are the 9 MCP tools registered in internal/mcp/server.go.
	expectedTools := []string{
		"list",
		"get",
		"next",
		"search",
		"context",
		"set",
		"validate",
		"graph",
		"status",
	}

	if len(m.Tools) != len(expectedTools) {
		t.Fatalf("got %d tools, want %d", len(m.Tools), len(expectedTools))
	}

	toolNames := make(map[string]bool)
	for _, tool := range m.Tools {
		if tool.Name == "" {
			t.Error("tool has empty name")
		}
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		toolNames[tool.Name] = true
	}

	for _, name := range expectedTools {
		if !toolNames[name] {
			t.Errorf("expected tool %q not found in manifest", name)
		}
	}
}

func TestManifestTemplate_VersionPlaceholder(t *testing.T) {
	m := loadTemplate(t)
	if m.Version != "VERSION_PLACEHOLDER" {
		t.Errorf("version = %q, want %q (should be placeholder for build-time substitution)", m.Version, "VERSION_PLACEHOLDER")
	}
}
