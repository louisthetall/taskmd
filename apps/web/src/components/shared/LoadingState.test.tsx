import { describe, it, expect } from "vitest";
import { render } from "@testing-library/react";
import { LoadingState } from "./LoadingState.tsx";

describe("LoadingState", () => {
  it("renders default skeleton when no variant provided", () => {
    const { container } = render(<LoadingState />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    // Default has centered spinner with py-12
    expect(container.querySelector(".py-12")).toBeInTheDocument();
  });

  it("renders default skeleton for variant='default'", () => {
    const { container } = render(<LoadingState variant="default" />);
    expect(container.querySelector(".py-12")).toBeInTheDocument();
  });

  it("renders table skeleton with header and 6 rows", () => {
    const { container } = render(<LoadingState variant="table" />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    // 6 rows, each a flex container with gap-4
    const rows = container.querySelectorAll(".flex.gap-4");
    expect(rows).toHaveLength(6);
  });

  it("renders board skeleton with 4 columns", () => {
    const { container } = render(<LoadingState variant="board" />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    const columns = container.querySelectorAll(".flex-1");
    expect(columns).toHaveLength(4);
  });

  it("renders graph skeleton with centered content", () => {
    const { container } = render(<LoadingState variant="graph" />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    expect(container.querySelector(".text-center")).toBeInTheDocument();
  });

  it("renders cards skeleton with 4 metric cards and 2 chart cards", () => {
    const { container } = render(<LoadingState variant="cards" />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    // 4 metric cards in the first grid
    const metricGrid = container.querySelector(".grid-cols-2");
    expect(metricGrid).toBeInTheDocument();
    const metricCards = metricGrid!.children;
    expect(metricCards).toHaveLength(4);
    // 2 chart cards in the second grid
    const chartGrid = container.querySelector(".grid-cols-1");
    expect(chartGrid).toBeInTheDocument();
    expect(chartGrid!.children).toHaveLength(2);
  });

  it("renders detail skeleton with header, info grid, and body", () => {
    const { container } = render(<LoadingState variant="detail" />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
    // Info grid has 3 items
    const infoGrid = container.querySelector(".grid-cols-2");
    expect(infoGrid).toBeInTheDocument();
    expect(infoGrid!.children).toHaveLength(3);
    // Body section with 3 text lines
    const bodySection = container.querySelector(".border-t");
    expect(bodySection).toBeInTheDocument();
    expect(bodySection!.children).toHaveLength(3);
  });
});
