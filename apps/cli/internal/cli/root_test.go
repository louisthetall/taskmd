package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestResolveRelativeToConfig_AbsolutePath(t *testing.T) {
	// Absolute paths should be returned unchanged.
	abs := "/absolute/path/to/tasks"
	result := resolveRelativeToConfig(abs)
	if result != abs {
		t.Errorf("expected %q, got %q", abs, result)
	}
}

func TestResolveRelativeToConfig_NoConfigFile(t *testing.T) {
	// When no config file is loaded, relative paths are returned unchanged.
	viper.Reset()
	rel := "my-tasks"
	result := resolveRelativeToConfig(rel)
	if result != rel {
		t.Errorf("expected %q, got %q", rel, result)
	}
}

func TestResolveRelativeToConfig_RelativePath(t *testing.T) {
	// When a config file is loaded, relative paths resolve against the config dir.
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".taskmd.yaml")
	if err := os.WriteFile(configPath, []byte("task-dir: tasks\n"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	viper.Reset()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	result := resolveRelativeToConfig("tasks")
	expected := filepath.Join(tmpDir, "tasks")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveRelativeToConfig_DotPath(t *testing.T) {
	// "." should resolve to the config file's directory.
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".taskmd.yaml")
	if err := os.WriteFile(configPath, []byte("task-dir: .\n"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	viper.Reset()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	result := resolveRelativeToConfig(".")
	expected := tmpDir
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
