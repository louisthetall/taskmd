package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetProjectInitFlags(tmpDir string) {
	projectInitForce = false
	projectInitStdout = false
	projectInitClaude = false
	projectInitGemini = false
	projectInitCodex = false
	projectInitNoSpec = false
	projectInitNoAgent = false
	projectInitNoTemplates = false
	projectInitTaskDir = tmpDir
	projectInitIDStrategy = ""
	projectInitIDPrefix = ""
	projectInitRoot = tmpDir
	projectInitIsTTY = func() bool { return false }
	taskDir = tmpDir
	// Reset cobra flag changed state so resolveInitIDStrategy sees them as unset
	if f := projectInitCmd.Flags().Lookup("id-strategy"); f != nil {
		f.Changed = false
	}
	if f := projectInitCmd.Flags().Lookup("id-prefix"); f != nil {
		f.Changed = false
	}
}

func TestProjectInit_DefaultWritesBothFiles(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create CLAUDE.md (default agent) in root
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	content, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	if !bytes.Equal(content, claudeTemplate) {
		t.Error("CLAUDE.md content does not match template")
	}

	// Should create TASKMD_SPEC.md in task dir (same as root in this test)
	specPath := filepath.Join(tmpDir, specFilename)
	content, err = os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", specFilename, err)
	}
	if !bytes.Equal(content, initSpecTemplate) {
		t.Error("TASKMD_SPEC.md content does not match template")
	}
}

func TestProjectInit_GeminiFlag(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitGemini = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create GEMINI.md
	geminiPath := filepath.Join(tmpDir, "GEMINI.md")
	content, err := os.ReadFile(geminiPath)
	if err != nil {
		t.Fatalf("failed to read GEMINI.md: %v", err)
	}
	if !bytes.Equal(content, geminiTemplate) {
		t.Error("GEMINI.md content does not match template")
	}

	// Should create TASKMD_SPEC.md
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created")
	}

	// Should NOT create CLAUDE.md
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("CLAUDE.md should not have been created when --gemini is specified")
	}
}

func TestProjectInit_MultipleAgentFlags(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitGemini = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create CLAUDE.md, GEMINI.md, and TASKMD_SPEC.md
	for _, name := range []string{"CLAUDE.md", "GEMINI.md", specFilename} {
		path := filepath.Join(tmpDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s should have been created", name)
		}
	}
}

func TestProjectInit_NoSpecFlag(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitNoSpec = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create CLAUDE.md
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should have been created")
	}

	// Should NOT create TASKMD_SPEC.md
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); err == nil {
		t.Error("TASKMD_SPEC.md should not have been created with --no-spec")
	}
}

func TestProjectInit_NoAgentFlag(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitNoAgent = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create TASKMD_SPEC.md
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created")
	}

	// Should NOT create CLAUDE.md
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("CLAUDE.md should not have been created with --no-agent")
	}
}

func TestProjectInit_NoSpecAndNoAgentAndNoTemplatesIsError(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitNoSpec = true
	projectInitNoAgent = true
	projectInitNoTemplates = true

	err := runProjectInit(projectInitCmd, []string{})
	if err == nil {
		t.Fatal("expected error when all --no-* flags are set")
	}

	if !strings.Contains(err.Error(), "nothing to do") {
		t.Errorf("error message %q should contain 'nothing to do'", err.Error())
	}
}

func TestProjectInit_ForceOverwritesExistingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitForce = true

	// Create existing files
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	specPath := filepath.Join(tmpDir, specFilename)
	os.WriteFile(claudePath, []byte("old claude"), 0644)
	os.WriteFile(specPath, []byte("old spec"), 0644)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were overwritten
	claudeContent, _ := os.ReadFile(claudePath)
	if !bytes.Equal(claudeContent, claudeTemplate) {
		t.Error("CLAUDE.md should have been overwritten")
	}

	specContent, _ := os.ReadFile(specPath)
	if !bytes.Equal(specContent, initSpecTemplate) {
		t.Error("TASKMD_SPEC.md should have been overwritten")
	}
}

