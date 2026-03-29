import { useState } from "react";
import { useFeed } from "../hooks/use-feed.ts";
import { useProject } from "../hooks/use-project.ts";
import { FeedView } from "../components/feed/FeedView.tsx";
import { LoadingState } from "../components/shared/LoadingState.tsx";
import { ErrorState } from "../components/shared/ErrorState.tsx";

const sourceOptions = ["all", "git", "worklog"] as const;
const sinceOptions = [
  { label: "All time", value: "" },
  { label: "24 hours", value: "1d" },
  { label: "7 days", value: "7d" },
  { label: "30 days", value: "30d" },
] as const;

export function FeedPage() {
  const { project } = useProject();
  const [source, setSource] = useState("all");
  const [since, setSince] = useState("7d");
  const [scope, setScope] = useState("");

  const { data, error, isLoading, mutate } = useFeed({
    source,
    since: since || undefined,
    scope: scope || undefined,
    limit: 50,
    project,
  });

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1.5">
          <label className="text-xs text-gray-500 dark:text-gray-400">
            Source
          </label>
          <select
            value={source}
            onChange={(e) => setSource(e.target.value)}
            className="text-sm rounded-md border border-gray-300 bg-white px-2 py-1 dark:bg-gray-800 dark:border-gray-600 dark:text-gray-200"
          >
            {sourceOptions.map((s) => (
              <option key={s} value={s}>
                {s === "all" ? "All" : s === "git" ? "Git" : "Worklog"}
              </option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-1.5">
          <label className="text-xs text-gray-500 dark:text-gray-400">
            Since
          </label>
          <select
            value={since}
            onChange={(e) => setSince(e.target.value)}
            className="text-sm rounded-md border border-gray-300 bg-white px-2 py-1 dark:bg-gray-800 dark:border-gray-600 dark:text-gray-200"
          >
            {sinceOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-1.5">
          <label className="text-xs text-gray-500 dark:text-gray-400">
            Scope
          </label>
          <input
            type="text"
            value={scope}
            onChange={(e) => setScope(e.target.value)}
            placeholder="e.g. cli"
            className="text-sm rounded-md border border-gray-300 bg-white px-2 py-1 w-28 dark:bg-gray-800 dark:border-gray-600 dark:text-gray-200 placeholder:text-gray-400 dark:placeholder:text-gray-500"
          />
        </div>
      </div>

      {isLoading && <LoadingState variant="table" />}
      {error && <ErrorState error={error} onRetry={() => mutate()} />}
      {!isLoading && !error && data && data.length === 0 && (
        <p className="text-sm text-gray-500 py-8 text-center">
          No recent activity
        </p>
      )}
      {!isLoading && !error && data && data.length > 0 && (
        <FeedView entries={data} />
      )}
    </div>
  );
}
