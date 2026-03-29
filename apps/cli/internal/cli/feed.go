package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/sdk/go/feed"
)

var (
	feedFormat string
	feedLimit  int
	feedSince  string
	feedScope  string
	feedSource string
)

// gitLogFunc is the function used to run git log.
// Override in tests to avoid running actual git commands.
var gitLogFunc = runGitLog

// gitShowFunc is the function used to run git show.
// Override in tests to avoid running actual git commands.
var gitShowFunc = runGitShow

var feedCmd = &cobra.Command{
	Use:        "feed",
	SuggestFor: []string{"activity", "log", "history"},
	Short:      "Show a chronological activity feed of task changes",
	Long: `Show a chronological activity feed of recent changes to task files.

Uses git log to detect task creation, modification, and renames,
presenting them as a time-ordered feed.

Examples:
  taskmd feed
  taskmd feed --since 7d
  taskmd feed --limit 10
  taskmd feed --scope cli
  taskmd feed --format json
  taskmd feed --source worklog
  taskmd feed --source git`,
	Args: cobra.NoArgs,
	RunE: runFeed,
}

func init() {
	rootCmd.AddCommand(feedCmd)

	feedCmd.Flags().StringVar(&feedFormat, "format", "text", "output format (text, json)")
	feedCmd.Flags().IntVar(&feedLimit, "limit", 20, "maximum number of entries to show")
	feedCmd.Flags().StringVar(&feedSince, "since", "", "show changes since (e.g. 2d, 1w, 2026-02-28)")
	feedCmd.Flags().StringVar(&feedScope, "scope", "", "filter to a tasks subdirectory; supports wildcards (e.g. cli, cli*)")
	feedCmd.Flags().StringVar(&feedSource, "source", "all", "filter by event source (all, git, worklog)")
}

func runFeed(_ *cobra.Command, _ []string) error {
	if err := ValidateFormat(feedFormat, []string{"text", "json"}); err != nil {
		return err
	}

	validSources := map[string]bool{"all": true, "git": true, "worklog": true}
	if !validSources[feedSource] {
		return fmt.Errorf("unsupported source: %q (supported: all, git, worklog)", feedSource)
	}

	flags := GetGlobalFlags()

	entries, err := feed.Query(feed.Options{
		TasksDir:  flags.TaskDir,
		Limit:     feedLimit,
		Since:     feedSince,
		Scope:     feedScope,
		Source:    feedSource,
		Verbose:   flags.Verbose,
		GitLogFn:  gitLogFunc,
		GitShowFn: gitShowFunc,
	})
	if err != nil {
		return fmt.Errorf("failed to read git history (is this a git repository?): %w", err)
	}

	if len(entries) == 0 {
		if feedFormat == "text" {
			fmt.Println("No recent task changes.")
		} else {
			fmt.Print("[]\n")
		}
		return nil
	}

	switch feedFormat {
	case "json":
		return WriteJSON(os.Stdout, entries)
	default:
		return writeFeedText(entries)
	}
}

func runGitLog(_ string, args []string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func runGitShow(hash, path string) (string, error) {
	cmd := exec.Command("git", "show", hash+":"+path)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeFeedText(entries []feed.FeedEntry) error {
	r := getRenderer()

	fmt.Println(formatDim("Recent task activity", r))
	fmt.Println()

	for i, entry := range entries {
		if i > 0 {
			fmt.Println()
		}

		if entry.Source == "worklog" {
			writeWorklogEntryText(entry, r)
			continue
		}

		date := formatDim(entry.Timestamp.Format("2006-01-02 15:04"), r)
		author := formatLabel(entry.Author, r)
		fmt.Printf("%s %s: %s\n", date, author, entry.Message)

		for _, f := range entry.Files {
			writeFileChangeText(f, r)
		}
	}

	return nil
}

func writeWorklogEntryText(entry feed.FeedEntry, r *lipgloss.Renderer) {
	date := formatDim(entry.Timestamp.Format("2006-01-02 15:04"), r)
	taskRef := ""
	if entry.TaskID != "" {
		taskRef = fmt.Sprintf(" (%s)", formatTaskID(entry.TaskID, r))
	}
	fmt.Printf("%s [Worklog]%s %s\n", date, taskRef, entry.Message)
}

func writeFileChangeText(f feed.FileChange, r *lipgloss.Renderer) {
	statusTag := fileStatusTag(f)
	taskRef := ""
	if f.TaskID != "" {
		taskRef = fmt.Sprintf(" (%s)", formatTaskID(f.TaskID, r))
	}

	summary := formatChangeSummary(f)
	if summary != "" {
		fmt.Printf("  %s %s%s: %s\n", statusTag, f.Path, taskRef, summary)
	} else {
		fmt.Printf("  %s %s%s\n", statusTag, f.Path, taskRef)
	}
}

// formatChangeSummary builds a compact one-line summary of field and subtask changes.
func formatChangeSummary(f feed.FileChange) string {
	var parts []string
	for _, fc := range f.FieldChanges {
		parts = append(parts, fmt.Sprintf("%s %s \u2192 %s", fc.Field, fc.OldValue, fc.NewValue))
	}
	done := 0
	undone := 0
	for _, sc := range f.SubtaskChanges {
		if sc.Done {
			done++
		} else {
			undone++
		}
	}
	if done > 0 {
		parts = append(parts, fmt.Sprintf("%d subtask(s) completed", done))
	}
	if undone > 0 {
		parts = append(parts, fmt.Sprintf("%d subtask(s) unchecked", undone))
	}
	return strings.Join(parts, ", ")
}

func fileStatusTag(fc feed.FileChange) string {
	if fc.TaskStatus == "completed" {
		return "[Completed]"
	}
	if fc.TaskStatus == "cancelled" {
		return "[Cancelled]"
	}
	switch fc.Status {
	case "created":
		return "[Added]"
	case "modified":
		return "[Modified]"
	case "renamed":
		return "[Renamed]"
	default:
		return "[?]"
	}
}
