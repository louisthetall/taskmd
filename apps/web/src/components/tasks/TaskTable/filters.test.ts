import { describe, it, expect } from "vitest";
import type { Task } from "../../../api/types.ts";
import { applyFilters, hasActiveFilters, defaultFilterState } from "./filters.ts";
import { STATUSES, PRIORITIES, EFFORTS, TYPES } from "./constants.ts";

function makeTask(overrides: Partial<Task> = {}): Task {
  return {
    id: "001",
    title: "Test task",
    status: "pending",
    priority: "medium",
    effort: "small",
    type: "feature",
    dependencies: null,
    tags: null,
    phase: "",
    group: "",
    owner: "",
    parent: "",
    created: "2026-01-01",
    body: "",
    file_path: "tasks/001-test.md",
    ...overrides,
  };
}

const sampleTasks: Task[] = [
  makeTask({ id: "001", status: "pending", priority: "high", type: "feature", tags: ["api"], effort: "small" }),
  makeTask({ id: "002", status: "in-progress", priority: "medium", type: "bug", tags: ["web", "api"], effort: "medium" }),
  makeTask({ id: "003", status: "completed", priority: "low", type: "chore", tags: ["docs"], effort: "large" }),
  makeTask({ id: "004", status: "blocked", priority: "critical", type: "feature", tags: null, effort: "" }),
  makeTask({ id: "005", status: "pending", priority: "", type: "", tags: null, effort: "" }),
];

describe("applyFilters", () => {
  it("returns all tasks with default filter state", () => {
    const result = applyFilters(sampleTasks, defaultFilterState());
    expect(result).toHaveLength(sampleTasks.length);
  });

  it("filters by single status", () => {
    const filters = { ...defaultFilterState(), selectedStatuses: new Set(["pending"]) };
    const result = applyFilters(sampleTasks, filters);
    expect(result.map((t) => t.id)).toEqual(["001", "005"]);
  });

  it("filters by in-review status", () => {
    const tasksWithReview = [
      ...sampleTasks,
      makeTask({ id: "006", status: "in-review", priority: "medium", type: "feature" }),
    ];
    const filters = { ...defaultFilterState(), selectedStatuses: new Set(["in-review"]) };
    const result = applyFilters(tasksWithReview, filters);
    expect(result.map((t) => t.id)).toEqual(["006"]);
  });

  it("filters by multiple statuses", () => {
    const filters = { ...defaultFilterState(), selectedStatuses: new Set(["pending", "blocked"]) };
    const result = applyFilters(sampleTasks, filters);
    expect(result.map((t) => t.id)).toEqual(["001", "004", "005"]);
  });

  it("filters by priority", () => {
    const filters = { ...defaultFilterState(), selectedPriorities: new Set(["high"]) };
    const result = applyFilters(sampleTasks, filters);
    // 001 has high priority, 005 has empty priority (passes through)
    expect(result.map((t) => t.id)).toEqual(["001", "005"]);
  });

  it("filters by type", () => {
    const filters = { ...defaultFilterState(), selectedTypes: new Set(["bug"]) };
    const result = applyFilters(sampleTasks, filters);
    // 002 is a bug, 005 has empty type (passes through)
    expect(result.map((t) => t.id)).toEqual(["002", "005"]);
  });

  it("filters by tags (OR among selected tags)", () => {
    const filters = { ...defaultFilterState(), selectedTags: new Set(["api"]) };
    const result = applyFilters(sampleTasks, filters);
    // 001 has [api], 002 has [web, api]
    expect(result.map((t) => t.id)).toEqual(["001", "002"]);
  });

  it("filters by multiple tags (OR logic)", () => {
    const filters = { ...defaultFilterState(), selectedTags: new Set(["docs", "web"]) };
    const result = applyFilters(sampleTasks, filters);
    // 002 has [web, api], 003 has [docs]
    expect(result.map((t) => t.id)).toEqual(["002", "003"]);
  });

  it("filters by effort", () => {
    const filters = { ...defaultFilterState(), selectedEffort: new Set(["small"]) };
    const result = applyFilters(sampleTasks, filters);
    expect(result.map((t) => t.id)).toEqual(["001"]);
  });

  it("applies intersection of multiple filters (status AND priority AND type)", () => {
    const filters = {
      ...defaultFilterState(),
      selectedStatuses: new Set(["pending", "in-progress"]),
      selectedPriorities: new Set(["medium", "high"]),
      selectedTypes: new Set(["feature", "bug"]),
    };
    const result = applyFilters(sampleTasks, filters);
    // 001: pending + high + feature -> yes
    // 002: in-progress + medium + bug -> yes
    // 005: pending + empty priority (passes) + empty type (passes) -> yes
    expect(result.map((t) => t.id)).toEqual(["001", "002", "005"]);
  });

  it("applies intersection of all filter criteria", () => {
    const filters = {
      selectedStatuses: new Set(["pending", "in-progress"]),
      selectedPriorities: new Set(["high", "medium"]),
      selectedTypes: new Set(["feature", "bug"]),
      selectedTags: new Set(["api"]),
      selectedEffort: new Set(["small"]),
      selectedPhases: new Set<string>(),
      globalFilter: "",
    };
    const result = applyFilters(sampleTasks, filters);
    // Only 001: pending + high + feature + has "api" tag + small effort
    expect(result.map((t) => t.id)).toEqual(["001"]);
  });

  it("filters by selected phases", () => {
    const tasksWithPhases = [
      makeTask({ id: "001", phase: "mvp" }),
      makeTask({ id: "002", phase: "v2" }),
      makeTask({ id: "003", phase: "" }),
    ];
    const filters = { ...defaultFilterState(), selectedPhases: new Set(["mvp"]) };
    const result = applyFilters(tasksWithPhases, filters);
    expect(result.map((t) => t.id)).toEqual(["001"]);
  });

  it("excludes tasks without phase when phase filter is active", () => {
    const tasksWithPhases = [
      makeTask({ id: "001", phase: "mvp" }),
      makeTask({ id: "002", phase: "" }),
    ];
    const filters = { ...defaultFilterState(), selectedPhases: new Set(["mvp"]) };
    const result = applyFilters(tasksWithPhases, filters);
    expect(result.map((t) => t.id)).toEqual(["001"]);
  });

  it("returns empty array when no tasks match", () => {
    const filters = { ...defaultFilterState(), selectedStatuses: new Set(["cancelled"]) };
    const result = applyFilters(sampleTasks, filters);
    expect(result).toEqual([]);
  });

  it("tasks without priority pass through priority filter", () => {
    const filters = { ...defaultFilterState(), selectedPriorities: new Set(["critical"]) };
    const result = applyFilters(sampleTasks, filters);
    // 004 has critical priority, 005 has empty priority (passes)
    expect(result.map((t) => t.id)).toEqual(["004", "005"]);
  });

  it("tasks without type pass through type filter", () => {
    const filters = { ...defaultFilterState(), selectedTypes: new Set(["chore"]) };
    const result = applyFilters(sampleTasks, filters);
    // 003 is chore, 005 has empty type (passes)
    expect(result.map((t) => t.id)).toEqual(["003", "005"]);
  });

  it("tasks without tags are excluded when tag filter is active", () => {
    const filters = { ...defaultFilterState(), selectedTags: new Set(["api"]) };
    const result = applyFilters(sampleTasks, filters);
    // 004 and 005 have no tags -> excluded
    expect(result.every((t) => t.tags !== null)).toBe(true);
  });

  it("tasks without effort are excluded when effort filter is active", () => {
    const filters = { ...defaultFilterState(), selectedEffort: new Set(["medium"]) };
    const result = applyFilters(sampleTasks, filters);
    expect(result.map((t) => t.id)).toEqual(["002"]);
  });
});

