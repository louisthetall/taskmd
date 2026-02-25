import * as vscode from "vscode";
import { ENUM_FIELDS } from "./schema";
import { findFrontmatterBounds } from "./frontmatter";
import { readScopes, scanTaskIds, type ScopeEntry, type TaskEntry } from "./config";

/** Result of resolving completions for a line in frontmatter. */
export interface CompletionResult {
  fieldName: string;
  values: readonly string[];
  /** Optional detail text per value. */
  details?: (string | undefined)[];
  /** Insert text for each value (includes leading space). */
  insertTexts: string[];
  /** Column range to replace: [startCol, endCol]. */
  replaceColumns: [number, number];
}

/**
 * Resolve completions for a given document text and cursor position.
 * Pure logic — no vscode dependency.
 */
export function resolveCompletions(
  lines: string[],
  cursorLine: number,
  cursorCol: number,
  frontmatterStartLine: number,
  frontmatterEndLine: number,
  scopes?: readonly ScopeEntry[],
  taskEntries?: readonly TaskEntry[]
): CompletionResult | undefined {
  if (cursorLine <= frontmatterStartLine || cursorLine >= frontmatterEndLine) {
    return undefined;
  }

  const lineText = lines[cursorLine];
  if (!lineText) return undefined;

  const beforeCursor = lineText.substring(0, cursorCol);

  // Check for enum field completions: "status: "
  const enumMatch = beforeCursor.match(/^(\w+):\s*$/);
  if (enumMatch) {
    const fieldName = enumMatch[1];
    const allowed = ENUM_FIELDS[fieldName];
    if (allowed) {
      const colonIndex = lineText.indexOf(":");
      return {
        fieldName,
        values: allowed,
        insertTexts: allowed.map((val) => ` ${val}`),
        replaceColumns: [colonIndex + 1, cursorCol],
      };
    }
  }

  // Check for touches scope completions
  if (scopes && scopes.length > 0) {
    const touchesResult = resolveTouchesCompletions(
      lines, cursorLine, cursorCol, frontmatterStartLine, scopes
    );
    if (touchesResult) return touchesResult;
  }

  // Check for task ID completions (dependencies, parent)
  if (taskEntries && taskEntries.length > 0) {
    const taskIdResult = resolveTaskIdCompletions(
      lines, cursorLine, cursorCol, frontmatterStartLine, taskEntries
    );
    if (taskIdResult) return taskIdResult;
  }

  return undefined;
}

/**
 * Resolve completions for touches field values (scope names).
 * Handles block array items (`  - value`) and inline arrays (`touches: [value, `).
 */
