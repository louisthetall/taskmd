import { describe, it, expect } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { GraphLegend } from "./GraphLegend.tsx";

describe("GraphLegend", () => {
  it("initially shows Legend button", () => {
    render(<GraphLegend />);

    expect(screen.getByText("Legend")).toBeInTheDocument();
  });

  it("opens panel with status items when Legend is clicked", () => {
    render(<GraphLegend />);

    fireEvent.click(screen.getByText("Legend"));

    expect(screen.getByText("Pending")).toBeInTheDocument();
    expect(screen.getByText("In Progress")).toBeInTheDocument();
    expect(screen.getByText("Completed")).toBeInTheDocument();
    expect(screen.getByText("Blocked")).toBeInTheDocument();
    expect(screen.getByText("Cancelled")).toBeInTheDocument();
  });

  it("shows priority items when open", () => {
    render(<GraphLegend />);

    fireEvent.click(screen.getByText("Legend"));

    expect(screen.getByText("Critical")).toBeInTheDocument();
    expect(screen.getByText("High")).toBeInTheDocument();
  });

  it("shows Depends on edge description when open", () => {
    render(<GraphLegend />);

    fireEvent.click(screen.getByText("Legend"));

    expect(screen.getByText("Depends on")).toBeInTheDocument();
  });

  it("closes panel when close button is clicked", () => {
    render(<GraphLegend />);

    fireEvent.click(screen.getByText("Legend"));
    expect(screen.getByText("Pending")).toBeInTheDocument();

    fireEvent.click(screen.getByLabelText("Close legend"));
    expect(screen.queryByText("Pending")).not.toBeInTheDocument();
  });
});
