package todos

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	goFuncRe     = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)
	jsFuncRe     = regexp.MustCompile(`(?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=|class\s+(\w+)|(\w+)\s*\([^)]*\)\s*\{)`)
	pyDefClassRe = regexp.MustCompile(`^(\s*)(?:def|class)\s+(\w+)`)
)

// DetectScope reads the file up to the given line and returns the enclosing
// scope name (function, method, or class). Returns "" if no scope is found
// or the language is unsupported.
func DetectScope(filePath string, line int, lang string) string {
	switch lang {
	case "go":
		return detectGoScope(filePath, line)
	case "javascript", "typescript":
		return detectJSScope(filePath, line)
	case "python":
		return detectPythonScope(filePath, line)
	default:
		return ""
	}
}

func detectGoScope(filePath string, targetLine int) string {
	lines, err := readLinesUpTo(filePath, targetLine)
	if err != nil {
		return ""
	}

	for i := len(lines) - 1; i >= 0; i-- {
		if m := goFuncRe.FindStringSubmatch(lines[i]); m != nil {
			return m[1]
		}
	}
	return ""
}

func detectJSScope(filePath string, targetLine int) string {
	lines, err := readLinesUpTo(filePath, targetLine)
	if err != nil {
		return ""
	}

	for i := len(lines) - 1; i >= 0; i-- {
		m := jsFuncRe.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		for _, g := range m[1:] {
			if g != "" {
				return g
			}
		}
	}
	return ""
}

func detectPythonScope(filePath string, targetLine int) string {
	lines, err := readLinesUpTo(filePath, targetLine)
	if err != nil || len(lines) == 0 {
		return ""
	}

	// Find the indentation of the target line
	targetIndent := len(lines[len(lines)-1]) - len(strings.TrimLeft(lines[len(lines)-1], " \t"))

	for i := len(lines) - 1; i >= 0; i-- {
		m := pyDefClassRe.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		scopeIndent := len(m[1])
		if scopeIndent < targetIndent {
			return m[2]
		}
	}
	return ""
}

func readLinesUpTo(filePath string, n int) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for i := 0; i < n && sc.Scan(); i++ {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}
