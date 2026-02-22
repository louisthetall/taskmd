package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetWebFlags() {
	webPort = 8080
	webDev = false
	webOpen = false
	webReadOnly = false
}

func TestWebStart_NonExistentDirectory(t *testing.T) {
	resetWebFlags()
	taskDir = "/nonexistent/path/that/does/not/exist"

	err := runWebStart(webStartCmd, nil)
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}

	if !strings.Contains(err.Error(), "not a valid directory") {
		t.Errorf("expected 'not a valid directory' error, got: %v", err)
	}
}

func TestWebStart_FileInsteadOfDirectory(t *testing.T) {
	resetWebFlags()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not-a-dir.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	taskDir = filePath

	err := runWebStart(webStartCmd, nil)
	if err == nil {
		t.Fatal("expected error when task-dir points to a file")
	}

	if !strings.Contains(err.Error(), "not a valid directory") {
		t.Errorf("expected 'not a valid directory' error, got: %v", err)
	}
}
