package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

const sampleGitLogOutput = `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
Alice
2026-02-28 10:30:00 +0000
chore: update task 042 status

M	tasks/cli/042-add-auth.md

bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
Bob
2026-02-27 14:00:00 +0000
feat: add new task 043

A	tasks/cli/043-new-feature.md
R100	tasks/old/010-rename-me.md	tasks/cli/010-renamed.md
`

func captureFeedOutput(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func resetFeedFlags() {
	feedFormat = "text"
	feedLimit = 20
	feedSince = ""
	feedScope = ""
}

// noopGitShow returns an error so enrichEntriesWithTaskStatus is a no-op.
func noopGitShow(_, _ string) (string, error) {
	return "", fmt.Errorf("not available")
}

func TestParseGitLogOutput(t *testing.T) {
	entries := parseGitLogOutput(sampleGitLogOutput)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// First entry
	e := entries[0]
	if e.Hash != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected hash: %s", e.Hash)
	}
	if e.Author != "Alice" {
		t.Errorf("unexpected author: %s", e.Author)
	}
	if e.Message != "chore: update task 042 status" {
		t.Errorf("unexpected message: %s", e.Message)
	}
	if len(e.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(e.Files))
	}
	if e.Files[0].Status != "modified" {
		t.Errorf("expected modified, got %s", e.Files[0].Status)
	}
	if e.Files[0].TaskID != "042" {
		t.Errorf("expected task ID 042, got %s", e.Files[0].TaskID)
	}

	// Second entry
	e2 := entries[1]
	if e2.Author != "Bob" {
		t.Errorf("unexpected author: %s", e2.Author)
	}
	if len(e2.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(e2.Files))
	}
	if e2.Files[0].Status != "created" {
		t.Errorf("expected created, got %s", e2.Files[0].Status)
	}
	if e2.Files[0].TaskID != "043" {
		t.Errorf("expected task ID 043, got %s", e2.Files[0].TaskID)
	}
	if e2.Files[1].Status != "renamed" {
		t.Errorf("expected renamed, got %s", e2.Files[1].Status)
	}
	if e2.Files[1].Path != "tasks/cli/010-renamed.md" {
		t.Errorf("expected renamed path, got %s", e2.Files[1].Path)
	}
}

func TestParseGitLogOutput_Empty(t *testing.T) {
	entries := parseGitLogOutput("")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}

	entries = parseGitLogOutput("   \n\n  ")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for whitespace, got %d", len(entries))
	}
}

func TestFeedCommand_PlainText(t *testing.T) {
	resetFeedFlags()
	noColor = true
	defer func() { noColor = false }()

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return sampleGitLogOutput, nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Recent task activity") {
		t.Error("expected header line in output")
	}
	if !strings.Contains(output, "2026-02-28 10:30") {
		t.Error("expected date in output")
	}
	if !strings.Contains(output, "Alice") {
		t.Error("expected author in output")
	}
	if !strings.Contains(output, "[Modified]") {
		t.Error("expected [Modified] status tag")
	}
	if !strings.Contains(output, "[Added] tasks/") {
		t.Error("expected [Added] file entry")
	}
	if !strings.Contains(output, "042") {
		t.Error("expected task ID 042 in output")
	}
}

func TestFeedCommand_JSON(t *testing.T) {
	resetFeedFlags()
	feedFormat = "json"

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return sampleGitLogOutput, nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entries []FeedEntry
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, output)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Author != "Alice" {
		t.Errorf("expected Alice, got %s", entries[0].Author)
	}
	if len(entries[1].Files) != 2 {
		t.Errorf("expected 2 files in second entry, got %d", len(entries[1].Files))
	}
}

func TestFeedCommand_Limit(t *testing.T) {
	resetFeedFlags()
	feedLimit = 5

	var capturedArgs []string
	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, args []string) (string, error) {
		capturedArgs = args
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	_, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, arg := range capturedArgs {
		if arg == "-5" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected -5 in git args, got: %v", capturedArgs)
	}
}

func TestFeedCommand_Since(t *testing.T) {
	resetFeedFlags()
	feedSince = "7d"

	var capturedArgs []string
	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, args []string) (string, error) {
		capturedArgs = args
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	_, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, arg := range capturedArgs {
		if arg == "--since=7.days.ago" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected --since=7.days.ago in git args, got: %v", capturedArgs)
	}
}

func TestFeedCommand_Scope(t *testing.T) {
	resetFeedFlags()
	feedScope = "cli"

	var capturedArgs []string
	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, args []string) (string, error) {
		capturedArgs = args
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	_, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The last arg should be the path glob containing the scope
	lastArg := capturedArgs[len(capturedArgs)-1]
	if !strings.Contains(lastArg, "cli") {
		t.Errorf("expected scope 'cli' in path glob, got: %s", lastArg)
	}
}

func TestFeedCommand_EmptyFeed(t *testing.T) {
	resetFeedFlags()

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "No recent task changes.") {
		t.Errorf("expected empty feed message, got: %s", output)
	}
}

func TestFeedCommand_EmptyFeed_JSON(t *testing.T) {
	resetFeedFlags()
	feedFormat = "json"

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.TrimSpace(output) != "[]" {
		t.Errorf("expected empty JSON array, got: %s", output)
	}
}

func TestFeedCommand_InvalidFormat(t *testing.T) {
	resetFeedFlags()
	feedFormat = "csv"

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return "", nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	err := runFeed(feedCmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %v", err)
	}
}

