import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { vi } from "vitest";

vi.mock("@xyflow/react", () => ({
  Handle: ({ type, position }: { type: string; position: string }) => (
    <div data-testid={`handle-${type}`} data-position={position} />
  ),
  Position: { Top: "top", Bottom: "bottom" },
  memo: (fn: any) => fn,
}));

import { TaskNode } from "./TaskNode.tsx";

describe("TaskNode", () => {
  it("renders task ID and label", () => {
    render(
      <TaskNode
        data={{ taskId: "042", label: "Implement feature", status: "pending" }}
      />,
    );

    expect(screen.getByText("042")).toBeInTheDocument();
    expect(screen.getByText("Implement feature")).toBeInTheDocument();
  });

  it("renders target and source handles", () => {
    render(
      <TaskNode
        data={{ taskId: "001", label: "Task", status: "pending" }}
      />,
    );

    const target = screen.getByTestId("handle-target");
    const source = screen.getByTestId("handle-source");

    expect(target).toHaveAttribute("data-position", "top");
    expect(source).toHaveAttribute("data-position", "bottom");
  });

  it("applies highlight ring when highlighted", () => {
    const { container } = render(
      <TaskNode
        data={{ taskId: "001", label: "Task", status: "pending", highlighted: true }}
      />,
    );

    const node = container.querySelector(".ring-2.ring-blue-500");
    expect(node).toBeInTheDocument();
  });

  it("applies dimmed opacity when dimmed", () => {
    const { container } = render(
      <TaskNode
        data={{ taskId: "001", label: "Task", status: "pending", dimmed: true }}
      />,
    );

    const node = container.querySelector(".opacity-40");
    expect(node).toBeInTheDocument();
  });

  it("shows priority border for critical priority", () => {
    const { container } = render(
      <TaskNode
        data={{ taskId: "001", label: "Task", status: "pending", priority: "critical" }}
      />,
    );

    const node = container.querySelector(".border-l-4.border-l-red-500");
    expect(node).toBeInTheDocument();
  });
});
