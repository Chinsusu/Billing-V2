"use client";

import { Bell, ChevronDown } from "lucide-react";

interface TopbarProps {
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  actions?: React.ReactNode;
}

export function Topbar({ title, breadcrumbs, meta, actions }: TopbarProps) {
  return (
    <header className="h-12 px-5 bg-white border-b border-gray-200 flex items-center gap-4 shrink-0">
      <div className="flex-1 min-w-0">
        {breadcrumbs.length > 0 && (
          <div className="flex items-center gap-1 text-[11px] text-gray-400 mb-0.5">
            {breadcrumbs.map((crumb, i) => (
              <span key={i} className="flex items-center gap-1">
                {i > 0 && <span className="opacity-40">›</span>}
                {crumb}
              </span>
            ))}
          </div>
        )}
        <div className="flex items-center gap-4.5">
          <h1 className="text-[15px] font-medium tracking-tight text-gray-900 m-0">{title}</h1>
          {meta}
        </div>
      </div>

      <div className="flex items-center gap-1.5 bg-gray-100 p-4 py-1 rounded-[3px] w-64 border border-transparent">
        <span className="text-gray-400 text-[12px]">⌕</span>
        <input
          placeholder="Search…"
          className="flex-1 bg-transparent border-0 outline-none text-[12px] text-gray-700 font-[inherit]"
        />
        <span className="text-[10px] text-gray-400 border border-gray-300 rounded px-1">⌘K</span>
      </div>

      <div className="flex items-center gap-1">
        {actions}
        <button className="h-8 flex items-center justify-center gap-1 text-gray-500 hover:text-gray-900 rounded-[3px] transition-colors border-0 bg-transparent cursor-pointer px-1 relative">
          <div className="relative flex items-center justify-center w-8 h-8 hover:bg-gray-100 rounded-full transition-colors">
            <Bell size={18} />
            <span className="absolute top-0 right-0 flex h-4 min-w-4 items-center justify-center rounded-full bg-[#D50C2D] px-1 text-[10px] font-medium text-white border-2 border-white shadow-sm">
              25
            </span>
          </div>
          <ChevronDown size={14} className="text-gray-400" />
        </button>
      </div>
    </header>
  );
}
