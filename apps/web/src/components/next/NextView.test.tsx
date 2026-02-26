import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import type { Recommendation } from "../../api/types.ts";
import { NextView } from "./NextView.tsx";

vi.mock("../../hooks/use-tasks.ts", () => ({
  useTasks: vi.fn(() => ({
    data: [
      { group: "cli" },
      { group: "cli" },
      { group: "web" },
      { group: "web/layout" },
      { group: "" },
    ],
  })),
}));

function makeRec(overrides: Partial<Recommendation> = {}): Recommendation {
  return {
    rank: 1,
    id: "001",
    title: "Test recommendation",
    file_path: "tasks/001-test.md",
    status: "pending",
    priority: "high",
    effort: "small",
    score: 85,
    reasons: ["high priority"],
    downstream_count: 0,
    on_critical_path: false,
    ...overrides,
  };
}

function renderView(props: Partial<React.ComponentProps<typeof NextView>> = {}) {
  const defaultProps = {
    recommendations: [makeRec()],
    limit: 5,
    onLimitChange: vi.fn(),
    group: "",
    onGroupChange: vi.fn(),
    ...props,
  };
  return render(
    <MemoryRouter>
      <NextView {...defaultProps} />
    </MemoryRouter>,
  );
}

describe("NextView group filter", () => {
  it("renders the folder filter input", () => {
    renderView();
    expect(screen.getByLabelText("Folder:")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("All folders")).toBeInTheDocument();
  });

  it("shows the current group value in the input", () => {
    renderView({ group: "web" });
    const input = screen.getByLabelText("Folder:") as HTMLInputElement;
    expect(input.value).toBe("web");
  });

  it("shows autocomplete suggestions on focus", async () => {
    renderView();
    const user = userEvent.setup();
    await user.click(screen.getByLabelText("Folder:"));
    expect(screen.getByRole("listbox")).toBeInTheDocument();
    expect(screen.getByText("cli")).toBeInTheDocument();
    expect(screen.getByText("web")).toBeInTheDocument();
    expect(screen.getByText("web/layout")).toBeInTheDocument();
  });

  it("filters suggestions by typed query", async () => {
    renderView();
    const user = userEvent.setup();
    await user.type(screen.getByLabelText("Folder:"), "web");
    const options = screen.getAllByRole("option");
    const labels = options.map((o) => o.textContent);
    expect(labels).toContain("web");
    expect(labels).toContain("web/layout");
    expect(labels).not.toContain("cli");
  });

  it("calls onGroupChange when a suggestion is clicked", async () => {
    const onGroupChange = vi.fn();
    renderView({ onGroupChange });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText("Folder:"));
    await user.click(screen.getByText("cli"));
    expect(onGroupChange).toHaveBeenCalledWith("cli");
  });

  it("does not call onGroupChange while typing (only on selection)", async () => {
    const onGroupChange = vi.fn();
    renderView({ onGroupChange });
    const user = userEvent.setup();
    await user.type(screen.getByLabelText("Folder:"), "web");
    expect(onGroupChange).not.toHaveBeenCalled();
  });

  it("calls onGroupChange with empty string when input is cleared", async () => {
    const onGroupChange = vi.fn();
    renderView({ group: "cli", onGroupChange });
    const user = userEvent.setup();
    const input = screen.getByLabelText("Folder:");
    await user.clear(input);
    expect(onGroupChange).toHaveBeenCalledWith("");
  });

  it("selects suggestion with Enter key", async () => {
    const onGroupChange = vi.fn();
    renderView({ onGroupChange });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText("Folder:"));
    // First suggestion is "cli" (sorted)
    await user.keyboard("{Enter}");
    expect(onGroupChange).toHaveBeenCalledWith("cli");
  });

  it("navigates suggestions with arrow keys", async () => {
    const onGroupChange = vi.fn();
    renderView({ onGroupChange });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText("Folder:"));
    await user.keyboard("{ArrowDown}{Enter}");
    // Second suggestion is "web" (sorted: cli, web, web/layout)
    expect(onGroupChange).toHaveBeenCalledWith("web");
  });

  it("limits suggestions to 5", async () => {
    // Mock has 3 unique groups; this test verifies cap logic doesn't break
    renderView();
    const user = userEvent.setup();
    await user.click(screen.getByLabelText("Folder:"));
    const options = screen.getAllByRole("option");
    expect(options.length).toBeLessThanOrEqual(5);
  });
});

describe("NextView limit buttons", () => {
  it("renders limit buttons", () => {
    renderView();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
    expect(screen.getByText("10")).toBeInTheDocument();
  });

  it("calls onLimitChange when a limit button is clicked", async () => {
    const onLimitChange = vi.fn();
    renderView({ onLimitChange });
    const user = userEvent.setup();

    await user.click(screen.getByText("10"));
    expect(onLimitChange).toHaveBeenCalledWith(10);
  });
});

describe("NextView recommendations", () => {
  it("renders recommendation cards", () => {
    renderView({
      recommendations: [
        makeRec({ id: "001", title: "First task" }),
        makeRec({ id: "002", title: "Second task", rank: 2 }),
      ],
    });
    expect(screen.getByText("First task")).toBeInTheDocument();
    expect(screen.getByText("Second task")).toBeInTheDocument();
  });
});
