package model

import (
	"strings"
	"testing"
)

func TestValidateVerifySteps(t *testing.T) {
	tests := []struct {
		name       string
		steps      []VerifyStep
		wantErrs   int
		wantSubstr string
	}{
		{
			name:     "valid bash step",
			steps:    []VerifyStep{{Type: "bash", Run: "make test"}},
			wantErrs: 0,
		},
		{
			name:     "valid assert step",
			steps:    []VerifyStep{{Type: "assert", Check: "file exists"}},
			wantErrs: 0,
		},
		{
			name:       "missing type field",
			steps:      []VerifyStep{{Run: "make test"}},
			wantErrs:   1,
			wantSubstr: "missing required field 'type'",
		},
		{
			name:       "bash step missing run",
			steps:      []VerifyStep{{Type: "bash"}},
			wantErrs:   1,
			wantSubstr: "bash step missing required field 'run'",
		},
		{
			name:       "assert step missing check",
			steps:      []VerifyStep{{Type: "assert"}},
			wantErrs:   1,
			wantSubstr: "assert step missing required field 'check'",
		},
		{
			name:     "empty slice",
			steps:    []VerifyStep{},
			wantErrs: 0,
		},
		{
			name: "multiple steps with mixed valid and invalid",
			steps: []VerifyStep{
				{Type: "bash", Run: "make test"},
				{Type: "bash"},
				{Type: "assert", Check: "ok"},
				{Run: "orphan"},
			},
			wantErrs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateVerifySteps(tt.steps)
			if len(errs) != tt.wantErrs {
				t.Errorf("ValidateVerifySteps() returned %d errors, want %d: %v", len(errs), tt.wantErrs, errs)
			}
			if tt.wantSubstr != "" {
				assertContainsSubstr(t, errs, tt.wantSubstr)
			}
		})
	}
}

func assertContainsSubstr(t *testing.T, errs []string, substr string) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e, substr) {
			return
		}
	}
	t.Errorf("expected error containing %q, got %v", substr, errs)
}

func TestTask_IsValid(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want bool
	}{
		{
			name: "both ID and Title set",
			task: Task{ID: "001", Title: "My Task"},
			want: true,
		},
		{
			name: "missing ID",
			task: Task{Title: "My Task"},
			want: false,
		},
		{
			name: "missing Title",
			task: Task{ID: "001"},
			want: false,
		},
		{
			name: "both empty",
			task: Task{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.IsValid(); got != tt.want {
				t.Errorf("Task.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_GetGroup(t *testing.T) {
	t.Run("with group set", func(t *testing.T) {
		task := Task{Group: "backend"}
		if got := task.GetGroup(); got != "backend" {
			t.Errorf("GetGroup() = %q, want %q", got, "backend")
		}
	})

	t.Run("with empty group", func(t *testing.T) {
		task := Task{}
		if got := task.GetGroup(); got != "" {
			t.Errorf("GetGroup() = %q, want empty string", got)
		}
	})
}

func TestStatus_IsResolved(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusCompleted, true},
		{StatusCancelled, true},
		{StatusPending, false},
		{StatusInProgress, false},
		{StatusInReview, false},
		{StatusBlocked, false},
		{Status("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsResolved(); got != tt.expected {
				t.Errorf("Status(%q).IsResolved() = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}
