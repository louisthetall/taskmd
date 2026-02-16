import {
  STATUSES,
  PRIORITIES,
  EFFORTS,
  STATUS_COLORS,
  PRIORITY_COLORS,
  EFFORT_COLORS,
} from "../tasks/TaskTable/constants.ts";
import { toggleInSet } from "../tasks/TaskTable/utils.ts";
import { TagAutocomplete } from "./TagAutocomplete.tsx";

const INACTIVE_PILL =
  "bg-gray-50 border border-gray-200 text-gray-400 hover:bg-gray-100 hover:text-gray-500 dark:bg-gray-800/50 dark:border-gray-700 dark:text-gray-500 dark:hover:bg-gray-700 dark:hover:text-gray-400";


interface PillRowProps {
  label: string;
  items: string[];
  selected: Set<string>;
  colors: Record<string, string>;
  onToggle: (value: string) => void;
}

function PillRow({ label, items, selected, colors, onToggle }: PillRowProps) {
  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-xs text-gray-500 dark:text-gray-400 font-medium">
        {label}:
      </span>
      {items.map((item) => {
        const active = selected.has(item);
        return (
          <button
            key={item}
            onClick={() => onToggle(item)}
            className={`px-2.5 py-1 text-xs rounded-full transition-colors duration-150 ${
              active ? (colors[item] ?? "") : INACTIVE_PILL
            }`}
          >
            {item}
          </button>
        );
      })}
    </div>
  );
}

export interface BoardFilterBarProps {
  groupBy: string;
  selectedStatuses: Set<string>;
  onStatusesChange: (next: Set<string>) => void;
  selectedPriorities: Set<string>;
  onPrioritiesChange: (next: Set<string>) => void;
  selectedEfforts: Set<string>;
  onEffortsChange: (next: Set<string>) => void;
  selectedTags: Set<string>;
  onTagsChange: (next: Set<string>) => void;
  availableTags: string[];
}

export function BoardFilterBar({
  groupBy,
  selectedStatuses,
  onStatusesChange,
  selectedPriorities,
  onPrioritiesChange,
  selectedEfforts,
  onEffortsChange,
  selectedTags,
  onTagsChange,
  availableTags,
}: BoardFilterBarProps) {
  return (
    <div className="mb-4 space-y-2">
      {groupBy !== "status" && (
        <PillRow
          label="Status"
          items={STATUSES}
          selected={selectedStatuses}
          colors={STATUS_COLORS}
          onToggle={(s) => onStatusesChange(toggleInSet(selectedStatuses, s))}
        />
      )}
      {groupBy !== "priority" && (
        <PillRow
          label="Priority"
          items={PRIORITIES}
          selected={selectedPriorities}
          colors={PRIORITY_COLORS}
          onToggle={(p) => onPrioritiesChange(toggleInSet(selectedPriorities, p))}
        />
      )}
      {groupBy !== "effort" && (
        <PillRow
          label="Effort"
          items={EFFORTS}
          selected={selectedEfforts}
          colors={EFFORT_COLORS}
          onToggle={(e) => onEffortsChange(toggleInSet(selectedEfforts, e))}
        />
      )}
      {groupBy !== "tag" && availableTags.length > 0 && (
        <TagAutocomplete
          availableTags={availableTags}
          selectedTags={selectedTags}
          onTagsChange={onTagsChange}
        />
      )}
    </div>
  );
}
