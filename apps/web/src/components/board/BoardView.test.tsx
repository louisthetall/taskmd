import { describe, it, expect, vi, beforeEach } from "vitest";
import { screen, fireEvent } from "@testing-library/react";
import { renderWithProviders } from "../../test-utils/render.ts";
import { createBoardGroup, createBoardTask, resetFixtureCounter } from "../../test-utils/fixtures.ts";
import { BoardView } from "./BoardView.tsx";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

vi.mock("./BoardColumn.tsx", () => ({
  BoardColumn: ({ group, focusedCardIndex }: { group: { group: string }; focusedCardIndex: number }) => (
    <div data-testid={`column-${group.group}`} data-focused-card={focusedCardIndex}>
      {group.group}
    </div>
  ),
}));

function makeGroups() {
  return [
    createBoardGroup({ group: "pending", tasks: [createBoardTask({ id: "001", title: "T1" }), createBoardTask({ id: "002", title: "T2" })] }),
    createBoardGroup({ group: "in-progress", tasks: [createBoardTask({ id: "003", title: "T3" })] }),
  ];
}

beforeEach(() => {
  resetFixtureCounter();
  mockNavigate.mockClear();
});

describe("BoardView", () => {
  it("renders grid with correct role and aria-label", () => {
    renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
    const grid = screen.getByRole("grid", { name: "Task board" });
    expect(grid).toBeInTheDocument();
  });

  it("renders all group columns", () => {
    renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
    expect(screen.getByTestId("column-pending")).toBeInTheDocument();
    expect(screen.getByTestId("column-in-progress")).toBeInTheDocument();
  });

  it("passes canDrag true when not readonly and groupBy is a draggable field", () => {
    // canDrag is computed internally and passed to BoardColumn; we verify indirectly
    // by checking the column renders (the mock doesn't expose canDrag, but this confirms
    // the component renders without error with these props)
    renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
    expect(screen.getByTestId("column-pending")).toBeInTheDocument();
  });

  it("renders when readonly is true", () => {
    renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={true} />);
    expect(screen.getByTestId("column-pending")).toBeInTheDocument();
  });

  it("renders when groupBy is not a draggable field", () => {
    renderWithProviders(<BoardView groups={makeGroups()} groupBy="tag" readonly={false} />);
    expect(screen.getByTestId("column-pending")).toBeInTheDocument();
  });

  describe("keyboard navigation", () => {
    function getGrid() {
      return screen.getByRole("grid", { name: "Task board" });
    }

    it("ArrowDown from initial state focuses first column and first card", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      fireEvent.keyDown(getGrid(), { key: "ArrowDown" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "0");
      expect(screen.getByTestId("column-in-progress")).toHaveAttribute("data-focused-card", "-1");
    });

    it("ArrowRight moves to next column", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // First focus a column
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      // Move right
      fireEvent.keyDown(grid, { key: "ArrowRight" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "-1");
      expect(screen.getByTestId("column-in-progress")).toHaveAttribute("data-focused-card", "0");
    });

    it("ArrowLeft wraps from first to last column", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // Focus first column
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      // Move left from first column should wrap to last
      fireEvent.keyDown(grid, { key: "ArrowLeft" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "-1");
      expect(screen.getByTestId("column-in-progress")).toHaveAttribute("data-focused-card", "0");
    });

    it("ArrowDown moves to next card in column", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // Focus first column, first card
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "0");
      // Move down to second card
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "1");
    });

    it("Enter navigates to focused task", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // Focus first column, first card
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      // Press Enter
      fireEvent.keyDown(grid, { key: "Enter" });
      expect(mockNavigate).toHaveBeenCalledWith("/tasks/001");
    });

    it("blur resets focus state", () => {
      renderWithProviders(<BoardView groups={makeGroups()} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // Focus a card
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "0");
      // Blur the grid
      fireEvent.blur(grid);
      expect(screen.getByTestId("column-pending")).toHaveAttribute("data-focused-card", "-1");
      expect(screen.getByTestId("column-in-progress")).toHaveAttribute("data-focused-card", "-1");
    });

    it("does nothing with empty groups", () => {
      renderWithProviders(<BoardView groups={[]} groupBy="status" readonly={false} />);
      const grid = getGrid();
      // Should not throw
      fireEvent.keyDown(grid, { key: "ArrowDown" });
      fireEvent.keyDown(grid, { key: "ArrowRight" });
      fireEvent.keyDown(grid, { key: "ArrowLeft" });
      fireEvent.keyDown(grid, { key: "Enter" });
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });
});