func TestProjectInit_ExistingFilesSkippedWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)

	// Create an existing CLAUDE.md
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	os.WriteFile(claudePath, []byte("existing claude"), 0644)

	// Capture stderr for the warning
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("expected no error when skipping existing files, got: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	// Should have a warning about skipping
	if !strings.Contains(stderrOutput, "Skipped") {
		t.Errorf("expected skip warning on stderr, got: %q", stderrOutput)
	}

	// Original file should be unchanged
	content, _ := os.ReadFile(claudePath)
	if string(content) != "existing claude" {
		t.Error("existing CLAUDE.md should not have been overwritten")
	}

	// TASKMD_SPEC.md should still be created (it didn't exist)
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created even though CLAUDE.md was skipped")
	}
}

func TestProjectInit_DirFlag(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "my-project")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	resetProjectInitFlags(subDir)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claudePath := filepath.Join(subDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should have been created in subdirectory")
	}

	specPath := filepath.Join(subDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created in subdirectory")
	}
}

func TestProjectInit_StdoutPrintsWithoutCreatingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitStdout = true

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain both agent template and spec template
	if !strings.Contains(output, string(claudeTemplate)) {
		t.Error("stdout output should contain Claude template")
	}
	if !strings.Contains(output, string(initSpecTemplate)) {
		t.Error("stdout output should contain spec template")
	}

	// No files should have been created
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("CLAUDE.md should not have been created with --stdout")
	}
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); err == nil {
		t.Error("TASKMD_SPEC.md should not have been created with --stdout")
	}
}

func TestProjectInit_CodexFlag(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitCodex = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create AGENTS.md and TASKMD_SPEC.md
	agentsPath := filepath.Join(tmpDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		t.Error("AGENTS.md should have been created")
	}

	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created")
	}

	// Should NOT create CLAUDE.md
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("CLAUDE.md should not have been created when --codex is specified")
	}
}

func TestProjectInit_AllAgentFlags(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitGemini = true
	projectInitCodex = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"CLAUDE.md", "GEMINI.md", "AGENTS.md", specFilename}
	for _, name := range expected {
		path := filepath.Join(tmpDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s should have been created", name)
		}
	}
}

func TestProjectInit_PartialSkipStillCreatesOthers(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitGemini = true

	// Create only CLAUDE.md as existing
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	os.WriteFile(claudePath, []byte("existing"), 0644)

	// Suppress stderr warnings
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// CLAUDE.md should be unchanged (skipped)
	content, _ := os.ReadFile(claudePath)
	if string(content) != "existing" {
		t.Error("existing CLAUDE.md should not have been overwritten")
	}

	// GEMINI.md should be created
	geminiPath := filepath.Join(tmpDir, "GEMINI.md")
	if _, err := os.Stat(geminiPath); os.IsNotExist(err) {
		t.Error("GEMINI.md should have been created")
	}

	// TASKMD_SPEC.md should be created
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created")
	}
}

// --- New tests for interactive init ---

func TestProjectInit_SeparateDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	rootDir := tmpDir
	taskDirPath := filepath.Join(tmpDir, "my-tasks")

	resetProjectInitFlags(taskDirPath)
	projectInitRoot = rootDir
	projectInitTaskDir = taskDirPath
	projectInitClaude = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Agent config should be in task directory
	claudePath := filepath.Join(taskDirPath, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should have been created in task directory")
	}

	// Agent config should NOT be in project root
	claudeInRoot := filepath.Join(rootDir, "CLAUDE.md")
	if _, err := os.Stat(claudeInRoot); err == nil {
		t.Error("CLAUDE.md should not have been created in project root")
	}

	// Spec should be in task directory
	specPath := filepath.Join(taskDirPath, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created in task directory")
	}

	// Spec should NOT be in project root
	specInRoot := filepath.Join(rootDir, specFilename)
	if _, err := os.Stat(specInRoot); err == nil {
		t.Error("TASKMD_SPEC.md should not have been created in project root")
	}
}

func TestProjectInit_ConfigFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configPath := filepath.Join(tmpDir, configFilename)
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", configFilename, err)
	}

	expected := "dir: " + tmpDir + "\n"
	if string(content) != expected {
		t.Errorf("config content = %q, want %q", string(content), expected)
	}
}

