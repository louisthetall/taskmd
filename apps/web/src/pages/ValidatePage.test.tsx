import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { ValidatePage } from "./ValidatePage.tsx";
import { createValidationResult } from "../test-utils/index.ts";

vi.mock("../hooks/use-validate.ts", () => ({
  useValidate: vi.fn(),
}));

vi.mock("../hooks/use-phase.tsx", () => ({
  usePhase: () => ({ phase: null, setPhase: vi.fn() }),
}));

vi.mock("../hooks/use-project.ts", () => ({
  useProject: () => ({ project: null, setProject: vi.fn() }),
}));

vi.mock("../components/validate/ValidateView.tsx", () => ({
  ValidateView: () => <div data-testid="validate-view">ValidateView</div>,
}));

import { useValidate } from "../hooks/use-validate.ts";
const mockUseValidate = vi.mocked(useValidate);

describe("ValidatePage", () => {
  it("renders loading state", () => {
    mockUseValidate.mockReturnValue({
      data: undefined,
      error: undefined,
      isLoading: true,
      mutate: vi.fn(),
      isValidating: false,
    });
    const { container } = render(<ValidatePage />);
    expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
  });

  it("renders error state", () => {
    mockUseValidate.mockReturnValue({
      data: undefined,
      error: new Error("Server error"),
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    render(<ValidatePage />);
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders ValidateView when data is available", () => {
    mockUseValidate.mockReturnValue({
      data: createValidationResult(),
      error: undefined,
      isLoading: false,
      mutate: vi.fn(),
      isValidating: false,
    });
    render(<ValidatePage />);
    expect(screen.getByTestId("validate-view")).toBeInTheDocument();
  });
});
