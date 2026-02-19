package todos

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// BlameInfo holds git blame metadata for a TODO item.
type BlameInfo struct {
	Author string `json:"author" yaml:"author"`
	Commit string `json:"commit" yaml:"commit"`
	Date   string `json:"date" yaml:"date"`
}

// TodoItem represents a single marker comment found in a source file.
type TodoItem struct {
	ID       string     `json:"id" yaml:"id"`
	FilePath string     `json:"file" yaml:"file"`
	Line     int        `json:"line" yaml:"line"`
	Column   int        `json:"column" yaml:"column"`
	Marker   string     `json:"tag" yaml:"tag"`
	Language string     `json:"language" yaml:"language"`
	Text     string     `json:"text" yaml:"text"`
	RawText  string     `json:"raw_text,omitempty" yaml:"raw_text,omitempty"`
	Scope    string     `json:"scope,omitempty" yaml:"scope,omitempty"`
	Blame    *BlameInfo `json:"blame,omitempty" yaml:"blame,omitempty"`
	Age      int        `json:"age,omitempty" yaml:"age,omitempty"`
}

// allMarkersRegex matches any known marker keyword. Used to detect
// non-filtered markers that should still break continuation.
var allMarkersRegex = buildRegex([]string{"TODO", "FIXME", "HACK", "XXX", "NOTE", "BUG", "OPTIMIZE"})

// buildRegex builds a regex that matches any of the given words as whole words.
func buildRegex(words []string) *regexp.Regexp {
	escaped := make([]string, len(words))
	for i, m := range words {
		escaped[i] = regexp.QuoteMeta(m)
	}
	pattern := fmt.Sprintf(`\b(%s)\b`, strings.Join(escaped, "|"))
	return regexp.MustCompile(pattern)
}

// buildMarkerRegex builds a regex that matches any of the given markers as whole words.
func buildMarkerRegex(markers []string) *regexp.Regexp {
	return buildRegex(markers)
}

// ParseFile reads a file and extracts TODO items from comments.
func ParseFile(path string, syntax *CommentSyntax, markers []string, rawText bool) ([]TodoItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseLines(bufio.NewScanner(f), path, syntax, markers, rawText), nil
}

type parserState int

const (
	stateNormal parserState = iota
	stateBlock
)

// parseResult bundles the return values from line handlers to simplify signatures.
type parseResult struct {
	items   []TodoItem
	current *TodoItem
	state   parserState
}

func parseLines(sc *bufio.Scanner, filePath string, syntax *CommentSyntax, markers []string, rawText bool) []TodoItem {
	re := buildMarkerRegex(markers)
	var items []TodoItem
	var current *TodoItem
	state := stateNormal
	lineNum := 0

	for sc.Scan() {
		lineNum++
		line := sc.Text()

		var r parseResult
		switch state {
		case stateNormal:
			r = handleNormalLine(line, lineNum, filePath, syntax, re, items, current, rawText)
		case stateBlock:
			r = handleBlockLine(line, lineNum, filePath, syntax, re, items, current, rawText)
		}
		items, current, state = r.items, r.current, r.state
	}

	if current != nil {
		items = append(items, *current)
	}

	return items
}

// finishCurrent appends the current item (if any) to items and returns nil.
func finishCurrent(items []TodoItem, current *TodoItem) []TodoItem {
	if current != nil {
		items = append(items, *current)
	}
	return items
}

func handleNormalLine(
	line string, lineNum int, filePath string,
	syntax *CommentSyntax, re *regexp.Regexp,
	items []TodoItem, current *TodoItem, rawText bool,
) parseResult {
	if syntax.BlockStart != "" {
		if r, ok := tryBlockStart(line, lineNum, filePath, syntax, re, items, current, rawText); ok {
			return r
		}
	}

	if r, ok := tryLineComment(line, lineNum, filePath, syntax, re, items, current, rawText); ok {
		return r
	}

	return parseResult{finishCurrent(items, current), nil, stateNormal}
}

func tryBlockStart(
	line string, lineNum int, filePath string,
	syntax *CommentSyntax, re *regexp.Regexp,
	items []TodoItem, current *TodoItem, rawText bool,
) (parseResult, bool) {
	idx := strings.Index(line, syntax.BlockStart)
	if idx < 0 {
		return parseResult{}, false
	}

	commentText := line[idx+len(syntax.BlockStart):]

	// Check if block closes on same line
	if endIdx := strings.Index(commentText, syntax.BlockEnd); endIdx >= 0 {
		items = finishCurrent(items, current)
		m := matchMarker(strings.TrimSpace(commentText[:endIdx]), re, filePath, lineNum, line)
		if m != nil && rawText {
			m.RawText = line
		}
		items = finishCurrent(items, m)
		return parseResult{items, nil, stateNormal}, true
	}

	// Block continues to next line
	items = finishCurrent(items, current)
	m := matchMarker(strings.TrimSpace(commentText), re, filePath, lineNum, line)
	if m != nil && rawText {
		m.RawText = line
	}
	return parseResult{items, m, stateBlock}, true
}

