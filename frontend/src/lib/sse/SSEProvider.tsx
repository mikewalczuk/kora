import { createContext, useContext, useEffect, useRef } from "react";
import type { ServerEvent } from "@/lib/api/generated/koraAPI.schemas";

type EventHandler = (event: ServerEvent) => void;

interface SSEContextValue {
  on: (handler: EventHandler) => () => void;
}

const SSEContext = createContext<SSEContextValue | null>(null);

const BASE_URL = import.meta.env.VITE_API_URL ?? "/api";
const SSE_URL = `${BASE_URL}/events`;
const RETRY_DELAY_MS = 3000;

export function SSEProvider({ children }: { children: React.ReactNode }) {
  const handlersRef = useRef<Set<EventHandler>>(new Set());

  useEffect(() => {
    let es: EventSource;
    let retryTimeout: ReturnType<typeof setTimeout>;

    function connect() {
      es = new EventSource(SSE_URL, { withCredentials: true });

      es.addEventListener("practice_ready", (e) => {
        try {
          const payload = JSON.parse((e as MessageEvent).data);
          const event: ServerEvent = { type: "practice_ready", ...payload };
          handlersRef.current.forEach((h) => h(event));
        } catch {}
      });

      es.onerror = () => {
        es.close();
        retryTimeout = setTimeout(connect, RETRY_DELAY_MS);
      };
    }

    connect();

    return () => {
      clearTimeout(retryTimeout);
      es?.close();
    };
  }, []);

  const on = (handler: EventHandler) => {
    handlersRef.current.add(handler);
    return () => {
      handlersRef.current.delete(handler);
    };
  };

  return <SSEContext.Provider value={{ on }}>{children}</SSEContext.Provider>;
}

export function useSSEEvents() {
  const ctx = useContext(SSEContext);
  if (!ctx) throw new Error("useSSEEvents must be used inside SSEProvider");
  return ctx;
}
