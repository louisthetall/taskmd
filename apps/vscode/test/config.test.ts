import { describe, it, expect, beforeEach, afterEach } from "vitest";
import * as fs from "fs";
import * as path from "path";
import * as os from "os";
import { findConfigFile, resolveTaskDir, isUnderTaskDir, readScopes, scanTaskIds } from "../src/config";

function makeTmpDir(): string {
  return fs.mkdtempSync(path.join(os.tmpdir(), "taskmd-test-"));
}

function cleanup(dir: string): void {
  fs.rmSync(dir, { recursive: true, force: true });
}

describe("findConfigFile", () => {
  let tmpDir: string;
  beforeEach(() => { tmpDir = makeTmpDir(); });
  afterEach(() => cleanup(tmpDir));

  it("finds config in the same directory", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "dir: tasks\n");
    const result = findConfigFile(tmpDir);
    expect(result).toBe(path.join(tmpDir, ".taskmd.yaml"));
  });

  it("finds config in a parent directory", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "dir: tasks\n");
    const subDir = path.join(tmpDir, "tasks", "cli");
    fs.mkdirSync(subDir, { recursive: true });
    const result = findConfigFile(subDir);
    expect(result).toBe(path.join(tmpDir, ".taskmd.yaml"));
  });

  it("returns null when no config exists", () => {
    const result = findConfigFile(tmpDir);
    expect(result).toBeNull();
  });
});

describe("resolveTaskDir", () => {
  let tmpDir: string;
  beforeEach(() => { tmpDir = makeTmpDir(); });
  afterEach(() => cleanup(tmpDir));

  it("uses task-dir from config", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: my-tasks\n");
    const filePath = path.join(tmpDir, "my-tasks", "001.md");
    const result = resolveTaskDir(filePath);
    expect(result).toBe(path.join(tmpDir, "my-tasks"));
  });

  it("uses dir from config (legacy key)", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "dir: work\n");
    const filePath = path.join(tmpDir, "work", "001.md");
    const result = resolveTaskDir(filePath);
    expect(result).toBe(path.join(tmpDir, "work"));
  });

  it("task-dir takes precedence over dir", () => {
    fs.writeFileSync(
      path.join(tmpDir, ".taskmd.yaml"),
      "task-dir: primary\ndir: secondary\n"
    );
    const filePath = path.join(tmpDir, "primary", "001.md");
    const result = resolveTaskDir(filePath);
    expect(result).toBe(path.join(tmpDir, "primary"));
  });

  it("defaults to tasks/ when config exists but has no dir key", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "workflow: solo\n");
    const filePath = path.join(tmpDir, "tasks", "001.md");
    const result = resolveTaskDir(filePath);
    expect(result).toBe(path.join(tmpDir, "tasks"));
  });

  it("returns null when no config exists", () => {
    const result = resolveTaskDir(path.join(tmpDir, "tasks", "001.md"));
    expect(result).toBeNull();
  });

  it("resolves relative paths against config directory", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: ./src/tasks\n");
    const filePath = path.join(tmpDir, "src", "tasks", "001.md");
    const result = resolveTaskDir(filePath);
    expect(result).toBe(path.resolve(tmpDir, "src/tasks"));
  });
});

describe("isUnderTaskDir", () => {
  let tmpDir: string;
  beforeEach(() => { tmpDir = makeTmpDir(); });
  afterEach(() => cleanup(tmpDir));

  it("returns true for file under task dir", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const filePath = path.join(tmpDir, "tasks", "001.md");
    expect(isUnderTaskDir(filePath)).toBe(true);
  });

  it("returns true for file in nested subdirectory", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const filePath = path.join(tmpDir, "tasks", "cli", "001.md");
    expect(isUnderTaskDir(filePath)).toBe(true);
  });

  it("returns false for file outside task dir", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const filePath = path.join(tmpDir, "docs", "readme.md");
    expect(isUnderTaskDir(filePath)).toBe(false);
  });

  it("returns false when no config exists", () => {
    const filePath = path.join(tmpDir, "tasks", "001.md");
    expect(isUnderTaskDir(filePath)).toBe(false);
  });

  it("works with default tasks/ when config has no dir", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "workflow: solo\n");
    expect(isUnderTaskDir(path.join(tmpDir, "tasks", "001.md"))).toBe(true);
    expect(isUnderTaskDir(path.join(tmpDir, "docs", "001.md"))).toBe(false);
  });
});

