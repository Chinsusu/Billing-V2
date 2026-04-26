"use client";

import { Filter, RotateCcw } from "lucide-react";
import { FormEvent, InputHTMLAttributes, ReactNode, SelectHTMLAttributes } from "react";

type AdminFilterTone = "default" | "loading" | "success" | "error";

interface AdminFilterBarProps {
  children: ReactNode;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onReset: () => void;
  statusText?: string;
  statusTone?: AdminFilterTone;
}

interface AdminFilterFieldProps {
  label: string;
  children: ReactNode;
}

interface AdminFilterInputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
}

interface AdminFilterSelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  options: Array<{ value: string; label: string }>;
}

const STATUS_CLASSES: Record<AdminFilterTone, string> = {
  default: "border-gray-200 bg-white text-gray-500",
  loading: "border-blue-200 bg-blue-50 text-blue-700",
  success: "border-emerald-200 bg-emerald-50 text-emerald-700",
  error: "border-amber-200 bg-amber-50 text-amber-700",
};

const INPUT_CLASS_NAME = [
  "h-10 w-full rounded-md border border-gray-200 bg-white px-3 text-[13px] text-gray-700 outline-none transition",
  "placeholder:text-gray-400 focus:border-[#D50C2D] focus:ring-2 focus:ring-[#D50C2D]/10",
].join(" ");

export function AdminFilterBar({
  children,
  onSubmit,
  onReset,
  statusText,
  statusTone = "default",
}: AdminFilterBarProps) {
  return (
    <form onSubmit={onSubmit} className="border-b border-gray-100 bg-gray-50/70 p-4">
      <div className="flex flex-col gap-3 xl:flex-row xl:items-end xl:justify-between">
        <div className="grid flex-1 gap-3 sm:grid-cols-2 xl:grid-cols-5">{children}</div>
        <div className="flex flex-col gap-2 sm:flex-row xl:flex-col xl:min-w-[148px]">
          <button
            type="submit"
            className="inline-flex h-10 items-center justify-center gap-2 rounded-md bg-[#D50C2D] px-4 text-[13px] font-medium text-white transition hover:bg-[#B3082A] cursor-pointer"
          >
            <Filter className="h-4 w-4" />
            Apply
          </button>
          <button
            type="button"
            onClick={onReset}
            className="inline-flex h-10 items-center justify-center gap-2 rounded-md border border-gray-200 bg-white px-4 text-[13px] font-medium text-gray-600 transition hover:bg-gray-50 cursor-pointer"
          >
            <RotateCcw className="h-4 w-4" />
            Reset
          </button>
        </div>
      </div>

      {statusText && (
        <div className={`mt-3 inline-flex items-center rounded-md border px-3 py-2 text-[12px] ${STATUS_CLASSES[statusTone]}`}>
          {statusText}
        </div>
      )}
    </form>
  );
}

export function AdminFilterField({ label, children }: AdminFilterFieldProps) {
  return (
    <label className="flex min-w-0 flex-col gap-1.5">
      <span className="text-[11px] font-medium uppercase tracking-wide text-gray-400">{label}</span>
      {children}
    </label>
  );
}

export function AdminFilterInput({ label, className = "", ...props }: AdminFilterInputProps) {
  return (
    <AdminFilterField label={label}>
      <input {...props} className={`${INPUT_CLASS_NAME} ${className}`.trim()} />
    </AdminFilterField>
  );
}

export function AdminFilterSelect({
  label,
  options,
  className = "",
  "aria-label": ariaLabel,
  ...props
}: AdminFilterSelectProps) {
  return (
    <AdminFilterField label={label}>
      <select {...props} aria-label={ariaLabel ?? label} className={`${INPUT_CLASS_NAME} ${className}`.trim()}>
        {options.map((option) => (
          <option key={option.value || "all"} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </AdminFilterField>
  );
}
