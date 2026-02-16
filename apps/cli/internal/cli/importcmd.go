package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/apps/cli/internal/sync"
	"github.com/driangle/taskmd/apps/cli/internal/sync/github"
)

var (
	importSource    string
	importProject   string
	importTokenEnv  string
	importUserEnv   string
	importBaseURL   string
	importOutDir    string
	importFilter    string
	importDryRun    bool
	importFormat    string
	importRepo      string
	importLabels    string
	importMilestone string
	importAssignee  string
	importURL       string
	importJQL       string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import tasks from external sources",
	Long: `Import fetches tasks from an external source (GitHub Issues, Jira, etc.)
and creates local markdown task files. This is a one-time onboarding tool
for populating your tasks/ directory.

When run without --source, an interactive wizard guides you through setup.

GitHub-specific shortcut flags:
  --repo        Alias for --project (owner/repo)
  --labels      Filter by labels (comma-separated)
  --milestone   Filter by milestone
  --assignee    Filter by assignee

Jira-specific shortcut flags:
  --url         Alias for --base-url (Jira instance URL)
  --jql         Jira Query Language filter

Examples:
  taskmd import
  taskmd import --source github --repo owner/repo
  taskmd import --source github --project owner/repo --token-env GITHUB_TOKEN
  taskmd import --source github --repo owner/repo --labels bug,critical --assignee alice
  taskmd import --source jira --project PROJ --url https://company.atlassian.net
  taskmd import --source jira --project PROJ --url https://company.atlassian.net --jql "assignee = currentUser()"
  taskmd import --source github --repo owner/repo --dry-run
  taskmd import --source github --repo owner/repo --format json
  taskmd import --source github --project owner/repo --token-env GITHUB_TOKEN --filter "state:open labels:bug"`,
	Args: cobra.NoArgs,
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(&importSource, "source", "", "source name (github, jira, etc.)")
	importCmd.Flags().StringVar(&importProject, "project", "", "project identifier (owner/repo for GitHub, project key for Jira)")
	importCmd.Flags().StringVar(&importTokenEnv, "token-env", "", "environment variable name for auth token")
	importCmd.Flags().StringVar(&importUserEnv, "user-env", "", "environment variable name for username (Jira)")
	importCmd.Flags().StringVar(&importBaseURL, "base-url", "", "API base URL (for Jira or GitHub Enterprise)")
	importCmd.Flags().StringVar(&importOutDir, "output-dir", "./tasks", "target directory for imported task files")
	importCmd.Flags().StringVar(&importFilter, "filter", "", `source-specific filters as key:value pairs (e.g. "state:open labels:bug")`)
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "preview import without writing files")
	importCmd.Flags().StringVar(&importFormat, "format", "table", "output format: table, json, yaml")

	// GitHub-specific shortcut flags
	importCmd.Flags().StringVar(&importRepo, "repo", "", "alias for --project (owner/repo for GitHub)")
	importCmd.Flags().StringVar(&importLabels, "labels", "", "filter by labels (comma-separated, GitHub)")
	importCmd.Flags().StringVar(&importMilestone, "milestone", "", "filter by milestone (GitHub)")
	importCmd.Flags().StringVar(&importAssignee, "assignee", "", "filter by assignee (GitHub)")

	// Jira-specific shortcut flags
	importCmd.Flags().StringVar(&importURL, "url", "", "alias for --base-url (Jira instance URL)")
	importCmd.Flags().StringVar(&importJQL, "jql", "", "Jira Query Language filter")
}

func runImport(_ *cobra.Command, _ []string) error {
	if err := ValidateFormat(importFormat, []string{"table", "json", "yaml"}); err != nil {
		return err
	}

	flags := GetGlobalFlags()

	var cfg sync.ImportConfig
	var err error

	if importSource != "" {
		cfg, err = buildImportConfigFromFlags()
	} else {
		cfg, err = runImportWizard()
	}
	if err != nil {
		return err
	}

	cfg.DryRun = importDryRun
	cfg.Verbose = flags.Verbose

	if !flags.Quiet && importFormat == "table" {
		fmt.Fprintf(os.Stderr, "Importing from %s...\n", cfg.SourceName)
		if cfg.DryRun {
			fmt.Fprintln(os.Stderr, "(dry-run mode: no files will be written)")
		}
	}

	result, err := sync.RunImport(cfg)
	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	return printImportResult(result, importFormat, flags.Quiet)
}