describe("readScopes", () => {
  let tmpDir: string;
  beforeEach(() => { tmpDir = makeTmpDir(); });
  afterEach(() => cleanup(tmpDir));

  it("returns scope names from config", () => {
    fs.writeFileSync(
      path.join(tmpDir, ".taskmd.yaml"),
      `scopes:\n  cli:\n    paths:\n      - "apps/cli/"\n  web:\n    paths:\n      - "apps/web/"\n`
    );
    const filePath = path.join(tmpDir, "tasks", "001.md");
    const scopes = readScopes(filePath);
    expect(scopes).toEqual([
      { name: "cli", description: undefined },
      { name: "web", description: undefined },
    ]);
  });

  it("includes scope descriptions when present", () => {
    fs.writeFileSync(
      path.join(tmpDir, ".taskmd.yaml"),
      `scopes:\n  cli:\n    paths:\n      - "apps/cli/"\n    description: "CLI application"\n`
    );
    const filePath = path.join(tmpDir, "tasks", "001.md");
    const scopes = readScopes(filePath);
    expect(scopes).toEqual([
      { name: "cli", description: "CLI application" },
    ]);
  });

  it("returns empty array when no scopes defined", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "dir: tasks\n");
    const filePath = path.join(tmpDir, "tasks", "001.md");
    expect(readScopes(filePath)).toEqual([]);
  });

  it("returns empty array when no config exists", () => {
    const filePath = path.join(tmpDir, "tasks", "001.md");
    expect(readScopes(filePath)).toEqual([]);
  });
});

describe("scanTaskIds", () => {
  let tmpDir: string;
  beforeEach(() => { tmpDir = makeTmpDir(); });
  afterEach(() => cleanup(tmpDir));

  it("scans task files and extracts id and title", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const tasksDir = path.join(tmpDir, "tasks");
    fs.mkdirSync(tasksDir, { recursive: true });
    fs.writeFileSync(
      path.join(tasksDir, "001-setup.md"),
      '---\nid: "001"\ntitle: "Setup project"\nstatus: pending\n---\n# Setup\n'
    );
    fs.writeFileSync(
      path.join(tasksDir, "002-auth.md"),
      '---\nid: "002"\ntitle: "Add auth"\n---\n# Auth\n'
    );
    const filePath = path.join(tasksDir, "001-setup.md");
    const entries = scanTaskIds(filePath);
    expect(entries).toHaveLength(2);
    expect(entries).toContainEqual({ id: "001", title: "Setup project" });
    expect(entries).toContainEqual({ id: "002", title: "Add auth" });
  });

  it("scans nested subdirectories", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const subDir = path.join(tmpDir, "tasks", "cli");
    fs.mkdirSync(subDir, { recursive: true });
    fs.writeFileSync(
      path.join(subDir, "010-feature.md"),
      '---\nid: "010"\ntitle: "CLI feature"\n---\n'
    );
    const entries = scanTaskIds(path.join(tmpDir, "tasks", "test.md"));
    expect(entries).toHaveLength(1);
    expect(entries[0]).toEqual({ id: "010", title: "CLI feature" });
  });

  it("skips files without frontmatter", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const tasksDir = path.join(tmpDir, "tasks");
    fs.mkdirSync(tasksDir, { recursive: true });
    fs.writeFileSync(path.join(tasksDir, "readme.md"), "# Just markdown\n");
    const entries = scanTaskIds(path.join(tasksDir, "readme.md"));
    expect(entries).toEqual([]);
  });

  it("skips files without id in frontmatter", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const tasksDir = path.join(tmpDir, "tasks");
    fs.mkdirSync(tasksDir, { recursive: true });
    fs.writeFileSync(
      path.join(tasksDir, "note.md"),
      '---\ntitle: "No ID"\nstatus: pending\n---\n'
    );
    const entries = scanTaskIds(path.join(tasksDir, "note.md"));
    expect(entries).toEqual([]);
  });

  it("handles unquoted id values", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: tasks\n");
    const tasksDir = path.join(tmpDir, "tasks");
    fs.mkdirSync(tasksDir, { recursive: true });
    fs.writeFileSync(
      path.join(tasksDir, "003.md"),
      '---\nid: 003\ntitle: Bare title\n---\n'
    );
    const entries = scanTaskIds(path.join(tasksDir, "003.md"));
    expect(entries).toHaveLength(1);
    expect(entries[0]).toEqual({ id: "003", title: "Bare title" });
  });

  it("returns empty when no config exists", () => {
    const entries = scanTaskIds(path.join(tmpDir, "tasks", "001.md"));
    expect(entries).toEqual([]);
  });

  it("returns empty when task dir does not exist", () => {
    fs.writeFileSync(path.join(tmpDir, ".taskmd.yaml"), "task-dir: nonexistent\n");
    const entries = scanTaskIds(path.join(tmpDir, "nonexistent", "001.md"));
    expect(entries).toEqual([]);
  });
});
