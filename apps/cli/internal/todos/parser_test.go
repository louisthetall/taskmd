package todos

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseLines_GoSingleLine(t *testing.T) {
	src := `package main

// TODO: refactor this function
func main() {
	// FIXME: handle error
	println("hello")
}
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "main.go", syntax, DefaultMarkers, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	assertItem(t, items[0], "main.go", 3, "TODO", "refactor this function")
	assertItem(t, items[1], "main.go", 5, "FIXME", "handle error")
}

func TestParseLines_GoBlockComment(t *testing.T) {
	src := `/*
 * TODO: implement caching
 * across all handlers
 */
func handle() {}
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "main.go", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "main.go", 2, "TODO", "implement caching across all handlers")
}

func TestParseLines_GoInlineBlockComment(t *testing.T) {
	src := `x := 1 /* TODO: remove this */ + 2
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "main.go", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "main.go", 1, "TODO", "remove this")
}

func TestParseLines_PythonHash(t *testing.T) {
	src := `# TODO: add logging
def foo():
    # HACK: workaround for library bug
    pass
`
	syntax := LookupSyntax(".py")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "app.py", syntax, DefaultMarkers, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	assertItem(t, items[0], "app.py", 1, "TODO", "add logging")
	assertItem(t, items[1], "app.py", 3, "HACK", "workaround for library bug")
}

func TestParseLines_PythonDocstring(t *testing.T) {
	src := `"""
TODO: document this module
"""
def foo():
    pass
`
	syntax := LookupSyntax(".py")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "app.py", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "app.py", 2, "TODO", "document this module")
}

func TestParseLines_HTML(t *testing.T) {
	src := `<html>
<!-- TODO: add meta tags -->
<body>
<!-- FIXME: broken layout
     needs responsive fix -->
</body>
</html>
`
	syntax := LookupSyntax(".html")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "index.html", syntax, DefaultMarkers, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	assertItem(t, items[0], "index.html", 2, "TODO", "add meta tags")
	assertItem(t, items[1], "index.html", 4, "FIXME", "broken layout needs responsive fix")
}

func TestParseLines_MultilineBlockContinuation(t *testing.T) {
	src := `/*
 * FIXME: the sorting algorithm
 * does not handle duplicates
 * correctly in all cases
 */
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "sort.go", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].Marker != "FIXME" {
		t.Errorf("expected marker FIXME, got %s", items[0].Marker)
	}
	if !strings.Contains(items[0].Text, "sorting algorithm") {
		t.Errorf("expected text to contain 'sorting algorithm', got %q", items[0].Text)
	}
	if !strings.Contains(items[0].Text, "duplicates") {
		t.Errorf("expected text to contain 'duplicates', got %q", items[0].Text)
	}
}

func TestParseLines_MultipleMarkersInFile(t *testing.T) {
	src := `// TODO: first thing
// FIXME: second thing
// HACK: third thing
// XXX: fourth thing
// NOTE: fifth thing
// BUG: sixth thing
// OPTIMIZE: seventh thing
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "all.go", syntax, DefaultMarkers, false)

	if len(items) != 7 {
		t.Fatalf("expected 7 items, got %d", len(items))
	}

	expectedMarkers := []string{"TODO", "FIXME", "HACK", "XXX", "NOTE", "BUG", "OPTIMIZE"}
	for i, m := range expectedMarkers {
		if items[i].Marker != m {
			t.Errorf("item %d: expected marker %s, got %s", i, m, items[i].Marker)
		}
	}
}

func TestParseLines_MarkerWithColon(t *testing.T) {
	src := `// TODO: implement this
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "t.go", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "t.go", 1, "TODO", "implement this")
}

func TestParseLines_MarkerWithParens(t *testing.T) {
	src := `// TODO(jsmith): implement auth
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "t.go", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "t.go", 1, "TODO", "implement auth")
}

func TestParseLines_MarkerFilter(t *testing.T) {
	src := `// TODO: do this
// FIXME: fix this
// HACK: hack this
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "t.go", syntax, []string{"TODO", "FIXME"}, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items (filtered), got %d", len(items))
	}

	assertItem(t, items[0], "t.go", 1, "TODO", "do this")
	assertItem(t, items[1], "t.go", 2, "FIXME", "fix this")
}

func TestParseLines_NoMarkers(t *testing.T) {
	src := `// This is just a normal comment
func foo() {}
`
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "t.go", syntax, DefaultMarkers, false)

	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestParseLines_EmptyFile(t *testing.T) {
	syntax := LookupSyntax(".go")
	items := parseLines(bufio.NewScanner(strings.NewReader("")), "t.go", syntax, DefaultMarkers, false)

	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestParseLines_ShellComment(t *testing.T) {
	src := `#!/bin/bash
# TODO: add error handling
echo "hello"
# NOTE: this is intentional
`
	syntax := LookupSyntax(".sh")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "run.sh", syntax, DefaultMarkers, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	assertItem(t, items[0], "run.sh", 2, "TODO", "add error handling")
	assertItem(t, items[1], "run.sh", 4, "NOTE", "this is intentional")
}

func TestParseLines_CSSBlockComment(t *testing.T) {
	src := `/* TODO: add responsive styles */
.container { width: 100%; }
`
	syntax := LookupSyntax(".css")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "style.css", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "style.css", 1, "TODO", "add responsive styles")
}

func TestParseLines_YAMLComment(t *testing.T) {
	src := `# TODO: add validation schema
name: test
`
	syntax := LookupSyntax(".yaml")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "config.yaml", syntax, DefaultMarkers, false)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	assertItem(t, items[0], "config.yaml", 1, "TODO", "add validation schema")
}

func TestParseLines_RustComment(t *testing.T) {
	src := `// TODO: implement Display trait
fn main() {
    /* FIXME: handle overflow */
    let x = 42;
}
`
	syntax := LookupSyntax(".rs")
	items := parseLines(bufio.NewScanner(strings.NewReader(src)), "main.rs", syntax, DefaultMarkers, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	assertItem(t, items[0], "main.rs", 1, "TODO", "implement Display trait")
	assertItem(t, items[1], "main.rs", 3, "FIXME", "handle overflow")
}

func TestParseFile_Integration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	content := `package test

// TODO: implement this
func foo() {}

// FIXME: broken
func bar() {}
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	syntax := LookupSyntax(".go")
	items, err := ParseFile(path, syntax, DefaultMarkers, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func assertItem(t *testing.T, item TodoItem, file string, line int, marker, text string) {
	t.Helper()
	if item.FilePath != file {
		t.Errorf("expected file %q, got %q", file, item.FilePath)
	}
	if item.Line != line {
		t.Errorf("expected line %d, got %d", line, item.Line)
	}
	if item.Marker != marker {
		t.Errorf("expected marker %q, got %q", marker, item.Marker)
	}
	if item.Text != text {
		t.Errorf("expected text %q, got %q", text, item.Text)
	}
}
