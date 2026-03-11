import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatsPage } from "./StatsPage.tsx";

vi.mock("../hooks/use-stats.ts", () => ({
  useStats: vi.fn(),
}));

vi.mock("../components/stats/StatsView.tsx", () => ({
  StatsView: ({ stats }: { stats: { total_tasks: number } }) => (
    <div data-testid="stats-view">Stats: {stats.total_tasks} tasks</div>
  ),
}));

import { useStats } from "../hooks/use-stats.ts";
const mockUseStats = vi.mocked(useStats);

describe("StatsPage", () => {
  it("renders loading state", () => {
    mockUseStats.mockReturnValue({
      data: undefined,
      error: undefined,
      isLoading: true,
      mutate: vi.fn(),
      isValidating: false,
    });
    const { container } = render(<StatsPage />);
    // LoadingState renders animated skeleton divs
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
  });

  it("renders error state", () => {
    mockUseStats.mockReturnValue({
      data: undefined,
      error: new Error("Server error"),
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    render(<StatsPage />);
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders empty state when total_tasks is 0", () => {
    mockUseStats.mockReturnValue({
      data: {
        total_tasks: 0,
        tasks_by_status: {},
        tasks_by_priority: {},
        tasks_by_effort: {},
        tasks_by_milestone: {},
        blocked_tasks_count: 0,
        critical_path_length: 0,
        max_dependency_depth: 0,
        avg_dependencies_per_task: 0,
        tags_by_count: [],
      },
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    render(<StatsPage />);
    expect(screen.getByText(/No tasks found/)).toBeInTheDocument();
  });

  it("renders StatsView when data is available", () => {
    mockUseStats.mockReturnValue({
      data: {
        total_tasks: 42,
        tasks_by_status: { pending: 10 },
        tasks_by_priority: { high: 5 },
        tasks_by_effort: { small: 20 },
        tasks_by_milestone: {},
        blocked_tasks_count: 2,
        critical_path_length: 3,
        max_dependency_depth: 4,
        avg_dependencies_per_task: 1.5,
        tags_by_count: [],
      },
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    render(<StatsPage />);
    expect(screen.getByTestId("stats-view")).toBeInTheDocument();
    expect(screen.getByText("Stats: 42 tasks")).toBeInTheDocument();
  });
});
