package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/sync"
)

func resetImportFlags() {
	importSource = ""
	importProject = ""
	importTokenEnv = ""
	importUserEnv = ""
	importBaseURL = ""
	importOutDir = "./tasks"
	importFilter = ""
	importDryRun = false
	importFormat = "table"
	importRepo = ""
	importLabels = ""
	importMilestone = ""
	importAssignee = ""
}

func TestImportCommand_NonInteractive(t *testing.T) {
	sourceName := "test-import-cli"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "CLI-1", Title: "Import task one", Status: "open", URL: "https://example.com/1"},
			{ExternalID: "CLI-2", Title: "Import task two", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	resetImportFlags()
	importSource = sourceName
	importOutDir = filepath.Join(tmpDir, "tasks")

	err := runImport(importCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were created
	entries, err := os.ReadDir(filepath.Join(tmpDir, "tasks"))
	if err != nil {
		t.Fatalf("failed to read tasks dir: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 task files, got %d", len(entries))
	}
}

func TestImportCommand_DryRun(t *testing.T) {
	sourceName := "test-import-cli-dry"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "DRY-1", Title: "Dry run task", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	resetImportFlags()
	importSource = sourceName
	importOutDir = filepath.Join(tmpDir, "tasks")
	importDryRun = true

	err := runImport(importCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No files should be created
	_, statErr := os.Stat(filepath.Join(tmpDir, "tasks"))
	if statErr == nil {
		t.Error("expected no tasks directory in dry-run mode")
	}
}

func TestImportCommand_JSONOutput(t *testing.T) {
	sourceName := "test-import-cli-json"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "JSON-1", Title: "JSON output task", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	resetImportFlags()
	importSource = sourceName
	importOutDir = filepath.Join(tmpDir, "tasks")
	importFormat = "json"

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runImport(importCmd, nil)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify valid JSON
	var data importResultData
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("invalid JSON output: %v\nOutput:\n%s", err, output)
	}

	if data.Summary.Created != 1 {
		t.Errorf("expected 1 created in summary, got %d", data.Summary.Created)
	}
	if len(data.Created) != 1 {
		t.Errorf("expected 1 created item, got %d", len(data.Created))
	}
	if data.Created[0].ExternalID != "JSON-1" {
		t.Errorf("expected external_id JSON-1, got %s", data.Created[0].ExternalID)
	}
}

func TestImportCommand_DuplicateSkip(t *testing.T) {
	sourceName := "test-import-cli-dedup"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "DUP-1", Title: "Already exists", Status: "open"},
			{ExternalID: "NEW-1", Title: "Brand new", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Create an existing task with external_id DUP-1
	tasksDir := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}
	existingContent := "---\nid: \"001\"\ntitle: \"Existing\"\nstatus: pending\nexternal_id: \"DUP-1\"\n---\n"
	if err := os.WriteFile(filepath.Join(tasksDir, "001-existing.md"), []byte(existingContent), 0644); err != nil {
		t.Fatal(err)
	}

	resetImportFlags()
	importSource = sourceName
	importOutDir = tasksDir
	importFormat = "json"

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runImport(importCmd, nil)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var data importResultData
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if data.Summary.Created != 1 {
		t.Errorf("expected 1 created, got %d", data.Summary.Created)
	}
	if data.Summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", data.Summary.Skipped)
	}
	if len(data.Skipped) != 1 || data.Skipped[0].ExternalID != "DUP-1" {
		t.Errorf("expected DUP-1 to be skipped, got: %+v", data.Skipped)
	}
}

func TestImportCommand_FilterParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name:     "single filter",
			input:    "state:open",
			expected: map[string]any{"state": "open"},
		},
		{
			name:     "multiple filters",
			input:    "state:open labels:bug",
			expected: map[string]any{"state": "open", "labels": "bug"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: map[string]any{},
		},
		{
			name:     "no colon",
			input:    "nocolon",
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseImportFilters(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d filters, got %d: %v", len(tt.expected), len(got), got)
				return
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("filter %q: expected %v, got %v", k, v, got[k])
				}
			}
		})
	}
}