func buildImportConfigFromFlags() (sync.ImportConfig, error) {
	project := importProject
	tokenEnv := importTokenEnv

	baseURL := importBaseURL
	userEnv := importUserEnv

	// GitHub-specific: --repo is an alias for --project
	if importSource == "github" {
		if project == "" && importRepo != "" {
			project = importRepo
		}
		if tokenEnv == "" {
			tokenEnv = "GITHUB_TOKEN"
		}
	}

	// Jira-specific: --url is an alias for --base-url, default env vars
	if importSource == "jira" {
		if baseURL == "" && importURL != "" {
			baseURL = importURL
		}
		if tokenEnv == "" {
			tokenEnv = "JIRA_TOKEN"
		}
		if userEnv == "" {
			userEnv = "JIRA_USER"
		}
	}

	srcCfg := sync.SourceConfig{
		Name:     importSource,
		Project:  project,
		TokenEnv: tokenEnv,
		UserEnv:  userEnv,
		BaseURL:  baseURL,
	}

	if importFilter != "" {
		srcCfg.Filters = parseImportFilters(importFilter)
	}

	// GitHub-specific: merge shortcut flags into filters
	if importSource == "github" {
		srcCfg.Filters = mergeGitHubShortcutFlags(srcCfg.Filters)
	}

	// Jira-specific: merge --jql into filters
	if importSource == "jira" {
		srcCfg.Filters = mergeJiraShortcutFlags(srcCfg.Filters)
	}

	return sync.ImportConfig{
		SourceName: importSource,
		SourceCfg:  srcCfg,
		OutputDir:  importOutDir,
		ScanDir:    ".",
	}, nil
}

// mergeGitHubShortcutFlags merges --labels, --milestone, --assignee flags
// into the filter map and defaults state to "open" if not set.
func mergeGitHubShortcutFlags(filters map[string]any) map[string]any {
	if filters == nil {
		filters = make(map[string]any)
	}

	if importLabels != "" {
		if _, ok := filters["labels"]; !ok {
			filters["labels"] = importLabels
		}
	}
	if importMilestone != "" {
		if _, ok := filters["milestone"]; !ok {
			filters["milestone"] = importMilestone
		}
	}
	if importAssignee != "" {
		if _, ok := filters["assignee"]; !ok {
			filters["assignee"] = importAssignee
		}
	}

	// Default to open issues for import
	if _, ok := filters["state"]; !ok {
		filters["state"] = "open"
	}

	return filters
}

// mergeJiraShortcutFlags merges --jql flag into the filter map.
func mergeJiraShortcutFlags(filters map[string]any) map[string]any {
	if importJQL == "" {
		return filters
	}
	if filters == nil {
		filters = make(map[string]any)
	}
	if _, ok := filters["jql"]; !ok {
		filters["jql"] = importJQL
	}
	return filters
}

func runImportWizard() (sync.ImportConfig, error) {
	sourceName, err := wizardSelectSource()
	if err != nil {
		return sync.ImportConfig{}, err
	}

	srcCfg, err := wizardSourceConfig(sourceName)
	if err != nil {
		return sync.ImportConfig{}, err
	}

	outDir, err := wizardOutputDir()
	if err != nil {
		return sync.ImportConfig{}, err
	}

	if err := wizardConfirm(); err != nil {
		return sync.ImportConfig{}, err
	}

	return sync.ImportConfig{
		SourceName: sourceName,
		SourceCfg:  srcCfg,
		OutputDir:  outDir,
		ScanDir:    ".",
	}, nil
}

func wizardSelectSource() (string, error) {
	names := sync.RegisteredNames()
	if len(names) == 0 {
		return "", fmt.Errorf("no sources registered")
	}

	options := make([]huh.Option[string], len(names))
	for i, n := range names {
		options[i] = huh.NewOption(n, n)
	}

	var sourceName string
	err := huh.NewSelect[string]().
		Title("Where are your tasks?").
		Options(options...).
		Value(&sourceName).
		Run()
	if err != nil {
		return "", fmt.Errorf("wizard cancelled: %w", err)
	}
	return sourceName, nil
}

