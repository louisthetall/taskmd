const PRIORITY_RANK: Record<string, number> = {
  critical: 0,
  high: 1,
  medium: 2,
  low: 3,
};
const UNRANKED = 4;

function priorityRank(priority: string | null | undefined): number {
  if (priority == null) return UNRANKED;
  return PRIORITY_RANK[priority] ?? UNRANKED;
}

export function comparePriority(a: string | null | undefined, b: string | null | undefined): number {
  return priorityRank(a) - priorityRank(b);
}

export function countUnmetDependencies(
  deps: string[] | null,
  taskStatusMap: Map<string, string> | undefined,
): number {
  if (!deps) return 0;
  if (!taskStatusMap) return deps.length;
  return deps.filter((id) => taskStatusMap.get(id) !== "completed").length;
}

export function compareBlocked(
  depsA: string[] | null,
  depsB: string[] | null,
  taskStatusMap: Map<string, string> | undefined,
): number {
  return countUnmetDependencies(depsA, taskStatusMap) - countUnmetDependencies(depsB, taskStatusMap);
}
