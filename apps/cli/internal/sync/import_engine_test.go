package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunImport_CreatesNewTasks(t *testing.T) {
	sourceName := "test-import-create"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{
			ExternalID:  "EXT-1",
			Title:       "First task",
			Description: "Task one body",
			Status:      "open",
			Priority:    "high",
			Labels:      []string{"bug"},
			Assignee:    "alice",
			URL:         "https://example.com/1",
		},
		{
			ExternalID:  "EXT-2",
			Title:       "Second task",
			Description: "Task two body",
			Status:      "closed",
		},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 2 {
		t.Fatalf("expected 2 created, got %d", len(result.Created))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}

	// Verify files were created
	for _, a := range result.Created {
		if _, err := os.Stat(a.FilePath); os.IsNotExist(err) {
			t.Errorf("expected file to exist: %s", a.FilePath)
		}
	}

	// Verify first file has correct content
	content, err := os.ReadFile(result.Created[0].FilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	s := string(content)
	if !strings.Contains(s, `title: "First task"`) {
		t.Errorf("expected title in file, got:\n%s", s)
	}
	if !strings.Contains(s, "status: pending") {
		t.Errorf("expected mapped status 'pending', got:\n%s", s)
	}
	if !strings.Contains(s, "priority: high") {
		t.Errorf("expected mapped priority 'high', got:\n%s", s)
	}
	if !strings.Contains(s, `external_id: "EXT-1"`) {
		t.Errorf("expected external_id in file, got:\n%s", s)
	}
	if !strings.Contains(s, "owner: alice") {
		t.Errorf("expected owner in file, got:\n%s", s)
	}
}

func TestRunImport_DryRun(t *testing.T) {
	sourceName := "test-import-dryrun"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{ExternalID: "EXT-1", Title: "A task", Status: "open"},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
		DryRun:     true,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 1 {
		t.Errorf("expected 1 created, got %d", len(result.Created))
	}

	// Verify no files were written
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		t.Error("expected output directory to not exist in dry-run mode")
	}

	// Verify action still has metadata
	if result.Created[0].FilePath != "" {
		t.Errorf("expected empty file path in dry-run, got %s", result.Created[0].FilePath)
	}
	if result.Created[0].LocalID == "" {
		t.Error("expected local ID to be assigned even in dry-run")
	}
}

func TestRunImport_DuplicateDetection(t *testing.T) {
	sourceName := "test-import-dedup"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{ExternalID: "EXISTING-1", Title: "Already imported", Status: "open"},
		{ExternalID: "NEW-1", Title: "New task", Status: "open"},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	// Create an existing task with external_id "EXISTING-1"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	existingContent := "---\nid: \"001\"\ntitle: \"Old task\"\nstatus: pending\nexternal_id: \"EXISTING-1\"\n---\n"
	if err := os.WriteFile(filepath.Join(outputDir, "001-old-task.md"), []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 1 {
		t.Fatalf("expected 1 created, got %d", len(result.Created))
	}
	if result.Created[0].ExternalID != "NEW-1" {
		t.Errorf("expected NEW-1 to be created, got %s", result.Created[0].ExternalID)
	}

	if len(result.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if result.Skipped[0].ExternalID != "EXISTING-1" {
		t.Errorf("expected EXISTING-1 to be skipped, got %s", result.Skipped[0].ExternalID)
	}
	if result.Skipped[0].Reason != "skipped_duplicate" {
		t.Errorf("expected reason 'skipped_duplicate', got %s", result.Skipped[0].Reason)
	}
}

func TestRunImport_SequentialIDs(t *testing.T) {
	sourceName := "test-import-seq"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{ExternalID: "EXT-1", Title: "First", Status: "open"},
		{ExternalID: "EXT-2", Title: "Second", Status: "open"},
		{ExternalID: "EXT-3", Title: "Third", Status: "open"},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 3 {
		t.Fatalf("expected 3 created, got %d", len(result.Created))
	}

	// IDs should be sequential
	ids := make(map[string]bool)
	for _, a := range result.Created {
		ids[a.LocalID] = true
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 unique IDs, got %d", len(ids))
	}

	// Check sequential order
	if result.Created[0].LocalID != "001" {
		t.Errorf("expected first ID to be 001, got %s", result.Created[0].LocalID)
	}
	if result.Created[1].LocalID != "002" {
		t.Errorf("expected second ID to be 002, got %s", result.Created[1].LocalID)
	}
	if result.Created[2].LocalID != "003" {
		t.Errorf("expected third ID to be 003, got %s", result.Created[2].LocalID)
	}
}

func TestRunImport_IDsRespectExisting(t *testing.T) {
	sourceName := "test-import-existing-ids"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{ExternalID: "NEW-1", Title: "New task", Status: "open"},
	})

	dir := t.TempDir()

	// Create existing tasks with IDs 010, 011
	cliDir := filepath.Join(dir, "tasks", "cli")
	createTaskFile(t, cliDir, "010", "Existing task ten")
	createTaskFile(t, cliDir, "011", "Existing task eleven")

	outputDir := filepath.Join(dir, "tasks", "imported")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 1 {
		t.Fatalf("expected 1 created, got %d", len(result.Created))
	}
	if result.Created[0].LocalID != "012" {
		t.Errorf("expected ID 012 (after existing 010, 011), got %s", result.Created[0].LocalID)
	}
}

