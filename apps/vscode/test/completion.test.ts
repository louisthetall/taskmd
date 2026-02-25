import { describe, it, expect } from "vitest";
import { findFrontmatterBounds } from "../src/frontmatter";
import { resolveCompletions } from "../src/completion";
import type { ScopeEntry, TaskEntry } from "../src/config";

const TEST_SCOPES: ScopeEntry[] = [
  { name: "cli", description: "CLI application" },
  { name: "web", description: "Web frontend" },
  { name: "cli/graph" },
];

const TEST_TASKS: TaskEntry[] = [
  { id: "001", title: "Setup project" },
  { id: "002", title: "Add authentication" },
  { id: "042", title: "Implement graphs" },
];

function getResult(
  text: string,
  line: number,
  character: number,
  scopes?: ScopeEntry[],
  taskEntries?: TaskEntry[]
) {
  const bounds = findFrontmatterBounds(text);
  if (!bounds) return undefined;
  const lines = text.split("\n");
  return resolveCompletions(lines, line, character, bounds.startLine, bounds.endLine, scopes, taskEntries);
}

describe("completion logic", () => {
  const doc = `---
id: "042"
title: "Test"
status:
priority:
effort:
type:
tags: []
---
# Body`;

  it("provides status completions", () => {
    const result = getResult(doc, 3, 8);
    expect(result).toBeDefined();
    expect(result!.values).toEqual([
      "pending",
      "in-progress",
      "completed",
      "in-review",
      "blocked",
      "cancelled",
    ]);
  });

  it("provides priority completions", () => {
    const result = getResult(doc, 4, 10);
    expect(result).toBeDefined();
    expect(result!.values).toEqual(["low", "medium", "high", "critical"]);
  });

  it("provides effort completions", () => {
    const result = getResult(doc, 5, 8);
    expect(result).toBeDefined();
    expect(result!.values).toEqual(["small", "medium", "large"]);
  });

  it("provides type completions", () => {
    const result = getResult(doc, 6, 6);
    expect(result).toBeDefined();
    expect(result!.values).toEqual(["feature", "bug", "improvement", "chore", "docs"]);
  });

  it("no completions for non-enum field", () => {
    expect(getResult(doc, 7, 6)).toBeUndefined();
  });

  it("no completions outside frontmatter", () => {
    expect(getResult(doc, 9, 5)).toBeUndefined();
  });

  it("no completions on opening delimiter line", () => {
    expect(getResult(doc, 0, 3)).toBeUndefined();
  });

  it("no completions on closing delimiter line", () => {
    expect(getResult(doc, 8, 3)).toBeUndefined();
  });

  it("no completions when cursor is in middle of value", () => {
    const docWithValue = `---
status: pend
---`;
    expect(getResult(docWithValue, 1, 12)).toBeUndefined();
  });

  it("no completions for document without frontmatter", () => {
    expect(getResult("# Just markdown", 0, 5)).toBeUndefined();
  });
});

describe("completion insert text", () => {
  it("includes leading space in insertText when no space after colon", () => {
    const doc = `---
status:
---`;
    // cursor right after colon: "status:" -> character 7
    const result = getResult(doc, 1, 7);
    expect(result).toBeDefined();
    expect(result!.insertTexts[0]).toBe(" pending");
    expect(result!.insertTexts).toEqual(
      result!.values.map((v) => ` ${v}`)
    );
  });

  it("includes leading space in insertText when space already present", () => {
    const doc = `---
status:
---`;
    // cursor after "status: " -> character 8
    const result = getResult(doc, 1, 8);
    expect(result).toBeDefined();
    expect(result!.insertTexts[0]).toBe(" pending");
  });

  it("replace range covers from after colon to cursor (no space)", () => {
    const doc = `---
status:
---`;
    const result = getResult(doc, 1, 7);
    expect(result).toBeDefined();
    // colon is at index 6, so replace starts at 7
    expect(result!.replaceColumns).toEqual([7, 7]);
  });

  it("replace range covers from after colon to cursor (with space)", () => {
    const doc = `---
status:
---`;
    const result = getResult(doc, 1, 8);
    expect(result).toBeDefined();
    // colon is at index 6, replace starts at 7, cursor at 8
    // replaces the existing space, then inserts " value"
    expect(result!.replaceColumns).toEqual([7, 8]);
  });

  it("applying insert text produces correct YAML (no space case)", () => {
    const line = "status:";
    const result = getResult(`---\n${line}\n---`, 1, 7);
    expect(result).toBeDefined();
    const [start, end] = result!.replaceColumns;
    const applied = line.substring(0, start) + result!.insertTexts[0] + line.substring(end);
    expect(applied).toBe("status: pending");
  });

  it("applying insert text produces correct YAML (space case)", () => {
    const line = "status: ";
    const result = getResult(`---\n${line}\n---`, 1, 8);
    expect(result).toBeDefined();
    const [start, end] = result!.replaceColumns;
    const applied = line.substring(0, start) + result!.insertTexts[0] + line.substring(end);
    expect(applied).toBe("status: pending");
  });
});

