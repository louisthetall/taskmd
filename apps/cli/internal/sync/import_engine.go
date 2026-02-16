package sync

import (
	"fmt"
	"os"
	"strings"

	"github.com/driangle/taskmd/apps/cli/internal/nextid"
	"github.com/driangle/taskmd/apps/cli/internal/scanner"
)

// ImportConfig holds configuration for a one-shot import run.
type ImportConfig struct {
	SourceName string
	SourceCfg  SourceConfig
	OutputDir  string
	ScanDir    string
	DryRun     bool
	Verbose    bool
}

// ImportAction describes a single import operation.
type ImportAction struct {
	ExternalID string
	LocalID    string
	FilePath   string
	Title      string
	Reason     string // "created" or "skipped_duplicate"
}

// ImportResult holds the outcome of an import run.
type ImportResult struct {
	Created []ImportAction
	Skipped []ImportAction
	Errors  []SyncError
}

// RunImport fetches tasks from an external source and writes them as local
// task files. Unlike RunSync, it performs no state tracking — it is a one-shot
// operation that detects duplicates via external_id in existing task files.
func RunImport(cfg ImportConfig) (*ImportResult, error) {
	externalTasks, err := fetchImportTasks(cfg)
	if err != nil {
		return nil, err
	}

	existingIDs, externalIDs, err := scanExistingTasks(cfg.ScanDir, cfg.Verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to scan existing tasks: %w", err)
	}

	fieldMap := resolveImportFieldMap(cfg.SourceName, cfg.SourceCfg.FieldMap)
	result := &ImportResult{}

	for _, ext := range externalTasks {
		action, ids := importTask(ext, cfg, fieldMap, existingIDs, externalIDs)
		existingIDs = ids

		switch action.Reason {
		case "skipped_duplicate":
			result.Skipped = append(result.Skipped, action.ImportAction)
		case "error":
			result.Errors = append(result.Errors, SyncError{
				ExternalID: ext.ExternalID, Title: ext.Title, Err: action.err,
			})
		default:
			result.Created = append(result.Created, action.ImportAction)
			externalIDs[ext.ExternalID] = true
		}
	}

	return result, nil
}

func fetchImportTasks(cfg ImportConfig) ([]ExternalTask, error) {
	src, err := GetSource(cfg.SourceName)
	if err != nil {
		return nil, err
	}
	if err := src.ValidateConfig(cfg.SourceCfg); err != nil {
		return nil, fmt.Errorf("invalid config for source %q: %w", cfg.SourceName, err)
	}
	return src.FetchTasks(cfg.SourceCfg)
}

func resolveImportFieldMap(sourceName string, userMap FieldMap) FieldMap {
	fm := defaultImportFieldMap()
	if sourceName == "jira" {
		fm = jiraDefaultFieldMap()
	}
	if len(userMap.Status) > 0 {
		fm.Status = userMap.Status
	}
	if len(userMap.Priority) > 0 {
		fm.Priority = userMap.Priority
	}
	return fm
}

// importActionInternal extends ImportAction with an error field for internal use.
type importActionInternal struct {
	ImportAction
	err error
}

func importTask(
	ext ExternalTask, cfg ImportConfig, fieldMap FieldMap,
	existingIDs []string, externalIDs map[string]bool,
) (importActionInternal, []string) {
	if externalIDs[ext.ExternalID] {
		return importActionInternal{ImportAction: ImportAction{
			ExternalID: ext.ExternalID, Title: ext.Title, Reason: "skipped_duplicate",
		}}, existingIDs
	}

	mapped := MapExternalTask(ext, fieldMap)
	mapped.Description = appendSourceURL(mapped.Description, ext.URL)

	newID := nextid.Calculate(existingIDs).NextID
	action := importActionInternal{ImportAction: ImportAction{
		ExternalID: ext.ExternalID, LocalID: newID, Title: ext.Title, Reason: "created",
	}}

	if !cfg.DryRun {
		filePath, err := WriteTaskFile(cfg.OutputDir, newID, mapped, ext.ExternalID, cfg.SourceName)
		if err != nil {
			action.Reason = "error"
			action.err = err
			return action, existingIDs
		}
		action.FilePath = filePath
	}

	return action, append(existingIDs, newID)
}

// defaultImportFieldMap returns sensible defaults for one-shot imports (GitHub-oriented).
func defaultImportFieldMap() FieldMap {
	return FieldMap{
		Status: map[string]string{
			"open":   "pending",
			"closed": "completed",
		},
		Priority: map[string]string{
			"critical": "critical",
			"high":     "high",
			"medium":   "medium",
			"low":      "low",
		},
		LabelsToTags:    true,
		AssigneeToOwner: true,
	}
}

// jiraDefaultFieldMap returns Jira-specific defaults for status and priority mappings.
func jiraDefaultFieldMap() FieldMap {
	return FieldMap{
		Status: map[string]string{
			"To Do":       "pending",
			"In Progress": "in-progress",
			"Done":        "completed",
			"open":        "pending",
			"closed":      "completed",
		},
		Priority: map[string]string{
			"Highest": "critical",
			"High":    "high",
			"Medium":  "medium",
			"Low":     "low",
			"Lowest":  "low",
		},
		LabelsToTags:    true,
		AssigneeToOwner: true,
	}
}

// scanExistingTasks scans the directory for task IDs and external_ids.
func scanExistingTasks(dir string, verbose bool) (ids []string, externalIDs map[string]bool, err error) {
	externalIDs = make(map[string]bool)

	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		return nil, externalIDs, nil
	}

	taskScanner := scanner.NewScanner(dir, verbose, nil)
	result, err := taskScanner.Scan()
	if err != nil {
		return nil, nil, err
	}

	ids = make([]string, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		ids = append(ids, t.ID)
		if t.ExternalID != "" {
			externalIDs[t.ExternalID] = true
		}
	}
	return ids, externalIDs, nil
}

// appendSourceURL appends a source reference to the task description.
func appendSourceURL(description, url string) string {
	if url == "" {
		return description
	}
	suffix := "\n\n---\nSource: " + url
	return strings.TrimRight(description, "\n") + suffix
}