func TestProjectInit_NonTTY_DefaultsClaude(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	// No agent flags set, non-TTY (default from reset)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should create CLAUDE.md (the non-TTY default)
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should have been created as non-TTY default")
	}

	// Should NOT create GEMINI.md or AGENTS.md
	geminiPath := filepath.Join(tmpDir, "GEMINI.md")
	if _, err := os.Stat(geminiPath); err == nil {
		t.Error("GEMINI.md should not have been created")
	}
	agentsPath := filepath.Join(tmpDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); err == nil {
		t.Error("AGENTS.md should not have been created")
	}
}

func TestProjectInit_ExistingConfig_SkippedWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true

	// Create existing config
	configPath := filepath.Join(tmpDir, configFilename)
	os.WriteFile(configPath, []byte("dir: ./old-tasks\n"), 0644)

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	// Should warn about skipping config
	if !strings.Contains(stderrOutput, "Skipped") || !strings.Contains(stderrOutput, configFilename) {
		t.Errorf("expected skip warning for %s, got: %q", configFilename, stderrOutput)
	}

	// Config should be unchanged
	content, _ := os.ReadFile(configPath)
	if string(content) != "dir: ./old-tasks\n" {
		t.Error("existing config should not have been overwritten")
	}
}

func TestProjectInit_ExistingConfig_OverwrittenWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitForce = true
	projectInitClaude = true

	// Create existing config
	configPath := filepath.Join(tmpDir, configFilename)
	os.WriteFile(configPath, []byte("dir: ./old-tasks\n"), 0644)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Config should be overwritten
	content, _ := os.ReadFile(configPath)
	expected := "dir: " + tmpDir + "\n"
	if string(content) != expected {
		t.Errorf("config content = %q, want %q", string(content), expected)
	}
}

func TestProjectInit_ExistingTaskDir_Graceful(t *testing.T) {
	tmpDir := t.TempDir()
	taskDirPath := filepath.Join(tmpDir, "tasks")
	os.MkdirAll(taskDirPath, 0755)

	resetProjectInitFlags(taskDirPath)
	projectInitRoot = tmpDir
	projectInitTaskDir = taskDirPath
	projectInitClaude = true

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	// Should note existing directory
	if !strings.Contains(stderrOutput, "already exists") {
		t.Errorf("expected 'already exists' note, got: %q", stderrOutput)
	}

	// Spec should still be created in the existing task dir
	specPath := filepath.Join(taskDirPath, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created in existing task directory")
	}
}

func TestProjectInit_CreatesTaskDir(t *testing.T) {
	tmpDir := t.TempDir()
	taskDirPath := filepath.Join(tmpDir, "new-tasks")

	resetProjectInitFlags(taskDirPath)
	projectInitRoot = tmpDir
	projectInitTaskDir = taskDirPath
	projectInitClaude = true

	// Verify task dir doesn't exist yet
	if _, err := os.Stat(taskDirPath); err == nil {
		t.Fatal("task directory should not exist before init")
	}

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Task directory should now exist
	info, err := os.Stat(taskDirPath)
	if err != nil {
		t.Fatalf("task directory should have been created: %v", err)
	}
	if !info.IsDir() {
		t.Error("task directory path should be a directory")
	}

	// Spec should be inside the new task directory
	specPath := filepath.Join(taskDirPath, specFilename)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Error("TASKMD_SPEC.md should have been created in new task directory")
	}
}

func TestProjectInit_Stdout_NoSideEffects(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitStdout = true
	projectInitClaude = true

	// Capture stdout
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No config file should have been created
	configPath := filepath.Join(tmpDir, configFilename)
	if _, err := os.Stat(configPath); err == nil {
		t.Error(".taskmd.yaml should not have been created with --stdout")
	}

	// No agent files should have been created
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("CLAUDE.md should not have been created with --stdout")
	}

	// No spec file should have been created
	specPath := filepath.Join(tmpDir, specFilename)
	if _, err := os.Stat(specPath); err == nil {
		t.Error("TASKMD_SPEC.md should not have been created with --stdout")
	}
}

func TestProjectInit_CreatesTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmplDir := filepath.Join(tmpDir, ".taskmd", "templates")
	info, err := os.Stat(tmplDir)
	if err != nil {
		t.Fatalf("expected .taskmd/templates/ to be created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected .taskmd/templates/ to be a directory")
	}

	// Check that built-in template files exist
	for _, name := range []string{"feature.md", "bug.md", "chore.md"} {
		path := filepath.Join(tmplDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected template %s to be created", name)
		}
	}
}

