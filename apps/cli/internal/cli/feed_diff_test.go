package cli

import (
	"testing"

	"github.com/driangle/taskmd/sdk/go/feed"
)

func TestExtractFrontmatterFields(t *testing.T) {
	content := "---\nid: 042\ntitle: \"Add Auth\"\nstatus: pending\npriority: medium\n---\n# Body"
	fields := feed.ExtractFrontmatterFields(content)

	if fields["id"] != "042" {
		t.Errorf("expected id=042, got %q", fields["id"])
	}
	if fields["title"] != "Add Auth" {
		t.Errorf("expected title=Add Auth, got %q", fields["title"])
	}
	if fields["status"] != "pending" {
		t.Errorf("expected status=pending, got %q", fields["status"])
	}
	if fields["priority"] != "medium" {
		t.Errorf("expected priority=medium, got %q", fields["priority"])
	}
}

func TestExtractFrontmatterFields_NoFrontmatter(t *testing.T) {
	fields := feed.ExtractFrontmatterFields("# Just markdown")
	if len(fields) != 0 {
		t.Errorf("expected no fields, got %d", len(fields))
	}
}

func TestExtractSubtasks(t *testing.T) {
	content := "---\nid: 042\n---\n# Task\n\n- [ ] Add tests\n- [x] Write docs\n- [ ] Deploy\n"
	subtasks := feed.ExtractSubtasks(content)

	if len(subtasks) != 3 {
		t.Fatalf("expected 3 subtasks, got %d", len(subtasks))
	}
	if subtasks["Add tests"] != false {
		t.Error("expected 'Add tests' unchecked")
	}
	if subtasks["Write docs"] != true {
		t.Error("expected 'Write docs' checked")
	}
	if subtasks["Deploy"] != false {
		t.Error("expected 'Deploy' unchecked")
	}
}

func TestExtractSubtasks_NoFrontmatter(t *testing.T) {
	content := "# Task\n\n- [x] Done\n- [ ] Not done\n"
	subtasks := feed.ExtractSubtasks(content)
	if len(subtasks) != 2 {
		t.Fatalf("expected 2 subtasks, got %d", len(subtasks))
	}
}

func TestAnalyzeDiff_StatusChange(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: pending\npriority: medium\n---\n# Task"
	newContent := "---\nid: 042\nstatus: in-progress\npriority: medium\n---\n# Task"

	fieldChanges, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 1 {
		t.Fatalf("expected 1 field change, got %d", len(fieldChanges))
	}
	if fieldChanges[0].Field != "status" {
		t.Errorf("expected field 'status', got %q", fieldChanges[0].Field)
	}
	if fieldChanges[0].OldValue != "pending" {
		t.Errorf("expected old value 'pending', got %q", fieldChanges[0].OldValue)
	}
	if fieldChanges[0].NewValue != "in-progress" {
		t.Errorf("expected new value 'in-progress', got %q", fieldChanges[0].NewValue)
	}
	if len(subtaskChanges) != 0 {
		t.Errorf("expected no subtask changes, got %d", len(subtaskChanges))
	}
}

func TestAnalyzeDiff_PriorityChange(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: pending\npriority: medium\n---\n# Task"
	newContent := "---\nid: 042\nstatus: pending\npriority: high\n---\n# Task"

	fieldChanges, _ := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 1 {
		t.Fatalf("expected 1 field change, got %d", len(fieldChanges))
	}
	if fieldChanges[0].Field != "priority" {
		t.Errorf("expected field 'priority', got %q", fieldChanges[0].Field)
	}
	if fieldChanges[0].OldValue != "medium" || fieldChanges[0].NewValue != "high" {
		t.Errorf("expected medium -> high, got %s -> %s", fieldChanges[0].OldValue, fieldChanges[0].NewValue)
	}
}

func TestAnalyzeDiff_SubtaskCompletion(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: in-progress\n---\n# Task\n\n- [ ] Add tests\n- [ ] Write docs\n"
	newContent := "---\nid: 042\nstatus: in-progress\n---\n# Task\n\n- [x] Add tests\n- [ ] Write docs\n"

	fieldChanges, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 0 {
		t.Errorf("expected no field changes, got %d", len(fieldChanges))
	}
	if len(subtaskChanges) != 1 {
		t.Fatalf("expected 1 subtask change, got %d", len(subtaskChanges))
	}
	if subtaskChanges[0].Text != "Add tests" {
		t.Errorf("expected subtask 'Add tests', got %q", subtaskChanges[0].Text)
	}
	if !subtaskChanges[0].Done {
		t.Error("expected subtask to be marked done")
	}
}

