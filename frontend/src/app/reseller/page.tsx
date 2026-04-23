"use client";

import { AuthGuard } from "@/components/layout/AuthGuard";
import { ResellerPortal } from "@/modules/reseller/ResellerPortal";

export default function ResellerPage() {
  return (
    <div className="flex flex-col h-full">
      <AuthGuard requiredPortal="reseller">
        {() => <ResellerPortal />}
      </AuthGuard>
    </div>
  );
}
