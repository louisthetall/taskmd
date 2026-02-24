import { describe, it, expect } from "vitest";
import { parseFrontmatter } from "../src/frontmatter";
import { validate, type ValidationIssue } from "../src/validate";

function getIssues(text: string): ValidationIssue[] {
  const parsed = parseFrontmatter(text);
  if (!parsed) throw new Error("Failed to parse frontmatter");
  return validate(parsed);
}

function issueFor(issues: ValidationIssue[], field: string): ValidationIssue | undefined {
  return issues.find((i) => i.field === field);
}

describe("validate: required fields", () => {
  it("reports error when id is missing", () => {
    const issues = getIssues(`---
title: "Test"
---`);
    const issue = issueFor(issues, "id");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("missing required field: id");
    expect(issue!.target).toBe("block");
  });

  it("reports error when title is missing", () => {
    const issues = getIssues(`---
id: "042"
---`);
    const issue = issueFor(issues, "title");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("missing required field: title");
  });

  it("reports error when id is empty", () => {
    const issues = getIssues(`---
id: ""
title: "Test"
---`);
    const issue = issueFor(issues, "id");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("required field 'id' is empty");
  });

  it("no errors for valid required fields", () => {
    const issues = getIssues(`---
id: "042"
title: "My Task"
---`);
    expect(issues.filter((i) => i.severity === "error")).toHaveLength(0);
  });
});

describe("validate: enum fields", () => {
  it("no error for valid status", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
status: in-progress
---`);
    expect(issueFor(issues, "status")).toBeUndefined();
  });

  it("error for invalid status", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
status: doing
---`);
    const issue = issueFor(issues, "status");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("invalid status");
    expect(issue!.message).toContain("doing");
  });

  it("allows empty status", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
status: ""
---`);
    expect(issueFor(issues, "status")).toBeUndefined();
  });

  it("error for invalid priority", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
priority: urgent
---`);
    const issue = issueFor(issues, "priority");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("invalid priority");
  });

  it("allows empty priority", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
priority: ""
---`);
    expect(issueFor(issues, "priority")).toBeUndefined();
  });

  it("error for invalid effort", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
effort: huge
---`);
    const issue = issueFor(issues, "effort");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("allows valid effort", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
effort: large
---`);
    expect(issueFor(issues, "effort")).toBeUndefined();
  });

  it("warning (not error) for invalid type", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
type: epic
---`);
    const issue = issueFor(issues, "type");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("warning");
    expect(issue!.message).toContain("invalid type");
  });

  it("allows valid type values", () => {
    for (const t of ["feature", "bug", "improvement", "chore", "docs"]) {
      const issues = getIssues(`---
id: "042"
title: "Test"
type: ${t}
---`);
      expect(issueFor(issues, "type")).toBeUndefined();
    }
  });
});

describe("validate: array fields", () => {
  it("error when dependencies is a scalar", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
dependencies: "041"
---`);
    const issue = issueFor(issues, "dependencies");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("should be an array");
  });

  it("no error when dependencies is an array", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
dependencies:
  - "041"
---`);
    expect(issueFor(issues, "dependencies")).toBeUndefined();
  });

  it("no error when tags is an empty array", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
tags: []
---`);
    expect(issueFor(issues, "tags")).toBeUndefined();
  });

  it("error when tags is a scalar", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
tags: backend
---`);
    const issue = issueFor(issues, "tags");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("error when touches is a scalar", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
touches: cli
---`);
    const issue = issueFor(issues, "touches");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("error when pr is a scalar", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
pr: "https://github.com/foo/bar/1"
---`);
    const issue = issueFor(issues, "pr");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });
});

describe("validate: date fields", () => {
  it("no error for valid date", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
created: "2026-02-14"
---`);
    expect(issueFor(issues, "created")).toBeUndefined();
  });

  it("error for invalid date format", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
created: "Feb 14, 2026"
---`);
    const issue = issueFor(issues, "created");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
    expect(issue!.message).toContain("YYYY-MM-DD");
  });

  it("allows unquoted date (parsed as Date object)", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
created: 2026-02-14
---`);
    expect(issueFor(issues, "created")).toBeUndefined();
  });
});

describe("validate: unknown fields", () => {
  it("warns on unknown field", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
stauts: pending
---`);
    const issue = issueFor(issues, "stauts");
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("warning");
    expect(issue!.target).toBe("key");
    expect(issue!.message).toContain("stauts");
  });

  it("no warning for valid fields", () => {
    const issues = getIssues(`---
id: "042"
title: "My Task"
status: pending
priority: high
tags:
  - cli
---`);
    const warnings = issues.filter((i) => i.severity === "warning" && i.message.includes("unknown"));
    expect(warnings).toHaveLength(0);
  });

  it("no warning for verify field", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
verify:
  - type: bash
    run: "echo hello"
---`);
    expect(issueFor(issues, "verify")).toBeUndefined();
  });

  it("warning message includes field name", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
foobar: something
---`);
    const issue = issueFor(issues, "foobar");
    expect(issue).toBeDefined();
    expect(issue!.message).toBe("unknown field 'foobar'");
  });
});

describe("validate: verify steps", () => {
  it("error when step missing type", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
verify:
  - run: "echo hello"
---`);
    const issue = issues.find((i) => i.message.includes("missing required field 'type'"));
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("error when bash step missing run", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
verify:
  - type: bash
---`);
    const issue = issues.find((i) => i.message.includes("bash step missing required field 'run'"));
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("error when assert step missing check", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
verify:
  - type: assert
---`);
    const issue = issues.find((i) => i.message.includes("assert step missing required field 'check'"));
    expect(issue).toBeDefined();
    expect(issue!.severity).toBe("error");
  });

  it("no error for valid verify steps", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
verify:
  - type: bash
    run: "go test ./..."
  - type: assert
    check: "all tests pass"
---`);
    const verifyIssues = issues.filter((i) => i.field === "verify");
    expect(verifyIssues).toHaveLength(0);
  });

  it("no error when verify is absent", () => {
    const issues = getIssues(`---
id: "042"
title: "Test"
---`);
    const verifyIssues = issues.filter((i) => i.field === "verify");
    expect(verifyIssues).toHaveLength(0);
  });
});
