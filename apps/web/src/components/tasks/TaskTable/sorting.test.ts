import { describe, it, expect } from "vitest";
import { countUnmetDependencies, compareBlocked, comparePriority } from "./sorting.ts";

describe("countUnmetDependencies", () => {
  it("returns 0 for null deps", () => {
    expect(countUnmetDependencies(null)).toBe(0);
  });

  it("returns 0 for empty deps", () => {
    expect(countUnmetDependencies([])).toBe(0);
  });

  it("counts all deps as unmet when no status map provided", () => {
    expect(countUnmetDependencies(["a", "b"])).toBe(2);
  });

  it("excludes completed deps", () => {
    const statusMap = new Map([
      ["a", "completed"],
      ["b", "in-progress"],
    ]);
    expect(countUnmetDependencies(["a", "b", "c"], statusMap)).toBe(2);
  });

  it("returns 0 when all deps are completed", () => {
    const statusMap = new Map([
      ["a", "completed"],
      ["b", "completed"],
    ]);
    expect(countUnmetDependencies(["a", "b"], statusMap)).toBe(0);
  });
});

describe("compareBlocked", () => {
  it("sorts task with fewer unmet deps first", () => {
    const statusMap = new Map([["a", "completed"]]);
    expect(compareBlocked(["a"], ["b"], statusMap)).toBeLessThan(0);
  });

  it("returns 0 for equal unmet counts", () => {
    expect(compareBlocked(["a"], ["b"])).toBe(0);
  });

  it("handles null deps", () => {
    expect(compareBlocked(null, ["a", "b"])).toBeLessThan(0);
  });
});

describe("comparePriority", () => {
  it("sorts critical before high", () => {
    expect(comparePriority("critical", "high")).toBeLessThan(0);
  });

  it("sorts high before medium", () => {
    expect(comparePriority("high", "medium")).toBeLessThan(0);
  });

  it("sorts medium before low", () => {
    expect(comparePriority("medium", "low")).toBeLessThan(0);
  });

  it("sorts low before unset", () => {
    expect(comparePriority("low", null)).toBeLessThan(0);
  });

  it("returns 0 for equal priorities", () => {
    expect(comparePriority("high", "high")).toBe(0);
  });

  it("treats null and undefined the same (unset)", () => {
    expect(comparePriority(null, undefined)).toBe(0);
  });

  it("sorts correctly when used with Array.sort", () => {
    const priorities = ["low", "critical", null, "medium", "high", undefined];
    const sorted = [...priorities].sort((a, b) => comparePriority(a, b));
    expect(sorted).toEqual(["critical", "high", "medium", "low", null, undefined]);
  });
});
