"use client";

import { useEffect, useRef, useState } from "react";

type ApiStatus = "idle" | "loading" | "success" | "error";

export interface ApiResource<T> {
  data: T | null;
  error: string | null;
  status: ApiStatus;
}

export function useApiResource<T>(load: () => Promise<T>, key = "default"): ApiResource<T> {
  const loadRef = useRef(load);
  const [state, setState] = useState<ApiResource<T>>({
    data: null,
    error: null,
    status: "idle",
  });

  useEffect(() => {
    loadRef.current = load;
  }, [load]);

  useEffect(() => {
    let active = true;
    setState((current) => ({ ...current, error: null, status: "loading" }));
    loadRef.current()
      .then((data) => {
        if (active) setState({ data, error: null, status: "success" });
      })
      .catch((error: unknown) => {
        if (!active) return;
        const message = error instanceof Error ? error.message : "API request failed.";
        setState({ data: null, error: message, status: "error" });
      });
    return () => {
      active = false;
    };
  }, [key]);

  return state;
}