func TestProjectInit_NoTemplatesFlag(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitNoTemplates = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmplDir := filepath.Join(tmpDir, ".taskmd", "templates")
	if _, err := os.Stat(tmplDir); err == nil {
		t.Error(".taskmd/templates/ should not have been created with --no-templates")
	}
}

func TestProjectInit_TemplatesNotOverwrittenWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true

	// Create an existing template
	tmplDir := filepath.Join(tmpDir, ".taskmd", "templates")
	os.MkdirAll(tmplDir, 0755)
	os.WriteFile(filepath.Join(tmplDir, "feature.md"), []byte("custom content"), 0644)

	// Suppress stderr
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Custom template should be unchanged
	content, _ := os.ReadFile(filepath.Join(tmplDir, "feature.md"))
	if string(content) != "custom content" {
		t.Error("existing template should not have been overwritten without --force")
	}
}

func TestProjectInit_TemplatesOverwrittenWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitForce = true

	// Create an existing template
	tmplDir := filepath.Join(tmpDir, ".taskmd", "templates")
	os.MkdirAll(tmplDir, 0755)
	os.WriteFile(filepath.Join(tmplDir, "feature.md"), []byte("custom content"), 0644)

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Template should be overwritten with built-in content
	content, _ := os.ReadFile(filepath.Join(tmplDir, "feature.md"))
	if string(content) == "custom content" {
		t.Error("template should have been overwritten with --force")
	}
	if !strings.Contains(string(content), "_template:") {
		t.Error("overwritten template should contain built-in template content")
	}
}

func TestProjectInit_EnsureTaskDir_PathIsFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file where the task directory should be
	filePath := filepath.Join(tmpDir, "tasks")
	if err := os.WriteFile(filePath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	resetProjectInitFlags(tmpDir)
	projectInitRoot = tmpDir
	projectInitTaskDir = filePath
	projectInitClaude = true

	err := runProjectInit(projectInitCmd, []string{})
	if err == nil {
		t.Fatal("expected error when task-dir path is a file")
	}

	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("expected 'not a directory' error, got: %v", err)
	}
}

// --- ID Strategy tests ---

func TestProjectInit_IDStrategy_ULID_ConfigAndTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitIDStrategy = "ulid"

	projectInitCmd.Flags().Set("id-strategy", "ulid")

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check config file includes id section
	configPath := filepath.Join(tmpDir, configFilename)
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	configStr := string(configContent)
	if !strings.Contains(configStr, "strategy: ulid") {
		t.Errorf("config should contain 'strategy: ulid', got: %s", configStr)
	}
	if !strings.Contains(configStr, "length: 9") {
		t.Errorf("config should contain 'length: 9', got: %s", configStr)
	}

	// Check CLAUDE.md has ULID examples
	claudePath := filepath.Join(tmpDir, "CLAUDE.md")
	claudeContent, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	claudeStr := string(claudeContent)
	if !strings.Contains(claudeStr, `id: "01h5a3mpk"`) {
		t.Errorf("CLAUDE.md should contain ULID example ID, got: %s", claudeStr)
	}
	if !strings.Contains(claudeStr, "01h5a3mpk-add-user-auth.md") {
		t.Errorf("CLAUDE.md should contain ULID example filename")
	}
	if strings.Contains(claudeStr, `id: "001"`) {
		t.Error("CLAUDE.md should not contain sequential example ID")
	}

	// Check spec has ULID section
	specPath := filepath.Join(tmpDir, specFilename)
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read spec: %v", err)
	}
	if !strings.Contains(string(specContent), "ULID") {
		t.Error("spec should contain ULID documentation section")
	}
}

func TestProjectInit_IDStrategy_Random_ConfigAndTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitIDStrategy = "random"

	projectInitCmd.Flags().Set("id-strategy", "random")

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configContent, _ := os.ReadFile(filepath.Join(tmpDir, configFilename))
	configStr := string(configContent)
	if !strings.Contains(configStr, "strategy: random") {
		t.Errorf("config should contain 'strategy: random', got: %s", configStr)
	}
	if !strings.Contains(configStr, "length: 6") {
		t.Errorf("config should contain 'length: 6', got: %s", configStr)
	}

	claudeContent, _ := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if !strings.Contains(string(claudeContent), `id: "a3f9x2"`) {
		t.Error("CLAUDE.md should contain random example ID")
	}
	if !strings.Contains(string(claudeContent), "a3f9x2-add-user-auth.md") {
		t.Error("CLAUDE.md should contain random example filename")
	}
}

