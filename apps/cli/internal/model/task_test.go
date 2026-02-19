package model

import "testing"

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
