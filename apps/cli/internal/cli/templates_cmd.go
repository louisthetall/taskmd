package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/apps/cli/internal/template"
)

var templatesFormat string

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage task templates",
	Long:  `Commands for listing and inspecting task templates used by the add command.`,
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available task templates",
	Long: `List all available task templates from project, user, and built-in sources.

Templates are discovered in precedence order: project (.taskmd/templates/) overrides
user (~/.taskmd/templates/) overrides built-in templates.

Examples:
  taskmd templates list
  taskmd templates list --format json
  taskmd templates list --format yaml`,
	Args: cobra.NoArgs,
	RunE: runTemplatesList,
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.AddCommand(templatesListCmd)

	templatesListCmd.Flags().StringVar(&templatesFormat, "format", "table", "output format (table, json, yaml)")
}

type templateListItem struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Source      string `json:"source" yaml:"source"`
}

func runTemplatesList(_ *cobra.Command, _ []string) error {
	if err := ValidateFormat(templatesFormat, []string{"table", "json", "yaml"}); err != nil {
		return err
	}

	projectRoot := resolveProjectRoot()
	userHome, _ := os.UserHomeDir()

	templates := template.Discover(projectRoot, userHome)

	items := make([]templateListItem, len(templates))
	for i, tmpl := range templates {
		items[i] = templateListItem{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Source:      tmpl.Source,
		}
	}

	switch templatesFormat {
	case "json":
		return WriteJSON(os.Stdout, items)
	case "yaml":
		return WriteYAML(os.Stdout, items)
	default:
		return outputTemplatesTable(items)
	}
}

func outputTemplatesTable(items []templateListItem) error {
	if len(items) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tSOURCE")
	fmt.Fprintln(w, "----\t-----------\t------")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.Description, item.Source)
	}
	return w.Flush()
}
