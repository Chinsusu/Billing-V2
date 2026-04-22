"use client";

import { useState } from "react";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";
import { Portal } from "@/lib/navigation/screens";

interface AppShellProps {
  portal: Portal;
  activeScreen: string;
  onSelectScreen: (id: string) => void;
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  actions?: React.ReactNode;
  children: React.ReactNode;
}

export function AppShell({
  portal, activeScreen, onSelectScreen,
  title, breadcrumbs, meta, actions, children,
}: AppShellProps) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-full bg-[#F5F6F7]">
      <Sidebar
        portal={portal}
        activeScreen={activeScreen}
        onSelect={onSelectScreen}
        collapsed={collapsed}
        onToggle={() => setCollapsed((c) => !c)}
      />
      <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <Topbar title={title} breadcrumbs={breadcrumbs} meta={meta} actions={actions} />
        <div className="flex-1 overflow-auto" style={{ scrollbarWidth: "none" }}>
          {children}
        </div>
      </main>
    </div>
  );
}