func TestProjectInit_IDStrategy_Prefixed_ConfigAndTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitIDStrategy = "prefixed"
	projectInitIDPrefix = "dr"

	projectInitCmd.Flags().Set("id-strategy", "prefixed")
	projectInitCmd.Flags().Set("id-prefix", "dr")

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configContent, _ := os.ReadFile(filepath.Join(tmpDir, configFilename))
	configStr := string(configContent)
	if !strings.Contains(configStr, "strategy: prefixed") {
		t.Errorf("config should contain 'strategy: prefixed', got: %s", configStr)
	}
	if !strings.Contains(configStr, "prefix: dr") {
		t.Errorf("config should contain 'prefix: dr', got: %s", configStr)
	}

	claudeContent, _ := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	claudeStr := string(claudeContent)
	if !strings.Contains(claudeStr, `id: "dr-001"`) {
		t.Errorf("CLAUDE.md should contain prefixed example ID, got relevant part missing")
	}
	if !strings.Contains(claudeStr, "dr-015-add-user-auth.md") {
		t.Error("CLAUDE.md should contain prefixed example filename")
	}

	specContent, _ := os.ReadFile(filepath.Join(tmpDir, specFilename))
	if !strings.Contains(string(specContent), "prefixed") {
		t.Error("spec should contain prefixed documentation section")
	}
}

func TestProjectInit_IDStrategy_Sequential_NoIDConfig(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	// No --id-strategy flag set, should default to sequential

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configContent, _ := os.ReadFile(filepath.Join(tmpDir, configFilename))
	configStr := string(configContent)
	// Sequential is default, so no id: section should be written
	if strings.Contains(configStr, "id:") {
		t.Errorf("config should not contain id: section for sequential strategy, got: %s", configStr)
	}

	// Templates should remain unchanged (sequential examples)
	claudeContent, _ := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if !bytes.Equal(claudeContent, claudeTemplate) {
		t.Error("CLAUDE.md should match raw template for sequential strategy")
	}
}

func TestProjectInit_IDStrategy_InvalidStrategy(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitIDStrategy = "invalid"

	projectInitCmd.Flags().Set("id-strategy", "invalid")

	err := runProjectInit(projectInitCmd, []string{})
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
	if !strings.Contains(err.Error(), "invalid --id-strategy") {
		t.Errorf("error should mention invalid strategy, got: %v", err)
	}
}

func TestProjectInit_IDStrategy_PrefixedWithoutPrefix_NonTTY(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitIDStrategy = "prefixed"

	projectInitCmd.Flags().Set("id-strategy", "prefixed")

	err := runProjectInit(projectInitCmd, []string{})
	if err == nil {
		t.Fatal("expected error for prefixed strategy without prefix in non-TTY mode")
	}
	if !strings.Contains(err.Error(), "id-prefix") {
		t.Errorf("error should mention id-prefix, got: %v", err)
	}
}

func TestProjectInit_IDStrategy_ULID_AllAgents(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitClaude = true
	projectInitGemini = true
	projectInitCodex = true
	projectInitIDStrategy = "ulid"

	projectInitCmd.Flags().Set("id-strategy", "ulid")

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All agent files should have ULID examples
	for _, name := range []string{"CLAUDE.md", "GEMINI.md", "AGENTS.md"} {
		content, err := os.ReadFile(filepath.Join(tmpDir, name))
		if err != nil {
			t.Fatalf("failed to read %s: %v", name, err)
		}
		if !strings.Contains(string(content), `id: "01h5a3mpk"`) {
			t.Errorf("%s should contain ULID example ID", name)
		}
	}
}

func TestProjectInit_IDStrategy_Stdout_ShowsReplacedContent(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitStdout = true
	projectInitClaude = true
	projectInitIDStrategy = "ulid"

	projectInitCmd.Flags().Set("id-strategy", "ulid")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runProjectInit(projectInitCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, `id: "01h5a3mpk"`) {
		t.Error("stdout should contain ULID example ID")
	}
	if strings.Contains(output, `id: "001"`) {
		t.Error("stdout should not contain sequential example ID")
	}
}

