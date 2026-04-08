import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { API_ERROR_TOAST_EVENT } from "@/api/constants";
import { cn } from "@/lib/utils";
import { useTheme } from "@/theme/ThemeContext";

type ToastType = "error" | "info";

interface ToastItem {
  id: number;
  message: string;
  type: ToastType;
}

const ToastContext = createContext<(msg: string, type?: ToastType) => void>(
  () => {}
);

export function ToastProvider({ children }: { children: ReactNode }) {
  const { mode } = useTheme();
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const showToast = useCallback((message: string, type: ToastType = "error") => {
    const id = Date.now();
    setToasts((prev) => [...prev, { id, message, type }]);
    window.setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 5000);
  }, []);

  useEffect(() => {
    const handler = (e: Event) => {
      const ce = e as CustomEvent<{ message?: string }>;
      const msg = ce.detail?.message;
      if (msg) showToast(msg, "error");
    };
    window.addEventListener(API_ERROR_TOAST_EVENT, handler);
    return () => window.removeEventListener(API_ERROR_TOAST_EVENT, handler);
  }, [showToast]);

  return (
    <ToastContext.Provider value={showToast}>
      {children}
      <div
        className="fixed bottom-24 right-4 z-[100] flex flex-col gap-2 max-w-sm pointer-events-none"
        aria-live="polite"
      >
        {toasts.map((t) => (
          <div
            key={t.id}
            className={cn(
              "pointer-events-auto rounded-xl px-4 py-3 text-sm font-medium shadow-lg border backdrop-blur-md",
              t.type === "error" &&
                (mode === "light"
                  ? "bg-red-100 text-red-900 border-red-300"
                  : "bg-red-950/90 text-red-100 border-red-500/40"),
              t.type === "info" &&
                (mode === "light"
                  ? "bg-white text-slate-800 border-slate-200"
                  : "bg-slate-900/90 text-slate-100 border-white/10")
            )}
          >
            {t.message}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  return useContext(ToastContext);
}
