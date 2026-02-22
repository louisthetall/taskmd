package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/apps/cli/internal/model"
	"github.com/driangle/taskmd/apps/cli/internal/scanner"
	"github.com/driangle/taskmd/apps/cli/internal/taskfile"
)

var (
	dedupDryRun bool
	dedupFormat string
)

var deduplicateCmd = &cobra.Command{
	Use:   "deduplicate [path]",
	Short: "Detect and resolve duplicate task IDs",
	Long: `Deduplicate finds tasks with colliding IDs and reassigns new IDs to resolve conflicts.

When multiple contributors create tasks on separate branches, IDs can collide after merge.
This command detects duplicates and assigns new IDs to the newer tasks (by created date).

For each collision:
  - The oldest task keeps its original ID
  - Newer tasks get reassigned a fresh ID
  - File is renamed to match the new ID
  - Cross-references (dependencies, parent) in all tasks are updated

Use --dry-run to preview changes without modifying files.

Examples:
  taskmd deduplicate
  taskmd deduplicate ./tasks
  taskmd deduplicate --dry-run
  taskmd deduplicate --format json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDeduplicate,
}

func init() {
	rootCmd.AddCommand(deduplicateCmd)

	deduplicateCmd.Flags().BoolVar(&dedupDryRun, "dry-run", false, "preview changes without modifying files")
	deduplicateCmd.Flags().StringVar(&dedupFormat, "format", "text", "output format (text, json)")
}

type reassignment struct {
	OldID       string `json:"old_id"`
	NewID       string `json:"new_id"`
	OldFilePath string `json:"old_file_path"`
	NewFilePath string `json:"new_file_path"`
	Title       string `json:"title"`
}

type deduplicateResult struct {
	DryRun        bool           `json:"dry_run"`
	Duplicates    int            `json:"duplicates"`
	Reassignments []reassignment `json:"reassignments"`
}

func runDeduplicate(cmd *cobra.Command, args []string) error {
	if err := ValidateFormat(dedupFormat, []string{"text", "json"}); err != nil {
		return err
	}

	flags := GetGlobalFlags()
	scanDir := ResolveScanDir(args)

	taskScanner := scanner.NewScanner(scanDir, flags.Verbose, flags.IgnoreDirs)
	scanResult, err := taskScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	reassignments, err := planReassignments(scanResult.Tasks)
	if err != nil {
		return err
	}

	if !dedupDryRun {
		if err := applyReassignments(reassignments, scanResult.Tasks); err != nil {
			return err
		}
	}

	result := deduplicateResult{
		DryRun:        dedupDryRun,
		Duplicates:    len(reassignments),
		Reassignments: reassignments,
	}

	return outputDeduplicateResult(result, flags.Quiet)
}

// planReassignments detects duplicate IDs and plans which tasks need new IDs.
func planReassignments(tasks []*model.Task) ([]reassignment, error) {
	idMap := buildIDMap(tasks)
	allIDs := collectAllIDs(tasks)
	cfg := resolveIDConfig()

	var reassignments []reassignment

	for id, group := range idMap {
		if len(group) < 2 {
			continue
		}

		sortByCreated(group)

		for _, task := range group[1:] {
			newID, err := generateID(allIDs, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to generate new ID for duplicate %q: %w", id, err)
			}

			allIDs = append(allIDs, newID)

			reassignments = append(reassignments, reassignment{
				OldID:       task.ID,
				NewID:       newID,
				OldFilePath: task.FilePath,
				NewFilePath: buildNewFilePath(task.FilePath, task.ID, newID),
				Title:       task.Title,
			})
		}
	}

	sort.Slice(reassignments, func(i, j int) bool {
		return reassignments[i].OldFilePath < reassignments[j].OldFilePath
	})

	return reassignments, nil
}

// buildIDMap groups tasks by their ID.
func buildIDMap(tasks []*model.Task) map[string][]*model.Task {
	m := make(map[string][]*model.Task)
	for _, t := range tasks {
		m[t.ID] = append(m[t.ID], t)
	}
	return m
}

// collectAllIDs returns a slice of all task IDs.
func collectAllIDs(tasks []*model.Task) []string {
	ids := make([]string, len(tasks))
	for i, t := range tasks {
		ids[i] = t.ID
	}
	return ids
}

// sortByCreated sorts tasks by Created date ascending (oldest first).
// Falls back to filepath alphabetical order for equal dates.
func sortByCreated(tasks []*model.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		if !tasks[i].Created.Equal(tasks[j].Created) {
			return tasks[i].Created.Before(tasks[j].Created)
		}
		return tasks[i].FilePath < tasks[j].FilePath
	})
}

// buildNewFilePath constructs the renamed file path by replacing the old ID prefix with the new ID.
func buildNewFilePath(oldPath, oldID, newID string) string {
	dir := filepath.Dir(oldPath)
	base := filepath.Base(oldPath)

	// Replace the ID prefix in the filename: "001-some-slug.md" → "abc123-some-slug.md"
	if strings.HasPrefix(base, oldID+"-") {
		base = newID + base[len(oldID):]
	} else if strings.HasPrefix(base, oldID+".") {
		base = newID + base[len(oldID):]
	} else {
		// Fallback: prepend new ID
		base = newID + "-" + base
	}

	return filepath.Join(dir, base)
}

// applyReassignments performs the actual file modifications.
func applyReassignments(reassignments []reassignment, allTasks []*model.Task) error {
	for _, r := range reassignments {
		// 1. Update the id field in the task's frontmatter.
		if err := taskfile.ReplaceID(r.OldFilePath, r.NewID); err != nil {
			return fmt.Errorf("failed to update ID in %s: %w", r.OldFilePath, err)
		}

		// 2. Rename the file.
		if r.OldFilePath != r.NewFilePath {
			if err := os.Rename(r.OldFilePath, r.NewFilePath); err != nil {
				return fmt.Errorf("failed to rename %s → %s: %w", r.OldFilePath, r.NewFilePath, err)
			}
		}
	}

	// 3. Update cross-references in all task files.
	//    Build a set of current file paths (accounting for renames).
	renamedPaths := make(map[string]string, len(reassignments))
	for _, r := range reassignments {
		renamedPaths[r.OldFilePath] = r.NewFilePath
	}

	for _, task := range allTasks {
		filePath := task.FilePath
		if newPath, ok := renamedPaths[filePath]; ok {
			filePath = newPath
		}

		for _, r := range reassignments {
			if err := taskfile.ReplaceReference(filePath, r.OldID, r.NewID); err != nil {
				return fmt.Errorf("failed to update references in %s: %w", filePath, err)
			}
		}
	}

	return nil
}

func outputDeduplicateResult(result deduplicateResult, quiet bool) error {
	if dedupFormat == "json" {
		return WriteJSON(os.Stdout, result)
	}

	if result.Duplicates == 0 {
		if !quiet {
			fmt.Println("No duplicate IDs found.")
		}
		return nil
	}

	prefix := ""
	if result.DryRun {
		prefix = "[dry-run] "
	}

	fmt.Printf("%sFound %d duplicate(s) to resolve:\n\n", prefix, result.Duplicates)
	for _, r := range result.Reassignments {
		fmt.Printf("  %s → %s  %s\n", r.OldID, r.NewID, r.Title)
		fmt.Printf("    %s → %s\n", r.OldFilePath, r.NewFilePath)
	}

	if result.DryRun {
		fmt.Println("\nNo changes made (dry-run mode).")
	} else {
		fmt.Printf("\nResolved %d duplicate(s).\n", result.Duplicates)
	}

	return nil
}
