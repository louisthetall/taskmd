import useSWR from "swr";
import { fetcher } from "../api/client.ts";
import type { FeedEntry } from "../api/types.ts";

interface UseFeedOptions {
  source?: string;
  since?: string;
  limit?: number;
  scope?: string;
  project?: string | null;
}

export function useFeed(options: UseFeedOptions = {}) {
  const params = new URLSearchParams();
  if (options.source && options.source !== "all")
    params.set("source", options.source);
  if (options.since) params.set("since", options.since);
  if (options.limit) params.set("limit", String(options.limit));
  if (options.scope) params.set("scope", options.scope);
  if (options.project) params.set("project", options.project);
  const qs = params.toString();
  return useSWR<FeedEntry[]>(`/api/feed${qs ? `?${qs}` : ""}`, fetcher);
}
