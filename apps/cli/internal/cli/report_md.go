package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/driangle/taskmd/apps/cli/internal/metrics"
	"github.com/driangle/taskmd/apps/cli/internal/model"
)

func outputReportMarkdown(data *reportData, w io.Writer) error {
	r := getReportRenderer()

	fmt.Fprintln(w, formatLabel("# Project Report", r))
	fmt.Fprintln(w)

	writeMarkdownSummary(data.Metrics, w, r)
	writeMarkdownGroups(data, w, r)
	writeMarkdownCriticalPath(data.CriticalPath, w, r)
	writeMarkdownBlockedTasks(data, w, r)

	if data.IncludeGraph {
		writeMarkdownGraph(data.GraphMermaid, w)
	}

	return nil
}

// getReportRenderer returns a renderer that disables color when writing to a file.
func getReportRenderer() *lipgloss.Renderer {
	if reportOut != "" {
		return getNoColorRenderer()
	}
	return getRenderer()
}

func writeMarkdownSummary(m *metrics.Metrics, w io.Writer, r *lipgloss.Renderer) {
	fmt.Fprintln(w, formatLabel("## Summary", r))
	fmt.Fprintln(w)
	fmt.Fprintf(w, "| Metric | Value |\n")
	fmt.Fprintf(w, "|--------|-------|\n")
	fmt.Fprintf(w, "| Total Tasks | %d |\n", m.TotalTasks)
	fmt.Fprintf(w, "| Blocked Tasks | %d |\n", m.BlockedTasksCount)
	fmt.Fprintf(w, "| Critical Path Length | %d |\n", m.CriticalPathLength)
	fmt.Fprintf(w, "| Avg Dependencies | %.1f |\n", m.AvgDependenciesPerTask)
	fmt.Fprintln(w)

	writeMarkdownStatusBreakdown(m, w, r)
	writeMarkdownPriorityBreakdown(m, w, r)
}

func writeMarkdownStatusBreakdown(m *metrics.Metrics, w io.Writer, r *lipgloss.Renderer) {
	statusOrder := []model.Status{
		model.StatusPending, model.StatusInProgress, model.StatusBlocked,
		model.StatusCompleted, model.StatusCancelled,
	}

	fmt.Fprintln(w, formatLabel("### By Status", r))
	fmt.Fprintln(w)
	for _, s := range statusOrder {
		if count, ok := m.TasksByStatus[s]; ok && count > 0 {
			fmt.Fprintf(w, "- %s: %d\n", formatStatus(string(s), r), count)
		}
	}
	fmt.Fprintln(w)
}

