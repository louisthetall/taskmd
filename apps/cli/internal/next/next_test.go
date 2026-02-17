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
