//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// --- Helpers ---

// mustParseJSON unmarshals JSON data into v or fails the test.
func mustParseJSON(t *testing.T, data string, v any) {
	t.Helper()
	if err := json.Unmarshal([]byte(data), v); err != nil {
		t.Fatalf("failed to parse JSON: %v\nraw:\n%s", err, data)
	}
}

// mustParseYAML unmarshals YAML data into v or fails the test.
func mustParseYAML(t *testing.T, data string, v any) {
	t.Helper()
	if err := yaml.Unmarshal([]byte(data), v); err != nil {
		t.Fatalf("failed to parse YAML: %v\nraw:\n%s", err, data)
	}
}

// writeTaskFull creates a task file with all fields configurable.
func writeTaskFull(t *testing.T, dir, filename, id, title, status, priority string, deps []string) {
	t.Helper()

	depsYAML := "[]"
	if len(deps) > 0 {
		parts := make([]string, len(deps))
		for i, d := range deps {
			parts[i] = fmt.Sprintf("%q", d)
		}
		depsYAML = "[" + strings.Join(parts, ", ") + "]"
	}

	content := fmt.Sprintf(`---
id: %q
title: %q
status: %s
priority: %s
effort: small
dependencies: %s
tags: ["e2e"]
created: 2026-01-01
---

# %s

Test task for e2e tests.
`, id, title, status, priority, depsYAML, title)

	path := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write task file %s: %v", path, err)
	}
}

// writeTaskWithBody creates a task file with a custom body for search tests.
func writeTaskWithBody(t *testing.T, dir, filename, id, title, status, priority, body string) {
	t.Helper()

	content := fmt.Sprintf(`---
id: %q
title: %q
status: %s
priority: %s
effort: small
dependencies: []
tags: ["e2e"]
created: 2026-01-01
---

# %s

%s
`, id, title, status, priority, title, body)

	path := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write task file %s: %v", path, err)
	}
}

// --- Workflow Tests ---

func TestWorkflow_AddThenList(t *testing.T) {
	dir := setupTaskDir(t)

	// Add a task via the CLI.
	addResult := mustRun(t, dir, "add", "My Workflow Task", "--format", "json")

	var added struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	mustParseJSON(t, addResult.Stdout, &added)

	if added.ID == "" {
		t.Fatal("expected add to return a non-empty ID")
	}
	if added.Title != "My Workflow Task" {
		t.Errorf("expected title 'My Workflow Task', got %q", added.Title)
	}

	// List tasks and verify the new task appears.
	listResult := mustRun(t, dir, "list", "--format", "json")

	var tasks []map[string]any
	mustParseJSON(t, listResult.Stdout, &tasks)

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in list, got %d", len(tasks))
	}
	if tasks[0]["id"] != added.ID {
		t.Errorf("expected list task ID %q, got %q", added.ID, tasks[0]["id"])
	}
	if tasks[0]["title"] != "My Workflow Task" {
		t.Errorf("expected list task title 'My Workflow Task', got %q", tasks[0]["title"])
	}
}

func TestWorkflow_AddSetList(t *testing.T) {
	dir := setupTaskDir(t)

	// Add a task.
	addResult := mustRun(t, dir, "add", "Status Change Task", "--format", "json")

	var added struct {
		ID string `json:"id"`
	}
	mustParseJSON(t, addResult.Stdout, &added)

	// Change its status.
	mustRun(t, dir, "set", added.ID, "--status", "in-progress")

	// Verify the status changed in list output.
	listResult := mustRun(t, dir, "list", "--format", "json")

	var tasks []map[string]any
	mustParseJSON(t, listResult.Stdout, &tasks)

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0]["status"] != "in-progress" {
		t.Errorf("expected status 'in-progress', got %q", tasks[0]["status"])
	}
}

