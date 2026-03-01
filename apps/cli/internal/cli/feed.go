package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	feedFormat string
	feedLimit  int
	feedSince  string
	feedScope  string
)

// gitLogFunc is the function used to run git log.
// Override in tests to avoid running actual git commands.
var gitLogFunc = runGitLog

// FeedEntry represents a single commit in the activity feed.
type FeedEntry struct {
	Hash      string       `json:"hash"`
	Author    string       `json:"author"`
	Timestamp time.Time    `json:"timestamp"`
	Message   string       `json:"message"`
	Files     []FileChange `json:"files"`
}

// FileChange represents a file changed in a commit.
type FileChange struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	TaskID string `json:"taskID,omitempty"`
}

var taskIDFromFilenameRegex = regexp.MustCompile(`(?:^|/)(\w+)-`)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Show a chronological activity feed of task changes",
	Long: `Show a chronological activity feed of recent changes to task files.

Uses git log to detect task creation, modification, and renames,
presenting them as a time-ordered feed.

Examples:
  taskmd feed
  taskmd feed --since 7d
  taskmd feed --limit 10
  taskmd feed --scope cli
  taskmd feed --format json`,
	Args: cobra.NoArgs,
	RunE: runFeed,
}

func init() {
	rootCmd.AddCommand(feedCmd)

	feedCmd.Flags().StringVar(&feedFormat, "format", "text", "output format (text, json)")
	feedCmd.Flags().IntVar(&feedLimit, "limit", 20, "maximum number of commits to show")
	feedCmd.Flags().StringVar(&feedSince, "since", "", "show changes since (e.g. 2d, 1w, 2026-02-28)")
	feedCmd.Flags().StringVar(&feedScope, "scope", "", "filter to a tasks subdirectory (e.g. cli)")
}

func runFeed(_ *cobra.Command, _ []string) error {
	if err := ValidateFormat(feedFormat, []string{"text", "json"}); err != nil {
		return err
	}

	flags := GetGlobalFlags()
	tasksDir := flags.TaskDir

	args := buildGitLogArgs(tasksDir, feedLimit, feedSince, feedScope)

	output, err := gitLogFunc(tasksDir, args)
	if err != nil {
		return fmt.Errorf("failed to read git history (is this a git repository?): %w", err)
	}

	entries := parseGitLogOutput(output)

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

func buildGitLogArgs(tasksDir string, limit int, since, scope string) []string {
	args := []string{
		"log",
		"--format=%H%n%an%n%ai%n%s",
		"--name-status",
		"--diff-filter=ACMR",
		fmt.Sprintf("-%d", limit),
	}

	if since != "" {
		args = append(args, "--since="+normalizeSince(since))
	}

	pathGlob := filepath.Join(tasksDir, "**", "*.md")
	if scope != "" {
		pathGlob = filepath.Join(tasksDir, scope, "**", "*.md")
	}
	args = append(args, "--", pathGlob)

	return args
}

// normalizeSince converts shorthand durations like "2d" or "1w" into
// git-compatible relative date strings like "2.days.ago" or "1.weeks.ago".
// Values that don't match the shorthand pattern are returned as-is,
// allowing absolute dates to pass through.
func normalizeSince(s string) string {
	unitMap := map[byte]string{
		'd': "days",
		'w': "weeks",
		'm': "months",
		'y': "years",
	}

	if len(s) < 2 {
		return s
	}

	unit := s[len(s)-1]
	numPart := s[:len(s)-1]

	word, ok := unitMap[unit]
	if !ok {
		return s
	}

	// Verify the numeric part is all digits
	for _, c := range numPart {
		if c < '0' || c > '9' {
			return s
		}
	}

	return numPart + "." + word + ".ago"
}

func runGitLog(_ string, args []string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func parseGitLogOutput(output string) []FeedEntry {
	if strings.TrimSpace(output) == "" {
		return nil
	}

	var entries []FeedEntry
	var current *FeedEntry
	lineAfterHash := 0

	s := bufio.NewScanner(strings.NewReader(output))
	for s.Scan() {
		line := s.Text()

		if len(line) == 40 && isHexString(line) {
			if current != nil {
				entries = append(entries, *current)
			}
			current = &FeedEntry{Hash: line}
			lineAfterHash = 1
			continue
		}

		if current == nil || len(line) == 0 {
			continue
		}

		lineAfterHash = parseEntryLine(current, line, lineAfterHash)
	}

	if current != nil {
		entries = append(entries, *current)
	}

	return entries
}

// parseEntryLine parses a single line within a commit entry and returns the
// updated line position counter.
func parseEntryLine(entry *FeedEntry, line string, pos int) int {
	switch pos {
	case 1:
		entry.Author = line
		return 2
	case 2:
		t, err := time.Parse("2006-01-02 15:04:05 -0700", line)
		if err == nil {
			entry.Timestamp = t
		}
		return 3
	case 3:
		entry.Message = line
		return 4
	default:
		if fc := parseFileChangeLine(line); fc != nil {
			entry.Files = append(entry.Files, *fc)
		}
		return pos
	}
}

func parseFileChangeLine(line string) *FileChange {
	parts := strings.Split(line, "\t")
	if len(parts) < 2 {
		return nil
	}

	statusCode := parts[0]
	var path, status string

	switch {
	case statusCode == "A":
		status = "created"
		path = parts[1]
	case statusCode == "M":
		status = "modified"
		path = parts[1]
	case strings.HasPrefix(statusCode, "R"):
		status = "renamed"
		if len(parts) >= 3 {
			path = parts[2] // new path
		} else {
			path = parts[1]
		}
	default:
		return nil
	}

	taskID := extractTaskIDFromPath(path)

	return &FileChange{
		Path:   path,
		Status: status,
		TaskID: taskID,
	}
}

func extractTaskIDFromPath(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	match := taskIDFromFilenameRegex.FindStringSubmatch(base)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func writeFeedText(entries []FeedEntry) error {
	r := getRenderer()

	fmt.Println(formatDim("[A] added  [M] modified  [R] renamed", r))
	fmt.Println()

	for i, entry := range entries {
		if i > 0 {
			fmt.Println()
		}

		date := formatDim(entry.Timestamp.Format("2006-01-02 15:04"), r)
		author := formatLabel(entry.Author, r)
		fmt.Printf("%s %s: %s\n", date, author, entry.Message)

		for _, f := range entry.Files {
			statusTag := fileStatusTag(f.Status)
			line := fmt.Sprintf("  %s %s", statusTag, f.Path)
			if f.TaskID != "" {
				line = fmt.Sprintf("  %s %s (%s)", statusTag, f.Path, formatTaskID(f.TaskID, r))
			}
			fmt.Println(line)
		}
	}

	return nil
}

func fileStatusTag(status string) string {
	switch status {
	case "created":
		return "[A]"
	case "modified":
		return "[M]"
	case "renamed":
		return "[R]"
	default:
		return "[?]"
	}
}
