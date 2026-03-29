import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useFeed } from "./use-feed.ts";

vi.mock("swr", () => ({
  default: vi.fn((key: string) => ({ data: undefined, error: undefined, isLoading: false, mutate: vi.fn(), key })),
}));

import useSWR from "swr";
const mockUseSWR = vi.mocked(useSWR);

describe("useFeed", () => {
  it("calls /api/feed with no params by default", () => {
    renderHook(() => useFeed());
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed", expect.any(Function));
  });

  it("includes source param when not 'all'", () => {
    renderHook(() => useFeed({ source: "git" }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed?source=git", expect.any(Function));
  });

  it("omits source param when 'all'", () => {
    renderHook(() => useFeed({ source: "all" }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed", expect.any(Function));
  });

  it("includes since param", () => {
    renderHook(() => useFeed({ since: "7d" }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed?since=7d", expect.any(Function));
  });

  it("includes limit param", () => {
    renderHook(() => useFeed({ limit: 10 }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed?limit=10", expect.any(Function));
  });

  it("includes scope param", () => {
    renderHook(() => useFeed({ scope: "cli" }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed?scope=cli", expect.any(Function));
  });

  it("includes project param", () => {
    renderHook(() => useFeed({ project: "proj" }));
    expect(mockUseSWR).toHaveBeenCalledWith("/api/feed?project=proj", expect.any(Function));
  });

  it("combines multiple params", () => {
    renderHook(() => useFeed({ source: "worklog", since: "1d", limit: 5, scope: "web" }));
    expect(mockUseSWR).toHaveBeenCalledWith(
      "/api/feed?source=worklog&since=1d&limit=5&scope=web",
      expect.any(Function),
    );
  });
});
