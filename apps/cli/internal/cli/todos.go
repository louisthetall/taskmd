package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/driangle/taskmd/apps/cli/internal/todos"
)

var (
	todosDir     string
	todosMarkers []string
	todosInclude []string
	todosExclude []string
	todosFormat  string
	todosRawText bool
	todosRich    bool
)

var todosCmd = &cobra.Command{
	Use:   "todos",
	Short: "Find TODO/FIXME comments in source code",
	Long:  `Commands for scanning source code files to find marker comments like TODO, FIXME, HACK, and more.`,
}

var todosListCmd = &cobra.Command{
	Use:   "list",
	Short: "List TODO/FIXME comments found in source files",
	Long: `Scan source code files recursively for marker comments (TODO, FIXME, HACK, XXX, NOTE, BUG, OPTIMIZE)
and display them with file path, line number, marker type, and comment text.

Respects .gitignore and skips common non-source directories (node_modules, .git, vendor, etc.).
Supports language-aware comment parsing for Go, JavaScript, TypeScript, Python, Ruby, Shell, CSS, HTML, Rust, YAML, and TOML.

Exclude patterns can also be configured in .taskmd.yaml under todos.exclude.
CLI --exclude flags are additive with config patterns (both are applied).

Table columns:
  Default: id, file, line, tag, text
  --rich:  id, file, line, tag, scope, age, author, text

Examples:
  taskmd todos list
  taskmd todos list --dir ./src
  taskmd todos list --marker TODO --marker FIXME
  taskmd todos list --include "*.go"
  taskmd todos list --exclude "*.test.go"
  taskmd todos list --format json
  taskmd todos list --rich`,
	Args: cobra.NoArgs,
	RunE: runTodosList,
}

func init() {
	rootCmd.AddCommand(todosCmd)
	todosCmd.AddCommand(todosListCmd)

	todosListCmd.Flags().StringVar(&todosDir, "dir", ".", "directory to scan for source code")
	todosListCmd.Flags().StringArrayVar(&todosMarkers, "marker", nil, "filter by marker type (e.g. --marker TODO --marker FIXME)")
	todosListCmd.Flags().StringArrayVar(&todosInclude, "include", nil, "include only files matching glob pattern")
	todosListCmd.Flags().StringArrayVar(&todosExclude, "exclude", nil, "exclude files matching glob pattern")
	todosListCmd.Flags().StringVar(&todosFormat, "format", "table", "output format (table, json, yaml)")
	todosListCmd.Flags().BoolVar(&todosRawText, "raw-text", false, "include original source line text in output")
	todosListCmd.Flags().BoolVar(&todosRich, "rich", false, "include scope and git blame information (slower)")
}

func runTodosList(_ *cobra.Command, _ []string) error {
	flags := GetGlobalFlags()

	markers := todosMarkers
	if len(markers) == 0 {
		markers = todos.DefaultMarkers
	}

	if err := validateMarkers(markers); err != nil {
		return err
	}

	if err := ValidateFormat(todosFormat, []string{"table", "json", "yaml"}); err != nil {
		return err
	}

	excludeGlobs := mergeConfigExcludes(todosExclude)

	items, err := todos.Scan(todos.ScanOptions{
		Dir:          todosDir,
		Markers:      markers,
		IncludeGlobs: todosInclude,
		ExcludeGlobs: excludeGlobs,
		Verbose:      flags.Verbose,
		RawText:      todosRawText,
	})
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if todosRich {
		todos.EnrichRich(items, todosDir)
	}

	switch todosFormat {
	case "json":
		if items == nil {
			items = []todos.TodoItem{}
		}
		return WriteJSON(os.Stdout, items)
	case "yaml":
		if items == nil {
			items = []todos.TodoItem{}
		}
		return WriteYAML(os.Stdout, items)
	default:
		cols := defaultColumns
		if todosRich {
			cols = richColumns
		}
		return outputTodosTable(items, cols, todosRich)
	}
}

