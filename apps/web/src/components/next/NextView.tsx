import { useMemo } from "react";
import type { Recommendation } from "../../api/types.ts";
import { RecommendationCard } from "./RecommendationCard.tsx";
import { FolderAutocomplete } from "./FolderAutocomplete.tsx";
import { useTasks } from "../../hooks/use-tasks.ts";

interface NextViewProps {
  recommendations: Recommendation[];
  limit: number;
  onLimitChange: (limit: number) => void;
  group: string;
  onGroupChange: (group: string) => void;
}

const LIMIT_OPTIONS = [3, 5, 10];

export function NextView({
  recommendations,
  limit,
  onLimitChange,
  group,
  onGroupChange,
}: NextViewProps) {
  const { data: tasks } = useTasks();

  const groups = useMemo(() => {
    if (!tasks) return [];
    const set = new Set<string>();
    for (const t of tasks) {
      if (t.group) set.add(t.group);
    }
    return Array.from(set).sort();
  }, [tasks]);

  return (
    <div>
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-4">
        <h2 className="text-lg font-semibold">Recommended Tasks</h2>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1">
            <label
              htmlFor="group-filter"
              className="text-sm text-gray-500 dark:text-gray-400"
            >
              Folder:
            </label>
            <FolderAutocomplete
              folders={groups}
              value={group}
              onChange={onGroupChange}
            />
          </div>
          <div className="flex items-center gap-1">
            <span className="text-sm text-gray-500 dark:text-gray-400 mr-2">
              Show:
            </span>
            {LIMIT_OPTIONS.map((n) => (
              <button
                key={n}
                onClick={() => onLimitChange(n)}
                className={`min-h-[44px] sm:min-h-0 inline-flex items-center px-3 py-1 text-sm rounded-md transition-colors ${
                  limit === n
                    ? "bg-gray-900 text-white dark:bg-white dark:text-gray-900"
                    : "text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-100 dark:hover:bg-gray-700"
                }`}
              >
                {n}
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className="space-y-3">
        {recommendations.map((rec) => (
          <RecommendationCard key={rec.id} rec={rec} />
        ))}
      </div>
    </div>
  );
}
