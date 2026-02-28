import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { DesktopNav, MobileMenu } from "./NavTabs.tsx";

const tabLabels = ["Tasks", "Next Up", "Board", "Tracks", "Graph", "Stats", "Validate"];

describe("DesktopNav", () => {
  function renderDesktopNav(onSearchOpen = vi.fn()) {
    return {
      onSearchOpen,
      ...render(
        <MemoryRouter initialEntries={["/tasks"]}>
          <DesktopNav onSearchOpen={onSearchOpen} />
        </MemoryRouter>,
      ),
    };
  }

  it("renders all 7 navigation tabs", () => {
    renderDesktopNav();
    for (const label of tabLabels) {
      expect(screen.getByRole("link", { name: label })).toBeInTheDocument();
    }
  });

  it("renders tabs with correct paths", () => {
    renderDesktopNav();
    expect(screen.getByRole("link", { name: "Tasks" })).toHaveAttribute("href", "/tasks");
    expect(screen.getByRole("link", { name: "Next Up" })).toHaveAttribute("href", "/next");
    expect(screen.getByRole("link", { name: "Board" })).toHaveAttribute("href", "/board");
    expect(screen.getByRole("link", { name: "Graph" })).toHaveAttribute("href", "/graph");
    expect(screen.getByRole("link", { name: "Stats" })).toHaveAttribute("href", "/stats");
    expect(screen.getByRole("link", { name: "Validate" })).toHaveAttribute("href", "/validate");
  });

  it("renders search button with aria-label", () => {
    renderDesktopNav();
    expect(screen.getByRole("button", { name: "Search tasks" })).toBeInTheDocument();
  });

  it("calls onSearchOpen when search button is clicked", async () => {
    const { onSearchOpen } = renderDesktopNav();
    await userEvent.click(screen.getByRole("button", { name: "Search tasks" }));
    expect(onSearchOpen).toHaveBeenCalledOnce();
  });

  it("renders Docs external link", () => {
    renderDesktopNav();
    const docsLink = screen.getByText(/Docs/);
    expect(docsLink).toHaveAttribute("target", "_blank");
    expect(docsLink).toHaveAttribute("rel", "noopener noreferrer");
  });

  it("renders GitHub external link with aria-label", () => {
    renderDesktopNav();
    const githubLink = screen.getByRole("link", { name: "GitHub repository" });
    expect(githubLink).toHaveAttribute("target", "_blank");
    expect(githubLink).toHaveAttribute("rel", "noopener noreferrer");
  });
});

describe("MobileMenu", () => {
  function renderMobileMenu() {
    return render(
      <MemoryRouter initialEntries={["/tasks"]}>
        <MobileMenu />
      </MemoryRouter>,
    );
  }

  it("renders all 7 navigation tabs", () => {
    renderMobileMenu();
    for (const label of tabLabels) {
      expect(screen.getByRole("link", { name: label })).toBeInTheDocument();
    }
  });

  it("renders Docs and GitHub external links", () => {
    renderMobileMenu();
    expect(screen.getByText(/Docs/)).toBeInTheDocument();
    expect(screen.getByText(/GitHub/)).toBeInTheDocument();
  });
});