func TestRunImport_FieldMapping(t *testing.T) {
	sourceName := "test-import-mapping"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{
			ExternalID: "EXT-1",
			Title:      "Bug fix",
			Status:     "open",
			Priority:   "high",
			Labels:     []string{"bug", "urgent"},
			Assignee:   "bob",
		},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(result.Created[0].FilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	s := string(content)

	if !strings.Contains(s, "status: pending") {
		t.Errorf("expected 'open' mapped to 'pending', got:\n%s", s)
	}
	if !strings.Contains(s, "priority: high") {
		t.Errorf("expected priority 'high', got:\n%s", s)
	}
	if !strings.Contains(s, `"bug"`) || !strings.Contains(s, `"urgent"`) {
		t.Errorf("expected labels mapped to tags, got:\n%s", s)
	}
	if !strings.Contains(s, "owner: bob") {
		t.Errorf("expected assignee mapped to owner, got:\n%s", s)
	}
}

func TestRunImport_SourceURLInBody(t *testing.T) {
	sourceName := "test-import-url"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{
			ExternalID:  "EXT-1",
			Title:       "Task with URL",
			Description: "Some description",
			Status:      "open",
			URL:         "https://github.com/owner/repo/issues/42",
		},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(result.Created[0].FilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !strings.Contains(string(content), "Source: https://github.com/owner/repo/issues/42") {
		t.Errorf("expected source URL in body, got:\n%s", string(content))
	}
}

func TestRunImport_UnknownSource(t *testing.T) {
	dir := t.TempDir()

	cfg := ImportConfig{
		SourceName: "nonexistent",
		SourceCfg:  SourceConfig{Name: "nonexistent"},
		OutputDir:  filepath.Join(dir, "tasks"),
		ScanDir:    dir,
	}

	_, err := RunImport(cfg)
	if err == nil {
		t.Fatal("expected error for unknown source")
	}
	if !strings.Contains(err.Error(), "unknown source") {
		t.Errorf("expected 'unknown source' in error, got: %v", err)
	}
}

func TestRunImport_FetchError(t *testing.T) {
	sourceName := "test-import-fetcherr"
	defer cleanupRegistry(sourceName)

	Register(&mockSource{
		name:     sourceName,
		fetchErr: fmt.Errorf("API rate limit exceeded"),
	})

	dir := t.TempDir()

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  filepath.Join(dir, "tasks"),
		ScanDir:    dir,
	}

	_, err := RunImport(cfg)
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
	if !strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Errorf("expected fetch error message, got: %v", err)
	}
}

