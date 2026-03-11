import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { StatsView } from "./StatsView.tsx";
import type { Stats } from "../../api/types.ts";

function makeStats(overrides: Partial<Stats> = {}): Stats {
  return {
    total_tasks: 42,
    tasks_by_status: { pending: 10, "in-progress": 5, completed: 27 },
    tasks_by_priority: { critical: 2, high: 8, medium: 20, low: 12 },
    tasks_by_effort: { small: 15, medium: 18, large: 9 },
    tasks_by_milestone: {},
    blocked_tasks_count: 3,
    critical_path_length: 7,
    max_dependency_depth: 4,
    avg_dependencies_per_task: 1.5,
    tags_by_count: [
      { tag: "backend", count: 12 },
      { tag: "frontend", count: 8 },
    ],
    ...overrides,
  };
}

function renderWithRouter(stats: Stats) {
  return render(
    <MemoryRouter>
      <StatsView stats={stats} />
    </MemoryRouter>,
  );
}

describe("StatsView", () => {
  it("renders all four metric cards with correct values", () => {
    renderWithRouter(makeStats());

    expect(screen.getByText("Total Tasks")).toBeInTheDocument();
    expect(screen.getByText("42")).toBeInTheDocument();

    expect(screen.getByText("Blocked")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();

    expect(screen.getByText("Critical Path")).toBeInTheDocument();
    expect(screen.getByText("7")).toBeInTheDocument();

    expect(screen.getByText("Avg Deps/Task")).toBeInTheDocument();
    expect(screen.getByText("1.5")).toBeInTheDocument();
  });

  it("renders breakdown cards with correct titles", () => {
    renderWithRouter(makeStats());
    expect(screen.getByText("By Status")).toBeInTheDocument();
    expect(screen.getByText("By Priority")).toBeInTheDocument();
    expect(screen.getByText("By Effort")).toBeInTheDocument();
  });

  it("renders status breakdown entries with counts", () => {
    renderWithRouter(makeStats());
    expect(screen.getByText("pending")).toBeInTheDocument();
    expect(screen.getByText("in-progress")).toBeInTheDocument();
    expect(screen.getByText("completed")).toBeInTheDocument();
  });

  it("filters out zero-count entries from breakdown cards", () => {
    renderWithRouter(makeStats({
      tasks_by_effort: { small: 5, medium: 0, large: 3 },
    }));
    // "medium" appears as a priority label, but as an effort entry with 0 count it should be filtered
    const effortCard = screen.getByText("By Effort").closest("div")!;
    // small and large should be present in the effort card
    expect(effortCard.textContent).toContain("small");
    expect(effortCard.textContent).toContain("large");
  });

  it("shows 'No data' for empty breakdown cards", () => {
    renderWithRouter(makeStats({
      tasks_by_effort: {},
    }));
    expect(screen.getByText("No data")).toBeInTheDocument();
  });

  it("renders breakdown links with correct href", () => {
    renderWithRouter(makeStats());
    const pendingLink = screen.getByRole("link", { name: "pending" });
    expect(pendingLink).toHaveAttribute("href", "/tasks?status=pending");

    const criticalLink = screen.getByRole("link", { name: "critical" });
    expect(criticalLink).toHaveAttribute("href", "/tasks?priority=critical");
  });

  it("renders tags with counts and links", () => {
    renderWithRouter(makeStats());
    expect(screen.getByText("Tags")).toBeInTheDocument();

    const backendLink = screen.getByRole("link", { name: "backend" });
    expect(backendLink).toHaveAttribute("href", "/tasks?tag=backend");

    const frontendLink = screen.getByRole("link", { name: "frontend" });
    expect(frontendLink).toHaveAttribute("href", "/tasks?tag=frontend");

    // Verify counts are rendered next to their tags
    const tagsSection = screen.getByText("Tags").closest("div")!;
    expect(tagsSection.textContent).toContain("backend");
    expect(tagsSection.textContent).toContain("8");
  });

  it("shows 'No tags found' when tags_by_count is empty", () => {
    renderWithRouter(makeStats({ tags_by_count: [] }));
    expect(screen.getByText("No tags found")).toBeInTheDocument();
  });

  it("shows 'No tags found' when tags_by_count is null-ish", () => {
    renderWithRouter(makeStats({ tags_by_count: undefined as unknown as Stats["tags_by_count"] }));
    expect(screen.getByText("No tags found")).toBeInTheDocument();
  });

  it("formats avg_dependencies_per_task with one decimal", () => {
    renderWithRouter(makeStats({ avg_dependencies_per_task: 2 }));
    expect(screen.getByText("2.0")).toBeInTheDocument();
  });
});