func TestWorkflow_AddThenGet(t *testing.T) {
	dir := setupTaskDir(t)

	// Add a task with specific fields.
	addResult := mustRun(t, dir, "add", "Detailed Task", "--priority", "high", "--tags", "cli,test", "--format", "json")

	var added struct {
		ID       string `json:"id"`
		Priority string `json:"priority"`
	}
	mustParseJSON(t, addResult.Stdout, &added)

	if added.Priority != "high" {
		t.Errorf("expected add priority 'high', got %q", added.Priority)
	}

	// Get the task by ID.
	getResult := mustRun(t, dir, "get", added.ID, "--format", "json")

	var got map[string]any
	mustParseJSON(t, getResult.Stdout, &got)

	if got["id"] != added.ID {
		t.Errorf("expected get ID %q, got %q", added.ID, got["id"])
	}
	if got["priority"] != "high" {
		t.Errorf("expected get priority 'high', got %q", got["priority"])
	}

	// Verify tags are present.
	tags, ok := got["tags"].([]any)
	if !ok {
		t.Fatalf("expected tags to be an array, got %T", got["tags"])
	}
	tagStrs := make([]string, len(tags))
	for i, tag := range tags {
		tagStrs[i] = fmt.Sprintf("%v", tag)
	}
	joined := strings.Join(tagStrs, ",")
	if !strings.Contains(joined, "cli") || !strings.Contains(joined, "test") {
		t.Errorf("expected tags to contain 'cli' and 'test', got %v", tagStrs)
	}
}

func TestWorkflow_AddNextWithDeps(t *testing.T) {
	dir := setupTaskDir(t)

	// Create two tasks: 001 depends on 002, 002 has no deps.
	writeTask(t, dir, "001-blocked.md", "001", "Blocked Task", "pending", []string{"002"})
	writeTask(t, dir, "002-ready.md", "002", "Ready Task", "pending", nil)

	// Next should recommend the unblocked task (002).
	nextResult := mustRun(t, dir, "next", "--format", "json")

	var recs []map[string]any
	mustParseJSON(t, nextResult.Stdout, &recs)

	if len(recs) == 0 {
		t.Fatal("expected at least one recommendation from next")
	}

	// 002 should appear (it's unblocked). 001 should not (it's blocked).
	ids := make([]string, len(recs))
	for i, rec := range recs {
		ids[i] = fmt.Sprintf("%v", rec["id"])
	}

	found002 := false
	for _, id := range ids {
		if id == "002" {
			found002 = true
		}
		if id == "001" {
			t.Error("expected blocked task 001 to NOT appear in next recommendations")
		}
	}
	if !found002 {
		t.Errorf("expected unblocked task 002 in next recommendations, got %v", ids)
	}
}

func TestWorkflow_GraphWithDependencies(t *testing.T) {
	dir := setupTaskDir(t)

	// Create a dependency chain: 001 -> 002 -> 003.
	writeTask(t, dir, "001-top.md", "001", "Top Task", "pending", []string{"002"})
	writeTask(t, dir, "002-mid.md", "002", "Mid Task", "pending", []string{"003"})
	writeTask(t, dir, "003-base.md", "003", "Base Task", "pending", nil)

	result := mustRun(t, dir, "graph", "--format", "json", "--all")

	var graphData struct {
		Nodes []map[string]any    `json:"nodes"`
		Edges []map[string]string `json:"edges"`
	}
	mustParseJSON(t, result.Stdout, &graphData)

	// Should have 3 nodes.
	if len(graphData.Nodes) != 3 {
		t.Errorf("expected 3 graph nodes, got %d", len(graphData.Nodes))
	}

	// Should have 2 edges (001->002, 002->003).
	if len(graphData.Edges) != 2 {
		t.Errorf("expected 2 graph edges, got %d", len(graphData.Edges))
	}

	// Verify edge structure. Edges go from dependency -> dependent
	// (i.e., "001 depends on 002" produces edge from:002 to:001).
	edgeSet := make(map[string]bool)
	for _, edge := range graphData.Edges {
		key := edge["from"] + "->" + edge["to"]
		edgeSet[key] = true
	}
	if !edgeSet["002->001"] {
		t.Errorf("expected edge 002->001 (001 depends on 002), got edges: %v", graphData.Edges)
	}
	if !edgeSet["003->002"] {
		t.Errorf("expected edge 003->002 (002 depends on 003), got edges: %v", graphData.Edges)
	}
}

