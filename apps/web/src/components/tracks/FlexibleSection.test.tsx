import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { FlexibleSection } from "./FlexibleSection.tsx";
import { createTrackTask } from "../../test-utils/fixtures.ts";

vi.mock("./TrackCard.tsx", () => ({
  TrackCard: ({ task }: { task: { title: string } }) => (
    <div data-testid="track-card">{task.title}</div>
  ),
}));

function renderSection(tasks: Parameters<typeof FlexibleSection>[0]["tasks"]) {
  return render(
    <MemoryRouter>
      <FlexibleSection tasks={tasks} />
    </MemoryRouter>,
  );
}

describe("FlexibleSection", () => {
  it("returns null when tasks is empty", () => {
    const { container } = renderSection([]);
    expect(container.firstChild).toBeNull();
  });

  it("renders heading with task count", () => {
    renderSection([createTrackTask(), createTrackTask()]);
    expect(screen.getByText("Flexible")).toBeInTheDocument();
    expect(screen.getByText("(2)")).toBeInTheDocument();
  });

  it("renders track cards for each task", () => {
    const tasks = [
      createTrackTask({ title: "Flex A" }),
      createTrackTask({ title: "Flex B" }),
    ];
    renderSection(tasks);
    expect(screen.getByText("Flex A")).toBeInTheDocument();
    expect(screen.getByText("Flex B")).toBeInTheDocument();
  });
});
