"use client";

import { createContext, useCallback, useContext, useEffect, useState, useRef } from "react";
import { createPortal } from "react-dom";
import { CheckCircle, XCircle, AlertTriangle, Info, X } from "lucide-react";

type ToastType = "success" | "error" | "warning" | "info";

interface Toast {
  id: string;
  type: ToastType;
  message: string;
}

interface ToastCtx {
  toast: (message: string, type?: ToastType) => void;
}

const Ctx = createContext<ToastCtx>({ toast: () => {} });

const ICONS = {
  success: <CheckCircle size={15} className="text-emerald-500 shrink-0" />,
  error: <XCircle size={15} className="text-red-500 shrink-0" />,
  warning: <AlertTriangle size={15} className="text-amber-500 shrink-0" />,
  info: <Info size={15} className="text-blue-500 shrink-0" />,
};

const BORDER = {
  success: "border-l-emerald-500",
  error: "border-l-red-500",
  warning: "border-l-amber-500",
  info: "border-l-blue-500",
};

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const timers = useRef<Record<string, ReturnType<typeof setTimeout>>>({});

  useEffect(() => { setMounted(true); }, []);

  const dismiss = useCallback((id: string) => {
    clearTimeout(timers.current[id]);
    delete timers.current[id];
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast = useCallback((message: string, type: ToastType = "info") => {
    const id = Math.random().toString(36).slice(2);
    setToasts((prev) => [...prev, { id, type, message }]);
    timers.current[id] = setTimeout(() => dismiss(id), 4000);
  }, [dismiss]);

  return (
    <Ctx.Provider value={{ toast }}>
      {children}
      {mounted && createPortal(
        <div className="fixed bottom-5 right-5 z-[100] flex flex-col gap-2 w-[340px]">
          {toasts.map((t) => (
            <div
              key={t.id}
              className={`flex items-start gap-3 bg-white border border-gray-200 border-l-4 ${BORDER[t.type]} rounded-md shadow-lg px-4 py-3`}
            >
              {ICONS[t.type]}
              <span className="text-[13px] text-gray-800 flex-1">{t.message}</span>
              <button
                onClick={() => dismiss(t.id)}
                className="w-5 h-5 flex items-center justify-center rounded hover:bg-gray-100 text-gray-400 border-0 bg-transparent cursor-pointer shrink-0"
              >
                <X size={12} />
              </button>
            </div>
          ))}
        </div>,
        document.body,
      )}
    </Ctx.Provider>
  );
}

export function useToast() {
  return useContext(Ctx);
}
