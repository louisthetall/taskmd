package next

import (
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/model"
)

func makeTask(id string, status model.Status, priority model.Priority, deps []string) *model.Task {
	return &model.Task{
		ID:           id,
		Title:        "Task " + id,
		Status:       status,
		Priority:     priority,
		Dependencies: deps,
	}
}

func makeTaskWithParent(id string, status model.Status, priority model.Priority, parent string) *model.Task {
	return &model.Task{
		ID:       id,
		Title:    "Task " + id,
		Status:   status,
		Priority: priority,
		Parent:   parent,
	}
}

func TestRecommend_ArchivedCompletedDepSatisfied(t *testing.T) {
	// Task 002 depends on 001, but 001 is archived and completed.
	// 002 should be actionable.
	tasks := []*model.Task{
		makeTask("002", model.StatusPending, model.PriorityHigh, []string{"001"}),
	}
	archived := []*model.Task{
		makeTask("001", model.StatusCompleted, model.PriorityHigh, nil),
	}

	recs, err := Recommend(tasks, Options{
		Limit:         10,
		ArchivedTasks: archived,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recs) != 1 {
		t.Fatalf("Expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].ID != "002" {
		t.Errorf("Expected task 002, got %s", recs[0].ID)
	}
}

func TestRecommend_ArchivedNonCompletedDepBlocks(t *testing.T) {
	// Task 002 depends on 001, which is archived but still pending.
	// 002 should be blocked.
	tasks := []*model.Task{
		makeTask("002", model.StatusPending, model.PriorityHigh, []string{"001"}),
	}
	archived := []*model.Task{
		makeTask("001", model.StatusPending, model.PriorityHigh, nil),
	}

	recs, err := Recommend(tasks, Options{
		Limit:         10,
		ArchivedTasks: archived,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(recs) != 0 {
		t.Errorf("Expected 0 recommendations (dep not completed), got %d", len(recs))
	}
}

func TestRecommend_ArchivedTasksNotRecommended(t *testing.T) {
	// Archived tasks should never appear in recommendations, even if actionable.
	tasks := []*model.Task{
		makeTask("002", model.StatusPending, model.PriorityHigh, nil),
	}
	archived := []*model.Task{
		makeTask("001", model.StatusPending, model.PriorityHigh, nil),
	}

	recs, err := Recommend(tasks, Options{
		Limit:         10,
		ArchivedTasks: archived,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, rec := range recs {
		if rec.ID == "001" {
			t.Error("Archived task 001 should not appear in recommendations")
		}
	}

	if len(recs) != 1 || recs[0].ID != "002" {
		t.Errorf("Expected only task 002, got %v", recs)
	}
}

func TestRecommend_ActiveTaskPrecedenceOverArchived(t *testing.T) {
	// If the same ID exists in both active and archived, active wins.
	tasks := []*model.Task{
		makeTask("001", model.StatusPending, model.PriorityHigh, nil),
		makeTask("002", model.StatusPending, model.PriorityMedium, []string{"001"}),
	}
	// Archived version has status=completed, but active version is pending.
	// Task 002 depends on 001 — since active 001 is pending, 002 should be blocked.
	archived := []*model.Task{
		makeTask("001", model.StatusCompleted, model.PriorityHigh, nil),
	}

	recs, err := Recommend(tasks, Options{
		Limit:         10,
		ArchivedTasks: archived,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 001 is active+pending → actionable. 002 depends on 001 (pending) → blocked.
	if len(recs) != 1 {
		t.Fatalf("Expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].ID != "001" {
		t.Errorf("Expected task 001, got %s", recs[0].ID)
	}
}

func TestScoreTask_LowChainDoesNotOutscoreMediumTask(t *testing.T) {
	// A low-priority task unblocking 5 low-priority tasks should NOT outscore
	// a standalone medium-priority task.
	criticalPath := map[string]bool{"low1": true}
	downstreamInfo := map[string]DownstreamInfo{
		"low1": {Count: 5, MaxPriority: model.PriorityLow},
		"med1": {Count: 0},
	}

	lowTask := &model.Task{ID: "low1", Priority: model.PriorityLow}
	medTask := &model.Task{ID: "med1", Priority: model.PriorityMedium}

	lowScore, _ := ScoreTask(lowTask, criticalPath, downstreamInfo)
	medScore, _ := ScoreTask(medTask, map[string]bool{}, downstreamInfo)

	if lowScore >= medScore {
		t.Errorf("Low-priority task with all-low downstream chain (score=%d) should not outscore standalone medium task (score=%d)",
			lowScore, medScore)
	}
}

func TestScoreTask_MixedChainGetsFullDownstreamBonus(t *testing.T) {
	// A low-priority task that unblocks a high-priority task should still get
	// a full downstream bonus (multiplier = 1.0).
	downstreamInfo := map[string]DownstreamInfo{
		"low1": {Count: 1, MaxPriority: model.PriorityHigh},
	}

	task := &model.Task{ID: "low1", Priority: model.PriorityLow}
	score, _ := ScoreTask(task, map[string]bool{}, downstreamInfo)

	// Expected: base low (10) + full downstream bonus (1 * 3 * 1.0 = 3) = 13
	expectedScore := ScorePriorityLow + 1*ScorePerDownstream
	if score != expectedScore {
		t.Errorf("Mixed chain score = %d, want %d (full downstream bonus for high-priority downstream)", score, expectedScore)
	}
}

func TestScoreTask_HighChainPreservesExistingBehavior(t *testing.T) {
	// A task on the critical path with high-priority downstream tasks should
	// get full bonuses (same as before the priority-aware change).
	criticalPath := map[string]bool{"t1": true}
	downstreamInfo := map[string]DownstreamInfo{
		"t1": {Count: 3, MaxPriority: model.PriorityCritical},
	}

	task := &model.Task{ID: "t1", Priority: model.PriorityHigh}
	score, _ := ScoreTask(task, criticalPath, downstreamInfo)

	// Expected: high (30) + critical path (15 * 1.0) + downstream (min(9,15) * 1.0 = 9) = 54
	expectedScore := ScorePriorityHigh + ScoreCriticalPath + min(3*ScorePerDownstream, ScoreDownstreamMax)
	if score != expectedScore {
		t.Errorf("High/critical chain score = %d, want %d", score, expectedScore)
	}
}

func TestRecommend_MediumTaskRanksAboveLowChain(t *testing.T) {
	// Integration test: an unblocked medium-priority task should rank higher than
	// a low-priority task whose entire downstream chain is low priority.
	//
	// low1 (low, no deps) -> low2 (low, depends on low1) -> low3 -> low4 -> low5
	// med1 (medium, no deps, standalone)
	tasks := []*model.Task{
		makeTask("low1", model.StatusPending, model.PriorityLow, nil),
		makeTask("low2", model.StatusPending, model.PriorityLow, []string{"low1"}),
		makeTask("low3", model.StatusPending, model.PriorityLow, []string{"low2"}),
		makeTask("low4", model.StatusPending, model.PriorityLow, []string{"low3"}),
		makeTask("low5", model.StatusPending, model.PriorityLow, []string{"low4"}),
		makeTask("med1", model.StatusPending, model.PriorityMedium, nil),
	}

	recs, err := Recommend(tasks, Options{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find positions
	var medRank, lowRank int
	for _, rec := range recs {
		if rec.ID == "med1" {
			medRank = rec.Rank
		}
		if rec.ID == "low1" {
			lowRank = rec.Rank
		}
	}

	if medRank == 0 {
		t.Fatal("med1 not found in recommendations")
	}
	if lowRank == 0 {
		t.Fatal("low1 not found in recommendations")
	}
	if medRank >= lowRank {
		t.Errorf("Medium task (rank %d) should rank above low task with all-low downstream chain (rank %d)", medRank, lowRank)
	}
}

func TestCalculateCriticalPathTasks_IgnoresCompletedDependencies(t *testing.T) {
	// Scenario: tasks with completed dependencies should not have inflated depth.
	//
	// Graph:
	//   A (completed, no deps)
	//   B (pending, depends on A)  — A is done, so B's remaining depth is 1
	//   C (pending, no deps)       — depth 1
	//   D (pending, depends on C)  — C is pending, real remaining chain depth 2
	//
	// The only real remaining dependency chain is C → D.
	// B should NOT be on the critical path because its dependency A is completed.
	tasks := []*model.Task{
		{ID: "A", Status: model.StatusCompleted, Dependencies: nil},
		{ID: "B", Status: model.StatusPending, Dependencies: []string{"A"}},
		{ID: "C", Status: model.StatusPending, Dependencies: nil},
		{ID: "D", Status: model.StatusPending, Dependencies: []string{"C"}},
	}

	taskMap := BuildTaskMap(tasks)
	criticalPath := CalculateCriticalPathTasks(tasks, taskMap)

	// C and D should be on the critical path (the only real remaining chain)
	if !criticalPath["C"] {
		t.Error("Expected task C to be on critical path")
	}
	if !criticalPath["D"] {
		t.Error("Expected task D to be on critical path")
	}

	// B should NOT be on the critical path — its dependency A is already completed
	if criticalPath["B"] {
		t.Error("Task B should NOT be on critical path: its dependency A is completed")
	}

	// A is completed and should not be on the critical path either
	if criticalPath["A"] {
		t.Error("Completed task A should NOT be on critical path")
	}
}

func TestCalculateCriticalPathTasks_PendingChainIsCritical(t *testing.T) {
	// When all tasks in a chain are pending, the longest chain is the critical path.
	//
	// Graph:
	//   001 (pending, no deps)         — depth 1
	//   002 (pending, depends on 001)  — depth 2
	//   003 (pending, depends on 002)  — depth 3
	//   004 (pending, no deps)         — depth 1
	//
	// Critical path: 001 → 002 → 003
	tasks := []*model.Task{
		{ID: "001", Status: model.StatusPending, Dependencies: nil},
		{ID: "002", Status: model.StatusPending, Dependencies: []string{"001"}},
		{ID: "003", Status: model.StatusPending, Dependencies: []string{"002"}},
		{ID: "004", Status: model.StatusPending, Dependencies: nil},
	}

	taskMap := BuildTaskMap(tasks)
	criticalPath := CalculateCriticalPathTasks(tasks, taskMap)

	for _, id := range []string{"001", "002", "003"} {
		if !criticalPath[id] {
			t.Errorf("Expected task %s to be on critical path", id)
		}
	}

	if criticalPath["004"] {
		t.Error("Task 004 should NOT be on critical path (shorter parallel path)")
	}
}

func TestCalculateCriticalPathTasks_MixedCompletedPendingChain(t *testing.T) {
	// A longer chain where early tasks are completed should have reduced effective depth.
	//
	// Graph:
	//   A (completed, no deps)
	//   B (completed, depends on A)
	//   C (pending, depends on B)    — B is done, so C's remaining depth is 1
	//   D (pending, no deps)         — depth 1
	//   E (pending, depends on D)    — depth 2
	//   F (pending, depends on E)    — depth 3
	//
	// Remaining chain D → E → F is longer than just C.
	// Critical path should be D → E → F only.
	tasks := []*model.Task{
		{ID: "A", Status: model.StatusCompleted, Dependencies: nil},
		{ID: "B", Status: model.StatusCompleted, Dependencies: []string{"A"}},
		{ID: "C", Status: model.StatusPending, Dependencies: []string{"B"}},
		{ID: "D", Status: model.StatusPending, Dependencies: nil},
		{ID: "E", Status: model.StatusPending, Dependencies: []string{"D"}},
		{ID: "F", Status: model.StatusPending, Dependencies: []string{"E"}},
	}

	taskMap := BuildTaskMap(tasks)
	criticalPath := CalculateCriticalPathTasks(tasks, taskMap)

	// D → E → F is the real critical path
	for _, id := range []string{"D", "E", "F"} {
		if !criticalPath[id] {
			t.Errorf("Expected task %s to be on critical path", id)
		}
	}

	// C should NOT be on critical path (only 1 remaining step, shorter than D→E→F)
	if criticalPath["C"] {
		t.Error("Task C should NOT be on critical path (shorter remaining chain)")
	}

	// Completed tasks should not be on critical path
	if criticalPath["A"] {
		t.Error("Completed task A should NOT be on critical path")
	}
	if criticalPath["B"] {
		t.Error("Completed task B should NOT be on critical path")
	}
}

func TestHasIncompleteChildren(t *testing.T) {
	tests := []struct {
		name     string
		task     *model.Task
		children []*model.Task
		expected bool
	}{
		{
			name:     "no children",
			task:     makeTask("P1", model.StatusPending, model.PriorityMedium, nil),
			children: nil,
			expected: false,
		},
		{
			name: "all children completed",
			task: makeTask("P1", model.StatusPending, model.PriorityMedium, nil),
			children: []*model.Task{
				makeTaskWithParent("C1", model.StatusCompleted, model.PriorityMedium, "P1"),
				makeTaskWithParent("C2", model.StatusCompleted, model.PriorityMedium, "P1"),
			},
			expected: false,
		},
		{
			name: "one pending child",
			task: makeTask("P1", model.StatusPending, model.PriorityMedium, nil),
			children: []*model.Task{
				makeTaskWithParent("C1", model.StatusCompleted, model.PriorityMedium, "P1"),
				makeTaskWithParent("C2", model.StatusPending, model.PriorityMedium, "P1"),
			},
			expected: true,
		},
		{
			name: "cancelled child counts as resolved",
			task: makeTask("P1", model.StatusPending, model.PriorityMedium, nil),
			children: []*model.Task{
				makeTaskWithParent("C1", model.StatusCancelled, model.PriorityMedium, "P1"),
			},
			expected: false,
		},
		{
			name: "in-progress child is incomplete",
			task: makeTask("P1", model.StatusPending, model.PriorityMedium, nil),
			children: []*model.Task{
				makeTaskWithParent("C1", model.StatusInProgress, model.PriorityMedium, "P1"),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childrenMap := BuildChildrenMap(tt.children)
			got := HasIncompleteChildren(tt.task, childrenMap)
			if got != tt.expected {
				t.Errorf("HasIncompleteChildren() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsActionable_WithChildren(t *testing.T) {
	parent := makeTask("P1", model.StatusPending, model.PriorityMedium, nil)
	child := makeTaskWithParent("C1", model.StatusPending, model.PriorityMedium, "P1")

	tasks := []*model.Task{parent, child}
	taskMap := BuildTaskMap(tasks)
	childrenMap := BuildChildrenMap(tasks)

	// Parent with incomplete child should not be actionable
	if IsActionable(parent, taskMap, childrenMap) {
		t.Error("Parent with incomplete child should not be actionable")
	}

	// Child itself should be actionable (no deps, no children of its own)
	if !IsActionable(child, taskMap, childrenMap) {
		t.Error("Child task should be actionable")
	}
}

func TestRecommend_ParentExcludedWhenChildrenIncomplete(t *testing.T) {
	parent := makeTask("P1", model.StatusPending, model.PriorityHigh, nil)
	child := makeTaskWithParent("C1", model.StatusPending, model.PriorityMedium, "P1")

	tasks := []*model.Task{parent, child}

	recs, err := Recommend(tasks, Options{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, rec := range recs {
		if rec.ID == "P1" {
			t.Error("Parent P1 should not be recommended while child C1 is incomplete")
		}
	}

	if len(recs) != 1 || recs[0].ID != "C1" {
		t.Errorf("Expected only child C1, got %v", recs)
	}
}

func TestRecommend_ParentIncludedWhenAllChildrenResolved(t *testing.T) {
	parent := makeTask("P1", model.StatusPending, model.PriorityHigh, nil)
	childCompleted := makeTaskWithParent("C1", model.StatusCompleted, model.PriorityMedium, "P1")
	childCancelled := makeTaskWithParent("C2", model.StatusCancelled, model.PriorityMedium, "P1")

	tasks := []*model.Task{parent, childCompleted, childCancelled}

	recs, err := Recommend(tasks, Options{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, rec := range recs {
		if rec.ID == "P1" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Parent P1 should be recommended when all children are resolved")
	}
}
