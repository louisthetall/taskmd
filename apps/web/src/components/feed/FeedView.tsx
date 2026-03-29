import { Link } from "react-router-dom";
import type {
  FeedEntry,
  FeedFileChange,
  FeedFieldChange,
  FeedSubtaskChange,
} from "../../api/types.ts";

interface FeedViewProps {
  entries: FeedEntry[];
}

export function FeedView({ entries }: FeedViewProps) {
  return (
    <div className="space-y-3">
      {entries.map((entry, i) => (
        <FeedEntryCard key={`${entry.timestamp}-${i}`} entry={entry} />
      ))}
    </div>
  );
}

function FeedEntryCard({ entry }: { entry: FeedEntry }) {
  const isWorklog = entry.source === "worklog";

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4 dark:bg-gray-800 dark:border-gray-700">
      <div className="flex items-start gap-3">
        <SourceBadge source={entry.source} />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <Timestamp value={entry.timestamp} />
            {entry.author && (
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                {entry.author}
              </span>
            )}
            {entry.taskID && (
              <Link
                to={`/tasks/${entry.taskID}`}
                className="text-xs font-mono px-1.5 py-0.5 rounded bg-gray-100 text-blue-600 hover:text-blue-800 hover:bg-gray-200 dark:bg-gray-700 dark:text-blue-400 dark:hover:text-blue-300 dark:hover:bg-gray-600"
              >
                {entry.taskID}
              </Link>
            )}
          </div>
          <p className="text-sm text-gray-900 dark:text-gray-100 mt-1">
            {entry.message}
          </p>
          {!isWorklog && entry.files && entry.files.length > 0 && (
            <div className="mt-2 space-y-1.5">
              {entry.files.map((file, j) => (
                <FileChangeItem key={`${file.path}-${j}`} file={file} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function SourceBadge({ source }: { source: string }) {
  if (source === "worklog") {
    return (
      <span className="mt-0.5 flex-shrink-0 inline-flex items-center justify-center w-7 h-7 rounded-full bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400">
        <svg
          className="w-3.5 h-3.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
          />
        </svg>
      </span>
    );
  }

  return (
    <span className="mt-0.5 flex-shrink-0 inline-flex items-center justify-center w-7 h-7 rounded-full bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
      <svg
        className="w-3.5 h-3.5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
        />
      </svg>
    </span>
  );
}

function Timestamp({ value }: { value: string }) {
  const date = new Date(value);
  const formatted = date.toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
  return (
    <time
      dateTime={value}
      className="text-xs text-gray-500 dark:text-gray-400"
    >
      {formatted}
    </time>
  );
}

function FileChangeItem({ file }: { file: FeedFileChange }) {
  return (
    <div className="text-xs">
      <div className="flex items-center gap-1.5">
        <FileStatusBadge status={file.status} taskStatus={file.taskStatus} />
        <span className="font-mono text-gray-600 dark:text-gray-400 truncate">
          {file.path}
        </span>
        {file.taskID && (
          <Link
            to={`/tasks/${file.taskID}`}
            className="font-mono text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 flex-shrink-0"
          >
            {file.taskID}
          </Link>
        )}
      </div>
      {file.fieldChanges && file.fieldChanges.length > 0 && (
        <div className="mt-1 flex flex-wrap gap-1 ml-6">
          {file.fieldChanges.map((fc, i) => (
            <FieldChangeBadge key={`${fc.field}-${i}`} change={fc} />
          ))}
        </div>
      )}
      {file.subtaskChanges && file.subtaskChanges.length > 0 && (
        <div className="mt-1 flex flex-wrap gap-1 ml-6">
          {file.subtaskChanges.map((sc, i) => (
            <SubtaskChangeBadge key={`${sc.text}-${i}`} change={sc} />
          ))}
        </div>
      )}
    </div>
  );
}

function FileStatusBadge({
  status,
  taskStatus,
}: {
  status: string;
  taskStatus?: string;
}) {
  if (taskStatus === "completed") {
    return (
      <span className="inline-flex px-1.5 py-0.5 rounded text-[10px] font-medium bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400">
        Completed
      </span>
    );
  }
  if (taskStatus === "cancelled") {
    return (
      <span className="inline-flex px-1.5 py-0.5 rounded text-[10px] font-medium bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400">
        Cancelled
      </span>
    );
  }

  const styles: Record<string, string> = {
    created:
      "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400",
    modified:
      "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400",
    renamed:
      "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400",
  };

  const labels: Record<string, string> = {
    created: "Added",
    modified: "Modified",
    renamed: "Renamed",
  };

  return (
    <span
      className={`inline-flex px-1.5 py-0.5 rounded text-[10px] font-medium ${styles[status] ?? "bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-400"}`}
    >
      {labels[status] ?? status}
    </span>
  );
}

function FieldChangeBadge({ change }: { change: FeedFieldChange }) {
  return (
    <span className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-[10px]">
      <span className="font-medium text-gray-600 dark:text-gray-300">
        {change.field}
      </span>
      <span className="text-gray-400 dark:text-gray-500">
        {change.oldValue || "(none)"}
      </span>
      <span className="text-gray-400 dark:text-gray-500">&rarr;</span>
      <span className="text-gray-700 dark:text-gray-200">
        {change.newValue || "(none)"}
      </span>
    </span>
  );
}

function SubtaskChangeBadge({ change }: { change: FeedSubtaskChange }) {
  return (
    <span
      className={`inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] ${
        change.done
          ? "bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400"
          : "bg-gray-50 text-gray-600 dark:bg-gray-700 dark:text-gray-400"
      }`}
    >
      <span>{change.done ? "\u2611" : "\u2610"}</span>
      <span className="truncate max-w-[200px]">{change.text}</span>
    </span>
  );
}