func TestFeedCommand_NonGitRepo(t *testing.T) {
	resetFeedFlags()

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return "", fmt.Errorf("not a git repository")
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = noopGitShow
	defer func() { gitShowFunc = oldGitShow }()

	err := runFeed(feedCmd, nil)
	if err == nil {
		t.Fatal("expected error for non-git repo")
	}
	if !strings.Contains(err.Error(), "git repository") {
		t.Errorf("expected git repository error, got: %v", err)
	}
}

func TestExtractTaskIDFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"tasks/cli/042-add-auth.md", "042"},
		{"tasks/043-new-feature.md", "043"},
		{"tasks/cli/01kjmg6sc-implement-feed.md", "01kjmg6sc"},
		{"README.md", ""},
	}

	for _, tt := range tests {
		got := extractTaskIDFromPath(tt.path)
		if got != tt.expected {
			t.Errorf("extractTaskIDFromPath(%q) = %q, want %q", tt.path, got, tt.expected)
		}
	}
}

func TestBuildGitLogArgs(t *testing.T) {
	args := buildGitLogArgs("tasks", 10, "7d", "cli")

	hasLimit := false
	hasSince := false
	hasScope := false
	for _, arg := range args {
		if arg == "-10" {
			hasLimit = true
		}
		if arg == "--since=7.days.ago" {
			hasSince = true
		}
		if strings.Contains(arg, "cli") {
			hasScope = true
		}
	}

	if !hasLimit {
		t.Error("expected -10 in args")
	}
	if !hasSince {
		t.Error("expected --since=7.days.ago in args")
	}
	if !hasScope {
		t.Error("expected cli scope in args")
	}
}

func TestNormalizeSince(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2d", "2.days.ago"},
		{"1w", "1.weeks.ago"},
		{"3m", "3.months.ago"},
		{"1y", "1.years.ago"},
		{"2026-02-28", "2026-02-28"},
		{"7 days ago", "7 days ago"},
		{"", ""},
		{"d", "d"},
	}

	for _, tt := range tests {
		got := normalizeSince(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeSince(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestExtractStatusFromContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "completed status",
			content:  "---\nid: 042\ntitle: Test\nstatus: completed\n---\n# Body",
			expected: "completed",
		},
		{
			name:     "cancelled status",
			content:  "---\nid: 043\ntitle: Test\nstatus: cancelled\n---\n# Body",
			expected: "cancelled",
		},
		{
			name:     "pending status",
			content:  "---\nid: 044\ntitle: Test\nstatus: pending\n---\n# Body",
			expected: "pending",
		},
		{
			name:     "no frontmatter",
			content:  "# Just a markdown file",
			expected: "",
		},
		{
			name:     "no status field",
			content:  "---\nid: 045\ntitle: Test\n---\n# Body",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractStatusFromContent(tt.content)
			if got != tt.expected {
				t.Errorf("extractStatusFromContent() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFeedCommand_CompletedStatus(t *testing.T) {
	resetFeedFlags()
	noColor = true
	defer func() { noColor = false }()

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return sampleGitLogOutput, nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = func(_, path string) (string, error) {
		if strings.Contains(path, "042") {
			return "---\nid: 042\ntitle: Add Auth\nstatus: completed\n---\n# Task", nil
		}
		return "", fmt.Errorf("not found")
	}
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "[Completed]") {
		t.Error("expected [Completed] tag for task 042")
	}
	if !strings.Contains(output, "[Added]") {
		t.Error("expected [Added] tag for task 043")
	}
}

func TestFeedCommand_CancelledStatus(t *testing.T) {
	resetFeedFlags()
	noColor = true
	defer func() { noColor = false }()

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return sampleGitLogOutput, nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = func(_, path string) (string, error) {
		if strings.Contains(path, "042") {
			return "---\nid: 042\ntitle: Add Auth\nstatus: cancelled\n---\n# Task", nil
		}
		return "", fmt.Errorf("not found")
	}
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "[Cancelled]") {
		t.Error("expected [Cancelled] tag for task 042")
	}
}

func TestFeedCommand_CompletedStatus_JSON(t *testing.T) {
	resetFeedFlags()
	feedFormat = "json"

	oldGitLog := gitLogFunc
	gitLogFunc = func(_ string, _ []string) (string, error) {
		return sampleGitLogOutput, nil
	}
	defer func() { gitLogFunc = oldGitLog }()

	oldGitShow := gitShowFunc
	gitShowFunc = func(_, path string) (string, error) {
		if strings.Contains(path, "042") {
			return "---\nid: 042\ntitle: Add Auth\nstatus: completed\n---\n# Task", nil
		}
		return "", fmt.Errorf("not found")
	}
	defer func() { gitShowFunc = oldGitShow }()

	output, err := captureFeedOutput(t, func() error {
		return runFeed(feedCmd, nil)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entries []FeedEntry
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entries[0].Files[0].TaskStatus != "completed" {
		t.Errorf("expected taskStatus 'completed', got %q", entries[0].Files[0].TaskStatus)
	}
}

func TestBuildGitLogArgs_NoOptionalFlags(t *testing.T) {
	args := buildGitLogArgs("tasks", 20, "", "")

	for _, arg := range args {
		if strings.HasPrefix(arg, "--since") {
			t.Error("did not expect --since when empty")
		}
	}

	// Last arg should be the path glob without a scope segment
	lastArg := args[len(args)-1]
	if strings.Contains(lastArg, "//") {
		t.Errorf("unexpected double slash in path: %s", lastArg)
	}
}
