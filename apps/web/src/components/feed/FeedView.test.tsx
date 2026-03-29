import { describe, it, expect } from "vitest";
import { screen } from "@testing-library/react";
import { renderWithProviders } from "../../test-utils/render.ts";
import { FeedView } from "./FeedView.tsx";
import type { FeedEntry } from "../../api/types.ts";

function createGitEntry(overrides: Partial<FeedEntry> = {}): FeedEntry {
  return {
    source: "git",
    hash: "abc123",
    author: "Alice",
    timestamp: "2026-03-01T10:00:00Z",
    message: "feat: add new feature",
    files: [
      {
        path: "tasks/cli/042-add-auth.md",
        status: "modified",
        taskID: "042",
        fieldChanges: [
          { field: "status", oldValue: "pending", newValue: "in-progress" },
        ],
      },
    ],
    ...overrides,
  };
}

function createWorklogEntry(overrides: Partial<FeedEntry> = {}): FeedEntry {
  return {
    source: "worklog",
    timestamp: "2026-03-01T09:00:00Z",
    message: "Started implementation of search",
    taskID: "015",
    ...overrides,
  };
}

describe("FeedView", () => {
  it("renders git entries with author and message", () => {
    renderWithProviders(<FeedView entries={[createGitEntry()]} />);
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("feat: add new feature")).toBeInTheDocument();
  });

  it("renders worklog entries with task ID link", () => {
    renderWithProviders(<FeedView entries={[createWorklogEntry()]} />);
    expect(screen.getByText("Started implementation of search")).toBeInTheDocument();
    expect(screen.getByText("015")).toBeInTheDocument();
  });

  it("renders field change badges", () => {
    renderWithProviders(<FeedView entries={[createGitEntry()]} />);
    expect(screen.getByText("status")).toBeInTheDocument();
    expect(screen.getByText("pending")).toBeInTheDocument();
    expect(screen.getByText("in-progress")).toBeInTheDocument();
  });

  it("renders file status badges", () => {
    renderWithProviders(<FeedView entries={[createGitEntry()]} />);
    expect(screen.getByText("Modified")).toBeInTheDocument();
  });

  it("renders completed task status", () => {
    const entry = createGitEntry({
      files: [
        {
          path: "tasks/cli/042-add-auth.md",
          status: "modified",
          taskID: "042",
          taskStatus: "completed",
        },
      ],
    });
    renderWithProviders(<FeedView entries={[entry]} />);
    expect(screen.getByText("Completed")).toBeInTheDocument();
  });

  it("renders subtask changes", () => {
    const entry = createGitEntry({
      files: [
        {
          path: "tasks/cli/042-add-auth.md",
          status: "modified",
          taskID: "042",
          subtaskChanges: [{ text: "Add tests", done: true }],
        },
      ],
    });
    renderWithProviders(<FeedView entries={[entry]} />);
    expect(screen.getByText("Add tests")).toBeInTheDocument();
  });

  it("renders mixed git and worklog entries", () => {
    renderWithProviders(
      <FeedView entries={[createGitEntry(), createWorklogEntry()]} />,
    );
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Started implementation of search")).toBeInTheDocument();
  });

  it("renders task ID links with correct href", () => {
    renderWithProviders(<FeedView entries={[createWorklogEntry()]} />);
    const link = screen.getByText("015");
    expect(link.closest("a")).toHaveAttribute("href", "/tasks/015");
  });

  it("renders created file status as Added", () => {
    const entry = createGitEntry({
      files: [
        {
          path: "tasks/cli/043-new.md",
          status: "created",
          taskID: "043",
        },
      ],
    });
    renderWithProviders(<FeedView entries={[entry]} />);
    expect(screen.getByText("Added")).toBeInTheDocument();
  });
});
