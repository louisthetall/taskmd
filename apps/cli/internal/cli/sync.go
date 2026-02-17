package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/driangle/taskmd/apps/cli/internal/sync"
	_ "github.com/driangle/taskmd/apps/cli/internal/sync/github" // register github sync source
	_ "github.com/driangle/taskmd/apps/cli/internal/sync/jira"   // register jira sync source
)

var (
	syncDryRun   bool
	syncSource   string
	syncConflict string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tasks from external sources",
	Long:  `Commands for syncing tasks with external sources (GitHub Issues, Jira, etc.).`,
}

var syncDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Pull tasks from external sources",
	Long: `Fetch tasks from configured external sources (GitHub Issues, Jira, etc.)
and create or update local markdown task files.

Configuration is read from .taskmd.yaml in the current directory.

Examples:
  taskmd sync down
  taskmd sync down --dry-run
  taskmd sync down --source github
  taskmd sync down --conflict remote
  taskmd sync down --conflict local`,
	Args: cobra.NoArgs,
	RunE: runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncDownCmd)

	syncDownCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "preview changes without writing files")
	syncDownCmd.Flags().StringVar(&syncSource, "source", "", "sync only the named source")
	syncDownCmd.Flags().StringVar(&syncConflict, "conflict", "skip", "conflict resolution strategy: skip, remote, local")
}

func runSync(_ *cobra.Command, _ []string) error {
	flags := GetGlobalFlags()

	switch syncConflict {
	case sync.ConflictSkip, sync.ConflictRemote, sync.ConflictLocal:
	default:
		return fmt.Errorf("invalid --conflict value %q: must be skip, remote, or local", syncConflict)
	}

	cfg, err := sync.LoadConfig(".")
	if err != nil {
		return err
	}

	engine := &sync.Engine{
		ConfigDir:        ".",
		Verbose:          flags.Verbose,
		DryRun:           syncDryRun,
		ConflictStrategy: syncConflict,
	}

	sources := cfg.Sources
	if syncSource != "" {
		sources = filterSources(cfg.Sources, syncSource)
		if len(sources) == 0 {
			return fmt.Errorf("source %q not found in config", syncSource)
		}
	}

	for _, srcCfg := range sources {
		if !flags.Quiet {
			fmt.Fprintf(os.Stderr, "Syncing from %s...\n", srcCfg.Name)
		}

		result, err := engine.RunSync(srcCfg)
		if err != nil {
			return fmt.Errorf("sync failed for %s: %w", srcCfg.Name, err)
		}

		printSyncResult(srcCfg.Name, result, flags.Quiet)
	}

	return nil
}

func filterSources(sources []sync.SourceConfig, name string) []sync.SourceConfig {
	for _, s := range sources {
		if s.Name == name {
			return []sync.SourceConfig{s}
		}
	}
	return nil
}

func printSyncResult(_ string, result *sync.SyncResult, quietMode bool) {
	if quietMode {
		return
	}

	if len(result.Created) > 0 {
		fmt.Printf("  Created %d task(s):\n", len(result.Created))
		for _, a := range result.Created {
			fmt.Printf("    + [%s] %s\n", a.LocalID, a.Title)
		}
	}

	if len(result.Updated) > 0 {
		fmt.Printf("  Updated %d task(s):\n", len(result.Updated))
		for _, a := range result.Updated {
			fmt.Printf("    ~ [%s] %s\n", a.LocalID, a.Title)
		}
	}

	if len(result.Conflicts) > 0 {
		fmt.Printf("  Conflicts %d task(s) (skipped, local changes detected):\n", len(result.Conflicts))
		for _, a := range result.Conflicts {
			fmt.Printf("    ! [%s] %s\n", a.LocalID, a.Title)
		}
	}

	if len(result.Skipped) > 0 {
		fmt.Printf("  Skipped %d task(s) (no changes)\n", len(result.Skipped))
	}

	if len(result.Errors) > 0 {
		fmt.Printf("  Errors %d task(s):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("    x [%s] %s: %v\n", e.ExternalID, e.Title, e.Err)
		}
	}

	total := len(result.Created) + len(result.Updated) + len(result.Skipped) + len(result.Conflicts)
	fmt.Printf("  Done: %d total, %d created, %d updated, %d skipped, %d conflicts\n",
		total, len(result.Created), len(result.Updated), len(result.Skipped), len(result.Conflicts))
}
