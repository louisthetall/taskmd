package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	projectInitForce   bool
	projectInitStdout  bool
	projectInitClaude  bool
	projectInitGemini  bool
	projectInitCodex   bool
	projectInitNoSpec  bool
	projectInitNoAgent bool
	projectInitTaskDir string
)

// projectInitRoot is the project root directory. Defaults to ".".
// Tests override this to t.TempDir() for parallel safety.
var projectInitRoot = "."

// projectInitIsTTY checks whether stdin is a terminal.
// Tests override this to return false.
var projectInitIsTTY = func() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

const configFilename = ".taskmd.yaml"

var projectInitCmd = &cobra.Command{
	Use:        "init",
	SuggestFor: []string{"setup", "create", "new"},
	Short:      "Initialize a taskmd project with config, task directory, and agent files",
	Long: `Initialize sets up a complete taskmd project in the current directory.

Creates a task directory, .taskmd.yaml config, agent configuration files, and
the taskmd specification document. When run interactively, prompts for any
values not provided via flags.

If a file already exists and --force is not set, it is skipped with a warning.

Examples:
  taskmd init                        # Interactive setup (prompts for missing info)
  taskmd init --task-dir ./tasks     # Set task directory, prompt for agents
  taskmd init --claude               # Claude agent, prompt for task directory
  taskmd init --task-dir ./tasks --claude  # Fully non-interactive
  taskmd init --claude --gemini      # Multiple agents
  taskmd init --no-spec              # Skip TASKMD_SPEC.md
  taskmd init --no-agent             # Skip agent configs
  taskmd init --force                # Overwrite existing files
  taskmd init --stdout               # Print all content to stdout`,
	Args: cobra.NoArgs,
	RunE: runProjectInit,
}

func init() {
	rootCmd.AddCommand(projectInitCmd)

	projectInitCmd.Flags().BoolVar(&projectInitForce, "force", false, "overwrite existing files")
	projectInitCmd.Flags().BoolVar(&projectInitStdout, "stdout", false, "print all content to stdout instead of writing files")
	projectInitCmd.Flags().BoolVar(&projectInitClaude, "claude", false, "initialize for Claude Code")
	projectInitCmd.Flags().BoolVar(&projectInitGemini, "gemini", false, "initialize for Gemini")
	projectInitCmd.Flags().BoolVar(&projectInitCodex, "codex", false, "initialize for Codex")
	projectInitCmd.Flags().BoolVar(&projectInitNoSpec, "no-spec", false, "skip generating TASKMD_SPEC.md")
	projectInitCmd.Flags().BoolVar(&projectInitNoAgent, "no-agent", false, "skip generating agent configuration files")
	projectInitCmd.Flags().StringVar(&projectInitTaskDir, "task-dir", "./tasks", "task directory path to create")
}

// fileToWrite represents a file that the init command will create.
type fileToWrite struct {
	filename string
	content  []byte
}

func runProjectInit(cmd *cobra.Command, _ []string) error {
	if projectInitNoSpec && projectInitNoAgent {
		return fmt.Errorf("--no-spec and --no-agent cannot both be set (nothing to do)")
	}

	root := projectInitRoot
	isTTY := projectInitIsTTY()
	quiet := GetGlobalFlags().Quiet

	// Resolve task directory: flag > prompt > default
	taskDirPath, err := resolveInitTaskDir(cmd, isTTY)
	if err != nil {
		return err
	}

	// Resolve agents: flags > prompt > default (Claude)
	resolveInitAgents(isTTY)

	// Collect files split by destination
	rootFiles, taskDirFiles := collectInitFiles()

	// --stdout mode: print everything and exit
	if projectInitStdout {
		allFiles := append(rootFiles, taskDirFiles...)
		return printFilesToStdout(allFiles)
	}

	taskDirAbs := taskDirPath
	if !filepath.IsAbs(taskDirAbs) {
		taskDirAbs = filepath.Join(root, taskDirPath)
	}

	return writeProjectFiles(root, taskDirAbs, taskDirPath, rootFiles, taskDirFiles, quiet)
}

// writeProjectFiles creates directories, config, and all init files.
func writeProjectFiles(root, taskDirAbs, taskDirPath string, rootFiles, taskDirFiles []fileToWrite, quiet bool) error {
	var createdPaths []string

	dirCreated, err := ensureTaskDir(taskDirAbs, quiet)
	if err != nil {
		return err
	}
	if dirCreated {
		abs, _ := filepath.Abs(taskDirAbs)
		createdPaths = append(createdPaths, abs+"/")
	}

	configCreated, err := writeConfigFile(root, taskDirPath, quiet)
	if err != nil {
		return err
	}
	if configCreated {
		abs, _ := filepath.Abs(filepath.Join(root, configFilename))
		createdPaths = append(createdPaths, abs)
	}

	rootCreated, err := writeInitFiles(root, rootFiles, quiet)
	if err != nil {
		return err
	}
	createdPaths = append(createdPaths, rootCreated...)

	tdCreated, err := writeInitFiles(taskDirAbs, taskDirFiles, quiet)
	if err != nil {
		return err
	}
	createdPaths = append(createdPaths, tdCreated...)

	if !quiet {
		printInitSummary(createdPaths)
	}

	return nil
}

// resolveInitTaskDir returns the task directory path.
// If --task-dir was explicitly provided, uses that.
// If TTY, prompts the user. Otherwise uses the default.
func resolveInitTaskDir(cmd *cobra.Command, isTTY bool) (string, error) {
	if cmd.Flags().Changed("task-dir") {
		return projectInitTaskDir, nil
	}

	if isTTY {
		value := projectInitTaskDir // default for the prompt
		err := huh.NewInput().
			Title("Task directory").
			Value(&value).
			Run()
		if err != nil {
			return "", fmt.Errorf("prompt cancelled: %w", err)
		}
		return value, nil
	}

	return projectInitTaskDir, nil
}

