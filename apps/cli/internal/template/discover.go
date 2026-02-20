package template

import (
	"os"
	"path/filepath"
	"strings"
)

// BuiltinTemplates is set by the CLI package to provide embedded templates.
// Each key is the template name, value is the raw markdown content.
var BuiltinTemplates map[string]string

// Discover finds all available templates across all sources.
// Returns templates deduplicated by name (project > user > built-in).
func Discover(projectRoot, userHome string) []Template {
	seen := make(map[string]bool)
	var templates []Template

	// 1. Project-level: .taskmd/templates/*.md
	projectDir := filepath.Join(projectRoot, ".taskmd", "templates")
	for _, tmpl := range loadFromDir(projectDir, "project") {
		if !seen[tmpl.Name] {
			seen[tmpl.Name] = true
			templates = append(templates, tmpl)
		}
	}

	// 2. User-level: ~/.taskmd/templates/*.md
	if userHome != "" {
		userDir := filepath.Join(userHome, ".taskmd", "templates")
		for _, tmpl := range loadFromDir(userDir, "user") {
			if !seen[tmpl.Name] {
				seen[tmpl.Name] = true
				templates = append(templates, tmpl)
			}
		}
	}

	// 3. Built-in templates
	for name, content := range BuiltinTemplates {
		if seen[name] {
			continue
		}
		tmpl, err := ParseTemplate(content, "built-in")
		if err != nil {
			continue
		}
		seen[name] = true
		templates = append(templates, tmpl)
	}

	return templates
}

// Resolve finds the first template matching a name across all sources.
func Resolve(name, projectRoot, userHome string) (Template, bool) {
	// 1. Project-level
	projectDir := filepath.Join(projectRoot, ".taskmd", "templates")
	if tmpl, ok := loadByName(projectDir, name, "project"); ok {
		return tmpl, true
	}

	// 2. User-level
	if userHome != "" {
		userDir := filepath.Join(userHome, ".taskmd", "templates")
		if tmpl, ok := loadByName(userDir, name, "user"); ok {
			return tmpl, true
		}
	}

	// 3. Built-in
	if content, ok := BuiltinTemplates[name]; ok {
		tmpl, err := ParseTemplate(content, "built-in")
		if err == nil {
			return tmpl, true
		}
	}

	return Template{}, false
}

// loadFromDir scans a directory for .md files and parses them as templates.
func loadFromDir(dir, source string) []Template {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var templates []Template
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		tmpl, err := ParseTemplate(string(content), source)
		if err != nil {
			continue
		}
		templates = append(templates, tmpl)
	}

	return templates
}

// loadByName looks for a template with a matching name in a directory.
func loadByName(dir, name, source string) (Template, bool) {
	for _, tmpl := range loadFromDir(dir, source) {
		if tmpl.Name == name {
			return tmpl, true
		}
	}
	return Template{}, false
}
