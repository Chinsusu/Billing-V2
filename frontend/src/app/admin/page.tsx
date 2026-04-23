"use client";

import { AuthGuard } from "@/components/layout/AuthGuard";
import { AdminPortal } from "@/modules/admin/AdminPortal";

export default function AdminPage() {
  return (
    <div className="flex flex-col h-full">
      <AuthGuard requiredPortal="admin">
        {() => <AdminPortal />}
      </AuthGuard>
    </div>
  );
}