func writeMarkdownPriorityBreakdown(m *metrics.Metrics, w io.Writer, r *lipgloss.Renderer) {
	priorityOrder := []model.Priority{
		model.PriorityCritical, model.PriorityHigh,
		model.PriorityMedium, model.PriorityLow,
	}

	hasAny := false
	for _, p := range priorityOrder {
		if count, ok := m.TasksByPriority[p]; ok && count > 0 {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return
	}

	fmt.Fprintln(w, formatLabel("### By Priority", r))
	fmt.Fprintln(w)
	for _, p := range priorityOrder {
		if count, ok := m.TasksByPriority[p]; ok && count > 0 {
			fmt.Fprintf(w, "- %s: %d\n", formatPriority(string(p), r), count)
		}
	}
	fmt.Fprintln(w)
}

func writeMarkdownGroups(data *reportData, w io.Writer, r *lipgloss.Renderer) {
	fmt.Fprintf(w, "%s\n", formatLabel(fmt.Sprintf("## Tasks by %s", capitalizeFirst(data.GroupBy)), r))
	fmt.Fprintln(w)

	for _, key := range data.GroupedTasks.Keys {
		tasks := data.GroupedTasks.Groups[key]
		coloredHeading := formatHeading(key, data.GroupBy, r)
		fmt.Fprintf(w, "### %s (%d)\n", coloredHeading, len(tasks))
		fmt.Fprintln(w)
		for _, t := range tasks {
			formattedID := formatTaskID(t.ID, r)
			line := fmt.Sprintf("- [%s] %s", formattedID, t.Title)
			if t.Priority != "" {
				line += fmt.Sprintf(" (priority: %s)", formatPriority(string(t.Priority), r))
			}
			fmt.Fprintln(w, line)
		}
		fmt.Fprintln(w)
	}
}

func writeMarkdownCriticalPath(cpTasks []reportTask, w io.Writer, r *lipgloss.Renderer) {
	fmt.Fprintln(w, formatLabel("## Critical Path", r))
	fmt.Fprintln(w)

	if len(cpTasks) == 0 {
		fmt.Fprintln(w, "No dependency chains found.")
		fmt.Fprintln(w)
		return
	}

	for i, t := range cpTasks {
		formattedID := formatTaskID(t.ID, r)
		formattedStatus := formatStatus(t.Status, r)
		fmt.Fprintf(w, "%d. [%s] %s (%s)\n", i+1, formattedID, t.Title, formattedStatus)
	}
	fmt.Fprintln(w)
}

func writeMarkdownBlockedTasks(data *reportData, w io.Writer, r *lipgloss.Renderer) {
	fmt.Fprintln(w, formatLabel("## Blocked Tasks", r))
	fmt.Fprintln(w)

	if len(data.BlockedTasks) == 0 {
		fmt.Fprintln(w, "No blocked tasks.")
		fmt.Fprintln(w)
		return
	}

	taskMap := make(map[string]reportTask)
	for _, t := range data.CriticalPath {
		taskMap[t.ID] = t
	}
	for _, t := range data.BlockedTasks {
		taskMap[t.ID] = t
	}
	// Also build from grouped tasks for status lookups
	allTaskMap := buildTaskMapFromGroups(data)

	for _, t := range data.BlockedTasks {
		formattedID := formatTaskID(t.ID, r)
		waitingOn := formatWaitingOnColored(t.Dependencies, allTaskMap, r)
		fmt.Fprintf(w, "- [%s] %s\n  Waiting on: %s\n", formattedID, t.Title, waitingOn)
	}
	fmt.Fprintln(w)
}

func buildTaskMapFromGroups(data *reportData) map[string]reportTask {
	m := make(map[string]reportTask)
	for _, key := range data.GroupedTasks.Keys {
		for _, t := range data.GroupedTasks.Groups[key] {
			m[t.ID] = reportTask{
				ID:     t.ID,
				Title:  t.Title,
				Status: string(t.Status),
			}
		}
	}
	return m
}

func formatWaitingOn(deps []string, taskMap map[string]reportTask) string {
	parts := make([]string, len(deps))
	for i, depID := range deps {
		if t, ok := taskMap[depID]; ok {
			parts[i] = fmt.Sprintf("%s (%s)", depID, t.Status)
		} else {
			parts[i] = depID
		}
	}
	return strings.Join(parts, ", ")
}

func formatWaitingOnColored(deps []string, taskMap map[string]reportTask, r *lipgloss.Renderer) string {
	parts := make([]string, len(deps))
	for i, depID := range deps {
		if t, ok := taskMap[depID]; ok {
			formattedID := formatTaskID(depID, r)
			formattedStatus := formatStatus(t.Status, r)
			parts[i] = fmt.Sprintf("%s (%s)", formattedID, formattedStatus)
		} else {
			parts[i] = formatTaskID(depID, r)
		}
	}
	return strings.Join(parts, ", ")
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func writeMarkdownGraph(mermaid string, w io.Writer) {
	fmt.Fprintln(w, "## Dependency Graph")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "```mermaid")
	fmt.Fprint(w, mermaid)
	fmt.Fprintln(w, "```")
	fmt.Fprintln(w)
}
