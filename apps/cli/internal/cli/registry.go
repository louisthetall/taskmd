package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// GlobalProjectEntry represents a registered project in the global config.
type GlobalProjectEntry struct {
	ID   string
	Name string
	Path string // Absolute, tilde-expanded
}

// globalConfig is the minimal YAML shape for reading the projects key.
type globalConfig struct {
	Projects []globalProjectYAML `yaml:"projects"`
}

type globalProjectYAML struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// LoadGlobalRegistry reads the projects list from ~/.taskmd.yaml (or $TASKMD_HOME_CONFIG).
// Returns an empty slice and nil error when the file does not exist or contains no projects.
// Does not validate that project paths exist on disk — callers decide that policy.
func LoadGlobalRegistry() ([]GlobalProjectEntry, error) {
	cfgPath, err := globalConfigPath()
	if err != nil {
		return nil, fmt.Errorf("resolve global config path: %w", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read global config %s: %w", cfgPath, err)
	}

	var cfg globalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse global config %s: %w", cfgPath, err)
	}

	if len(cfg.Projects) == 0 {
		return nil, nil
	}

	return parseProjectEntries(cfg.Projects, filepath.Dir(cfgPath))
}

// parseProjectEntries converts raw YAML entries into validated GlobalProjectEntry values.
func parseProjectEntries(projects []globalProjectYAML, cfgDir string) ([]GlobalProjectEntry, error) {
	entries := make([]GlobalProjectEntry, 0, len(projects))
	var errs []error

	for i, p := range projects {
		entry, err := resolveProjectEntry(p, i, cfgDir)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entries = append(entries, entry)
	}

	// Check for duplicate IDs.
	seen := make(map[string]int, len(entries))
	for i, e := range entries {
		if prev, ok := seen[e.ID]; ok {
			errs = append(errs, fmt.Errorf("projects[%d]: duplicate id %q (first seen at index %d)", i, e.ID, prev))
		} else {
			seen[e.ID] = i
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return entries, nil
}

// resolveProjectEntry normalizes a single YAML project entry into a GlobalProjectEntry.
func resolveProjectEntry(p globalProjectYAML, index int, cfgDir string) (GlobalProjectEntry, error) {
	if p.Path == "" {
		return GlobalProjectEntry{}, fmt.Errorf("projects[%d]: path is required", index)
	}

	resolved, err := resolvePath(p.Path, cfgDir)
	if err != nil {
		return GlobalProjectEntry{}, fmt.Errorf("projects[%d]: %w", index, err)
	}

	id := p.ID
	if id == "" {
		id = filepath.Base(resolved)
	}

	name := p.Name
	if name == "" {
		name = id
	}

	return GlobalProjectEntry{ID: id, Name: name, Path: resolved}, nil
}

// globalConfigPath returns the path to the global config file.
// Checks $TASKMD_HOME_CONFIG first, then falls back to ~/.taskmd.yaml.
func globalConfigPath() (string, error) {
	if env := os.Getenv("TASKMD_HOME_CONFIG"); env != "" {
		return expandTilde(env)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".taskmd.yaml"), nil
}

// resolvePath expands tilde and makes relative paths absolute against baseDir.
func resolvePath(path, baseDir string) (string, error) {
	expanded, err := expandTilde(path)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	return filepath.Clean(filepath.Join(baseDir, expanded)), nil
}

// expandTilde replaces a leading ~ with the user's home directory.
func expandTilde(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expand ~: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
