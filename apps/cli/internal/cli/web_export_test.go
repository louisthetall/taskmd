package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetWebExportFlags() {
	webExportOutput = "./taskmd-export"
	webExportBasePath = "/"
}

func TestWebExport_NonExistentDirectory(t *testing.T) {
	resetWebExportFlags()
	taskDir = "/nonexistent/path/that/does/not/exist"

	err := runWebExport(webExportCmd, nil)
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}

	if !strings.Contains(err.Error(), "not a valid directory") {
		t.Errorf("expected 'not a valid directory' error, got: %v", err)
	}
}

func TestWebExport_FileInsteadOfDirectory(t *testing.T) {
	resetWebExportFlags()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not-a-dir.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	taskDir = filePath

	err := runWebExport(webExportCmd, nil)
	if err == nil {
		t.Fatal("expected error when task-dir points to a file")
	}

	if !strings.Contains(err.Error(), "not a valid directory") {
		t.Errorf("expected 'not a valid directory' error, got: %v", err)
	}
}