// mergeConfigExcludes combines CLI --exclude flags with todos.exclude from .taskmd.yaml.
func mergeConfigExcludes(cliExcludes []string) []string {
	configExcludes := viper.GetStringSlice("todos.exclude")
	if len(configExcludes) == 0 {
		return cliExcludes
	}
	if len(cliExcludes) == 0 {
		return configExcludes
	}
	merged := make([]string, 0, len(cliExcludes)+len(configExcludes))
	merged = append(merged, cliExcludes...)
	merged = append(merged, configExcludes...)
	return merged
}

func validateMarkers(markers []string) error {
	for _, m := range markers {
		upper := strings.ToUpper(m)
		if !slices.Contains(todos.DefaultMarkers, upper) {
			return fmt.Errorf("invalid marker %q: must be one of %s", m, strings.Join(todos.DefaultMarkers, ", "))
		}
	}
	return nil
}

// todoColumn defines a table column with its header, separator, and value extractor.
type todoColumn struct {
	header string
	sep    string
	value  func(item todos.TodoItem, r *lipgloss.Renderer) string
}

var todoColumnDefs = map[string]todoColumn{
	"id": {"ID", "--------", func(item todos.TodoItem, r *lipgloss.Renderer) string {
		id := item.ID
		if len(id) > 8 {
			id = id[:8]
		}
		return formatDim(id, r)
	}},
	"file":  {"FILE", "----------", func(item todos.TodoItem, r *lipgloss.Renderer) string { return formatDim(item.FilePath, r) }},
	"line":  {"LINE", "----", func(item todos.TodoItem, _ *lipgloss.Renderer) string { return fmt.Sprintf("%d", item.Line) }},
	"tag":   {"TAG", "--------", func(item todos.TodoItem, r *lipgloss.Renderer) string { return formatMarker(item.Marker, r) }},
	"text":  {"TEXT", "----------", func(item todos.TodoItem, _ *lipgloss.Renderer) string { return item.Text }},
	"scope": {"SCOPE", "----------", func(item todos.TodoItem, _ *lipgloss.Renderer) string { return item.Scope }},
	"age": {"AGE", "---", func(item todos.TodoItem, _ *lipgloss.Renderer) string {
		if item.Age > 0 {
			return fmt.Sprintf("%dd", item.Age)
		}
		return ""
	}},
	"author": {"AUTHOR", "----------", func(item todos.TodoItem, _ *lipgloss.Renderer) string {
		if item.Blame != nil {
			return item.Blame.Author
		}
		return ""
	}},
}

var (
	defaultColumns = []string{"id", "file", "line", "tag", "text"}
	richColumns    = []string{"id", "file", "line", "tag", "scope", "age", "author", "text"}
)

func outputTodosTable(items []todos.TodoItem, columns []string, rich bool) error {
	if len(items) == 0 {
		fmt.Println("No TODO comments found")
		return nil
	}

	r := getRenderer()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Header
	headers := make([]string, len(columns))
	seps := make([]string, len(columns))
	for i, name := range columns {
		col := todoColumnDefs[name]
		headers[i] = col.header
		seps[i] = col.sep
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Rows
	for _, item := range items {
		vals := make([]string, len(columns))
		for i, name := range columns {
			vals[i] = todoColumnDefs[name].value(item, r)
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}

	w.Flush()

	fmt.Fprintf(os.Stderr, "\nFound %d comment(s)\n", len(items))

	if !rich {
		fmt.Fprintf(os.Stderr, "Use --rich to add: scope, age, author columns\n")
	}

	return nil
}

func formatMarker(marker string, r *lipgloss.Renderer) string {
	switch strings.ToUpper(marker) {
	case "TODO":
		return r.NewStyle().Foreground(lipgloss.Color("3")).Render(marker) // Yellow
	case "FIXME", "BUG", "XXX":
		return r.NewStyle().Foreground(lipgloss.Color("1")).Render(marker) // Red
	case "HACK":
		return r.NewStyle().Foreground(lipgloss.Color("5")).Render(marker) // Magenta
	case "NOTE":
		return r.NewStyle().Foreground(lipgloss.Color("6")).Render(marker) // Cyan
	case "OPTIMIZE":
		return r.NewStyle().Foreground(lipgloss.Color("4")).Render(marker) // Blue
	default:
		return marker
	}
}