func tryLineComment(
	line string, lineNum int, filePath string,
	syntax *CommentSyntax, re *regexp.Regexp,
	items []TodoItem, current *TodoItem, rawText bool,
) (parseResult, bool) {
	for _, prefix := range syntax.LinePrefix {
		idx := strings.Index(line, prefix)
		if idx < 0 {
			continue
		}

		commentText := strings.TrimSpace(line[idx+len(prefix):])
		if m := matchMarker(commentText, re, filePath, lineNum, line); m != nil {
			if rawText {
				m.RawText = line
			}
			items = finishCurrent(items, current)
			return parseResult{items, m, stateNormal}, true
		}

		// Continue the current item if the comment line doesn't start a new marker
		if current != nil && canContinue(current.Line, lineNum, commentText) {
			current.Text += " " + commentText
			if rawText {
				current.RawText += "\n" + line
			}
			return parseResult{items, current, stateNormal}, true
		}

		return parseResult{finishCurrent(items, current), nil, stateNormal}, true
	}
	return parseResult{}, false
}

// canContinue checks whether a comment line should extend the current item.
func canContinue(startLine, currentLine int, text string) bool {
	return currentLine-startLine <= 5 && !allMarkersRegex.MatchString(text)
}

func handleBlockLine(
	line string, lineNum int, filePath string,
	syntax *CommentSyntax, re *regexp.Regexp,
	items []TodoItem, current *TodoItem, rawText bool,
) parseResult {
	if endIdx := strings.Index(line, syntax.BlockEnd); endIdx >= 0 {
		return closeBlock(line[:endIdx], lineNum, filePath, re, items, current, line, rawText)
	}

	trimmed := stripBlockPrefix(line)
	if trimmed == "" {
		if current != nil && rawText {
			current.RawText += "\n" + line
		}
		return parseResult{items, current, stateBlock}
	}

	if m := matchMarker(trimmed, re, filePath, lineNum, line); m != nil {
		if rawText {
			m.RawText = line
		}
		items = finishCurrent(items, current)
		return parseResult{items, m, stateBlock}
	}
	if current != nil {
		current.Text += " " + trimmed
		if rawText {
			current.RawText += "\n" + line
		}
	}
	return parseResult{items, current, stateBlock}
}

func closeBlock(
	beforeEnd string, lineNum int, filePath string,
	re *regexp.Regexp, items []TodoItem, current *TodoItem, sourceLine string, rawText bool,
) parseResult {
	text := strings.TrimSpace(beforeEnd)
	if current != nil && text != "" {
		current.Text += " " + text
		if rawText {
			current.RawText += "\n" + sourceLine
		}
	} else if text != "" {
		if m := matchMarker(text, re, filePath, lineNum, sourceLine); m != nil {
			if rawText {
				m.RawText = sourceLine
			}
			items = finishCurrent(items, current)
			current = m
		}
	} else if current != nil && rawText {
		current.RawText += "\n" + sourceLine
	}
	items = finishCurrent(items, current)
	return parseResult{items, nil, stateNormal}
}

// stripBlockPrefix removes leading whitespace and asterisks common in block comments.
func stripBlockPrefix(line string) string {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimLeft(trimmed, "* ")
	return strings.TrimSpace(trimmed)
}

// matchMarker checks if text contains a marker and returns a TodoItem if found.
// sourceLine is the original unprocessed line used to compute the column offset.
func matchMarker(text string, re *regexp.Regexp, filePath string, line int, sourceLine string) *TodoItem {
	loc := re.FindStringIndex(text)
	if loc == nil {
		return nil
	}

	marker := text[loc[0]:loc[1]]
	rest := strings.TrimSpace(text[loc[1]:])
	rest = stripMarkerPrefix(rest)

	col := strings.Index(sourceLine, marker)
	if col < 0 {
		col = 0
	}

	return &TodoItem{
		FilePath: filePath,
		Line:     line,
		Column:   col,
		Marker:   marker,
		Text:     rest,
	}
}

// stripMarkerPrefix removes optional (author): or : prefix after a marker.
func stripMarkerPrefix(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '(' {
		if end := strings.Index(s, ")"); end >= 0 {
			s = s[end+1:]
		}
	}
	s = strings.TrimLeft(s, ": ")
	return strings.TrimSpace(s)
}
