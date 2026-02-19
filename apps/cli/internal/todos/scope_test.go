package todos

import (
	"os"
	"path/filepath"
	"testing"
)

func writeScopeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDetectScope_GoFunc(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "main.go", `package main

func doWork() {
	// TODO: implement
}
`)
	scope := DetectScope(path, 4, "go")
	if scope != "doWork" {
		t.Errorf("expected scope 'doWork', got %q", scope)
	}
}

func TestDetectScope_GoMethod(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "main.go", `package main

func (s *Server) Handle() {
	// TODO: implement
}
`)
	scope := DetectScope(path, 4, "go")
	if scope != "Handle" {
		t.Errorf("expected scope 'Handle', got %q", scope)
	}
}

func TestDetectScope_GoNoScope(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "main.go", `package main

// TODO: package-level todo
var x = 1
`)
	scope := DetectScope(path, 3, "go")
	if scope != "" {
		t.Errorf("expected empty scope, got %q", scope)
	}
}

func TestDetectScope_JSFunction(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.js", `function handleClick() {
  // TODO: implement
}
`)
	scope := DetectScope(path, 2, "javascript")
	if scope != "handleClick" {
		t.Errorf("expected scope 'handleClick', got %q", scope)
	}
}

func TestDetectScope_JSConst(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.js", `const processData = () => {
  // TODO: implement
}
`)
	scope := DetectScope(path, 2, "javascript")
	if scope != "processData" {
		t.Errorf("expected scope 'processData', got %q", scope)
	}
}

func TestDetectScope_JSClass(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.ts", `class UserService {
  // TODO: add methods
}
`)
	scope := DetectScope(path, 2, "typescript")
	if scope != "UserService" {
		t.Errorf("expected scope 'UserService', got %q", scope)
	}
}

func TestDetectScope_PythonDef(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.py", `def process():
    # TODO: implement
    pass
`)
	scope := DetectScope(path, 2, "python")
	if scope != "process" {
		t.Errorf("expected scope 'process', got %q", scope)
	}
}

func TestDetectScope_PythonClass(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.py", `class MyClass:
    # TODO: add init
    pass
`)
	scope := DetectScope(path, 2, "python")
	if scope != "MyClass" {
		t.Errorf("expected scope 'MyClass', got %q", scope)
	}
}

func TestDetectScope_PythonNested(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "app.py", `class MyClass:
    def method(self):
        # TODO: implement
        pass
`)
	scope := DetectScope(path, 3, "python")
	if scope != "method" {
		t.Errorf("expected scope 'method', got %q", scope)
	}
}

func TestDetectScope_PythonEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "empty.py", "")

	scope := DetectScope(path, 1, "python")
	if scope != "" {
		t.Errorf("expected empty scope for empty file, got %q", scope)
	}
}

func TestDetectScope_UnsupportedLanguage(t *testing.T) {
	dir := t.TempDir()
	path := writeScopeFile(t, dir, "style.css", `/* TODO: fix colors */`)

	scope := DetectScope(path, 1, "css")
	if scope != "" {
		t.Errorf("expected empty scope for unsupported language, got %q", scope)
	}
}
