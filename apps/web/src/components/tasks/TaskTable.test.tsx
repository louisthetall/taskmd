import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import type { Task } from "../../api/types.ts";
import { TaskTable } from "./TaskTable.tsx";

function makeTask(overrides: Partial<Task> = {}): Task {
  return {
    id: "001",
    title: "Test task",
    status: "pending",
    priority: "medium",
    effort: "small",
    type: "feature",
    dependencies: null,
    tags: null,
    milestone: "",
    group: "",
    owner: "",
    parent: "",
    created: "2026-01-01",
    body: "",
    file_path: "tasks/001-test.md",
    ...overrides,
  };
}

const tasks: Task[] = [
  makeTask({ id: "001", status: "pending", priority: "high", type: "feature", tags: ["api"] }),
  makeTask({ id: "002", status: "in-progress", priority: "medium", type: "bug", tags: ["web"] }),
  makeTask({ id: "003", status: "completed", priority: "low", type: "chore", tags: ["docs"] }),
];

function renderWithRouter(
  ui: React.ReactElement,
  { initialEntries = ["/"] }: { initialEntries?: string[] } = {},
) {
  return render(
    <MemoryRouter initialEntries={initialEntries}>{ui}</MemoryRouter>,
  );
}

describe("TaskTable URL sync", () => {
  it("initializes filters from URL search params", () => {
    renderWithRouter(
      <TaskTable
        tasks={tasks}
        initialStatuses={["pending"]}
        initialPriorities={["high"]}
      />,
    );
    // Only task 001 matches pending + high
    expect(screen.getByText("Showing 1 of 3 tasks")).toBeInTheDocument();
  });

  it("initializes tag filters from props", () => {
    renderWithRouter(
      <TaskTable tasks={tasks} initialTags={["api"]} />,
    );
    // Only task 001 has the "api" tag
    expect(screen.getByText("Showing 1 of 3 tasks")).toBeInTheDocument();
  });

  it("initializes effort filters from props", () => {
    const tasksWithEffort = [
      makeTask({ id: "001", status: "pending", effort: "small" }),
      makeTask({ id: "002", status: "pending", effort: "large" }),
    ];
    renderWithRouter(
      <TaskTable tasks={tasksWithEffort} initialEffort={["small"]} />,
    );
    expect(screen.getByText("Showing 1 of 2 tasks")).toBeInTheDocument();
  });

  it("initializes type filters from props", () => {
    renderWithRouter(
      <TaskTable tasks={tasks} initialTypes={["bug"]} />,
    );
    // Only task 002 has type "bug"
    expect(screen.getByText("Showing 1 of 3 tasks")).toBeInTheDocument();
  });

  it("shows all tasks with no initial filters", () => {
    renderWithRouter(<TaskTable tasks={tasks} />);
    expect(screen.getByText("Showing 3 of 3 tasks")).toBeInTheDocument();
  });
});

describe("TaskTable clearFilters", () => {
  it("shows Clear filters button when filters are active", () => {
    renderWithRouter(
      <TaskTable tasks={tasks} initialStatuses={["pending"]} />,
    );
    expect(screen.getByText("Clear filters")).toBeInTheDocument();
  });

  it("does not show Clear filters button with default filters", () => {
    renderWithRouter(<TaskTable tasks={tasks} />);
    expect(screen.queryByText("Clear filters")).not.toBeInTheDocument();
  });

  it("resets all filters when Clear filters is clicked", async () => {
    const user = userEvent.setup();
    renderWithRouter(
      <TaskTable tasks={tasks} initialStatuses={["pending"]} />,
    );
    // Initially filtered to 1 task
    expect(screen.getByText("Showing 1 of 3 tasks")).toBeInTheDocument();

    await user.click(screen.getByText("Clear filters"));

    // After clearing, all tasks are shown
    expect(screen.getByText("Showing 3 of 3 tasks")).toBeInTheDocument();
    // Clear filters button should disappear
    expect(screen.queryByText("Clear filters")).not.toBeInTheDocument();
  });
});

describe("TaskTable status toggle", () => {
  it("toggles status filter when clicking a status button", async () => {
    const user = userEvent.setup();
    renderWithRouter(<TaskTable tasks={tasks} />);
    expect(screen.getByText("Showing 3 of 3 tasks")).toBeInTheDocument();

    // Click "pending" to deselect it
    await user.click(screen.getByRole("button", { name: "pending" }));
    // Tasks 002 (in-progress) and 003 (completed) remain
    expect(screen.getByText("Showing 2 of 3 tasks")).toBeInTheDocument();

    // Click "pending" again to re-select it
    await user.click(screen.getByRole("button", { name: "pending" }));
    expect(screen.getByText("Showing 3 of 3 tasks")).toBeInTheDocument();
  });
});

describe("TaskTable priority toggle", () => {
  it("toggles priority filter when clicking a priority button", async () => {
    const user = userEvent.setup();
    renderWithRouter(<TaskTable tasks={tasks} />);

    // Click "high" to deselect it - task 001 has high priority
    await user.click(screen.getByRole("button", { name: "high" }));
    expect(screen.getByText("Showing 2 of 3 tasks")).toBeInTheDocument();
  });
});
