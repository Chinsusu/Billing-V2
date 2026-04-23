"use client";

import { AuthGuard } from "@/components/layout/AuthGuard";
import { ClientPortal } from "@/modules/client/ClientPortal";

export default function ClientPage() {
  return (
    <div className="flex flex-col h-full">
      <AuthGuard requiredPortal="client">
        {() => <ClientPortal />}
      </AuthGuard>
    </div>
  );
}
