import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";
import { useSSEEvents } from "./SSEProvider";
import type { ServerEvent } from "@/lib/api/generated/koraAPI.schemas";

interface Toast {
  id: string;
  practiceId: string;
}

export function Toaster() {
  const { on } = useSSEEvents();
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    return on((event: ServerEvent) => {
      if (event.type !== "practice_ready") return;
      const toast: Toast = {
        id: crypto.randomUUID(),
        practiceId: event.practiceId,
      };
      setToasts((prev) => [...prev, toast]);
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== toast.id));
      }, 5000);
    });
  }, [on]);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 flex flex-col gap-2 z-50">
      {toasts.map((t) => (
        <Link
          key={t.id}
          to="/practice/$practiceId"
          params={{ practiceId: t.practiceId }}
          className="flex items-center justify-between gap-4 bg-gray-900 text-white text-sm px-4 py-3 rounded-lg shadow-lg hover:bg-gray-700 transition-colors"
        >
          <span>Practice ready</span>
          <span className="text-violet-400 text-xs">View →</span>
        </Link>
      ))}
    </div>
  );
}