describe("touches scope completions", () => {
  it("provides scope completions for block array items under touches", () => {
    const doc = `---
id: "042"
touches:
  -
---`;
    const result = getResult(doc, 3, 4, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("touches");
    expect(result!.values).toEqual(["cli", "web", "cli/graph"]);
  });

  it("provides scope completions for dash without trailing space", () => {
    const doc = "---\ntouches:\n  -\n---";
    const result = getResult(doc, 2, 3, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("touches");
  });

  it("includes scope descriptions as details", () => {
    const doc = `---
touches:
  -
---`;
    const result = getResult(doc, 2, 4, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.details).toEqual(["CLI application", "Web frontend", undefined]);
  });

  it("replaces partial text in block array item", () => {
    const doc = `---
touches:
  - cl
---`;
    const result = getResult(doc, 2, 6, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.replaceColumns).toEqual([4, 6]);
    expect(result!.insertTexts[0]).toBe("cli");
  });

  it("no scope completions when scopes are empty", () => {
    const doc = `---
touches:
  -
---`;
    const result = getResult(doc, 2, 4, []);
    expect(result).toBeUndefined();
  });

  it("no scope completions when scopes are not provided", () => {
    const doc = `---
touches:
  -
---`;
    const result = getResult(doc, 2, 4);
    expect(result).toBeUndefined();
  });

  it("no scope completions for array items under other fields", () => {
    const doc = `---
tags:
  -
---`;
    const result = getResult(doc, 2, 4, TEST_SCOPES);
    expect(result).toBeUndefined();
  });

  it("provides scope completions for inline array", () => {
    const doc = `---
touches: [
---`;
    const result = getResult(doc, 1, 11, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("touches");
    expect(result!.values).toEqual(["cli", "web", "cli/graph"]);
  });

  it("provides scope completions after comma in inline array", () => {
    const doc = `---
touches: [cli,
---`;
    const result = getResult(doc, 1, 16, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.values).toEqual(["cli", "web", "cli/graph"]);
  });

  it("enum completions still work when scopes are provided", () => {
    const doc = `---
status:
---`;
    const result = getResult(doc, 1, 8, TEST_SCOPES);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("status");
    expect(result!.values).toContain("pending");
  });
});

describe("task ID completions", () => {
  it("provides task ID completions for dependencies block array", () => {
    const doc = `---
dependencies:
  -
---`;
    const result = getResult(doc, 2, 4, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("dependencies");
    expect(result!.values).toEqual(["001", "002", "042"]);
  });

  it("includes task titles as details", () => {
    const doc = `---
dependencies:
  -
---`;
    const result = getResult(doc, 2, 4, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.details).toEqual(["Setup project", "Add authentication", "Implement graphs"]);
  });

  it("wraps task IDs in quotes for dependencies block items", () => {
    const doc = `---
dependencies:
  -
---`;
    const result = getResult(doc, 2, 4, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.insertTexts[0]).toBe('"001"');
    expect(result!.insertTexts[1]).toBe('"002"');
  });

  it("handles dependencies block item with existing quote", () => {
    const doc = `---
dependencies:
  - "
---`;
    const result = getResult(doc, 2, 5, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("dependencies");
    // replaceColumns should cover the opening quote at position 4
    expect(result!.replaceColumns[0]).toBe(4);
    expect(result!.insertTexts[0]).toBe('"001"');
  });

  it("provides task ID completions for dependencies inline array", () => {
    const doc = `---
dependencies: ["
---`;
    const result = getResult(doc, 1, 17, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("dependencies");
    expect(result!.values).toEqual(["001", "002", "042"]);
  });

  it("provides completions after comma in dependencies inline array", () => {
    const doc = `---
dependencies: ["001", "
---`;
    const result = getResult(doc, 1, 24, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("dependencies");
  });

  it("provides task ID completions for parent field", () => {
    const doc = `---
parent:
---`;
    const result = getResult(doc, 1, 8, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("parent");
    expect(result!.values).toEqual(["001", "002", "042"]);
  });

  it("wraps task IDs in quotes for parent field", () => {
    const doc = `---
parent:
---`;
    const result = getResult(doc, 1, 8, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.insertTexts[0]).toBe('"001"');
  });

  it("handles parent field with existing quote", () => {
    const doc = `---
parent: "
---`;
    const result = getResult(doc, 1, 9, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.replaceColumns[0]).toBe(8);
    expect(result!.insertTexts[0]).toBe('"001"');
  });

  it("no task ID completions for non-ID fields", () => {
    const doc = `---
tags:
  -
---`;
    const result = getResult(doc, 2, 4, [], TEST_TASKS);
    expect(result).toBeUndefined();
  });

  it("no task ID completions when task entries are empty", () => {
    const doc = `---
dependencies:
  -
---`;
    const result = getResult(doc, 2, 4, [], []);
    expect(result).toBeUndefined();
  });

  it("no task ID completions when task entries not provided", () => {
    const doc = `---
dependencies:
  -
---`;
    const result = getResult(doc, 2, 4);
    expect(result).toBeUndefined();
  });

  it("enum completions still work with task entries provided", () => {
    const doc = `---
status:
---`;
    const result = getResult(doc, 1, 8, [], TEST_TASKS);
    expect(result).toBeDefined();
    expect(result!.fieldName).toBe("status");
    expect(result!.values).toContain("pending");
  });
});
