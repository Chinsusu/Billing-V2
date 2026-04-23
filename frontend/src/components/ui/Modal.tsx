"use client";

import { useEffect, useRef } from "react";
import { createPortal } from "react-dom";
import { X } from "lucide-react";

interface ModalProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: React.ReactNode;
  width?: "sm" | "md" | "lg";
}

const WIDTH = { sm: "max-w-sm", md: "max-w-md", lg: "max-w-lg" };

export function Modal({ open, onClose, title, description, children, width = "md" }: ModalProps) {
  const overlayRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;

  return createPortal(
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
      onMouseDown={(e) => { if (e.target === overlayRef.current) onClose(); }}
    >
      <div className={`bg-white rounded-lg shadow-xl w-full ${WIDTH[width]} flex flex-col max-h-[90vh]`}>
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-gray-200 shrink-0">
          <div>
            <h2 className="text-[14px] font-semibold text-gray-900 m-0">{title}</h2>
            {description && <p className="text-[12px] text-gray-500 mt-0.5">{description}</p>}
          </div>
          <button
            onClick={onClose}
            className="w-7 h-7 flex items-center justify-center rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors border-0 bg-transparent cursor-pointer"
          >
            <X size={15} />
          </button>
        </div>
        {/* Body */}
        <div className="overflow-auto flex-1 px-5 py-4">{children}</div>
      </div>
    </div>,
    document.body,
  );
}
