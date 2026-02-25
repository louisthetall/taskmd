import * as path from "path";
import * as fs from "fs";
import { parse } from "yaml";

const CONFIG_FILENAME = ".taskmd.yaml";
const DEFAULT_TASK_DIR = "tasks";

export interface ScopeDefinition {
  paths?: string[];
  description?: string;
}

interface TaskmdConfig {
  "task-dir"?: string;
  dir?: string;
  scopes?: Record<string, ScopeDefinition | null>;
}

/**
 * Find `.taskmd.yaml` by walking up from the given directory.
 * Returns the absolute path to the config file, or null if not found.
 */
export function findConfigFile(startDir: string): string | null {
  let current = startDir;
  while (true) {
    const candidate = path.join(current, CONFIG_FILENAME);
    if (fs.existsSync(candidate)) {
      return candidate;
    }
    const parent = path.dirname(current);
    if (parent === current) break; // reached filesystem root
    current = parent;
  }
  return null;
}

/**
 * Parse a `.taskmd.yaml` config file. Returns null on any error.
 */
function readConfig(configPath: string): TaskmdConfig | null {
  try {
    const content = fs.readFileSync(configPath, "utf-8");
    return (parse(content) as TaskmdConfig | null) ?? null;
  } catch {
    return null;
  }
}

/**
 * Read the task directory from a `.taskmd.yaml` config file.
 * Supports both `task-dir` and `dir` keys (task-dir takes precedence).
 * Returns the raw value, or null if neither key is set.
 */
function readTaskDirFromConfig(configPath: string): string | null {
  const config = readConfig(configPath);
  if (!config) return null;
  return config["task-dir"] ?? config.dir ?? null;
}

/**
 * Resolve the absolute task directory for a given file path.
 *
 * Logic mirrors the Go CLI:
 * 1. Walk up from the file to find `.taskmd.yaml`
 * 2. Read `task-dir` or `dir` from the config
 * 3. Resolve relative to the config file's directory
 * 4. If no config or no dir key, use `<project-root>/tasks` as default
 *    (where project-root is the config file's directory, or the workspace root)
 */
export function resolveTaskDir(filePath: string): string | null {
  const fileDir = path.dirname(filePath);
  const configPath = findConfigFile(fileDir);

  if (configPath) {
    const projectRoot = path.dirname(configPath);
    const taskDirValue = readTaskDirFromConfig(configPath);

    if (taskDirValue) {
      // Resolve relative to project root (where .taskmd.yaml lives)
      return path.resolve(projectRoot, taskDirValue);
    }

    // Config exists but no dir key — default to <project-root>/tasks
    return path.join(projectRoot, DEFAULT_TASK_DIR);
  }

  // No config found — can't determine task directory
  return null;
}

/**
 * Check whether a file path falls under the resolved task directory.
 */
export function isUnderTaskDir(filePath: string): boolean {
  const taskDir = resolveTaskDir(filePath);
  if (!taskDir) return false;

  const normalizedFile = path.resolve(filePath);
  const normalizedTaskDir = path.resolve(taskDir);

  return normalizedFile.startsWith(normalizedTaskDir + path.sep);
}

/** A task entry with its ID and title. */
export interface TaskEntry {
  id: string;
  title: string;
}

/**
 * Scan the task directory for all task files and extract their IDs and titles.
 * Uses lightweight regex to read `id:` and `title:` from frontmatter.
 */
export function scanTaskIds(filePath: string): TaskEntry[] {
  const taskDir = resolveTaskDir(filePath);
  if (!taskDir || !fs.existsSync(taskDir)) return [];

  let files: string[];
  try {
    files = fs.readdirSync(taskDir, { recursive: true }) as string[];
  } catch {
    return [];
  }

  const entries: TaskEntry[] = [];
  for (const rel of files) {
    if (!String(rel).endsWith(".md")) continue;
    const fullPath = path.join(taskDir, String(rel));
    try {
      const content = fs.readFileSync(fullPath, "utf-8");
      const fmMatch = content.match(/^---\n([\s\S]*?)\n---/);
      if (!fmMatch) continue;
      const fm = fmMatch[1];
      const idMatch = fm.match(/^id:\s*"?([^"\n]+)"?\s*$/m);
      const titleMatch = fm.match(/^title:\s*"?([^"\n]+)"?\s*$/m);
      if (idMatch) {
        entries.push({
          id: idMatch[1].trim(),
          title: titleMatch ? titleMatch[1].trim() : "",
        });
      }
    } catch {
      // skip unreadable files
    }
  }
  return entries;
}

/** A scope entry with its name and optional description. */
export interface ScopeEntry {
  name: string;
  description?: string;
}

/**
 * Read scope definitions from `.taskmd.yaml` for a given file path.
 * Returns an array of scope entries, or an empty array if no scopes are defined.
 */
export function readScopes(filePath: string): ScopeEntry[] {
  const fileDir = path.dirname(filePath);
  const configPath = findConfigFile(fileDir);
  if (!configPath) return [];

  const config = readConfig(configPath);
  if (!config?.scopes) return [];

  return Object.entries(config.scopes)
    .filter(([, def]) => def !== null)
    .map(([name, def]) => ({
      name,
      description: def?.description,
    }));
}
