import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Highlight } from "./Highlight.tsx";

describe("Highlight", () => {
  it("renders plain text when query is empty", () => {
    render(<Highlight text="Hello world" query="" />);
    expect(screen.getByText("Hello world")).toBeInTheDocument();
    expect(screen.queryByRole("mark")).not.toBeInTheDocument();
  });

  it("renders plain text when query does not match", () => {
    render(<Highlight text="Hello world" query="xyz" />);
    expect(screen.getByText("Hello world")).toBeInTheDocument();
    expect(screen.queryByRole("mark")).not.toBeInTheDocument();
  });

  it("highlights matching substring", () => {
    const { container } = render(<Highlight text="Hello world" query="world" />);
    const mark = container.querySelector("mark");
    expect(mark).not.toBeNull();
    expect(mark!.textContent).toBe("world");
  });

  it("is case insensitive", () => {
    const { container } = render(<Highlight text="Hello World" query="hello" />);
    const mark = container.querySelector("mark");
    expect(mark).not.toBeNull();
    expect(mark!.textContent).toBe("Hello");
  });

  it("preserves original casing in highlighted text", () => {
    const { container } = render(<Highlight text="FooBar" query="foobar" />);
    const mark = container.querySelector("mark");
    expect(mark!.textContent).toBe("FooBar");
  });

  it("highlights only the first occurrence", () => {
    const { container } = render(<Highlight text="test a test" query="test" />);
    const marks = container.querySelectorAll("mark");
    expect(marks).toHaveLength(1);
    expect(marks[0].textContent).toBe("test");
  });

  it("renders text before and after the match", () => {
    const { container } = render(<Highlight text="abc def ghi" query="def" />);
    expect(container.textContent).toBe("abc def ghi");
    const mark = container.querySelector("mark");
    expect(mark!.textContent).toBe("def");
  });

  it("handles match at the start of text", () => {
    const { container } = render(<Highlight text="Hello world" query="Hello" />);
    const mark = container.querySelector("mark");
    expect(mark!.textContent).toBe("Hello");
    expect(container.textContent).toBe("Hello world");
  });

  it("handles match at the end of text", () => {
    const { container } = render(<Highlight text="Hello world" query="world" />);
    const mark = container.querySelector("mark");
    expect(mark!.textContent).toBe("world");
    expect(container.textContent).toBe("Hello world");
  });

  it("handles full text match", () => {
    const { container } = render(<Highlight text="exact" query="exact" />);
    const mark = container.querySelector("mark");
    expect(mark!.textContent).toBe("exact");
  });
});
