package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/apps/cli/internal/model"
	"github.com/driangle/taskmd/apps/cli/internal/scanner"
	"github.com/driangle/taskmd/apps/cli/internal/taskfile"
	"github.com/driangle/taskmd/apps/cli/internal/verify"
)

var (
	setTaskID     string
	setStatus     string
	setPriority   string
	setEffort     string
	setType       string
	setOwner      string
	setParent     string
	setDone       bool
	setDryRun     bool
	setVerify     bool
	setAddTags    []string
	setRemoveTags []string
	setAddPRs     []string
	setRemovePRs  []string
)

var setCmd = &cobra.Command{
	Use:        "set [task-id]",
	SuggestFor: []string{"edit", "modify", "change"},
	Short:      "Set a task's frontmatter fields",
	Long: `Set modifies frontmatter fields (status, priority, effort, tags) of a task file.

The task is identified by a positional argument or --task-id (exact match only).

Examples:
  taskmd set cli-049 --status completed
  taskmd set cli-049 --priority high --effort large
  taskmd set cli-049 --done
  taskmd set cli-049 --add-tag backend --add-tag api
  taskmd set cli-049 --remove-tag deprecated
  taskmd set --task-id cli-049 --status completed   # --task-id also works`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSet,
}

// Deprecated: use "set" instead.
var updateCmd = &cobra.Command{
	Use:        "update [task-id]",
	Short:      "Update a task's frontmatter fields (deprecated: use 'set')",
	Args:       cobra.MaximumNArgs(1),
	RunE:       runSet,
	Hidden:     true,
	Deprecated: "use 'set' instead",
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(updateCmd)

	for _, cmd := range []*cobra.Command{setCmd, updateCmd} {
		cmd.Flags().StringVar(&setTaskID, "task-id", "", "task ID to update (required)")
		cmd.Flags().StringVar(&setStatus, "status", "", "new status (pending, in-progress, completed, in-review, blocked, cancelled)")
		cmd.Flags().StringVar(&setPriority, "priority", "", "new priority (low, medium, high, critical)")
		cmd.Flags().StringVar(&setEffort, "effort", "", "new effort (small, medium, large)")
		cmd.Flags().StringVar(&setType, "type", "", "work type (feature, bug, improvement, chore, docs)")
		cmd.Flags().StringVar(&setOwner, "owner", "", "owner/assignee of the task")
		cmd.Flags().StringVar(&setParent, "parent", "", "parent task ID (use empty string to clear)")
		cmd.Flags().BoolVar(&setDone, "done", false, "mark task as completed (alias for --status completed)")
		cmd.Flags().BoolVar(&setDryRun, "dry-run", false, "preview changes without writing to disk")
		cmd.Flags().BoolVar(&setVerify, "verify", false, "run verification checks before completing a task")
		cmd.Flags().StringArrayVar(&setAddTags, "add-tag", nil, "add a tag (repeatable)")
		cmd.Flags().StringArrayVar(&setRemoveTags, "remove-tag", nil, "remove a tag (repeatable)")
		cmd.Flags().StringArrayVar(&setAddPRs, "add-pr", nil, "add a PR URL (repeatable)")
		cmd.Flags().StringArrayVar(&setRemovePRs, "remove-pr", nil, "remove a PR URL (repeatable)")
	}
}

func resolveSetTaskID(cmd *cobra.Command, args []string) (string, error) {
	positional := ""
	if len(args) > 0 {
		positional = args[0]
	}
	flagVal := setTaskID

	if positional != "" && flagVal != "" && positional != flagVal {
		return "", fmt.Errorf("conflicting task ID: positional argument %q and --task-id %q differ", positional, flagVal)
	}

	if positional != "" {
		return positional, nil
	}
	if flagVal != "" {
		return flagVal, nil
	}
	return "", fmt.Errorf("task ID required: provide as positional argument or --task-id flag")
}

func runSet(cmd *cobra.Command, args []string) error {
	taskID, err := resolveSetTaskID(cmd, args)
	if err != nil {
		return err
	}

	req, err := buildSetRequest(cmd)
	if err != nil {
		return err
	}

	flags := GetGlobalFlags()
	scanDir := ResolveScanDir(nil)

	taskScanner := scanner.NewScanner(scanDir, flags.Verbose, flags.IgnoreDirs)
	result, err := taskScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	debugLog("scan directory: %s", scanDir)
	debugLog("found %d task(s)", len(result.Tasks))

	task := findExactMatch(taskID, result.Tasks)
	if task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if err := runSetVerification(task, req); err != nil {
		return err
	}

	changes := buildChangeLog(task, req)

	if setDryRun {
		printSetConfirmation(task, changes)
		r := getRenderer()
		fmt.Println("\n" + formatWarning("Dry run — no changes made.", r))
		return nil
	}

	if err := taskfile.UpdateTaskFile(task.FilePath, req); err != nil {
		return err
	}

	printSetConfirmation(task, changes)
	return nil
}

// setStringField sets a pointer field if the value is non-empty.
func setStringField(dst **string, value string) {
	if value != "" {
		*dst = &value
	}
}

