import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { TagAutocomplete } from "./TagAutocomplete.tsx";

// jsdom does not implement scrollIntoView
Element.prototype.scrollIntoView = vi.fn();

function defaultProps(overrides: Partial<React.ComponentProps<typeof TagAutocomplete>> = {}) {
  return {
    availableTags: ["backend", "frontend", "api", "database"],
    selectedTags: new Set<string>(),
    onTagsChange: vi.fn(),
    ...overrides,
  };
}

describe("TagAutocomplete", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("renders Tags label and input", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    expect(screen.getByText("Tags:")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Add tag...")).toBeInTheDocument();
  });

  it("shows selected tags as chips", () => {
    render(<TagAutocomplete {...defaultProps({ selectedTags: new Set(["backend", "api"]) })} />);
    expect(screen.getByText("backend")).toBeInTheDocument();
    expect(screen.getByText("api")).toBeInTheDocument();
  });

  it("shows suggestions dropdown on input focus", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    expect(screen.getByRole("listbox")).toBeInTheDocument();
    expect(screen.getAllByRole("option")).toHaveLength(4);
  });

  it("filters suggestions by query", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    fireEvent.change(input, { target: { value: "back" } });
    const options = screen.getAllByRole("option");
    expect(options).toHaveLength(1);
    expect(options[0]).toHaveTextContent("backend");
  });

  it("clicking a suggestion selects it", () => {
    const props = defaultProps();
    render(<TagAutocomplete {...props} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    fireEvent.click(screen.getByText("api"));
    expect(props.onTagsChange).toHaveBeenCalledOnce();
    const mockFn = props.onTagsChange as ReturnType<typeof vi.fn>;
    const result = mockFn.mock.calls[0][0] as Set<string>;
    expect(result.has("api")).toBe(true);
  });

  it("remove button on chip calls onTagsChange", () => {
    const props = defaultProps({ selectedTags: new Set(["backend"]) });
    render(<TagAutocomplete {...props} />);
    fireEvent.click(screen.getByLabelText("Remove backend filter"));
    expect(props.onTagsChange).toHaveBeenCalledOnce();
  });

  it("ArrowDown and ArrowUp navigate suggestions", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    // First option is active by default
    expect(screen.getByRole("option", { name: "backend" })).toHaveAttribute("aria-selected", "true");

    fireEvent.keyDown(input, { key: "ArrowDown" });
    expect(screen.getByRole("option", { name: "frontend" })).toHaveAttribute("aria-selected", "true");

    fireEvent.keyDown(input, { key: "ArrowUp" });
    expect(screen.getByRole("option", { name: "backend" })).toHaveAttribute("aria-selected", "true");
  });

  it("Enter selects the active suggestion", () => {
    const props = defaultProps();
    render(<TagAutocomplete {...props} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    fireEvent.keyDown(input, { key: "ArrowDown" });
    fireEvent.keyDown(input, { key: "Enter" });
    expect(props.onTagsChange).toHaveBeenCalledOnce();
    const enterMockFn = props.onTagsChange as ReturnType<typeof vi.fn>;
    const result = enterMockFn.mock.calls[0][0] as Set<string>;
    expect(result.has("frontend")).toBe(true);
  });

  it("Escape closes dropdown", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    expect(screen.getByRole("listbox")).toBeInTheDocument();
    fireEvent.keyDown(input, { key: "Escape" });
    expect(screen.queryByRole("listbox")).not.toBeInTheDocument();
  });

  it("does not show already-selected tags in suggestions", () => {
    render(<TagAutocomplete {...defaultProps({ selectedTags: new Set(["backend", "api"]) })} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    const options = screen.getAllByRole("option");
    expect(options).toHaveLength(2);
    expect(options[0]).toHaveTextContent("frontend");
    expect(options[1]).toHaveTextContent("database");
  });

  it("shows 'No matching tags' when query has no matches", () => {
    render(<TagAutocomplete {...defaultProps()} />);
    const input = screen.getByPlaceholderText("Add tag...");
    fireEvent.focus(input);
    fireEvent.change(input, { target: { value: "zzzzz" } });
    expect(screen.queryByRole("listbox")).not.toBeInTheDocument();
    expect(screen.getByText("No matching tags")).toBeInTheDocument();
  });
});
