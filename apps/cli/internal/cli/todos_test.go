package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/driangle/taskmd/apps/cli/internal/todos"
)

func resetTodosFlags() {
	todosDir = "."
	todosMarkers = nil
	todosInclude = nil
	todosExclude = nil
	todosFormat = "table"
	todosRawText = false
	todosRich = false
	noColor = true
}

func createTodosTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeTodosTestFile(t, filepath.Join(dir, "main.go"), `package main

// TODO: implement main logic
func main() {}

// FIXME: handle error case
func process() error { return nil }
`)

	writeTodosTestFile(t, filepath.Join(dir, "app.py"), `# HACK: workaround for upstream bug
import os
`)

	writeTodosTestFile(t, filepath.Join(dir, "style.css"), `/* NOTE: using hardcoded values */
.container { width: 100%; }
`)

	return dir
}

func writeTodosTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func captureTodosTableOutput(t *testing.T, items []todos.TodoItem) string {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = wErr

	err := outputTodosTable(items, defaultColumns, false)
	if err != nil {
		w.Close()
		wErr.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		t.Fatalf("outputTodosTable failed: %v", err)
	}

	w.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	// drain stderr too
	var stderrBuf bytes.Buffer
	stderrBuf.ReadFrom(rErr)
	return buf.String()
}

func TestTodosList_TableOutput(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)
	todosDir = dir

	// Scan and capture output
	items, err := todos.Scan(todos.ScanOptions{Dir: dir})
	if err != nil {
		t.Fatal(err)
	}

	output := captureTodosTableOutput(t, items)

	if !strings.Contains(output, "FILE") || !strings.Contains(output, "LINE") {
		t.Error("expected header with FILE and LINE")
	}
	if !strings.Contains(output, "TAG") || !strings.Contains(output, "TEXT") {
		t.Error("expected header with TAG and TEXT")
	}
	if !strings.Contains(output, "ID") {
		t.Error("expected header with ID")
	}
	if !strings.Contains(output, "TODO") {
		t.Error("expected TODO marker in output")
	}
	if !strings.Contains(output, "FIXME") {
		t.Error("expected FIXME marker in output")
	}
}

func TestTodosList_TableOutputEmpty(t *testing.T) {
	resetTodosFlags()

	output := captureTodosTableOutput(t, nil)
	if !strings.Contains(output, "No TODO comments found") {
		t.Error("expected 'No TODO comments found' message")
	}
}

func TestTodosList_JSONOutput(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)

	items, err := todos.Scan(todos.ScanOptions{Dir: dir})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = WriteJSON(os.Stdout, items)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v\noutput: %s", err, buf.String())
	}

	if len(parsed) == 0 {
		t.Fatal("expected items in JSON output")
	}

	// Verify fields are present
	for _, item := range parsed {
		if item.FilePath == "" {
			t.Error("expected non-empty file path")
		}
		if item.Line == 0 {
			t.Error("expected non-zero line number")
		}
		if item.Marker == "" {
			t.Error("expected non-empty marker")
		}
	}
}

func TestTodosList_YAMLOutput(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)

	items, err := todos.Scan(todos.ScanOptions{Dir: dir})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = WriteYAML(os.Stdout, items)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("WriteYAML failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "file:") || !strings.Contains(output, "line:") {
		t.Error("expected YAML with file and line fields")
	}
	if !strings.Contains(output, "tag:") || !strings.Contains(output, "text:") {
		t.Error("expected YAML with tag and text fields")
	}
}

func TestTodosList_MarkerFilter(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)

	items, err := todos.Scan(todos.ScanOptions{
		Dir:     dir,
		Markers: []string{"TODO"},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range items {
		if item.Marker != "TODO" {
			t.Errorf("expected only TODO markers, got %s", item.Marker)
		}
	}
}

func TestTodosList_InvalidMarker(t *testing.T) {
	resetTodosFlags()

	err := validateMarkers([]string{"INVALID"})
	if err == nil {
		t.Fatal("expected error for invalid marker")
	}

	if !strings.Contains(err.Error(), "invalid marker") {
		t.Errorf("expected 'invalid marker' in error, got: %s", err.Error())
	}
}

func TestTodosList_ValidMarkers(t *testing.T) {
	resetTodosFlags()

	err := validateMarkers(todos.DefaultMarkers)
	if err != nil {
		t.Fatalf("expected no error for valid markers, got: %v", err)
	}
}

func TestTodosList_RunCommand(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)
	todosDir = dir
	todosFormat = "json"

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) == 0 {
		t.Fatal("expected items from runTodosList")
	}
}

