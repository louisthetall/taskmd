import { STATUS_COLORS, PRIORITY_COLORS, TYPE_COLORS } from "./constants.ts";

export function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={`px-2 py-0.5 text-xs font-medium rounded-full whitespace-nowrap ${STATUS_COLORS[status] ?? "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"}`}
    >
      {status}
    </span>
  );
}

export function PriorityBadge({ priority }: { priority: string }) {
  return (
    <span
      className={`px-2 py-0.5 text-xs font-medium rounded-full ${PRIORITY_COLORS[priority] ?? "bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400"}`}
    >
      {priority}
    </span>
  );
}

export function TypeBadge({ type: taskType }: { type: string }) {
  return (
    <span
      className={`px-2 py-0.5 text-xs font-medium rounded-full ${TYPE_COLORS[taskType] ?? "bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400"}`}
    >
      {taskType}
    </span>
  );
}

export function MilestoneBadge({ milestone }: { milestone: string }) {
  return (
    <span
      className="px-2 py-0.5 text-xs font-medium rounded-full bg-teal-100 text-teal-700 ring-1 ring-teal-300 dark:bg-teal-900/30 dark:text-teal-300 dark:ring-teal-700"
    >
      {milestone}
    </span>
  );
}

export function BlockedStatusBadge({
  dependencies,
  taskStatusMap,
}: {
  dependencies: string[] | null;
  taskStatusMap?: Map<string, string>;
}) {
  const unmetDeps = dependencies?.filter(
    (id) => !taskStatusMap || taskStatusMap.get(id) !== "completed",
  ) ?? [];
  const blockedByCount = unmetDeps.length;
  const isBlocked = blockedByCount > 0;

  if (!isBlocked) {
    return (
      <span
        className="px-2 py-0.5 text-xs font-medium rounded-full bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300 inline-flex items-center gap-1"
        aria-label="Task is ready to work on"
      >
        <span aria-hidden="true">✓</span>
        <span className="hidden sm:inline">Ready</span>
      </span>
    );
  }

  const tooltipText = `Blocked by: ${unmetDeps.join(", ")}`;

  return (
    <span
      className="px-2 py-0.5 text-xs font-medium rounded-full bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300 inline-flex items-center gap-1 cursor-help"
      title={tooltipText}
      aria-label={tooltipText}
    >
      <span aria-hidden="true">⚠</span>
      <span>
        <span className="hidden sm:inline">Blocked </span>({blockedByCount})
      </span>
    </span>
  );
}
