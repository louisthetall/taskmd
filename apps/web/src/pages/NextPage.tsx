import { useSearchParams } from "react-router-dom";
import { useNext } from "../hooks/use-next.ts";
import { NextView } from "../components/next/NextView.tsx";
import { LoadingState } from "../components/shared/LoadingState.tsx";
import { ErrorState } from "../components/shared/ErrorState.tsx";

export function NextPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const limit = Number(searchParams.get("limit")) || 5;
  const group = searchParams.get("group") ?? "";
  const { data, error, isLoading, mutate } = useNext(limit, group || undefined);

  function setLimit(n: number) {
    setSearchParams(
      (prev) => {
        prev.set("limit", String(n));
        return prev;
      },
      { replace: true },
    );
  }

  function setGroup(g: string) {
    setSearchParams(
      (prev) => {
        if (g) {
          prev.set("group", g);
        } else {
          prev.delete("group");
        }
        return prev;
      },
      { replace: true },
    );
  }

  if (isLoading) return <LoadingState variant="cards" />;
  if (error) return <ErrorState error={error} onRetry={() => mutate()} />;
  if (!data) return null;

  if (data.length === 0) {
    return (
      <div>
        <NextView
          recommendations={[]}
          limit={limit}
          onLimitChange={setLimit}
          group={group}
          onGroupChange={setGroup}
        />
        <p className="text-sm text-gray-500 py-8 text-center">
          No actionable tasks found. All tasks are either completed or blocked.
        </p>
      </div>
    );
  }

  return (
    <NextView
      recommendations={data}
      limit={limit}
      onLimitChange={setLimit}
      group={group}
      onGroupChange={setGroup}
    />
  );
}