func TestTodosList_EmptyDirectory(t *testing.T) {
	resetTodosFlags()
	dir := t.TempDir()
	todosDir = dir
	todosFormat = "json"

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) != 0 {
		t.Fatalf("expected 0 items for empty dir, got %d", len(parsed))
	}
}

func TestTodosList_InvalidFormat(t *testing.T) {
	resetTodosFlags()
	todosDir = t.TempDir()
	todosFormat = "xml"

	err := runTodosList(nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected 'unsupported format' error, got: %s", err.Error())
	}
}

func TestMergeConfigExcludes_ConfigOnly(t *testing.T) {
	viper.Set("todos.exclude", []string{"*_test.go", "*.generated.*"})
	defer viper.Set("todos.exclude", nil)

	result := mergeConfigExcludes(nil)

	if len(result) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(result))
	}
	if result[0] != "*_test.go" || result[1] != "*.generated.*" {
		t.Errorf("unexpected patterns: %v", result)
	}
}

func TestMergeConfigExcludes_CLIOnly(t *testing.T) {
	viper.Set("todos.exclude", nil)

	result := mergeConfigExcludes([]string{"vendor/*"})

	if len(result) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(result))
	}
	if result[0] != "vendor/*" {
		t.Errorf("unexpected pattern: %v", result)
	}
}

func TestMergeConfigExcludes_BothMerged(t *testing.T) {
	viper.Set("todos.exclude", []string{"*_test.go"})
	defer viper.Set("todos.exclude", nil)

	result := mergeConfigExcludes([]string{"vendor/*"})

	if len(result) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(result))
	}
	if result[0] != "vendor/*" {
		t.Errorf("expected CLI pattern first, got %s", result[0])
	}
	if result[1] != "*_test.go" {
		t.Errorf("expected config pattern second, got %s", result[1])
	}
}

