import { useCallback } from "react";
import { useSearchParams } from "react-router-dom";

const PHASE_PARAM = "phase";

export function usePhase() {
  const [searchParams, setSearchParams] = useSearchParams();
  const phase = searchParams.get(PHASE_PARAM);

  const setPhase = useCallback(
    (next: string | null) => {
      setSearchParams(
        (prev) => {
          const updated = new URLSearchParams(prev);
          if (next) {
            updated.set(PHASE_PARAM, next);
          } else {
            updated.delete(PHASE_PARAM);
          }
          return updated;
        },
        { replace: false },
      );
    },
    [setSearchParams],
  );

  return { phase, setPhase };
}