func TestRunImport_EmptyResult(t *testing.T) {
	sourceName := "test-import-empty"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{})

	dir := t.TempDir()

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  filepath.Join(dir, "tasks"),
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 0 {
		t.Errorf("expected 0 created, got %d", len(result.Created))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

func TestRunImport_JiraStatusMapping(t *testing.T) {
	// Register mock under "jira" to trigger Jira-specific field mappings
	defer cleanupRegistry("jira")

	Register(&mockSource{
		name: "jira",
		tasks: []ExternalTask{
			{ExternalID: "PROJ-1", Title: "Todo task", Status: "To Do"},
			{ExternalID: "PROJ-2", Title: "In progress task", Status: "In Progress"},
			{ExternalID: "PROJ-3", Title: "Done task", Status: "Done"},
		},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: "jira",
		SourceCfg:  SourceConfig{Name: "jira"},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 3 {
		t.Fatalf("expected 3 created, got %d", len(result.Created))
	}

	statusTests := []struct {
		file     string
		expected string
	}{
		{result.Created[0].FilePath, "status: pending"},
		{result.Created[1].FilePath, "status: in-progress"},
		{result.Created[2].FilePath, "status: completed"},
	}

	for _, st := range statusTests {
		content, err := os.ReadFile(st.file)
		if err != nil {
			t.Fatalf("failed to read %s: %v", st.file, err)
		}
		if !strings.Contains(string(content), st.expected) {
			t.Errorf("expected %q in %s, got:\n%s", st.expected, st.file, string(content))
		}
	}
}

func TestRunImport_JiraPriorityMapping(t *testing.T) {
	// Register mock under "jira" to trigger Jira-specific field mappings
	defer cleanupRegistry("jira")

	Register(&mockSource{
		name: "jira",
		tasks: []ExternalTask{
			{ExternalID: "PROJ-1", Title: "Highest pri", Status: "To Do", Priority: "Highest"},
			{ExternalID: "PROJ-2", Title: "High pri", Status: "To Do", Priority: "High"},
			{ExternalID: "PROJ-3", Title: "Medium pri", Status: "To Do", Priority: "Medium"},
			{ExternalID: "PROJ-4", Title: "Low pri", Status: "To Do", Priority: "Low"},
			{ExternalID: "PROJ-5", Title: "Lowest pri", Status: "To Do", Priority: "Lowest"},
		},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: "jira",
		SourceCfg:  SourceConfig{Name: "jira"},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Created) != 5 {
		t.Fatalf("expected 5 created, got %d", len(result.Created))
	}

	priorityTests := []struct {
		file     string
		expected string
	}{
		{result.Created[0].FilePath, "priority: critical"},
		{result.Created[1].FilePath, "priority: high"},
		{result.Created[2].FilePath, "priority: medium"},
		{result.Created[3].FilePath, "priority: low"},
		{result.Created[4].FilePath, "priority: low"},
	}

	for _, pt := range priorityTests {
		content, err := os.ReadFile(pt.file)
		if err != nil {
			t.Fatalf("failed to read %s: %v", pt.file, err)
		}
		if !strings.Contains(string(content), pt.expected) {
			t.Errorf("expected %q in %s, got:\n%s", pt.expected, pt.file, string(content))
		}
	}
}

func TestRunImport_NonJiraUsesDefaultFieldMap(t *testing.T) {
	sourceName := "test-import-default-map"
	defer cleanupRegistry(sourceName)

	setupMockSource(sourceName, []ExternalTask{
		{ExternalID: "GH-1", Title: "Open issue", Status: "open", Priority: "high"},
	})

	dir := t.TempDir()
	outputDir := filepath.Join(dir, "tasks")

	cfg := ImportConfig{
		SourceName: sourceName,
		SourceCfg:  SourceConfig{Name: sourceName},
		OutputDir:  outputDir,
		ScanDir:    dir,
	}

	result, err := RunImport(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(result.Created[0].FilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	s := string(content)
	if !strings.Contains(s, "status: pending") {
		t.Errorf("expected 'open' mapped to 'pending', got:\n%s", s)
	}
	if !strings.Contains(s, "priority: high") {
		t.Errorf("expected priority 'high', got:\n%s", s)
	}
}

func TestAppendSourceURL(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		url      string
		expected string
	}{
		{
			name:     "with description and URL",
			desc:     "Some description",
			url:      "https://example.com",
			expected: "Some description\n\n---\nSource: https://example.com",
		},
		{
			name:     "empty URL",
			desc:     "Description",
			url:      "",
			expected: "Description",
		},
		{
			name:     "empty description with URL",
			desc:     "",
			url:      "https://example.com",
			expected: "\n\n---\nSource: https://example.com",
		},
		{
			name:     "description with trailing newlines",
			desc:     "Body text\n\n",
			url:      "https://example.com",
			expected: "Body text\n\n---\nSource: https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendSourceURL(tt.desc, tt.url)
			if got != tt.expected {
				t.Errorf("appendSourceURL(%q, %q) = %q, want %q", tt.desc, tt.url, got, tt.expected)
			}
		})
	}
}
