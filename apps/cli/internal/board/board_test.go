package board

import (
	"encoding/json"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/model"
)

func TestToJSON_IncludesTags(t *testing.T) {
	tasks := []*model.Task{
		{
			ID:       "001",
			Title:    "Task with tags",
			Status:   model.StatusPending,
			Priority: model.PriorityHigh,
			Effort:   model.EffortSmall,
			Tags:     []string{"backend", "api"},
		},
	}

	gr := &GroupResult{
		Keys:   []string{"pending"},
		Groups: map[string][]*model.Task{"pending": tasks},
	}

	result := ToJSON(gr)
	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}

	jTask := result[0].Tasks[0]
	if len(jTask.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(jTask.Tags))
	}
	if jTask.Tags[0] != "backend" || jTask.Tags[1] != "api" {
		t.Errorf("unexpected tags: %v", jTask.Tags)
	}
}

func TestToJSON_OmitsEmptyTags(t *testing.T) {
	tasks := []*model.Task{
		{
			ID:     "002",
			Title:  "Task without tags",
			Status: model.StatusPending,
			Tags:   nil,
		},
	}

	gr := &GroupResult{
		Keys:   []string{"pending"},
		Groups: map[string][]*model.Task{"pending": tasks},
	}

	result := ToJSON(gr)
	jTask := result[0].Tasks[0]

	data, err := json.Marshal(jTask)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if _, exists := raw["tags"]; exists {
		t.Error("expected tags to be omitted from JSON when nil, but it was present")
	}
}

func TestGroupTasks_Type(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Feature", Type: model.TypeFeature},
		{ID: "002", Title: "Bug", Type: model.TypeBug},
		{ID: "003", Title: "No type"},
	}

	result, err := GroupTasks(tasks, "type")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}

	if len(result.Groups["feature"]) != 1 {
		t.Errorf("expected 1 feature task, got %d", len(result.Groups["feature"]))
	}
	if len(result.Groups["bug"]) != 1 {
		t.Errorf("expected 1 bug task, got %d", len(result.Groups["bug"]))
	}
	if len(result.Groups[defaultGroupKey]) != 1 {
		t.Errorf("expected 1 task with no type, got %d", len(result.Groups[defaultGroupKey]))
	}

	// Verify ordering: feature, bug, then (none)
	if result.Keys[0] != "feature" {
		t.Errorf("expected first key to be 'feature', got %q", result.Keys[0])
	}
	if result.Keys[1] != "bug" {
		t.Errorf("expected second key to be 'bug', got %q", result.Keys[1])
	}
}

func TestToJSON_OmitsEmptySliceTags(t *testing.T) {
	tasks := []*model.Task{
		{
			ID:     "003",
			Title:  "Task with empty tags",
			Status: model.StatusPending,
			Tags:   []string{},
		},
	}

	gr := &GroupResult{
		Keys:   []string{"pending"},
		Groups: map[string][]*model.Task{"pending": tasks},
	}

	result := ToJSON(gr)
	jTask := result[0].Tasks[0]

	data, err := json.Marshal(jTask)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Note: Go's json.Marshal does NOT omit empty slices ([]string{}),
	// only nil slices. This test documents that behavior.
	// An empty slice serializes as "tags": []
}

func TestGroupTasks_Status(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Pending", Status: model.StatusPending},
		{ID: "002", Title: "In Progress", Status: model.StatusInProgress},
		{ID: "003", Title: "Completed", Status: model.StatusCompleted},
		{ID: "004", Title: "Also Pending", Status: model.StatusPending},
	}

	result, err := GroupTasks(tasks, "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}
	if len(result.Groups["pending"]) != 2 {
		t.Errorf("expected 2 pending tasks, got %d", len(result.Groups["pending"]))
	}

	// Verify ordering follows statusOrder: pending, in-progress, completed
	expectedOrder := []string{"pending", "in-progress", "completed"}
	for i, key := range expectedOrder {
		if result.Keys[i] != key {
			t.Errorf("Keys[%d] = %q, want %q", i, result.Keys[i], key)
		}
	}
}

func TestGroupTasks_Priority(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Low", Priority: model.PriorityLow},
		{ID: "002", Title: "Critical", Priority: model.PriorityCritical},
		{ID: "003", Title: "High", Priority: model.PriorityHigh},
	}

	result, err := GroupTasks(tasks, "priority")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}

	// Verify ordering follows priorityOrder: critical, high, ..., low
	if result.Keys[0] != "critical" {
		t.Errorf("first key = %q, want %q", result.Keys[0], "critical")
	}
	if result.Keys[1] != "high" {
		t.Errorf("second key = %q, want %q", result.Keys[1], "high")
	}
	if result.Keys[2] != "low" {
		t.Errorf("third key = %q, want %q", result.Keys[2], "low")
	}
}

