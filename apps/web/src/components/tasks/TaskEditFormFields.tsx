import {
  STATUSES,
  PRIORITIES,
  EFFORTS,
  TYPES,
} from "./TaskTable/constants.ts";

interface FieldGridProps {
  status: string;
  onStatusChange: (v: string) => void;
  priority: string;
  onPriorityChange: (v: string) => void;
  effort: string;
  onEffortChange: (v: string) => void;
  taskType: string;
  onTaskTypeChange: (v: string) => void;
  inputClasses: string;
}

export function FieldGrid({
  status, onStatusChange,
  priority, onPriorityChange,
  effort, onEffortChange,
  taskType, onTaskTypeChange,
  inputClasses,
}: FieldGridProps) {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Status
        </label>
        <select value={status} onChange={(e) => onStatusChange(e.target.value)} className={inputClasses}>
          {STATUSES.map((s) => (
            <option key={s} value={s}>{s}</option>
          ))}
        </select>
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Priority
        </label>
        <select value={priority} onChange={(e) => onPriorityChange(e.target.value)} className={inputClasses}>
          <option value="">-</option>
          {PRIORITIES.map((p) => (
            <option key={p} value={p}>{p}</option>
          ))}
        </select>
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Effort
        </label>
        <select value={effort} onChange={(e) => onEffortChange(e.target.value)} className={inputClasses}>
          <option value="">-</option>
          {EFFORTS.map((e) => (
            <option key={e} value={e}>{e}</option>
          ))}
        </select>
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Type
        </label>
        <select value={taskType} onChange={(e) => onTaskTypeChange(e.target.value)} className={inputClasses}>
          <option value="">-</option>
          {TYPES.map((ty) => (
            <option key={ty} value={ty}>{ty}</option>
          ))}
        </select>
      </div>
    </div>
  );
}

interface MetadataFieldsProps {
  milestone: string;
  onMilestoneChange: (v: string) => void;
  owner: string;
  onOwnerChange: (v: string) => void;
  parent: string;
  onParentChange: (v: string) => void;
  tags: string;
  onTagsChange: (v: string) => void;
  inputClasses: string;
}

export function MetadataFields({
  milestone, onMilestoneChange,
  owner, onOwnerChange,
  parent, onParentChange,
  tags, onTagsChange,
  inputClasses,
}: MetadataFieldsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Milestone
        </label>
        <input
          type="text"
          value={milestone}
          onChange={(e) => onMilestoneChange(e.target.value)}
          placeholder="e.g. v1.0"
          className={inputClasses}
        />
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Owner
        </label>
        <input
          type="text"
          value={owner}
          onChange={(e) => onOwnerChange(e.target.value)}
          placeholder="e.g. alice"
          className={inputClasses}
        />
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Parent
        </label>
        <input
          type="text"
          value={parent}
          onChange={(e) => onParentChange(e.target.value)}
          placeholder="e.g. 045"
          className={inputClasses}
        />
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          Tags (comma-separated)
        </label>
        <input
          type="text"
          value={tags}
          onChange={(e) => onTagsChange(e.target.value)}
          placeholder="e.g. backend, api, feature"
          className={inputClasses}
        />
      </div>
    </div>
  );
}
