//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTaskWithTouches creates a task file with a touches field.
func writeTaskWithTouches(t *testing.T, dir, filename, id, title, status, priority string, touches []string) {
	t.Helper()

	touchesYAML := "[]"
	if len(touches) > 0 {
		touchesYAML = "["
		for i, s := range touches {
			if i > 0 {
				touchesYAML += ", "
			}
			touchesYAML += fmt.Sprintf("%q", s)
		}
		touchesYAML += "]"
	}

	content := fmt.Sprintf(`---
id: %q
title: %q
status: %s
priority: %s
effort: small
touches: %s
created: 2026-01-01
---

# %s
`, id, title, status, priority, touchesYAML, title)

	path := filepath.Join(dir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write task file %s: %v", path, err)
	}
}

type nextRec struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func TestNext_ScopeFlag(t *testing.T) {
	dir := setupTaskDir(t)

	writeTaskWithTouches(t, dir, "001-web.md", "001", "Web feature", "pending", "high", []string{"web", "api"})
	writeTaskWithTouches(t, dir, "002-cli.md", "002", "CLI feature", "pending", "medium", []string{"cli"})
	writeTaskWithTouches(t, dir, "003-web.md", "003", "Web styling", "pending", "low", []string{"web"})
	writeTask(t, dir, "004-generic.md", "004", "No scope task", "pending", nil)

	t.Run("scope filters to matching tasks", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "web", "--format", "json")

		var recs []nextRec
		if err := json.Unmarshal([]byte(result.Stdout), &recs); err != nil {
			t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, result.Stdout)
		}

		if len(recs) != 2 {
			t.Fatalf("Expected 2 tasks for scope 'web', got %d", len(recs))
		}

		ids := map[string]bool{}
		for _, r := range recs {
			ids[r.ID] = true
		}
		if !ids["001"] || !ids["003"] {
			t.Errorf("Expected tasks 001 and 003, got %v", ids)
		}
	})

	t.Run("scope with no matches shows message", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "nonexistent")

		if !strings.Contains(result.Stdout, `No actionable tasks found for scope "nonexistent"`) {
			t.Errorf("Expected scope-specific message, got: %s", result.Stdout)
		}
	})

	t.Run("without scope returns all tasks", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--format", "json", "--limit", "10")

		var recs []nextRec
		if err := json.Unmarshal([]byte(result.Stdout), &recs); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if len(recs) != 4 {
			t.Errorf("Without scope, expected 4 tasks, got %d", len(recs))
		}
	})

	t.Run("scope combined with filter", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "web", "--filter", "priority=high", "--format", "json")

		var recs []nextRec
		if err := json.Unmarshal([]byte(result.Stdout), &recs); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if len(recs) != 1 || recs[0].ID != "001" {
			t.Errorf("Expected only task 001 (web + high), got %v", recs)
		}
	})

	t.Run("scope shows label in table output", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "web")

		if !strings.Contains(result.Stdout, "scope: web") {
			t.Errorf("Expected scope label in table output, got: %s", result.Stdout)
		}
	})
}

func TestNext_ScopeExpansion(t *testing.T) {
	dir := setupTaskDir(t)

	// 001: no touches, blocks 002
	// 002: touches web, depends on 001 (blocked)
	// 003: touches cli, unrelated
	// 004: touches web, no deps (actionable)
	writeTask(t, dir, "001-db.md", "001", "Setup database", "pending", nil)
	writeTaskWithTouches(t, dir, "002-web.md", "002", "Web dashboard", "pending", "medium", []string{"web"})
	// Manually write 002 with dependency since writeTaskWithTouches doesn't support deps
	content002 := `---
id: "002"
title: "Web dashboard"
status: pending
priority: medium
effort: small
touches: ["web"]
dependencies: ["001"]
created: 2026-01-01
---

# Web dashboard
`
	if err := os.WriteFile(filepath.Join(dir, "002-web.md"), []byte(content002), 0o644); err != nil {
		t.Fatalf("failed to write 002: %v", err)
	}
	writeTaskWithTouches(t, dir, "003-cli.md", "003", "CLI feature", "pending", "low", []string{"cli"})
	writeTaskWithTouches(t, dir, "004-web.md", "004", "Web styling", "pending", "low", []string{"web"})

	t.Run("scope expands to include blocking dependencies", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "web", "--format", "json", "--limit", "10")

		var recs []nextRec
		if err := json.Unmarshal([]byte(result.Stdout), &recs); err != nil {
			t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, result.Stdout)
		}

		ids := map[string]bool{}
		for _, r := range recs {
			ids[r.ID] = true
		}

		if !ids["001"] {
			t.Errorf("Expected task 001 (blocks web task) to be included, got %v", ids)
		}
		if !ids["004"] {
			t.Errorf("Expected task 004 (directly touches web) to be included, got %v", ids)
		}
		if ids["003"] {
			t.Errorf("Task 003 (cli scope) should not appear in web results")
		}
	})

	t.Run("exact skips dependency expansion", func(t *testing.T) {
		result := mustRun(t, dir, "next", "--scope", "web", "--exact", "--format", "json", "--limit", "10")

		var recs []nextRec
		if err := json.Unmarshal([]byte(result.Stdout), &recs); err != nil {
			t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, result.Stdout)
		}

		ids := map[string]bool{}
		for _, r := range recs {
			ids[r.ID] = true
		}

		if ids["001"] {
			t.Errorf("With --exact, task 001 (no touches) should not appear")
		}
		// Only task 004 should appear (touches web, actionable)
		// Task 002 touches web but is blocked
		if len(recs) != 1 || recs[0].ID != "004" {
			t.Errorf("Expected only task 004, got %v", ids)
		}
	})
}