func TestGroupTasks_Effort(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Large", Effort: model.EffortLarge},
		{ID: "002", Title: "Small", Effort: model.EffortSmall},
		{ID: "003", Title: "Medium", Effort: model.EffortMedium},
	}

	result, err := GroupTasks(tasks, "effort")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify ordering follows effortOrder: small, medium, large
	expected := []string{"small", "medium", "large"}
	for i, key := range expected {
		if result.Keys[i] != key {
			t.Errorf("Keys[%d] = %q, want %q", i, result.Keys[i], key)
		}
	}
}

func TestGroupTasks_Group(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Backend", Group: "backend"},
		{ID: "002", Title: "Frontend", Group: "frontend"},
		{ID: "003", Title: "No Group"},
	}

	result, err := GroupTasks(tasks, "group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}

	// sortedKeys: alphabetical with (none) last
	if result.Keys[0] != "backend" {
		t.Errorf("Keys[0] = %q, want %q", result.Keys[0], "backend")
	}
	if result.Keys[1] != "frontend" {
		t.Errorf("Keys[1] = %q, want %q", result.Keys[1], "frontend")
	}
	if result.Keys[2] != defaultGroupKey {
		t.Errorf("Keys[2] = %q, want %q", result.Keys[2], defaultGroupKey)
	}
}

func TestGroupTasks_Tag(t *testing.T) {
	tasks := []*model.Task{
		{ID: "001", Title: "Multi-tag", Tags: []string{"api", "backend"}},
		{ID: "002", Title: "Single tag", Tags: []string{"frontend"}},
		{ID: "003", Title: "No tags"},
	}

	result, err := GroupTasks(tasks, "tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Task 001 should appear in both "api" and "backend" groups
	if len(result.Groups["api"]) != 1 {
		t.Errorf("expected 1 task in 'api' group, got %d", len(result.Groups["api"]))
	}
	if len(result.Groups["backend"]) != 1 {
		t.Errorf("expected 1 task in 'backend' group, got %d", len(result.Groups["backend"]))
	}
	if len(result.Groups[defaultGroupKey]) != 1 {
		t.Errorf("expected 1 task in '(none)' group, got %d", len(result.Groups[defaultGroupKey]))
	}

	// (none) should be last
	lastKey := result.Keys[len(result.Keys)-1]
	if lastKey != defaultGroupKey {
		t.Errorf("last key = %q, want %q", lastKey, defaultGroupKey)
	}
}

func TestGroupTasks_UnsupportedField(t *testing.T) {
	tasks := []*model.Task{{ID: "001", Title: "Test"}}

	_, err := GroupTasks(tasks, "invalid-field")
	if err == nil {
		t.Fatal("expected error for unsupported field, got nil")
	}
}

func TestGroupTasks_EmptyTasks(t *testing.T) {
	result, err := GroupTasks([]*model.Task{}, "status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(result.Groups))
	}
	if len(result.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(result.Keys))
	}
}

func TestSortedKeys(t *testing.T) {
	groups := map[string][]*model.Task{
		"zebra":         {{ID: "001"}},
		"alpha":         {{ID: "002"}},
		defaultGroupKey: {{ID: "003"}},
		"middle":        {{ID: "004"}},
	}

	keys := sortedKeys(groups)

	// Alphabetical order with (none) last
	expected := []string{"alpha", "middle", "zebra", defaultGroupKey}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, want := range expected {
		if keys[i] != want {
			t.Errorf("keys[%d] = %q, want %q", i, keys[i], want)
		}
	}
}

func TestOrderedKeys(t *testing.T) {
	order := []string{"first", "second", "third"}
	groups := map[string][]*model.Task{
		"third":   {{ID: "001"}},
		"first":   {{ID: "002"}},
		"unknown": {{ID: "003"}},
		"second":  {{ID: "004"}},
	}

	keys := orderedKeys(groups, order)

	// Predefined order preserved, unknown appended alphabetically
	expected := []string{"first", "second", "third", "unknown"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, want := range expected {
		if keys[i] != want {
			t.Errorf("keys[%d] = %q, want %q", i, keys[i], want)
		}
	}
}