func TestWorkflow_BoardGroupsByStatus(t *testing.T) {
	dir := setupTaskDir(t)

	writeTask(t, dir, "001-pending.md", "001", "Pending Task", "pending", nil)
	writeTask(t, dir, "002-progress.md", "002", "Progress Task", "in-progress", nil)
	writeTask(t, dir, "003-done.md", "003", "Done Task", "completed", nil)

	result := mustRun(t, dir, "board", "--format", "json")

	var groups []struct {
		Group string           `json:"group"`
		Count int              `json:"count"`
		Tasks []map[string]any `json:"tasks"`
	}
	mustParseJSON(t, result.Stdout, &groups)

	// Build a map of group -> task IDs for easy assertion.
	groupMap := make(map[string][]string)
	for _, g := range groups {
		for _, task := range g.Tasks {
			groupMap[g.Group] = append(groupMap[g.Group], fmt.Sprintf("%v", task["id"]))
		}
	}

	assertContains := func(group, id string) {
		t.Helper()
		for _, taskID := range groupMap[group] {
			if taskID == id {
				return
			}
		}
		t.Errorf("expected task %s in group %q, got %v", id, group, groupMap[group])
	}

	assertContains("pending", "001")
	assertContains("in-progress", "002")
	assertContains("completed", "003")
}

func TestWorkflow_ValidateWellFormed(t *testing.T) {
	dir := setupTaskDir(t)

	writeTask(t, dir, "001-alpha.md", "001", "Alpha Task", "pending", nil)
	writeTask(t, dir, "002-beta.md", "002", "Beta Task", "in-progress", nil)

	result := mustRun(t, dir, "validate", "--format", "json")

	var validation struct {
		Errors    int   `json:"errors"`
		Warnings  int   `json:"warnings"`
		TaskCount int   `json:"task_count"`
		Issues    []any `json:"issues"`
	}
	mustParseJSON(t, result.Stdout, &validation)

	if validation.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", validation.Errors)
	}
	if validation.TaskCount != 2 {
		t.Errorf("expected task_count 2, got %d", validation.TaskCount)
	}
}

func TestWorkflow_JSONOutputParseable(t *testing.T) {
	dir := setupTaskDir(t)

	writeTask(t, dir, "001-json.md", "001", "JSON Test Task", "pending", nil)
	writeTask(t, dir, "002-json.md", "002", "JSON Test Task 2", "in-progress", []string{"001"})

	// Each subtest parses JSON output from a different command.
	t.Run("list", func(t *testing.T) {
		result := mustRun(t, dir, "list", "--format", "json")
		var v []any
		mustParseJSON(t, result.Stdout, &v)
		if len(v) != 2 {
			t.Errorf("expected 2 items, got %d", len(v))
		}
	})

	t.Run("get", func(t *testing.T) {
		result := mustRun(t, dir, "get", "001", "--format", "json")
		var v map[string]any
		mustParseJSON(t, result.Stdout, &v)
		if v["id"] != "001" {
			t.Errorf("expected id '001', got %q", v["id"])
		}
	})

	t.Run("next", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--format", "json")
		var v []any
		mustParseJSON(t, result.Stdout, &v)
	})

	t.Run("graph", func(t *testing.T) {
		result := mustRun(t, dir, "graph", "--format", "json", "--all")
		var v map[string]any
		mustParseJSON(t, result.Stdout, &v)
		if _, ok := v["nodes"]; !ok {
			t.Error("expected 'nodes' key in graph JSON")
		}
		if _, ok := v["edges"]; !ok {
			t.Error("expected 'edges' key in graph JSON")
		}
	})

	t.Run("board", func(t *testing.T) {
		result := mustRun(t, dir, "board", "--format", "json")
		var v []any
		mustParseJSON(t, result.Stdout, &v)
	})

	t.Run("validate", func(t *testing.T) {
		result := mustRun(t, dir, "validate", "--format", "json")
		var v map[string]any
		mustParseJSON(t, result.Stdout, &v)
		if _, ok := v["task_count"]; !ok {
			t.Error("expected 'task_count' key in validate JSON")
		}
	})
}

