package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/driangle/taskmd/sdk/go/next"
	"github.com/driangle/taskmd/sdk/go/scanner"
)

// Recommendation is re-exported from the shared package.
type Recommendation = next.Recommendation

var (
	nextFormat    string
	nextLimit     int
	nextFilters   []string
	nextQuickWins bool
	nextCritical  bool
	nextScope     string
	nextExact     bool
	nextPhase     string
)

var nextCmd = &cobra.Command{
	Use:        "next",
	SuggestFor: []string{"pick", "suggest", "what"},
	Short:      "Recommend what task to work on next",
	Long: `Next analyzes all tasks and recommends the best ones to work on next.

Tasks are scored based on priority, critical path position, downstream impact,
effort, and phase ordering (from .taskmd.yaml). Only actionable tasks
(pending or in-progress with all dependencies completed) are shown.

Output formats: table (default), json, yaml

Examples:
  taskmd next
  taskmd next ./tasks
  taskmd next --limit 3
  taskmd next --filter tag=cli
  taskmd next --filter priority=high --format json
  taskmd next --quick-wins
  taskmd next --critical --limit 1
  taskmd next --scope web/graph
  taskmd next --scope web/graph --exact
  taskmd next --phase v0.2`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNext,
}

func init() {
	rootCmd.AddCommand(nextCmd)

	nextCmd.Flags().StringVar(&nextFormat, "format", "table", "output format (table, json, yaml)")
	nextCmd.Flags().IntVar(&nextLimit, "limit", 5, "maximum number of recommendations")
	nextCmd.Flags().StringArrayVar(&nextFilters, "filter", []string{}, "filter tasks (e.g., --filter tag=cli)")
	nextCmd.Flags().BoolVar(&nextQuickWins, "quick-wins", false, "show only quick wins (effort: small)")
	nextCmd.Flags().BoolVar(&nextCritical, "critical", false, "show only critical path tasks")
	nextCmd.Flags().StringVar(&nextScope, "scope", "", "filter by scope; supports wildcards (e.g. cli, cli*)")
	nextCmd.Flags().BoolVar(&nextExact, "exact", false, "disable dependency expansion for --scope (only direct matches)")
	nextCmd.Flags().StringVar(&nextPhase, "phase", "", "filter by phase")
}

func runNext(cmd *cobra.Command, args []string) error {
	flags := GetGlobalFlags()
	scanDir := ResolveScanDir(args)

	taskScanner := scanner.NewScanner(scanDir, flags.Verbose, flags.IgnoreDirs)
	result, err := taskScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	allTasks := result.Tasks
	makeFilePathsRelative(allTasks, scanDir)

	warnDuplicateIDs(allTasks)

	archivedTasks, err := taskScanner.ScanArchive()
	if err != nil {
		return fmt.Errorf("archive scan failed: %w", err)
	}

	phaseOrder := loadPhaseOrder()

	recs, err := next.Recommend(allTasks, next.Options{
		Limit:         nextLimit,
		Filters:       nextFilters,
		QuickWins:     nextQuickWins,
		Critical:      nextCritical,
		Scope:         nextScope,
		ScopeExact:    nextExact,
		ArchivedTasks: archivedTasks,
		Phase:         nextPhase,
		PhaseOrder:    phaseOrder,
	})
	if err != nil {
		return err
	}

	switch nextFormat {
	case "json":
		return outputNextJSON(recs)
	case "yaml":
		return outputNextYAML(recs)
	case "table":
		return outputNextTable(recs)
	default:
		return ValidateFormat(nextFormat, []string{"table", "json", "yaml"})
	}
}

// loadPhaseOrder reads phase names from the viper config, preserving order.
func loadPhaseOrder() []string {
	raw := viper.Get("phases")
	if raw == nil {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	var names []string
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if name, ok := m["name"].(string); ok {
			names = append(names, name)
		}
	}
	return names
}

func outputNextJSON(recs []Recommendation) error {
	return WriteJSON(os.Stdout, recs)
}

func outputNextYAML(recs []Recommendation) error {
	return WriteYAML(os.Stdout, recs)
}

func outputNextTable(recs []Recommendation) error {
	r := getRenderer()

	if len(recs) == 0 {
		if nextScope != "" {
			fmt.Printf("No actionable tasks found for scope %q.\n", nextScope)
		} else if nextQuickWins {
			fmt.Println("No quick wins available.")
		} else if nextCritical {
			fmt.Println("No critical path tasks available.")
		} else {
			fmt.Println("No actionable tasks found.")
		}
		return nil
	}

	label := "Recommended tasks:"
	if nextScope != "" {
		label = fmt.Sprintf("Recommended tasks (scope: %s):", nextScope)
	}
	if nextQuickWins {
		label = "Recommended quick wins:"
	}
	if nextCritical {
		label = "Recommended critical path tasks:"
	}
	fmt.Println(formatLabel(label, r))
	fmt.Println()

	tw := NewTableWriter()
	tw.AddHeader([]string{"#", "ID", "Title", "Priority", "Effort", "File", "Reason"})
	tw.AddSeparator()

	for _, rec := range recs {
		rank := fmt.Sprintf("%d", rec.Rank)
		reason := strings.Join(rec.Reasons, ", ")
		plain := []string{rank, rec.ID, rec.Title, rec.Priority, rec.Effort, rec.FilePath, reason}
		colored := []string{
			rank,
			formatTaskID(rec.ID, r),
			rec.Title,
			formatPriority(rec.Priority, r),
			formatEffort(rec.Effort, r),
			formatDim(rec.FilePath, r),
			reason,
		}
		tw.AddRow(plain, colored)
	}

	tw.Flush(os.Stdout)
	return nil
}
