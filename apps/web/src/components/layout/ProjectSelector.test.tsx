import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

const mockSetProject = vi.fn();

vi.mock("../../hooks/use-projects.ts", () => ({ useProjects: vi.fn() }));
vi.mock("../../hooks/use-project.ts", () => ({
  useProject: () => ({ project: null, setProject: mockSetProject }),
}));

import { useProjects } from "../../hooks/use-projects.ts";
const mockUseProjects = vi.mocked(useProjects);

import { ProjectSelector } from "./ProjectSelector.tsx";

describe("ProjectSelector", () => {
  beforeEach(() => {
    mockSetProject.mockClear();
  });

  it("returns null when projects data is undefined", () => {
    mockUseProjects.mockReturnValue({
      data: undefined,
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    const { container } = render(<ProjectSelector />);
    expect(container.innerHTML).toBe("");
  });

  it("returns null when projects array is empty", () => {
    mockUseProjects.mockReturnValue({
      data: [],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    const { container } = render(<ProjectSelector />);
    expect(container.innerHTML).toBe("");
  });

  it("renders select with project options when projects exist", () => {
    mockUseProjects.mockReturnValue({
      data: [
        { id: "proj-1", name: "Project One" },
        { id: "proj-2", name: "Project Two" },
      ],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    render(<ProjectSelector />);
    expect(screen.getByRole("combobox")).toBeInTheDocument();
    expect(screen.getByText("Project One")).toBeInTheDocument();
    expect(screen.getByText("Project Two")).toBeInTheDocument();
  });

  it("renders (local) as default option", () => {
    mockUseProjects.mockReturnValue({
      data: [{ id: "proj-1", name: "Project One" }],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    render(<ProjectSelector />);
    expect(screen.getByText("(local)")).toBeInTheDocument();
  });

  it("calls setProject with selected value on change", () => {
    mockUseProjects.mockReturnValue({
      data: [
        { id: "proj-1", name: "Project One" },
        { id: "proj-2", name: "Project Two" },
      ],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    render(<ProjectSelector />);
    fireEvent.change(screen.getByRole("combobox"), {
      target: { value: "proj-2" },
    });
    expect(mockSetProject).toHaveBeenCalledWith("proj-2");
  });

  it("calls setProject with null when selecting empty value", () => {
    mockUseProjects.mockReturnValue({
      data: [{ id: "proj-1", name: "Project One" }],
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    } as any);
    render(<ProjectSelector />);
    fireEvent.change(screen.getByRole("combobox"), {
      target: { value: "" },
    });
    expect(mockSetProject).toHaveBeenCalledWith(null);
  });
});
