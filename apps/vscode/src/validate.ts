import type { ParsedFrontmatter, FrontmatterField } from "./frontmatter";
import {
  REQUIRED_FIELDS,
  ENUM_FIELDS,
  STRING_ARRAY_FIELDS,
  DATE_FIELDS,
  KNOWN_FIELDS,
} from "./schema";

export type Severity = "error" | "warning";
export type Target = "key" | "value" | "block";

export interface ValidationIssue {
  severity: Severity;
  message: string;
  field: string;
  target: Target;
}

const DATE_REGEX = /^\d{4}-\d{2}-\d{2}$/;

/**
 * Validate parsed frontmatter against the taskmd schema.
 * Returns a list of issues found.
 */
export function validate(frontmatter: ParsedFrontmatter): ValidationIssue[] {
  const issues: ValidationIssue[] = [];

  checkUnknownFields(frontmatter, issues);
  checkRequiredFields(frontmatter, issues);
  checkEnumFields(frontmatter, issues);
  checkArrayFields(frontmatter, issues);
  checkDateFields(frontmatter, issues);
  checkVerifySteps(frontmatter, issues);

  return issues;
}

function checkUnknownFields(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  for (const [name] of fm.fields) {
    if (!KNOWN_FIELDS.has(name)) {
      issues.push({
        severity: "warning",
        message: `unknown field '${name}'`,
        field: name,
        target: "key",
      });
    }
  }
}

function checkRequiredFields(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  for (const name of REQUIRED_FIELDS) {
    const field = fm.fields.get(name);
    if (!field) {
      issues.push({
        severity: "error",
        message: `missing required field: ${name}`,
        field: name,
        target: "block",
      });
    } else if (field.value === null || field.value === undefined || field.value === "") {
      issues.push({
        severity: "error",
        message: `required field '${name}' is empty`,
        field: name,
        target: "value",
      });
    }
  }
}

function checkEnumFields(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  for (const [name, allowed] of Object.entries(ENUM_FIELDS)) {
    const field = fm.fields.get(name);
    if (!field) continue;

    const val = field.value;
    // Empty is always allowed
    if (val === null || val === undefined || val === "") continue;

    if (typeof val !== "string" || !allowed.includes(val)) {
      const severity: Severity = name === "type" ? "warning" : "error";
      issues.push({
        severity,
        message: `invalid ${name}: '${val}' (valid values: ${allowed.join(", ")})`,
        field: name,
        target: "value",
      });
    }
  }
}

function checkArrayFields(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  for (const name of STRING_ARRAY_FIELDS) {
    const field = fm.fields.get(name);
    if (!field) continue;

    const val = field.value;
    // null/undefined means field present but no value, which is fine (treated as empty array)
    if (val === null || val === undefined) continue;

    if (!Array.isArray(val)) {
      issues.push({
        severity: "error",
        message: `'${name}' should be an array, got ${typeof val}`,
        field: name,
        target: "value",
      });
    }
  }
}

function checkDateFields(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  for (const name of DATE_FIELDS) {
    const field = fm.fields.get(name);
    if (!field) continue;

    const val = field.value;
    if (val === null || val === undefined || val === "") continue;

    // The YAML parser may parse dates as Date objects — that's valid
    if (val instanceof Date) continue;

    if (typeof val !== "string" || !DATE_REGEX.test(val)) {
      issues.push({
        severity: "error",
        message: `'${name}' must be a valid date in YYYY-MM-DD format`,
        field: name,
        target: "value",
      });
    }
  }
}

function checkVerifySteps(
  fm: ParsedFrontmatter,
  issues: ValidationIssue[]
): void {
  const field = fm.fields.get("verify");
  if (!field) return;

  const val = field.value;
  if (!Array.isArray(val)) return;

  for (let i = 0; i < val.length; i++) {
    const step = val[i];
    if (!step || typeof step !== "object") continue;

    const rec = step as Record<string, unknown>;

    if (!rec.type || rec.type === "") {
      issues.push({
        severity: "error",
        message: `verify[${i}]: missing required field 'type'`,
        field: "verify",
        target: "value",
      });
      continue;
    }

    if (rec.type === "bash" && (!rec.run || rec.run === "")) {
      issues.push({
        severity: "error",
        message: `verify[${i}]: bash step missing required field 'run'`,
        field: "verify",
        target: "value",
      });
    }

    if (rec.type === "assert" && (!rec.check || rec.check === "")) {
      issues.push({
        severity: "error",
        message: `verify[${i}]: assert step missing required field 'check'`,
        field: "verify",
        target: "value",
      });
    }
  }
}