describe("hasActiveFilters", () => {
  it("returns false for default filter state", () => {
    expect(hasActiveFilters(defaultFilterState())).toBe(false);
  });

  it("returns true when a status is deselected", () => {
    const filters = { ...defaultFilterState(), selectedStatuses: new Set(STATUSES.slice(0, 3)) };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns true when a priority is deselected", () => {
    const filters = { ...defaultFilterState(), selectedPriorities: new Set(PRIORITIES.slice(0, 2)) };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns true when a type is deselected", () => {
    const filters = { ...defaultFilterState(), selectedTypes: new Set(TYPES.slice(0, 3)) };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns true when tags are selected", () => {
    const filters = { ...defaultFilterState(), selectedTags: new Set(["api"]) };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns true when effort is selected", () => {
    const filters = { ...defaultFilterState(), selectedEffort: new Set(["small"]) };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns true when global filter has text", () => {
    const filters = { ...defaultFilterState(), globalFilter: "search" };
    expect(hasActiveFilters(filters)).toBe(true);
  });

  it("returns false when all statuses/priorities/types selected and nothing else active", () => {
    const filters = {
      selectedStatuses: new Set(STATUSES),
      selectedPriorities: new Set(PRIORITIES),
      selectedTypes: new Set(TYPES),
      selectedTags: new Set<string>(),
      selectedEffort: new Set(EFFORTS),
      selectedPhases: new Set<string>(),
      globalFilter: "",
    };
    expect(hasActiveFilters(filters)).toBe(false);
  });
});

describe("defaultFilterState", () => {
  it("has all statuses selected", () => {
    expect(defaultFilterState().selectedStatuses).toEqual(new Set(STATUSES));
  });

  it("includes in-review in default statuses", () => {
    expect(defaultFilterState().selectedStatuses.has("in-review")).toBe(true);
  });

  it("has all priorities selected", () => {
    expect(defaultFilterState().selectedPriorities).toEqual(new Set(PRIORITIES));
  });

  it("has all types selected", () => {
    expect(defaultFilterState().selectedTypes).toEqual(new Set(TYPES));
  });

  it("has no tags selected", () => {
    expect(defaultFilterState().selectedTags.size).toBe(0);
  });

  it("has all efforts selected", () => {
    expect(defaultFilterState().selectedEffort).toEqual(new Set(EFFORTS));
  });

  it("has empty global filter", () => {
    expect(defaultFilterState().globalFilter).toBe("");
  });
});
