package sync

import "strings"

// MappedTask holds the result of mapping an ExternalTask to taskmd fields.
type MappedTask struct {
	Title       string
	Description string
	Status      string
	Priority    string
	Owner       string
	Tags        []string
	URL         string
}

// MapExternalTask converts an ExternalTask to taskmd fields using the FieldMap.
func MapExternalTask(ext ExternalTask, fm FieldMap) MappedTask {
	m := MappedTask{
		Title:       ext.Title,
		Description: ext.Description,
		URL:         ext.URL,
	}

	m.Status = mapField(ext.Status, fm.Status, "pending")
	m.Priority = mapField(ext.Priority, fm.Priority, "")

	if fm.AssigneeToOwner {
		m.Owner = ext.Assignee
	}

	if fm.LabelsToTags {
		for _, label := range ext.Labels {
			m.Tags = append(m.Tags, normalizeLabel(label))
		}
	}

	return m
}

// normalizeLabel lowercases and replaces spaces with hyphens.
func normalizeLabel(label string) string {
	return strings.ReplaceAll(strings.ToLower(label), " ", "-")
}

func mapField(value string, mapping map[string]string, fallback string) string {
	if mapping != nil {
		if mapped, ok := mapping[value]; ok {
			return mapped
		}
	}
	if value == "" {
		return fallback
	}
	return fallback
}
