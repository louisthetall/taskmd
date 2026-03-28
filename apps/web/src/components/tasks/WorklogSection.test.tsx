import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { WorklogSection } from "./WorklogSection.tsx";
import { createWorklogEntry } from "../../test-utils/fixtures.ts";

describe("WorklogSection", () => {
  it("renders heading with entry count", () => {
    const entries = [
      createWorklogEntry(),
      createWorklogEntry({ timestamp: "2026-01-16T09:00:00Z" }),
    ];
    render(<WorklogSection entries={entries} />);

    expect(screen.getByText("Worklog")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });

  it("renders multiple entries with formatted timestamps", () => {
    const entries = [
      createWorklogEntry({ timestamp: "2026-01-15T10:30:00Z", content: "First entry." }),
      createWorklogEntry({ timestamp: "2026-01-16T14:00:00Z", content: "Second entry." }),
    ];
    render(<WorklogSection entries={entries} />);

    const times = screen.getAllByRole("time");
    expect(times).toHaveLength(2);
    expect(times[0]).toHaveTextContent(new Date("2026-01-15T10:30:00Z").toLocaleString());
    expect(times[1]).toHaveTextContent(new Date("2026-01-16T14:00:00Z").toLocaleString());
  });

  it("renders markdown content", () => {
    const entries = [
      createWorklogEntry({ content: "This is **bold** text." }),
    ];
    render(<WorklogSection entries={entries} />);

    const bold = screen.getByText("bold");
    expect(bold.tagName).toBe("STRONG");
  });

  it("renders empty list with count of 0", () => {
    render(<WorklogSection entries={[]} />);

    expect(screen.getByText("Worklog")).toBeInTheDocument();
    expect(screen.getByText("0")).toBeInTheDocument();
    expect(screen.queryAllByRole("time")).toHaveLength(0);
  });

  it("renders GFM features like strikethrough", () => {
    const entries = [
      createWorklogEntry({ content: "This is ~~removed~~ text." }),
    ];
    render(<WorklogSection entries={entries} />);

    const deleted = screen.getByText("removed");
    expect(deleted.tagName).toBe("DEL");
  });
});
