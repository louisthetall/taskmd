import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatusBadge, PriorityBadge, TypeBadge, BlockedStatusBadge } from "./Badges.tsx";
import { STATUS_COLORS, PRIORITY_COLORS, TYPE_COLORS } from "./constants.ts";

describe("StatusBadge", () => {
  it.each(Object.keys(STATUS_COLORS))("renders '%s' with correct color classes", (status) => {
    const { container } = render(<StatusBadge status={status} />);
    const badge = container.querySelector("span")!;
    expect(badge).toHaveTextContent(status);
    for (const cls of STATUS_COLORS[status].split(" ")) {
      expect(badge.className).toContain(cls);
    }
  });

  it("falls back to gray for unknown status", () => {
    const { container } = render(<StatusBadge status="unknown" />);
    const badge = container.querySelector("span")!;
    expect(badge).toHaveTextContent("unknown");
    expect(badge.className).toContain("bg-gray-100");
  });
});

describe("PriorityBadge", () => {
  it.each(Object.keys(PRIORITY_COLORS))("renders '%s' with correct color classes", (priority) => {
    const { container } = render(<PriorityBadge priority={priority} />);
    const badge = container.querySelector("span")!;
    expect(badge).toHaveTextContent(priority);
    for (const cls of PRIORITY_COLORS[priority].split(" ")) {
      expect(badge.className).toContain(cls);
    }
  });

  it("falls back to gray for unknown priority", () => {
    const { container } = render(<PriorityBadge priority="unknown" />);
    const badge = container.querySelector("span")!;
    expect(badge.className).toContain("bg-gray-100");
  });
});

describe("TypeBadge", () => {
  it.each(Object.keys(TYPE_COLORS))("renders '%s' with correct color classes", (type) => {
    const { container } = render(<TypeBadge type={type} />);
    const badge = container.querySelector("span")!;
    expect(badge).toHaveTextContent(type);
    for (const cls of TYPE_COLORS[type].split(" ")) {
      expect(badge.className).toContain(cls);
    }
  });

  it("falls back to gray for unknown type", () => {
    const { container } = render(<TypeBadge type="unknown" />);
    const badge = container.querySelector("span")!;
    expect(badge.className).toContain("bg-gray-100");
  });
});

describe("BlockedStatusBadge", () => {
  it("renders Ready badge when dependencies is null", () => {
    render(<BlockedStatusBadge dependencies={null} />);
    expect(screen.getByText("Ready")).toBeInTheDocument();
    expect(screen.getByText("✓")).toBeInTheDocument();
    expect(screen.getByLabelText("Task is ready to work on")).toBeInTheDocument();
  });

  it("renders Ready badge when dependencies is empty array", () => {
    render(<BlockedStatusBadge dependencies={[]} />);
    expect(screen.getByText("Ready")).toBeInTheDocument();
  });

  it("renders Blocked badge with count for single dependency", () => {
    render(<BlockedStatusBadge dependencies={["005"]} />);
    expect(screen.getByText("(1)")).toBeInTheDocument();
    expect(screen.getByText("⚠")).toBeInTheDocument();
  });

  it("renders Blocked badge with count for multiple dependencies", () => {
    render(<BlockedStatusBadge dependencies={["005", "010", "015"]} />);
    expect(screen.getByText("(3)")).toBeInTheDocument();
  });

  it("shows tooltip with blocked-by IDs", () => {
    render(<BlockedStatusBadge dependencies={["005", "010"]} />);
    const badge = screen.getByLabelText("Blocked by: 005, 010");
    expect(badge).toHaveAttribute("title", "Blocked by: 005, 010");
  });

  it("applies green styling for Ready state", () => {
    const { container } = render(<BlockedStatusBadge dependencies={null} />);
    const badge = container.querySelector("span")!;
    expect(badge.className).toContain("bg-green-100");
  });

  it("applies amber styling for Blocked state", () => {
    const { container } = render(<BlockedStatusBadge dependencies={["001"]} />);
    const badge = container.querySelector("span")!;
    expect(badge.className).toContain("bg-amber-100");
  });
});
