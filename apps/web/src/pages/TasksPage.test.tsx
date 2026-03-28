import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { TasksPage } from "./TasksPage.tsx";
import { createTask } from "../test-utils/index.ts";

vi.mock("../hooks/use-tasks.ts", () => ({
  useTasks: vi.fn(),
}));

vi.mock("../hooks/use-phase.tsx", () => ({
  usePhase: () => ({ phase: null, setPhase: vi.fn() }),
}));

vi.mock("../hooks/use-project.ts", () => ({
  useProject: () => ({ project: null, setProject: vi.fn() }),
}));

vi.mock("../components/tasks/TaskTable.tsx", () => ({
  TaskTable: ({ tasks }: { tasks: unknown[] }) => (
    <div data-testid="task-table">Tasks: {tasks.length}</div>
  ),
}));

import { useTasks } from "../hooks/use-tasks.ts";
const mockUseTasks = vi.mocked(useTasks);

function renderWithRouter() {
  return render(
    <MemoryRouter>
      <TasksPage />
    </MemoryRouter>,
  );
}

describe("TasksPage", () => {
  it("renders loading state", () => {
    mockUseTasks.mockReturnValue({
      data: undefined,
      error: undefined,
      isLoading: true,
      mutate: vi.fn(),
      isValidating: false,
    });
    const { container } = renderWithRouter();
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
  });

  it("renders error state", () => {
    mockUseTasks.mockReturnValue({
      data: undefined,
      error: new Error("Server error"),
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    renderWithRouter();
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders empty state when data is empty array", () => {
    mockUseTasks.mockReturnValue({
      data: [],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    renderWithRouter();
    expect(screen.getByText(/No tasks found/)).toBeInTheDocument();
  });

  it("renders TaskTable when data is available", () => {
    mockUseTasks.mockReturnValue({
      data: [createTask(), createTask()],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    renderWithRouter();
    expect(screen.getByTestId("task-table")).toBeInTheDocument();
    expect(screen.getByText("Tasks: 2")).toBeInTheDocument();
  });
});
