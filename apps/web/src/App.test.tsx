import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

vi.mock("./hooks/use-live-reload.ts", () => ({ useLiveReload: vi.fn() }));
vi.mock("./components/layout/Shell.tsx", () => ({
  Shell: ({ children }: { children: React.ReactNode }) => <div data-testid="shell">{children}</div>,
}));
vi.mock("./pages/TasksPage.tsx", () => ({ TasksPage: () => <div data-testid="tasks-page" /> }));
vi.mock("./pages/TaskDetailPage.tsx", () => ({ TaskDetailPage: () => <div data-testid="task-detail-page" /> }));
vi.mock("./pages/BoardPage.tsx", () => ({ BoardPage: () => <div data-testid="board-page" /> }));
vi.mock("./pages/GraphPage.tsx", () => ({ GraphPage: () => <div data-testid="graph-page" /> }));
vi.mock("./pages/NextPage.tsx", () => ({ NextPage: () => <div data-testid="next-page" /> }));
vi.mock("./pages/TracksPage.tsx", () => ({ TracksPage: () => <div data-testid="tracks-page" /> }));
vi.mock("./pages/StatsPage.tsx", () => ({ StatsPage: () => <div data-testid="stats-page" /> }));
vi.mock("./pages/ValidatePage.tsx", () => ({ ValidatePage: () => <div data-testid="validate-page" /> }));
vi.mock("./pages/PhasesPage.tsx", () => ({ PhasesPage: () => <div data-testid="phases-page" /> }));

import App from "./App.tsx";

describe("App", () => {
  it("renders shell wrapper", () => {
    render(
      <MemoryRouter initialEntries={["/tasks"]}>
        <App />
      </MemoryRouter>,
    );
    expect(screen.getByTestId("shell")).toBeInTheDocument();
  });

  it("renders tasks page at /tasks", () => {
    render(
      <MemoryRouter initialEntries={["/tasks"]}>
        <App />
      </MemoryRouter>,
    );
    expect(screen.getByTestId("tasks-page")).toBeInTheDocument();
  });

  it("redirects / to /tasks", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <App />
      </MemoryRouter>,
    );
    expect(screen.getByTestId("tasks-page")).toBeInTheDocument();
  });

  it("renders board page at /board", () => {
    render(
      <MemoryRouter initialEntries={["/board"]}>
        <App />
      </MemoryRouter>,
    );
    expect(screen.getByTestId("board-page")).toBeInTheDocument();
  });

  it("renders stats page at /stats", () => {
    render(
      <MemoryRouter initialEntries={["/stats"]}>
        <App />
      </MemoryRouter>,
    );
    expect(screen.getByTestId("stats-page")).toBeInTheDocument();
  });
});