export function resolveTouchesCompletions(
  lines: string[],
  cursorLine: number,
  cursorCol: number,
  frontmatterStartLine: number,
  scopes: readonly ScopeEntry[]
): CompletionResult | undefined {
  const lineText = lines[cursorLine];
  if (!lineText) return undefined;

  const beforeCursor = lineText.substring(0, cursorCol);
  const scopeNames = scopes.map((s) => s.name);
  const scopeDetails = scopes.map((s) => s.description);

  // Block array item: "  - value" or "  - " under touches field
  const blockMatch = beforeCursor.match(/^(\s+-\s*)(\S*)$/);
  if (blockMatch) {
    const parentField = findParentField(lines, cursorLine, frontmatterStartLine);
    if (parentField === "touches") {
      const prefixEnd = blockMatch[1].length;
      return {
        fieldName: "touches",
        values: scopeNames,
        details: scopeDetails,
        insertTexts: scopeNames.map((name) => name),
        replaceColumns: [prefixEnd, cursorCol],
      };
    }
  }

  // Inline array: "touches: [val1, " or "touches: ["
  const inlineMatch = beforeCursor.match(/^touches:\s*\[(?:.*,\s*)?(\S*)$/);
  if (inlineMatch) {
    const partial = inlineMatch[1] ?? "";
    const replaceStart = cursorCol - partial.length;
    return {
      fieldName: "touches",
      values: scopeNames,
      details: scopeDetails,
      insertTexts: scopeNames.map((name) => name),
      replaceColumns: [replaceStart, cursorCol],
    };
  }

  return undefined;
}

const TASK_ID_FIELDS = ["dependencies", "parent"];

/**
 * Resolve completions for task ID fields (dependencies, parent).
 * Handles block array items, inline arrays (dependencies), and scalar values (parent).
 */
export function resolveTaskIdCompletions(
  lines: string[],
  cursorLine: number,
  cursorCol: number,
  frontmatterStartLine: number,
  taskEntries: readonly TaskEntry[]
): CompletionResult | undefined {
  const lineText = lines[cursorLine];
  if (!lineText) return undefined;

  const beforeCursor = lineText.substring(0, cursorCol);
  const ids = taskEntries.map((e) => e.id);
  const titles = taskEntries.map((e) => e.title || undefined);

  // Block array item: "  - value" under dependencies field
  const blockMatch = beforeCursor.match(/^(\s+-\s*"?)(\S*)$/);
  if (blockMatch) {
    const parentField = findParentField(lines, cursorLine, frontmatterStartLine);
    if (parentField === "dependencies") {
      const prefix = blockMatch[1];
      const hasQuote = prefix.endsWith('"');
      const replaceStart = hasQuote ? prefix.length - 1 : prefix.length;
      return {
        fieldName: "dependencies",
        values: ids,
        details: titles,
        insertTexts: ids.map((id) => `"${id}"`),
        replaceColumns: [replaceStart, cursorCol],
      };
    }
  }

  // Inline array: `dependencies: ["` or `dependencies: ["xxx", "`
  const inlineMatch = beforeCursor.match(/^dependencies:\s*\[(?:.*,\s*)?"?(\S*)$/);
  if (inlineMatch) {
    const partial = inlineMatch[1] ?? "";
    // Check if there's an opening quote before the partial
    const beforePartial = beforeCursor.substring(0, cursorCol - partial.length);
    const hasQuote = beforePartial.endsWith('"');
    const replaceStart = hasQuote ? cursorCol - partial.length - 1 : cursorCol - partial.length;
    return {
      fieldName: "dependencies",
      values: ids,
      details: titles,
      insertTexts: ids.map((id) => `"${id}"`),
      replaceColumns: [replaceStart, cursorCol],
    };
  }

  // Parent field: `parent: ` (scalar value)
  const parentMatch = beforeCursor.match(/^parent:\s*"?(\S*)$/);
  if (parentMatch) {
    const partial = parentMatch[1] ?? "";
    const beforePartial = beforeCursor.substring(0, cursorCol - partial.length);
    const hasQuote = beforePartial.endsWith('"');
    const replaceStart = hasQuote ? cursorCol - partial.length - 1 : cursorCol - partial.length;
    return {
      fieldName: "parent",
      values: ids,
      details: titles,
      insertTexts: ids.map((id) => `"${id}"`),
      replaceColumns: [replaceStart, cursorCol],
    };
  }

  return undefined;
}

/**
 * Walk backwards from cursorLine to find the YAML field name that owns
 * the current block array items.
 */
function findParentField(
  lines: string[],
  cursorLine: number,
  frontmatterStartLine: number
): string | null {
  for (let i = cursorLine - 1; i > frontmatterStartLine; i--) {
    const line = lines[i];
    // A top-level field line: "fieldname:" at column 0
    const fieldMatch = line.match(/^(\w[\w-]*):/);
    if (fieldMatch) return fieldMatch[1];
  }
  return null;
}

const TASK_ID_FIELD_SET = new Set(TASK_ID_FIELDS);

export class TaskmdCompletionProvider implements vscode.CompletionItemProvider {
  private cachedTaskEntries: TaskEntry[] | null = null;

  invalidateTaskCache(): void {
    this.cachedTaskEntries = null;
  }

  private getTaskEntries(filePath: string): TaskEntry[] {
    if (!this.cachedTaskEntries) {
      this.cachedTaskEntries = scanTaskIds(filePath);
    }
    return this.cachedTaskEntries;
  }

  provideCompletionItems(
    document: vscode.TextDocument,
    position: vscode.Position
  ): vscode.CompletionItem[] | undefined {
    const text = document.getText();
    const bounds = findFrontmatterBounds(text);
    if (!bounds) return undefined;

    const filePath = document.uri.fsPath;
    const scopes = readScopes(filePath);
    const taskEntries = this.getTaskEntries(filePath);
    const lines = text.split("\n");
    const result = resolveCompletions(
      lines,
      position.line,
      position.character,
      bounds.startLine,
      bounds.endLine,
      scopes,
      taskEntries
    );
    if (!result) return undefined;

    const replaceRange = new vscode.Range(
      position.line, result.replaceColumns[0],
      position.line, result.replaceColumns[1]
    );

    return result.values.map((val, i) => {
      const kind = TASK_ID_FIELD_SET.has(result.fieldName)
        ? vscode.CompletionItemKind.Reference
        : result.fieldName === "touches"
          ? vscode.CompletionItemKind.Value
          : vscode.CompletionItemKind.EnumMember;
      const item = new vscode.CompletionItem(val, kind);
      item.insertText = result.insertTexts[i];
      item.range = replaceRange;
      item.detail = result.details?.[i] ?? `taskmd ${result.fieldName}`;
      return item;
    });
  }
}