func buildSetRequest(cmd *cobra.Command) (taskfile.UpdateRequest, error) {
	if setDone && cmd.Flags().Changed("status") {
		return taskfile.UpdateRequest{}, fmt.Errorf("--done and --status are mutually exclusive")
	}

	if setDone {
		flags := GetGlobalFlags()
		if flags.Workflow == "pr-review" {
			setStatus = string(model.StatusInReview)
		} else {
			setStatus = string(model.StatusCompleted)
		}
	}

	var req taskfile.UpdateRequest

	setStringField(&req.Status, setStatus)
	setStringField(&req.Priority, setPriority)
	setStringField(&req.Effort, setEffort)
	setStringField(&req.Type, setType)
	setStringField(&req.Owner, setOwner)

	if cmd.Flags().Changed("parent") {
		req.Parent = &setParent
	}

	if len(setAddTags) > 0 {
		req.AddTags = setAddTags
	}
	if len(setRemoveTags) > 0 {
		req.RemTags = setRemoveTags
	}
	if len(setAddPRs) > 0 {
		req.AddPRs = setAddPRs
	}
	if len(setRemovePRs) > 0 {
		req.RemPRs = setRemovePRs
	}

	if err := validateSetEnums(req); err != nil {
		return taskfile.UpdateRequest{}, err
	}

	if !hasUpdates(req) {
		return taskfile.UpdateRequest{}, fmt.Errorf("nothing to update: provide --status, --priority, --effort, --type, --owner, --parent, --done, --add-tag, --remove-tag, --add-pr, or --remove-pr")
	}

	return req, nil
}

func hasUpdates(req taskfile.UpdateRequest) bool {
	hasScalar := req.Status != nil || req.Priority != nil || req.Effort != nil ||
		req.Type != nil || req.Owner != nil || req.Parent != nil
	hasTags := len(req.AddTags) > 0 || len(req.RemTags) > 0
	hasPRs := len(req.AddPRs) > 0 || len(req.RemPRs) > 0
	return hasScalar || hasTags || hasPRs
}

type changeEntry struct {
	field    string
	oldValue string
	newValue string
}

func buildChangeLog(task *model.Task, req taskfile.UpdateRequest) []changeEntry {
	oldValues := map[string]string{
		"status":   string(task.Status),
		"priority": string(task.Priority),
		"effort":   string(task.Effort),
		"type":     string(task.Type),
		"owner":    task.Owner,
		"parent":   task.Parent,
	}

	var changes []changeEntry

	if req.Status != nil {
		changes = append(changes, changeEntry{field: "status", oldValue: oldValues["status"], newValue: *req.Status})
	}
	if req.Priority != nil {
		changes = append(changes, changeEntry{field: "priority", oldValue: oldValues["priority"], newValue: *req.Priority})
	}
	if req.Effort != nil {
		changes = append(changes, changeEntry{field: "effort", oldValue: oldValues["effort"], newValue: *req.Effort})
	}
	if req.Type != nil {
		changes = append(changes, changeEntry{field: "type", oldValue: oldValues["type"], newValue: *req.Type})
	}
	if req.Owner != nil {
		changes = append(changes, changeEntry{field: "owner", oldValue: oldValues["owner"], newValue: *req.Owner})
	}
	if req.Parent != nil {
		changes = append(changes, changeEntry{field: "parent", oldValue: oldValues["parent"], newValue: *req.Parent})
	}

	if len(req.AddTags) > 0 || len(req.RemTags) > 0 {
		newTags := taskfile.ComputeNewTags(task.Tags, req.AddTags, req.RemTags)
		changes = append(changes, changeEntry{
			field:    "tags",
			oldValue: "[" + strings.Join(task.Tags, ", ") + "]",
			newValue: "[" + strings.Join(newTags, ", ") + "]",
		})
	}

	if len(req.AddPRs) > 0 || len(req.RemPRs) > 0 {
		newPRs := taskfile.ComputeNewTags(task.PRs, req.AddPRs, req.RemPRs)
		changes = append(changes, changeEntry{
			field:    "pr",
			oldValue: "[" + strings.Join(task.PRs, ", ") + "]",
			newValue: "[" + strings.Join(newPRs, ", ") + "]",
		})
	}

	return changes
}

func printSetConfirmation(task *model.Task, changes []changeEntry) {
	r := getRenderer()
	fmt.Printf("Updated task %s (%s):\n", formatTaskID(task.ID, r), task.Title)
	for _, c := range changes {
		old := c.oldValue
		if old == "" {
			old = "(unset)"
		}
		fmt.Printf("  %s: %s -> %s\n",
			formatLabel(c.field, r),
			formatDim(old, r),
			colorizeFieldValue(c.field, c.newValue, r),
		)
	}
}

func colorizeFieldValue(field, value string, r *lipgloss.Renderer) string {
	switch field {
	case "status":
		return formatStatus(value, r)
	case "priority":
		return formatPriority(value, r)
	case "effort":
		return formatEffort(value, r)
	default:
		return value
	}
}

// runSetVerification runs verify checks if --verify is set and status is being set to completed.
func runSetVerification(task *model.Task, req taskfile.UpdateRequest) error {
	if !setVerify {
		return nil
	}
	isCompleting := req.Status != nil && *req.Status == string(model.StatusCompleted)
	if !isCompleting {
		return nil
	}
	if len(task.Verify) == 0 {
		return nil
	}

	if errs := model.ValidateVerifySteps(task.Verify); len(errs) > 0 {
		return fmt.Errorf("invalid verify steps:\n  %s", strings.Join(errs, "\n  "))
	}

	flags := GetGlobalFlags()
	projectRoot := resolveProjectRoot()
	opts := verify.Options{
		ProjectRoot: projectRoot,
		Timeout:     60 * time.Second,
		Verbose:     flags.Verbose,
		LogFunc: func(format string, args ...any) {
			if !flags.Quiet {
				fmt.Fprintf(os.Stderr, format+"\n", args...)
			}
		},
	}

	vResult := verify.Run(task.Verify, opts)
	printVerifyTable(vResult)

	if vResult.HasFailures() {
		return fmt.Errorf("verification failed: %d check(s) failed — status change aborted", vResult.Failed)
	}
	return nil
}