func wizardSourceConfig(sourceName string) (sync.SourceConfig, error) {
	cfg := sync.SourceConfig{Name: sourceName}

	if sourceName == "github" {
		return wizardGitHubConfig(cfg)
	}

	if err := promptInput("Project identifier", projectHint(sourceName), "", &cfg.Project); err != nil {
		return cfg, err
	}
	if err := promptInput("Auth token env var", "Name of the environment variable holding your API token", "GITHUB_TOKEN", &cfg.TokenEnv); err != nil {
		return cfg, err
	}
	if sourceName == "jira" {
		if err := promptInput("Base URL", "Your Jira instance URL (e.g. https://company.atlassian.net)", "", &cfg.BaseURL); err != nil {
			return cfg, err
		}
		if err := promptInput("Username env var", "Name of the environment variable holding your Jira username", "JIRA_USER", &cfg.UserEnv); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func wizardGitHubConfig(cfg sync.SourceConfig) (sync.SourceConfig, error) {
	// Auto-detect repo from git remote
	detected := github.DetectRepo()
	placeholder := "owner/repo"
	if detected != "" {
		cfg.Project = detected
		placeholder = detected
	}

	if err := promptInput("Repository", "GitHub repository (e.g. owner/repo)", placeholder, &cfg.Project); err != nil {
		return cfg, err
	}

	cfg.TokenEnv = "GITHUB_TOKEN"
	if err := promptInput("Auth token env var", "Name of the environment variable holding your GitHub token", "GITHUB_TOKEN", &cfg.TokenEnv); err != nil {
		return cfg, err
	}

	filters, err := wizardGitHubFilters()
	if err != nil {
		return cfg, err
	}
	cfg.Filters = filters

	return cfg, nil
}

func wizardGitHubFilters() (map[string]any, error) {
	var filterChoice string
	err := huh.NewSelect[string]().
		Title("Which issues to import?").
		Options(
			huh.NewOption("All open issues", "open"),
			huh.NewOption("Filter by label", "label"),
			huh.NewOption("Filter by assignee", "assignee"),
			huh.NewOption("Filter by milestone", "milestone"),
			huh.NewOption("All issues (open + closed)", "all"),
		).
		Value(&filterChoice).
		Run()
	if err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	filters := map[string]any{"state": "open"}

	switch filterChoice {
	case "all":
		filters["state"] = "all"
	case "label":
		var labels string
		if err := promptInput("Labels", "Comma-separated labels to filter by", "", &labels); err != nil {
			return nil, err
		}
		if labels != "" {
			filters["labels"] = labels
		}
	case "assignee":
		var assignee string
		if err := promptInput("Assignee", "GitHub username", "", &assignee); err != nil {
			return nil, err
		}
		if assignee != "" {
			filters["assignee"] = assignee
		}
	case "milestone":
		var milestone string
		if err := promptInput("Milestone", "Milestone name or number", "", &milestone); err != nil {
			return nil, err
		}
		if milestone != "" {
			filters["milestone"] = milestone
		}
	}

	return filters, nil
}

func wizardOutputDir() (string, error) {
	outDir := importOutDir
	if err := promptInput("Output directory", "", "", &outDir); err != nil {
		return "", err
	}
	return outDir, nil
}

func wizardConfirm() error {
	var confirm bool
	err := huh.NewConfirm().Title("Import tasks?").Value(&confirm).Run()
	if err != nil {
		return fmt.Errorf("wizard cancelled: %w", err)
	}
	if !confirm {
		return fmt.Errorf("import cancelled by user")
	}
	return nil
}

func promptInput(title, description, placeholder string, value *string) error {
	input := huh.NewInput().Title(title).Value(value)
	if description != "" {
		input = input.Description(description)
	}
	if placeholder != "" {
		input = input.Placeholder(placeholder)
	}
	if err := input.Run(); err != nil {
		return fmt.Errorf("wizard cancelled: %w", err)
	}
	return nil
}

func projectHint(source string) string {
	switch source {
	case "github":
		return "GitHub repository (e.g. owner/repo)"
	case "jira":
		return "Jira project key (e.g. PROJ)"
	default:
		return "Project identifier"
	}
}

// parseImportFilters parses "key:value key2:value2" into a map.
func parseImportFilters(raw string) map[string]any {
	filters := make(map[string]any)
	for _, pair := range strings.Fields(raw) {
		key, value, ok := strings.Cut(pair, ":")
		if ok && key != "" {
			filters[key] = value
		}
	}
	return filters
}

// importResultData is the structured representation for JSON/YAML output.
type importResultData struct {
	Created []importActionData `json:"created" yaml:"created"`
	Skipped []importActionData `json:"skipped" yaml:"skipped"`
	Errors  []importErrorData  `json:"errors,omitempty" yaml:"errors,omitempty"`
	Summary importSummary      `json:"summary" yaml:"summary"`
}

type importActionData struct {
	ExternalID string `json:"external_id" yaml:"external_id"`
	LocalID    string `json:"local_id,omitempty" yaml:"local_id,omitempty"`
	FilePath   string `json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Title      string `json:"title" yaml:"title"`
	Reason     string `json:"reason" yaml:"reason"`
}

type importErrorData struct {
	ExternalID string `json:"external_id" yaml:"external_id"`
	Title      string `json:"title" yaml:"title"`
	Error      string `json:"error" yaml:"error"`
}

type importSummary struct {
	Total   int `json:"total" yaml:"total"`
	Created int `json:"created" yaml:"created"`
	Skipped int `json:"skipped" yaml:"skipped"`
	Errors  int `json:"errors" yaml:"errors"`
}

func printImportResult(result *sync.ImportResult, format string, quietMode bool) error {
	data := buildImportResultData(result)

	switch format {
	case "json":
		return WriteJSON(os.Stdout, data)
	case "yaml":
		return WriteYAML(os.Stdout, data)
	default:
		return printImportTable(result, data.Summary, quietMode)
	}
}

func buildImportResultData(result *sync.ImportResult) importResultData {
	data := importResultData{
		Created: make([]importActionData, len(result.Created)),
		Skipped: make([]importActionData, len(result.Skipped)),
	}

	for i, a := range result.Created {
		data.Created[i] = importActionData{
			ExternalID: a.ExternalID,
			LocalID:    a.LocalID,
			FilePath:   a.FilePath,
			Title:      a.Title,
			Reason:     a.Reason,
		}
	}
	for i, a := range result.Skipped {
		data.Skipped[i] = importActionData{
			ExternalID: a.ExternalID,
			Title:      a.Title,
			Reason:     a.Reason,
		}
	}
	for _, e := range result.Errors {
		data.Errors = append(data.Errors, importErrorData{
			ExternalID: e.ExternalID,
			Title:      e.Title,
			Error:      e.Err.Error(),
		})
	}

	data.Summary = importSummary{
		Total:   len(result.Created) + len(result.Skipped) + len(result.Errors),
		Created: len(result.Created),
		Skipped: len(result.Skipped),
		Errors:  len(result.Errors),
	}

	return data
}

func printImportTable(result *sync.ImportResult, summary importSummary, quietMode bool) error {
	if quietMode {
		return nil
	}

	if len(result.Created) > 0 {
		fmt.Printf("  Created %d task(s):\n", len(result.Created))
		for _, a := range result.Created {
			fmt.Printf("    + [%s] %s\n", a.LocalID, a.Title)
		}
	}

	if len(result.Skipped) > 0 {
		fmt.Printf("  Skipped %d task(s) (duplicate external_id):\n", len(result.Skipped))
		for _, a := range result.Skipped {
			fmt.Printf("    - [%s] %s\n", a.ExternalID, a.Title)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Printf("  Errors %d task(s):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("    x [%s] %s: %v\n", e.ExternalID, e.Title, e.Err)
		}
	}

	fmt.Printf("  Done: %d total, %d created, %d skipped, %d errors\n",
		summary.Total, summary.Created, summary.Skipped, summary.Errors)

	return nil
}
