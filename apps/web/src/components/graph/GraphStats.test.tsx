import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { GraphStats } from "./GraphStats.tsx";
import { createGraphData, createGraphNode } from "./../../test-utils/fixtures.ts";

describe("GraphStats", () => {
  it("shows visible and total task counts", () => {
    const data = createGraphData();

    render(<GraphStats data={data} visibleCount={1} />);

    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
    expect(screen.getByText(/tasks/)).toBeInTheDocument();
  });

  it("shows blocked count when blocked nodes exist", () => {
    const data = createGraphData({
      nodes: [
        createGraphNode({ id: "001", status: "blocked" }),
        createGraphNode({ id: "002", status: "pending" }),
      ],
    });

    render(<GraphStats data={data} visibleCount={2} />);

    expect(screen.getByText("1 blocked")).toBeInTheDocument();
  });

  it("does not show blocked when no blocked nodes", () => {
    const data = createGraphData({
      nodes: [
        createGraphNode({ id: "001", status: "pending" }),
        createGraphNode({ id: "002", status: "completed" }),
      ],
    });

    render(<GraphStats data={data} visibleCount={2} />);

    expect(screen.queryByText(/blocked/)).not.toBeInTheDocument();
  });

  it("shows circular dependency warning when cycles present", () => {
    const data = createGraphData({
      cycles: [["001", "002"]],
    });

    render(<GraphStats data={data} visibleCount={2} />);

    expect(screen.getByText("Circular dependencies detected")).toBeInTheDocument();
  });

  it("does not show cycle warning when no cycles", () => {
    const data = createGraphData({ cycles: [] });

    render(<GraphStats data={data} visibleCount={2} />);

    expect(screen.queryByText("Circular dependencies detected")).not.toBeInTheDocument();
  });
});
