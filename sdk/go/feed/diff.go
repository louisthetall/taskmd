package feed

import (
	"regexp"
	"sort"
	"strings"
)

var subtaskRegex = regexp.MustCompile(`^- \[([ xX])\] (.+)$`)

// AnalyzeDiff compares old and new task file content, returning detected
// frontmatter field changes and subtask checkbox toggles.
func AnalyzeDiff(oldContent, newContent string) ([]FieldChange, []SubtaskChange) {
	oldFields := ExtractFrontmatterFields(oldContent)
	newFields := ExtractFrontmatterFields(newContent)

	var fieldChanges []FieldChange

	allKeys := make(map[string]struct{})
	for k := range oldFields {
		allKeys[k] = struct{}{}
	}
	for k := range newFields {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		oldVal := oldFields[key]
		newVal := newFields[key]
		if oldVal != newVal {
			fieldChanges = append(fieldChanges, FieldChange{
				Field:    key,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	oldSubtasks := ExtractSubtasks(oldContent)
	newSubtasks := ExtractSubtasks(newContent)

	var subtaskChanges []SubtaskChange
	for text, newDone := range newSubtasks {
		oldDone, exists := oldSubtasks[text]
		if exists && oldDone != newDone {
			subtaskChanges = append(subtaskChanges, SubtaskChange{
				Text: text,
				Done: newDone,
			})
		}
	}

	sort.Slice(subtaskChanges, func(i, j int) bool {
		return subtaskChanges[i].Text < subtaskChanges[j].Text
	})

	return fieldChanges, subtaskChanges
}

// ExtractFrontmatterFields parses YAML frontmatter (between --- delimiters)
// into a map of field name to value. Only handles simple key: value lines.
func ExtractFrontmatterFields(content string) map[string]string {
	fields := make(map[string]string)

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return fields
	}

	frontmatter := parts[1]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		value = strings.Trim(value, `"'`)
		fields[key] = value
	}

	return fields
}

// ExtractSubtasks finds all markdown checkbox lines (- [ ] or - [x]) in the
// body (after frontmatter) and returns a map of subtask text to checked state.
func ExtractSubtasks(content string) map[string]bool {
	subtasks := make(map[string]bool)

	body := content
	parts := strings.SplitN(content, "---", 3)
	if len(parts) >= 3 {
		body = parts[2]
	}

	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		match := subtaskRegex.FindStringSubmatch(line)
		if match != nil {
			checked := match[1] == "x" || match[1] == "X"
			text := match[2]
			subtasks[text] = checked
		}
	}

	return subtasks
}
