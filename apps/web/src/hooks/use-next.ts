import useSWR from "swr";
import { fetcher } from "../api/client.ts";
import type { Recommendation } from "../api/types.ts";

export function useNext(limit: number = 5, group?: string) {
  const params = new URLSearchParams();
  params.set("limit", String(limit));
  if (group) {
    params.set("filter", `group=${group}`);
  }
  return useSWR<Recommendation[]>(`/api/next?${params}`, fetcher);
}