func TestWorkflow_YAMLOutputParseable(t *testing.T) {
	dir := setupTaskDir(t)

	writeTask(t, dir, "001-yaml.md", "001", "YAML Test Task", "pending", nil)

	t.Run("list", func(t *testing.T) {
		result := mustRun(t, dir, "list", "--format", "yaml")
		var v []any
		mustParseYAML(t, result.Stdout, &v)
		if len(v) != 1 {
			t.Errorf("expected 1 item, got %d", len(v))
		}
	})

	t.Run("get", func(t *testing.T) {
		result := mustRun(t, dir, "get", "001", "--format", "yaml")
		var v map[string]any
		mustParseYAML(t, result.Stdout, &v)
		if v["id"] != "001" {
			t.Errorf("expected id '001', got %q", v["id"])
		}
	})

	t.Run("next", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--format", "yaml")
		var v []any
		mustParseYAML(t, result.Stdout, &v)
	})
}

func TestWorkflow_FlagWiring(t *testing.T) {
	dir := setupTaskDir(t)

	// Seed tasks with varying statuses and priorities.
	writeTaskFull(t, dir, "001-pending-low.md", "001", "Low Pending", "pending", "low", nil)
	writeTaskFull(t, dir, "002-pending-crit.md", "002", "Critical Pending", "pending", "critical", nil)
	writeTaskFull(t, dir, "003-progress-med.md", "003", "Medium Progress", "in-progress", "medium", nil)
	writeTaskFull(t, dir, "004-done-high.md", "004", "High Done", "completed", "high", nil)

	t.Run("list_filter_status", func(t *testing.T) {
		result := mustRun(t, dir, "list", "--filter", "status=pending", "--format", "json")

		var tasks []map[string]any
		mustParseJSON(t, result.Stdout, &tasks)

		for _, task := range tasks {
			if task["status"] != "pending" {
				t.Errorf("expected only pending tasks, got status %q for task %v", task["status"], task["id"])
			}
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 pending tasks, got %d", len(tasks))
		}
	})

	t.Run("list_sort_priority", func(t *testing.T) {
		result := mustRun(t, dir, "list", "--sort", "priority", "--format", "json")

		var tasks []map[string]any
		mustParseJSON(t, result.Stdout, &tasks)

		if len(tasks) < 2 {
			t.Fatalf("expected at least 2 tasks, got %d", len(tasks))
		}
		// Critical should come first when sorted by priority.
		if tasks[0]["priority"] != "critical" {
			t.Errorf("expected first task sorted by priority to be 'critical', got %q", tasks[0]["priority"])
		}
	})

	t.Run("add_priority_critical", func(t *testing.T) {
		addDir := setupTaskDir(t)
		result := mustRun(t, addDir, "add", "Critical Task", "--priority", "critical", "--format", "json")

		var added struct {
			Priority string `json:"priority"`
		}
		mustParseJSON(t, result.Stdout, &added)

		if added.Priority != "critical" {
			t.Errorf("expected priority 'critical', got %q", added.Priority)
		}
	})

	t.Run("add_status_in_progress", func(t *testing.T) {
		addDir := setupTaskDir(t)
		result := mustRun(t, addDir, "add", "In Progress Task", "--status", "in-progress", "--format", "json")

		var added struct {
			Status string `json:"status"`
		}
		mustParseJSON(t, result.Stdout, &added)

		if added.Status != "in-progress" {
			t.Errorf("expected status 'in-progress', got %q", added.Status)
		}
	})

	t.Run("next_limit", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--limit", "1", "--format", "json")

		var recs []any
		mustParseJSON(t, result.Stdout, &recs)

		if len(recs) > 1 {
			t.Errorf("expected at most 1 recommendation with --limit 1, got %d", len(recs))
		}
	})

	t.Run("graph_exclude_status", func(t *testing.T) {
		result := mustRun(t, dir, "graph", "--exclude-status", "completed", "--format", "json")

		var graphData struct {
			Nodes []map[string]any `json:"nodes"`
		}
		mustParseJSON(t, result.Stdout, &graphData)

		for _, node := range graphData.Nodes {
			if node["status"] == "completed" {
				t.Errorf("expected no completed tasks in graph with --exclude-status completed, found node %v", node["id"])
			}
		}
	})
}

