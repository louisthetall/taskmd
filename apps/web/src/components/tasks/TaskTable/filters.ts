import type { Task } from "../../../api/types.ts";
import { STATUSES, PRIORITIES, EFFORTS, TYPES } from "./constants.ts";

export interface FilterState {
  selectedStatuses: Set<string>;
  selectedPriorities: Set<string>;
  selectedTypes: Set<string>;
  selectedTags: Set<string>;
  selectedEffort: Set<string>;
  globalFilter: string;
}

export function applyFilters(tasks: Task[], filters: FilterState): Task[] {
  return tasks.filter((task) => {
    if (!filters.selectedStatuses.has(task.status)) return false;
    if (task.priority && !filters.selectedPriorities.has(task.priority))
      return false;
    if (task.type && !filters.selectedTypes.has(task.type)) return false;
    if (filters.selectedTags.size > 0) {
      if (!task.tags || !task.tags.some((t) => filters.selectedTags.has(t)))
        return false;
    }
    if (filters.selectedEffort.size < EFFORTS.length) {
      if (!task.effort || !filters.selectedEffort.has(task.effort)) return false;
    }
    return true;
  });
}

export function hasActiveFilters(filters: FilterState): boolean {
  return (
    filters.selectedStatuses.size !== STATUSES.length ||
    filters.selectedPriorities.size !== PRIORITIES.length ||
    filters.selectedTypes.size !== TYPES.length ||
    filters.selectedTags.size > 0 ||
    filters.selectedEffort.size !== EFFORTS.length ||
    filters.globalFilter !== ""
  );
}

export function defaultFilterState(): FilterState {
  return {
    selectedStatuses: new Set(STATUSES),
    selectedPriorities: new Set(PRIORITIES),
    selectedTypes: new Set(TYPES),
    selectedTags: new Set<string>(),
    selectedEffort: new Set(EFFORTS),
    globalFilter: "",
  };
}