func TestBuildIDConfigYAML(t *testing.T) {
	tests := []struct {
		name     string
		cfg      idStrategyConfig
		expected string
	}{
		{
			name:     "sequential produces no output",
			cfg:      idStrategyConfig{strategy: "sequential"},
			expected: "",
		},
		{
			name:     "ulid",
			cfg:      idStrategyConfig{strategy: "ulid"},
			expected: "id:\n  strategy: ulid\n  length: 9\n",
		},
		{
			name:     "random",
			cfg:      idStrategyConfig{strategy: "random"},
			expected: "id:\n  strategy: random\n  length: 6\n",
		},
		{
			name:     "prefixed",
			cfg:      idStrategyConfig{strategy: "prefixed", prefix: "cli"},
			expected: "id:\n  strategy: prefixed\n  prefix: cli\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildIDConfigYAML(tt.cfg)
			if got != tt.expected {
				t.Errorf("buildIDConfigYAML() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetIDStrategyExamples(t *testing.T) {
	tests := []struct {
		name            string
		cfg             idStrategyConfig
		wantID          string
		wantFilename    string
		wantFilePattern string
	}{
		{"sequential", idStrategyConfig{strategy: "sequential"}, "001", "015-add-user-auth.md", "NNN-descriptive-title.md"},
		{"ulid", idStrategyConfig{strategy: "ulid"}, "01h5a3mpk", "01h5a3mpk-add-user-auth.md", "ID-descriptive-title.md"},
		{"random", idStrategyConfig{strategy: "random"}, "a3f9x2", "a3f9x2-add-user-auth.md", "ID-descriptive-title.md"},
		{"prefixed", idStrategyConfig{strategy: "prefixed", prefix: "dr"}, "dr-001", "dr-015-add-user-auth.md", "DR-NNN-descriptive-title.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := getIDStrategyExamples(tt.cfg)
			if ex.exampleID != tt.wantID {
				t.Errorf("exampleID = %q, want %q", ex.exampleID, tt.wantID)
			}
			if ex.exampleFilename != tt.wantFilename {
				t.Errorf("exampleFilename = %q, want %q", ex.exampleFilename, tt.wantFilename)
			}
			if ex.filePattern != tt.wantFilePattern {
				t.Errorf("filePattern = %q, want %q", ex.filePattern, tt.wantFilePattern)
			}
		})
	}
}

func TestApplyIDStrategyReplacements(t *testing.T) {
	input := []byte(`id: "001"` + "\n" + `Files follow the pattern ` + "`NNN-descriptive-title.md`" + ` (e.g., ` + "`015-add-user-auth.md`" + `).`)

	examples := getIDStrategyExamples(idStrategyConfig{strategy: "ulid"})
	result := string(applyIDStrategyReplacements(input, examples))

	if !strings.Contains(result, `id: "01h5a3mpk"`) {
		t.Error("should replace example ID")
	}
	if !strings.Contains(result, "01h5a3mpk-add-user-auth.md") {
		t.Error("should replace example filename")
	}
	if !strings.Contains(result, "ID-descriptive-title.md") {
		t.Error("should replace file pattern")
	}
}

func TestApplyIDStrategyReplacements_Sequential_NoOp(t *testing.T) {
	input := []byte(`id: "001"` + "\n" + "015-add-user-auth.md")
	examples := getIDStrategyExamples(idStrategyConfig{strategy: "sequential"})
	result := applyIDStrategyReplacements(input, examples)
	if !bytes.Equal(result, input) {
		t.Error("sequential strategy should return unchanged content")
	}
}

func TestProjectInit_NoSpecNoAgentStillCreatesTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	resetProjectInitFlags(tmpDir)
	projectInitNoSpec = true
	projectInitNoAgent = true

	err := runProjectInit(projectInitCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmplDir := filepath.Join(tmpDir, ".taskmd", "templates")
	if _, err := os.Stat(tmplDir); os.IsNotExist(err) {
		t.Error("expected .taskmd/templates/ to be created even with --no-spec and --no-agent")
	}
}