func TestWorkflow_SearchFlags(t *testing.T) {
	dir := setupTaskDir(t)

	// Create tasks with unique keywords in body for search tests.
	writeTaskWithBody(t, dir, "001-auth-high.md", "001", "Auth login feature", "pending", "high",
		"Implement JWT authentication for the login endpoint.")
	writeTaskWithBody(t, dir, "002-auth-low.md", "002", "Auth docs update", "pending", "low",
		"Document the authentication API endpoints.")
	writeTaskWithBody(t, dir, "003-deploy.md", "003", "Deploy pipeline", "in-progress", "medium",
		"Set up CI/CD pipeline with authentication checks.")
	writeTaskWithBody(t, dir, "004-auth-done.md", "004", "Auth refactor", "completed", "critical",
		"Refactored the authentication module for performance.")

	t.Run("search_filter_status", func(t *testing.T) {
		result := mustRun(t, dir, "search", "authentication", "--filter", "status=pending", "--format", "json")

		var results []map[string]any
		mustParseJSON(t, result.Stdout, &results)

		for _, r := range results {
			if r["status"] != "pending" {
				t.Errorf("expected only pending results, got status %q for task %v", r["status"], r["id"])
			}
		}
		if len(results) != 2 {
			t.Errorf("expected 2 pending authentication results, got %d", len(results))
		}
	})

	t.Run("search_sort_priority", func(t *testing.T) {
		result := mustRun(t, dir, "search", "authentication", "--sort", "priority", "--format", "json")

		var results []map[string]any
		mustParseJSON(t, result.Stdout, &results)

		if len(results) < 2 {
			t.Fatalf("expected at least 2 results, got %d", len(results))
		}
		// Critical should come first when sorted by priority.
		if results[0]["priority"] != "critical" {
			t.Errorf("expected first result to be critical priority, got %q", results[0]["priority"])
		}
	})

	t.Run("search_limit", func(t *testing.T) {
		result := mustRun(t, dir, "search", "authentication", "--limit", "1", "--format", "json")

		var results []map[string]any
		mustParseJSON(t, result.Stdout, &results)

		if len(results) != 1 {
			t.Errorf("expected 1 result with --limit 1, got %d", len(results))
		}
	})

	t.Run("search_combined", func(t *testing.T) {
		result := mustRun(t, dir, "search", "authentication", "--filter", "status=pending", "--sort", "priority", "--limit", "2", "--format", "json")

		var results []map[string]any
		mustParseJSON(t, result.Stdout, &results)

		if len(results) > 2 {
			t.Errorf("expected at most 2 results with --limit 2, got %d", len(results))
		}
		for _, r := range results {
			if r["status"] != "pending" {
				t.Errorf("expected only pending results, got status %q", r["status"])
			}
		}
		// With sort by priority, high should come before low.
		if len(results) == 2 && results[0]["priority"] != "high" {
			t.Errorf("expected first result to be high priority, got %q", results[0]["priority"])
		}
	})
}
