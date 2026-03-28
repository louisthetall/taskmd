import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

const mockFitView = vi.fn();
const mockGetViewport = vi.fn(() => ({ x: 0, y: 0, zoom: 1 }));
const mockSetViewport = vi.fn();

vi.mock("@xyflow/react", () => ({
  useReactFlow: () => ({
    fitView: mockFitView,
    getViewport: mockGetViewport,
    setViewport: mockSetViewport,
  }),
}));

import { GraphSearch } from "./GraphSearch.tsx";

describe("GraphSearch", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders search input with placeholder", () => {
    render(<GraphSearch query="" onQueryChange={vi.fn()} matchedNodeIds={new Set()} />);
    expect(screen.getByPlaceholderText("Search tasks...")).toBeInTheDocument();
  });

  it("typing calls onQueryChange", () => {
    const onQueryChange = vi.fn();
    render(<GraphSearch query="" onQueryChange={onQueryChange} matchedNodeIds={new Set()} />);
    fireEvent.change(screen.getByPlaceholderText("Search tasks..."), { target: { value: "test" } });
    expect(onQueryChange).toHaveBeenCalledWith("test");
  });

  it("shows match count when query is non-empty", () => {
    render(<GraphSearch query="task" onQueryChange={vi.fn()} matchedNodeIds={new Set(["1", "2", "3"])} />);
    expect(screen.getByText("3 found")).toBeInTheDocument();
  });

  it("shows clear button when query is non-empty", () => {
    render(<GraphSearch query="task" onQueryChange={vi.fn()} matchedNodeIds={new Set(["1"])} />);
    expect(screen.getByLabelText("Clear search")).toBeInTheDocument();
  });

  it("clicking clear calls onQueryChange with empty string", () => {
    const onQueryChange = vi.fn();
    render(<GraphSearch query="task" onQueryChange={onQueryChange} matchedNodeIds={new Set(["1"])} />);
    fireEvent.click(screen.getByLabelText("Clear search"));
    expect(onQueryChange).toHaveBeenCalledWith("");
  });

  it("does not show match count when query is empty", () => {
    render(<GraphSearch query="" onQueryChange={vi.fn()} matchedNodeIds={new Set()} />);
    expect(screen.queryByText("found")).not.toBeInTheDocument();
    expect(screen.queryByLabelText("Clear search")).not.toBeInTheDocument();
  });
});