// resolveInitAgents sets agent flags via prompt if none were explicitly set.
func resolveInitAgents(isTTY bool) {
	// If any agent flag is set, respect it
	if projectInitClaude || projectInitGemini || projectInitCodex || projectInitNoAgent {
		return
	}

	if isTTY {
		promptAgentSelection()
		return
	}

	// Non-TTY with no agent flags: default to Claude
	projectInitClaude = true
}

// promptAgentSelection shows a multi-select for agent configs.
func promptAgentSelection() {
	options := []huh.Option[string]{
		huh.NewOption("Claude Code", "claude").Selected(true),
		huh.NewOption("Gemini", "gemini"),
		huh.NewOption("Codex", "codex"),
	}

	var selected []string
	err := huh.NewMultiSelect[string]().
		Title("Which AI assistants do you use?").
		Options(options...).
		Value(&selected).
		Run()
	if err != nil || len(selected) == 0 {
		// Cancelled or nothing selected: default to Claude
		projectInitClaude = true
		return
	}

	for _, s := range selected {
		switch s {
		case "claude":
			projectInitClaude = true
		case "gemini":
			projectInitGemini = true
		case "codex":
			projectInitCodex = true
		}
	}
}

// collectInitFiles returns files split into root (agent configs) and task dir (spec).
func collectInitFiles() (rootFiles, taskDirFiles []fileToWrite) {
	if !projectInitNoAgent {
		agents := getProjectInitAgents()
		for _, agent := range agents {
			rootFiles = append(rootFiles, fileToWrite{
				filename: agent.filename,
				content:  agent.template,
			})
		}
	}

	if !projectInitNoSpec {
		taskDirFiles = append(taskDirFiles, fileToWrite{
			filename: specFilename,
			content:  specTemplate,
		})
	}

	return rootFiles, taskDirFiles
}

func getProjectInitAgents() []agentConfig {
	var agents []agentConfig

	// If no agent flags specified, default to Claude
	if !projectInitClaude && !projectInitGemini && !projectInitCodex {
		projectInitClaude = true
	}

	if projectInitClaude {
		agents = append(agents, agentConfig{
			name:     "Claude Code",
			filename: "CLAUDE.md",
			template: claudeTemplate,
		})
	}

	if projectInitGemini {
		agents = append(agents, agentConfig{
			name:     "Gemini",
			filename: "GEMINI.md",
			template: geminiTemplate,
		})
	}

	if projectInitCodex {
		agents = append(agents, agentConfig{
			name:     "Codex",
			filename: "AGENTS.md",
			template: codexTemplate,
		})
	}

	return agents
}

// ensureTaskDir creates the task directory if it doesn't exist.
func ensureTaskDir(path string, quiet bool) (created bool, err error) {
	info, statErr := os.Stat(path)
	if statErr == nil {
		if !info.IsDir() {
			return false, fmt.Errorf("not a directory: %s", path)
		}
		if !quiet {
			fmt.Fprintf(os.Stderr, "Task directory already exists: %s\n", path)
		}
		return false, nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return false, fmt.Errorf("failed to create task directory: %w", err)
	}
	return true, nil
}

// writeConfigFile writes .taskmd.yaml to the project root.
func writeConfigFile(root, taskDirPath string, quiet bool) (created bool, err error) {
	configPath := filepath.Join(root, configFilename)

	if !projectInitForce {
		if _, err := os.Stat(configPath); err == nil {
			if !quiet {
				abs, _ := filepath.Abs(configPath)
				fmt.Fprintf(os.Stderr, "Skipped %s (already exists, use --force to overwrite)\n", abs)
			}
			return false, nil
		}
	}

	content := fmt.Sprintf("dir: %s\n", taskDirPath)
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("failed to write %s: %w", configFilename, err)
	}
	return true, nil
}

// writeInitFiles writes files to a directory, returning created paths.
func writeInitFiles(dir string, files []fileToWrite, quiet bool) ([]string, error) {
	var created []string

	for _, f := range files {
		absPath, skipped, err := writeInitFile(dir, f)
		if err != nil {
			return created, err
		}
		if skipped {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Skipped %s (already exists, use --force to overwrite)\n", absPath)
			}
			continue
		}
		created = append(created, absPath)
	}

	return created, nil
}

func writeInitFile(targetDir string, f fileToWrite) (absPath string, skipped bool, err error) {
	outputPath := filepath.Join(targetDir, f.filename)
	absPath, err = filepath.Abs(outputPath)
	if err != nil {
		absPath = outputPath
	}

	if !projectInitForce {
		if _, err := os.Stat(outputPath); err == nil {
			return absPath, true, nil
		}
	}

	if err := os.WriteFile(outputPath, f.content, 0644); err != nil {
		return absPath, false, fmt.Errorf("failed to write %s: %w", f.filename, err)
	}

	return absPath, false, nil
}

// printInitSummary prints the list of created files and next steps.
func printInitSummary(createdPaths []string) {
	if len(createdPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Nothing to create (everything already exists).")
		return
	}

	fmt.Println("\nCreated:")
	for _, p := range createdPaths {
		fmt.Printf("  %s\n", p)
	}
	fmt.Println("\nYou're ready! Try:")
	fmt.Println("  taskmd add \"My first task\"")
	fmt.Println("  taskmd list")
	fmt.Println("  taskmd web start --open")
}

func printFilesToStdout(files []fileToWrite) error {
	for i, f := range files {
		if i > 0 {
			fmt.Print("\n---\n")
			fmt.Printf("# %s\n", f.filename)
			fmt.Print("---\n\n")
		}
		fmt.Print(string(f.content))
	}
	return nil
}
