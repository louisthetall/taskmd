package web

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// mockStaticFS creates a minimal embedded FS for testing.
func mockStaticFS() fstest.MapFS {
	return fstest.MapFS{
		"static/dist/index.html": &fstest.MapFile{
			Data: []byte(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>taskmd</title>
  <script type="module" crossorigin src="/assets/index-abc123.js"></script>
  <link rel="stylesheet" crossorigin href="/assets/index-abc123.css">
</head>
<body>
  <div id="root"></div>
</body>
</html>`),
		},
		"static/dist/assets/index-abc123.js": &fstest.MapFile{
			Data: []byte(`console.log("app")`),
		},
		"static/dist/assets/index-abc123.css": &fstest.MapFile{
			Data: []byte(`body { margin: 0; }`),
		},
		"static/dist/favicon.ico": &fstest.MapFile{
			Data: []byte("icon"),
		},
	}
}

// exportWithMockFS runs Export after temporarily replacing StaticFiles.
// Since StaticFiles() returns an empty FS in non-embed builds, we
// use ExportWithFS which accepts an explicit FS for testing.
func exportWithMockFS(t *testing.T, cfg ExportConfig) error {
	t.Helper()
	return ExportWithFS(cfg, mockStaticFS())
}

func TestExport_OutputStructure(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test-1.0.0",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check core files exist
	expected := []string{
		"index.html",
		"404.html",
		"assets/index-abc123.js",
		"assets/index-abc123.css",
		"api/config.json",
		"api/tasks.json",
		"api/tasks/001.json",
		"api/tasks/002.json",
		"api/tasks/001/worklog.json",
		"api/tasks/002/worklog.json",
		"api/board/status.json",
		"api/board/priority.json",
		"api/board/effort.json",
		"api/board/type.json",
		"api/board/group.json",
		"api/board/tag.json",
		"api/graph.json",
		"api/stats.json",
		"api/next.json",
		"api/tracks.json",
		"api/validate.json",
	}

	for _, f := range expected {
		path := filepath.Join(outDir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", f, err)
		}
	}
}

func TestExport_TasksJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "api", "tasks.json"))
	if err != nil {
		t.Fatalf("failed to read tasks.json: %v", err)
	}

	var tasks []map[string]any
	if err := json.Unmarshal(data, &tasks); err != nil {
		t.Fatalf("invalid JSON in tasks.json: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestExport_TaskDetailJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "api", "tasks", "001.json"))
	if err != nil {
		t.Fatalf("failed to read task detail: %v", err)
	}

	var detail map[string]any
	if err := json.Unmarshal(data, &detail); err != nil {
		t.Fatalf("invalid JSON in task detail: %v", err)
	}

	if detail["id"] != "001" {
		t.Errorf("expected id '001', got %v", detail["id"])
	}

	if _, ok := detail["body"]; !ok {
		t.Error("expected body field in task detail")
	}
}

func TestExport_BoardJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	groupByValues := []string{"status", "priority", "effort", "type", "group", "tag"}
	for _, gb := range groupByValues {
		path := filepath.Join(outDir, "api", "board", gb+".json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read board/%s.json: %v", gb, err)
			continue
		}

		var groups []map[string]any
		if err := json.Unmarshal(data, &groups); err != nil {
			t.Errorf("invalid JSON in board/%s.json: %v", gb, err)
		}
	}
}

func TestExport_ConfigJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test-1.0.0",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "api", "config.json"))
	if err != nil {
		t.Fatalf("failed to read config.json: %v", err)
	}

	var config ConfigResponse
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("invalid JSON in config.json: %v", err)
	}

	if !config.ReadOnly {
		t.Error("expected readonly to be true in exported config")
	}

	if config.Version != "test-1.0.0" {
		t.Errorf("expected version 'test-1.0.0', got %q", config.Version)
	}
}

func TestExport_IndexHTML_FetchInterceptor(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	html := string(data)

	if !strings.Contains(html, "window.fetch") {
		t.Error("expected fetch interceptor in index.html")
	}

	if !strings.Contains(html, "window.EventSource") {
		t.Error("expected EventSource stub in index.html")
	}

	// Asset paths should be relative
	if strings.Contains(html, `="/assets/`) {
		t.Error("expected absolute asset paths to be rewritten to relative")
	}

	if !strings.Contains(html, `="./assets/`) {
		t.Error("expected relative asset paths in index.html")
	}
}

func TestExport_IndexHTML_BasePath(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/demo/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	html := string(data)

	if !strings.Contains(html, `<base href="/demo/">`) {
		t.Error("expected <base> tag with /demo/ in index.html")
	}
}

func TestExport_IndexHTML_NoBasePathForRoot(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	if strings.Contains(string(data), "<base") {
		t.Error("expected no <base> tag when base-path is /")
	}
}

func TestExport_SPARouteFiles(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Top-level SPA routes
	routes := []string{"tasks", "board", "graph", "next", "stats", "validate", "tracks"}
	for _, route := range routes {
		path := filepath.Join(outDir, route, "index.html")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected SPA fallback at %s/index.html: %v", route, err)
		}
	}

	// Per-task routes
	for _, id := range []string{"001", "002"} {
		path := filepath.Join(outDir, "tasks", id, "index.html")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected task route at tasks/%s/index.html: %v", id, err)
		}
	}

	// 404.html
	if _, err := os.Stat(filepath.Join(outDir, "404.html")); err != nil {
		t.Error("expected 404.html to exist")
	}
}

