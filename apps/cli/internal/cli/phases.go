package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/driangle/taskmd/sdk/go/model"
	"github.com/driangle/taskmd/sdk/go/scanner"
	"github.com/driangle/taskmd/sdk/go/validator"
)

var phasesFormat string

var phasesCmd = &cobra.Command{
	Use:        "phases",
	SuggestFor: []string{"milestones", "phase"},
	Short:      "List project phases with progress stats",
	Long: `Phases displays all configured phases from .taskmd.yaml with summary stats:
- Task count and completion percentage per phase
- Status breakdown per phase
- Due dates

Phases are configured in .taskmd.yaml under the "phases" key.

Examples:
  taskmd phases
  taskmd phases ./tasks
  taskmd phases --format json
  taskmd phases --format yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPhases,
}

func init() {
	rootCmd.AddCommand(phasesCmd)

	phasesCmd.Flags().StringVar(&phasesFormat, "format", "table", "output format (table, json, yaml)")
}

// PhaseSummary holds computed stats for a single phase.
type PhaseSummary struct {
	ID       string         `json:"id" yaml:"id"`
	Name     string         `json:"name" yaml:"name"`
	Tasks    int            `json:"tasks" yaml:"tasks"`
	Done     int            `json:"done" yaml:"done"`
	Progress string         `json:"progress" yaml:"progress"`
	Due      string         `json:"due,omitempty" yaml:"due,omitempty"`
	ByStatus map[string]int `json:"by_status" yaml:"by_status"`
}

func runPhases(cmd *cobra.Command, args []string) error {
	flags := GetGlobalFlags()
	scanDir := ResolveScanDir(args)

	taskScanner := scanner.NewScanner(scanDir, flags.Verbose, flags.IgnoreDirs)
	result, err := taskScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	tasks := result.Tasks

	if flags.Verbose && len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "\nWarning: encountered %d errors during scan:\n", len(result.Errors))
		for _, scanErr := range result.Errors {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", scanErr.FilePath, scanErr.Error)
		}
		fmt.Fprintln(os.Stderr)
	}

	phases := parsePhasesConfig(viper.Get("phases"))
	if len(phases) == 0 {
		fmt.Fprintln(os.Stderr, "No phases configured. Add phases to your .taskmd.yaml file:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  phases:")
		fmt.Fprintln(os.Stderr, "    - id: mvp")
		fmt.Fprintln(os.Stderr, "      name: \"MVP\"")
		fmt.Fprintln(os.Stderr, "      due: 2026-06-01")
		return nil
	}

	summaries := computePhaseSummaries(phases, tasks)
	warnOrphanedPhases(phases, tasks)

	switch phasesFormat {
	case "json":
		return WriteJSON(os.Stdout, summaries)
	case "yaml":
		return WriteYAML(os.Stdout, summaries)
	case "table":
		return outputPhasesTable(summaries)
	default:
		return ValidateFormat(phasesFormat, []string{"table", "json", "yaml"})
	}
}

func computePhaseSummaries(phases []validator.PhaseConfig, tasks []*model.Task) []PhaseSummary {
	summaries := make([]PhaseSummary, 0, len(phases))
	for _, phase := range phases {
		due := ""
		if !phase.Due.IsZero() {
			due = phase.Due.Format("2006-01-02")
		}
		summary := PhaseSummary{
			ID:       phase.ID,
			Name:     phase.Name,
			Due:      due,
			ByStatus: make(map[string]int),
		}
		for _, task := range tasks {
			if task.Phase != phase.ID {
				continue
			}
			summary.ByStatus[string(task.Status)]++
			if task.Status == model.StatusCancelled {
				continue
			}
			summary.Tasks++
			if task.Status == model.StatusCompleted {
				summary.Done++
			}
		}
		if summary.Tasks > 0 {
			pct := float64(summary.Done) / float64(summary.Tasks) * 100
			summary.Progress = fmt.Sprintf("%.0f%%", pct)
		} else {
			summary.Progress = "0%"
		}
		summaries = append(summaries, summary)
	}
	return summaries
}

func warnOrphanedPhases(phases []validator.PhaseConfig, tasks []*model.Task) {
	knownIDs := make(map[string]bool, len(phases))
	for _, p := range phases {
		knownIDs[p.ID] = true
	}
	orphaned := make(map[string]int)
	for _, task := range tasks {
		if task.Phase != "" && !knownIDs[task.Phase] {
			orphaned[task.Phase]++
		}
	}
	sortedPhases := make([]string, 0, len(orphaned))
	for phase := range orphaned {
		sortedPhases = append(sortedPhases, phase)
	}
	sort.Strings(sortedPhases)
	for _, phase := range sortedPhases {
		fmt.Fprintf(os.Stderr, "Warning: %d task(s) reference undefined phase %q\n", orphaned[phase], phase)
	}
}

func outputPhasesTable(summaries []PhaseSummary) error {
	tw := NewTableWriter()
	tw.AddHeader([]string{"ID", "Name", "Tasks", "Done", "Progress", "Due"})
	tw.AddSeparator()
	for _, s := range summaries {
		due := s.Due
		if due == "" {
			due = "-"
		}
		row := []string{
			s.ID,
			s.Name,
			fmt.Sprintf("%d", s.Tasks),
			fmt.Sprintf("%d", s.Done),
			s.Progress,
			due,
		}
		tw.AddRow(row, row)
	}
	tw.Flush(os.Stdout)
	return nil
}