func TestMergeConfigExcludes_NeitherSet(t *testing.T) {
	viper.Set("todos.exclude", nil)

	result := mergeConfigExcludes(nil)

	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestTodosList_ConfigExcludePattern(t *testing.T) {
	resetTodosFlags()
	dir := t.TempDir()
	todosDir = dir
	todosFormat = "json"

	writeTodosTestFile(t, filepath.Join(dir, "main.go"), `package main
// TODO: keep this
`)
	writeTodosTestFile(t, filepath.Join(dir, "main_test.go"), `package main
// TODO: exclude this via config
`)

	viper.Set("todos.exclude", []string{"*_test.go"})
	defer viper.Set("todos.exclude", nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 item (test file excluded by config), got %d", len(parsed))
	}
	if parsed[0].FilePath != "main.go" {
		t.Errorf("expected main.go, got %s", parsed[0].FilePath)
	}
}

func TestTodosList_ConfigAndCLIExcludeCombine(t *testing.T) {
	resetTodosFlags()
	dir := t.TempDir()
	todosDir = dir
	todosFormat = "json"
	todosExclude = []string{"*.py"}

	writeTodosTestFile(t, filepath.Join(dir, "main.go"), `package main
// TODO: keep this
`)
	writeTodosTestFile(t, filepath.Join(dir, "main_test.go"), `package main
// TODO: exclude via config
`)
	writeTodosTestFile(t, filepath.Join(dir, "app.py"), `# TODO: exclude via CLI flag
`)

	viper.Set("todos.exclude", []string{"*_test.go"})
	defer viper.Set("todos.exclude", nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 item (test+py excluded), got %d", len(parsed))
	}
	if parsed[0].FilePath != "main.go" {
		t.Errorf("expected main.go, got %s", parsed[0].FilePath)
	}
}

func TestTodosList_ConfigExcludePathPattern(t *testing.T) {
	resetTodosFlags()
	dir := t.TempDir()
	todosDir = dir
	todosFormat = "json"

	writeTodosTestFile(t, filepath.Join(dir, "main.go"), `package main
// TODO: keep this
`)
	writeTodosTestFile(t, filepath.Join(dir, "sub", "deep.go"), `package sub
// TODO: exclude via path pattern
`)

	viper.Set("todos.exclude", []string{"sub/*.go"})
	defer viper.Set("todos.exclude", nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 item (sub/*.go excluded by config), got %d", len(parsed))
	}
	if parsed[0].FilePath != "main.go" {
		t.Errorf("expected main.go, got %s", parsed[0].FilePath)
	}
}

func TestTodosList_NoConfigExcludeUnchangedBehavior(t *testing.T) {
	resetTodosFlags()
	dir := t.TempDir()
	todosDir = dir
	todosFormat = "json"

	// Ensure no config excludes are set
	viper.Set("todos.exclude", nil)

	writeTodosTestFile(t, filepath.Join(dir, "main.go"), `package main
// TODO: one
`)
	writeTodosTestFile(t, filepath.Join(dir, "main_test.go"), `package main
// TODO: two
`)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 items (no config excludes), got %d", len(parsed))
	}
}

func TestTodosList_JSONOutputNewFields(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)

	items, err := todos.Scan(todos.ScanOptions{Dir: dir})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = WriteJSON(os.Stdout, items)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that new fields are present in JSON
	if !strings.Contains(output, `"id"`) {
		t.Error("expected 'id' field in JSON output")
	}
	if !strings.Contains(output, `"column"`) {
		t.Error("expected 'column' field in JSON output")
	}
	if !strings.Contains(output, `"language"`) {
		t.Error("expected 'language' field in JSON output")
	}
	if !strings.Contains(output, `"tag"`) {
		t.Error("expected 'tag' field in JSON output")
	}

	// Verify fields parse correctly
	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	for _, item := range parsed {
		if item.ID == "" {
			t.Error("expected non-empty ID")
		}
		if len(item.ID) != 12 {
			t.Errorf("expected 12-char ID, got %d: %q", len(item.ID), item.ID)
		}
		if item.Language == "" {
			t.Error("expected non-empty language")
		}
	}
}

func TestTodosList_RawTextFlag(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)
	todosDir = dir
	todosFormat = "json"
	todosRawText = true

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, `"raw_text"`) {
		t.Error("expected 'raw_text' field in JSON output when --raw-text is set")
	}

	var parsed []todos.TodoItem
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	for _, item := range parsed {
		if item.RawText == "" {
			t.Errorf("expected non-empty raw_text for %s:%d", item.FilePath, item.Line)
		}
	}
}

func TestTodosList_NoRawTextByDefault(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)
	todosDir = dir
	todosFormat = "json"

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTodosList(nil, nil)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runTodosList failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if strings.Contains(output, `"raw_text"`) {
		t.Error("did not expect 'raw_text' in JSON output by default")
	}
}

func TestTodosList_RichTableOutput(t *testing.T) {
	resetTodosFlags()
	dir := createTodosTestDir(t)
	todosDir = dir

	items, err := todos.Scan(todos.ScanOptions{Dir: dir})
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = wErr

	err = outputTodosTable(items, richColumns, true)
	w.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("outputTodosTable with rich failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	var stderrBuf bytes.Buffer
	stderrBuf.ReadFrom(rErr)
	output := buf.String()

	if !strings.Contains(output, "SCOPE") {
		t.Error("expected SCOPE column in rich table output")
	}
	if !strings.Contains(output, "AGE") {
		t.Error("expected AGE column in rich table output")
	}
	if !strings.Contains(output, "AUTHOR") {
		t.Error("expected AUTHOR column in rich table output")
	}
}