func TestExport_NoEmbeddedAssets(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	// Use empty FS to simulate no embedded assets
	err := ExportWithFS(ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	}, fstest.MapFS{})

	if err == nil {
		t.Fatal("expected error when no embedded assets")
	}

	if !strings.Contains(err.Error(), "no embedded web assets") {
		t.Errorf("expected 'no embedded web assets' error, got: %v", err)
	}
}

func TestExport_CustomOutput(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "custom-dir", "nested")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify output is in the custom directory
	if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
		t.Error("expected index.html in custom output directory")
	}
}

func TestExport_GraphJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "api", "graph.json"))
	if err != nil {
		t.Fatalf("failed to read graph.json: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON in graph.json: %v", err)
	}

	if _, ok := result["nodes"]; !ok {
		t.Error("expected 'nodes' in graph.json")
	}
	if _, ok := result["edges"]; !ok {
		t.Error("expected 'edges' in graph.json")
	}
}

func TestExport_StatsJSON(t *testing.T) {
	taskDir := createTestTaskDir(t)
	outDir := filepath.Join(t.TempDir(), "export")

	err := exportWithMockFS(t, ExportConfig{
		OutputDir: outDir,
		ScanDir:   taskDir,
		BasePath:  "/",
		Version:   "test",
	})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "api", "stats.json"))
	if err != nil {
		t.Fatalf("failed to read stats.json: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON in stats.json: %v", err)
	}

	totalTasks, ok := result["total_tasks"].(float64)
	if !ok {
		t.Fatal("expected total_tasks in stats.json")
	}
	if int(totalTasks) != 2 {
		t.Errorf("expected 2 total tasks, got %d", int(totalTasks))
	}
}

func TestPatchIndexHTML_RelativePaths(t *testing.T) {
	input := `<head><script src="/assets/app.js"></script></head>`
	result := patchIndexHTML(input, "/")

	if strings.Contains(result, `"/assets/`) {
		t.Error("expected absolute asset paths to be rewritten")
	}
}

func TestPatchIndexHTML_BasePathTrailingSlash(t *testing.T) {
	input := `<head></head>`
	result := patchIndexHTML(input, "/demo")

	if !strings.Contains(result, `<base href="/demo/">`) {
		t.Error("expected base-path to get trailing slash")
	}
}
