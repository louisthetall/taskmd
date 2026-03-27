import { describe, it, expect, vi, beforeEach } from "vitest";
import { fetcher, updateTask, ApiRequestError } from "./client.ts";

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

beforeEach(() => {
  mockFetch.mockReset();
});

describe("fetcher", () => {
  it("returns parsed JSON on success", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ id: "1", title: "Test" }),
    });
    const result = await fetcher("/api/tasks");
    expect(result).toEqual({ id: "1", title: "Test" });
    expect(mockFetch).toHaveBeenCalledWith("/api/tasks");
  });

  it("throws on non-200 status", async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 404,
      statusText: "Not Found",
      json: () => Promise.reject(new Error("no body")),
    });
    await expect(fetcher("/api/tasks/999")).rejects.toThrow(
      "API error: 404 Not Found",
    );
  });
});

describe("ApiRequestError", () => {
  it("has name, message, and details", () => {
    const err = new ApiRequestError("bad input", ["field is required"]);
    expect(err.name).toBe("ApiRequestError");
    expect(err.message).toBe("bad input");
    expect(err.details).toEqual(["field is required"]);
  });

  it("defaults details to empty array", () => {
    const err = new ApiRequestError("fail");
    expect(err.details).toEqual([]);
  });
});

describe("updateTask", () => {
  it("sends PUT with JSON body and returns parsed response", async () => {
    const task = { id: "42", title: "Updated" };
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(task),
    });

    const result = await updateTask("42", { title: "Updated" });
    expect(result).toEqual(task);
    expect(mockFetch).toHaveBeenCalledWith("/api/tasks/42", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title: "Updated" }),
    });
  });

  it("throws ApiRequestError with details on error response", async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: () =>
        Promise.resolve({
          error: "validation failed",
          details: ["title is required"],
        }),
    });

    try {
      await updateTask("42", {});
      expect.fail("should have thrown");
    } catch (err) {
      expect(err).toBeInstanceOf(ApiRequestError);
      expect((err as ApiRequestError).message).toBe("validation failed");
      expect((err as ApiRequestError).details).toEqual(["title is required"]);
    }
  });

  it("throws ApiRequestError with fallback when response body is not JSON", async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 500,
      json: () => Promise.reject(new Error("not json")),
    });

    try {
      await updateTask("42", {});
      expect.fail("should have thrown");
    } catch (err) {
      expect(err).toBeInstanceOf(ApiRequestError);
      expect((err as ApiRequestError).message).toBe("HTTP 500");
    }
  });
});