func TestAnalyzeDiff_MultipleSubtaskCompletions(t *testing.T) {
	oldContent := "---\nid: 042\n---\n# Task\n\n- [ ] Add tests\n- [ ] Write docs\n- [ ] Deploy\n"
	newContent := "---\nid: 042\n---\n# Task\n\n- [x] Add tests\n- [ ] Write docs\n- [x] Deploy\n"

	_, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(subtaskChanges) != 2 {
		t.Fatalf("expected 2 subtask changes, got %d", len(subtaskChanges))
	}
	// Sorted by text
	if subtaskChanges[0].Text != "Add tests" {
		t.Errorf("expected first subtask 'Add tests', got %q", subtaskChanges[0].Text)
	}
	if subtaskChanges[1].Text != "Deploy" {
		t.Errorf("expected second subtask 'Deploy', got %q", subtaskChanges[1].Text)
	}
}

func TestAnalyzeDiff_MixedChanges(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: pending\npriority: medium\n---\n# Task\n\n- [ ] Add tests\n"
	newContent := "---\nid: 042\nstatus: in-progress\npriority: high\n---\n# Task\n\n- [x] Add tests\n"

	fieldChanges, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 2 {
		t.Fatalf("expected 2 field changes, got %d", len(fieldChanges))
	}
	if len(subtaskChanges) != 1 {
		t.Fatalf("expected 1 subtask change, got %d", len(subtaskChanges))
	}

	// Field changes are sorted by key name
	if fieldChanges[0].Field != "priority" {
		t.Errorf("expected first field change 'priority', got %q", fieldChanges[0].Field)
	}
	if fieldChanges[1].Field != "status" {
		t.Errorf("expected second field change 'status', got %q", fieldChanges[1].Field)
	}
}

func TestAnalyzeDiff_NoChanges(t *testing.T) {
	content := "---\nid: 042\nstatus: pending\n---\n# Task\n\n- [ ] Add tests\n"

	fieldChanges, subtaskChanges := feed.AnalyzeDiff(content, content)

	if len(fieldChanges) != 0 {
		t.Errorf("expected no field changes, got %d", len(fieldChanges))
	}
	if len(subtaskChanges) != 0 {
		t.Errorf("expected no subtask changes, got %d", len(subtaskChanges))
	}
}

func TestAnalyzeDiff_NoFrontmatter(t *testing.T) {
	oldContent := "# Just markdown"
	newContent := "# Updated markdown"

	fieldChanges, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 0 {
		t.Errorf("expected no field changes, got %d", len(fieldChanges))
	}
	if len(subtaskChanges) != 0 {
		t.Errorf("expected no subtask changes, got %d", len(subtaskChanges))
	}
}

func TestAnalyzeDiff_SubtaskUnchecked(t *testing.T) {
	oldContent := "---\nid: 042\n---\n# Task\n\n- [x] Add tests\n"
	newContent := "---\nid: 042\n---\n# Task\n\n- [ ] Add tests\n"

	_, subtaskChanges := feed.AnalyzeDiff(oldContent, newContent)

	if len(subtaskChanges) != 1 {
		t.Fatalf("expected 1 subtask change, got %d", len(subtaskChanges))
	}
	if subtaskChanges[0].Done {
		t.Error("expected subtask to be unchecked")
	}
}

func TestAnalyzeDiff_NewFieldAdded(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: pending\n---\n# Task"
	newContent := "---\nid: 042\nstatus: pending\npriority: high\n---\n# Task"

	fieldChanges, _ := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 1 {
		t.Fatalf("expected 1 field change, got %d", len(fieldChanges))
	}
	if fieldChanges[0].Field != "priority" {
		t.Errorf("expected field 'priority', got %q", fieldChanges[0].Field)
	}
	if fieldChanges[0].OldValue != "" {
		t.Errorf("expected empty old value, got %q", fieldChanges[0].OldValue)
	}
	if fieldChanges[0].NewValue != "high" {
		t.Errorf("expected new value 'high', got %q", fieldChanges[0].NewValue)
	}
}

func TestAnalyzeDiff_StatusToCompleted(t *testing.T) {
	oldContent := "---\nid: 042\nstatus: in-progress\n---\n# Task"
	newContent := "---\nid: 042\nstatus: completed\n---\n# Task"

	fieldChanges, _ := feed.AnalyzeDiff(oldContent, newContent)

	if len(fieldChanges) != 1 {
		t.Fatalf("expected 1 field change, got %d", len(fieldChanges))
	}
	if fieldChanges[0].Field != "status" {
		t.Errorf("expected field 'status', got %q", fieldChanges[0].Field)
	}
	if fieldChanges[0].NewValue != "completed" {
		t.Errorf("expected new value 'completed', got %q", fieldChanges[0].NewValue)
	}
}