func TestImportCommand_OutputDirFlag(t *testing.T) {
	sourceName := "test-import-cli-outdir"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "DIR-1", Title: "Custom dir task", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	customDir := filepath.Join(tmpDir, "custom", "output")

	resetImportFlags()
	importSource = sourceName
	importOutDir = customDir

	err := runImport(importCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(customDir)
	if err != nil {
		t.Fatalf("failed to read custom dir: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 task file in custom dir, got %d", len(entries))
	}
}

func TestImportCommand_FilterFlagPopulatesConfig(t *testing.T) {
	sourceName := "test-import-cli-filter-cfg"
	defer sync.Unregister(sourceName)

	sync.Register(&cliMockSource{
		name: sourceName,
		tasks: []sync.ExternalTask{
			{ExternalID: "F-1", Title: "Filtered", Status: "open"},
		},
	})

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	resetImportFlags()
	importSource = sourceName
	importProject = "owner/repo"
	importTokenEnv = "GITHUB_TOKEN"
	importFilter = "state:open labels:bug"
	importOutDir = filepath.Join(tmpDir, "tasks")

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters == nil {
		t.Fatal("expected filters to be set")
	}
	if cfg.SourceCfg.Filters["state"] != "open" {
		t.Errorf("expected state filter 'open', got %v", cfg.SourceCfg.Filters["state"])
	}
	if cfg.SourceCfg.Filters["labels"] != "bug" {
		t.Errorf("expected labels filter 'bug', got %v", cfg.SourceCfg.Filters["labels"])
	}
}

func TestImportCommand_InvalidFormat(t *testing.T) {
	resetImportFlags()
	importSource = "something"
	importFormat = "invalid"

	err := runImport(importCmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestImportCommand_RepoFlagAliasesProject(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "myorg/myrepo"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Project != "myorg/myrepo" {
		t.Errorf("expected project=myorg/myrepo, got %q", cfg.SourceCfg.Project)
	}
}

func TestImportCommand_ProjectTakesPrecedenceOverRepo(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importProject = "explicit/project"
	importRepo = "fallback/repo"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Project != "explicit/project" {
		t.Errorf("expected project=explicit/project, got %q", cfg.SourceCfg.Project)
	}
}

func TestImportCommand_RepoFlagIgnoredForNonGitHub(t *testing.T) {
	resetImportFlags()
	importSource = "jira"
	importRepo = "should-be-ignored"
	importProject = "PROJ"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Project != "PROJ" {
		t.Errorf("expected project=PROJ, got %q", cfg.SourceCfg.Project)
	}
}

func TestImportCommand_GitHubDefaultTokenEnv(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.TokenEnv != "GITHUB_TOKEN" {
		t.Errorf("expected default token env=GITHUB_TOKEN, got %q", cfg.SourceCfg.TokenEnv)
	}
}

func TestImportCommand_GitHubDefaultStateOpen(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters == nil {
		t.Fatal("expected filters to be set")
	}
	if cfg.SourceCfg.Filters["state"] != "open" {
		t.Errorf("expected default state=open, got %v", cfg.SourceCfg.Filters["state"])
	}
}

func TestImportCommand_GitHubStateNotOverridden(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importFilter = "state:closed"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters["state"] != "closed" {
		t.Errorf("expected state=closed (from --filter), got %v", cfg.SourceCfg.Filters["state"])
	}
}

func TestImportCommand_LabelsShortcutFlag(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importLabels = "bug,critical"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters["labels"] != "bug,critical" {
		t.Errorf("expected labels=bug,critical, got %v", cfg.SourceCfg.Filters["labels"])
	}
}

func TestImportCommand_MilestoneShortcutFlag(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importMilestone = "v1.0"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters["milestone"] != "v1.0" {
		t.Errorf("expected milestone=v1.0, got %v", cfg.SourceCfg.Filters["milestone"])
	}
}

func TestImportCommand_AssigneeShortcutFlag(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importAssignee = "alice"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters["assignee"] != "alice" {
		t.Errorf("expected assignee=alice, got %v", cfg.SourceCfg.Filters["assignee"])
	}
}

func TestImportCommand_FilterFlagTakesPrecedenceOverShortcut(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importFilter = "labels:from-filter"
	importLabels = "from-shortcut"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// --filter labels take precedence over --labels shortcut
	if cfg.SourceCfg.Filters["labels"] != "from-filter" {
		t.Errorf("expected labels=from-filter (from --filter), got %v", cfg.SourceCfg.Filters["labels"])
	}
}

func TestImportCommand_AllShortcutFlags(t *testing.T) {
	resetImportFlags()
	importSource = "github"
	importRepo = "owner/repo"
	importLabels = "bug"
	importMilestone = "v2.0"
	importAssignee = "bob"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SourceCfg.Filters["labels"] != "bug" {
		t.Errorf("expected labels=bug, got %v", cfg.SourceCfg.Filters["labels"])
	}
	if cfg.SourceCfg.Filters["milestone"] != "v2.0" {
		t.Errorf("expected milestone=v2.0, got %v", cfg.SourceCfg.Filters["milestone"])
	}
	if cfg.SourceCfg.Filters["assignee"] != "bob" {
		t.Errorf("expected assignee=bob, got %v", cfg.SourceCfg.Filters["assignee"])
	}
	if cfg.SourceCfg.Filters["state"] != "open" {
		t.Errorf("expected default state=open, got %v", cfg.SourceCfg.Filters["state"])
	}
}

func TestImportCommand_NonGitHubNoDefaultState(t *testing.T) {
	resetImportFlags()
	importSource = "jira"
	importProject = "PROJ"
	importTokenEnv = "JIRA_TOKEN"

	cfg, err := buildImportConfigFromFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-github sources should not get default state filter
	if cfg.SourceCfg.Filters != nil {
		t.Errorf("expected nil filters for non-github source, got %v", cfg.SourceCfg.Filters)
	}
}
