import { describe, it, expect, vi } from "vitest";
import { screen, fireEvent } from "@testing-library/react";

const mockNavigate = vi.fn();

vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});
vi.mock("@xyflow/react", () => ({
  ReactFlow: ({ nodes, edges, onNodeClick, children }: any) => (
    <div data-testid="react-flow" data-nodes={nodes.length} data-edges={edges.length}>
      {nodes.map((n: any) => (
        <div key={n.id} data-testid={`node-${n.id}`} onClick={(e) => onNodeClick?.(e, n)}>
          {n.data?.label}
        </div>
      ))}
      {children}
    </div>
  ),
  Controls: () => <div data-testid="controls" />,
  Background: () => <div data-testid="background" />,
  BackgroundVariant: { Dots: "dots" },
}));
vi.mock("../../hooks/use-theme.ts", () => ({
  useTheme: () => ({ theme: "light", toggle: vi.fn() }),
}));
vi.mock("./TaskNode.tsx", () => ({ TaskNode: () => null }));

import { GraphView } from "./GraphView.tsx";
import { renderWithProviders } from "../../test-utils/render.ts";

describe("GraphView", () => {
  it("shows 'No tasks to display' when nodes is empty", () => {
    renderWithProviders(<GraphView nodes={[]} edges={[]} />);
    expect(screen.getByText("No tasks to display")).toBeInTheDocument();
  });

  it("renders ReactFlow when nodes exist", () => {
    const nodes = [
      { id: "1", position: { x: 0, y: 0 }, data: { label: "Task A", taskId: "001" } },
      { id: "2", position: { x: 100, y: 0 }, data: { label: "Task B", taskId: "002" } },
    ];
    const edges = [{ id: "e1-2", source: "1", target: "2" }];
    renderWithProviders(<GraphView nodes={nodes} edges={edges} />);
    expect(screen.getByTestId("react-flow")).toBeInTheDocument();
    expect(screen.getByTestId("react-flow")).toHaveAttribute("data-nodes", "2");
    expect(screen.getByTestId("react-flow")).toHaveAttribute("data-edges", "1");
  });

  it("renders controls and background", () => {
    const nodes = [{ id: "1", position: { x: 0, y: 0 }, data: { label: "Task A", taskId: "001" } }];
    renderWithProviders(<GraphView nodes={nodes} edges={[]} />);
    expect(screen.getByTestId("controls")).toBeInTheDocument();
    expect(screen.getByTestId("background")).toBeInTheDocument();
  });

  it("navigates to task detail on node click", () => {
    const nodes = [{ id: "1", position: { x: 0, y: 0 }, data: { label: "Task A", taskId: "001" } }];
    renderWithProviders(<GraphView nodes={nodes} edges={[]} />);
    fireEvent.click(screen.getByTestId("node-1"));
    expect(mockNavigate).toHaveBeenCalledWith("/tasks/001");
  });
});
