import {
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  useReactTable,
  type SortingState,
} from "@tanstack/react-table";
import { useState, useMemo, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import type { Task } from "../../api/types.ts";
import { STATUSES, PRIORITIES, EFFORTS, TYPES } from "./TaskTable/constants.ts";
import { FilterBar } from "./TaskTable/FilterBar.tsx";
import { createTaskColumns } from "./TaskTable/columns.tsx";
import { toggleInSet } from "./TaskTable/utils.ts";
import { applyFilters, hasActiveFilters as checkActiveFilters } from "./TaskTable/filters.ts";
import { MobileCardList } from "./TaskTable/MobileCardList.tsx";

interface TaskTableProps {
  tasks: Task[];
  initialTags?: string[];
  initialStatuses?: string[];
  initialPriorities?: string[];
  initialEffort?: string[];
}

export function TaskTable({ tasks, initialTags, initialStatuses, initialPriorities, initialEffort }: TaskTableProps) {
  const [, setSearchParams] = useSearchParams();
  const [sorting, setSorting] = useState<SortingState>([]);
  const [globalFilter, setGlobalFilter] = useState("");
  const [selectedStatuses, setSelectedStatuses] = useState<Set<string>>(
    () => initialStatuses && initialStatuses.length > 0 ? new Set(initialStatuses) : new Set(STATUSES),
  );
  const [selectedPriorities, setSelectedPriorities] = useState<Set<string>>(
    () => initialPriorities && initialPriorities.length > 0 ? new Set(initialPriorities) : new Set(PRIORITIES),
  );
  const [selectedTypes, setSelectedTypes] = useState<Set<string>>(
    new Set(TYPES),
  );
  const [selectedTags, setSelectedTags] = useState<Set<string>>(
    () => new Set(initialTags),
  );
  const [selectedEffort, setSelectedEffort] = useState<Set<string>>(
    () => initialEffort && initialEffort.length > 0 ? new Set(initialEffort) : new Set(EFFORTS),
  );

  const filterState = { selectedStatuses, selectedPriorities, selectedTypes, selectedTags, selectedEffort, globalFilter };
  const hasActiveFilters = checkActiveFilters(filterState);

  const syncFiltersToUrl = useCallback(
    (updates: { tag?: Set<string>; status?: Set<string>; priority?: Set<string>; effort?: Set<string> }) => {
      setSearchParams(
        (prev) => {
          for (const [param, values] of Object.entries(updates)) {
            prev.delete(param);
            if (values) {
              values.forEach((v: string) => prev.append(param, v));
            }
          }
          return prev;
        },
        { replace: true },
      );
    },
    [setSearchParams],
  );

  function clearFilters() {
    setSelectedStatuses(new Set(STATUSES));
    setSelectedPriorities(new Set(PRIORITIES));
    setSelectedTypes(new Set(TYPES));
    setSelectedTags(new Set());
    setSelectedEffort(new Set(EFFORTS));
    syncFiltersToUrl({ tag: new Set(), status: new Set(), priority: new Set(), effort: new Set() });
    setGlobalFilter("");
  }

  function toggleTag(tag: string) {
    setSelectedTags((prev) => {
      const next = toggleInSet(prev, tag);
      syncFiltersToUrl({ tag: next });
      return next;
    });
  }

  const filteredTasks = useMemo(
    () => applyFilters(tasks, filterState),
    [tasks, selectedStatuses, selectedPriorities, selectedTypes, selectedTags, selectedEffort],
  );

  const columns = useMemo(
    () => createTaskColumns(selectedTags, toggleTag),
    [selectedTags],
  );

  const table = useReactTable({
    data: filteredTasks,
    columns,
    state: { sorting, globalFilter },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  const rows = table.getRowModel().rows;
  const visibleCount = rows.length;

  return (
    <div>
      <FilterBar
        globalFilter={globalFilter}
        onGlobalFilterChange={setGlobalFilter}
        selectedStatuses={selectedStatuses}
        onToggleStatus={(s) =>
          setSelectedStatuses((prev) => {
            const next = toggleInSet(prev, s);
            syncFiltersToUrl({ status: next.size === STATUSES.length ? new Set() : next });
            return next;
          })
        }
        selectedPriorities={selectedPriorities}
        onTogglePriority={(p) =>
          setSelectedPriorities((prev) => {
            const next = toggleInSet(prev, p);
            syncFiltersToUrl({ priority: next.size === PRIORITIES.length ? new Set() : next });
            return next;
          })
        }
        selectedEffort={selectedEffort}
        onToggleEffort={(e) =>
          setSelectedEffort((prev) => {
            const next = toggleInSet(prev, e);
            syncFiltersToUrl({ effort: next.size === EFFORTS.length ? new Set() : next });
            return next;
          })
        }
        selectedTypes={selectedTypes}
        onToggleType={(ty) =>
          setSelectedTypes((prev) => toggleInSet(prev, ty))
        }
        selectedTags={selectedTags}
        onRemoveTag={toggleTag}
        onClearFilters={clearFilters}
        hasActiveFilters={hasActiveFilters}
      />
      <MobileCardList rows={rows} onClearFilters={clearFilters} />
      {/* Desktop table */}
      <div className="hidden sm:block overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800">
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id}>
                {hg.headers.map((header) => (
                  <th
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className={`px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider cursor-pointer select-none ${(header.column.columnDef.meta as Record<string, string>)?.className ?? ""}`}
                  >
                    <div className="flex items-center gap-1">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                      {{ asc: " ^", desc: " v" }[
                        header.column.getIsSorted() as string
                      ] ?? ""}
                    </div>
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
            {visibleCount === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-4 py-8 text-center text-sm text-gray-500 dark:text-gray-400"
                >
                  No tasks match your filters.{" "}
                  <button
                    onClick={clearFilters}
                    className="text-blue-600 hover:underline dark:text-blue-400"
                  >
                    Clear filters
                  </button>
                </td>
              </tr>
            ) : (
              rows.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className={`px-4 py-3 text-sm ${(cell.column.columnDef.meta as Record<string, string>)?.className ?? ""}`}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      <p className="mt-2 text-xs text-gray-400">
        Showing {visibleCount} of {tasks.length} tasks
      </p>
    </div>
  );
}
