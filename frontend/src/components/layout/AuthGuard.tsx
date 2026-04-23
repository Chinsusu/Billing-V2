"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getSession, type Session } from "@/lib/auth/mockAuth";

interface AuthGuardProps {
  requiredPortal: "admin" | "reseller" | "client";
  children: (session: Session) => React.ReactNode;
}

export function AuthGuard({ requiredPortal, children }: AuthGuardProps) {
  const router = useRouter();
  const [session, setSession] = useState<Session | null | "loading">("loading");

  useEffect(() => {
    const s = getSession();
    if (!s) { router.replace("/login"); return; }
    if (s.portal !== requiredPortal) { router.replace(`/${s.portal}`); return; }
    setSession(s);
  }, [router, requiredPortal]);

  if (session === "loading") {
    return (
      <div className="flex h-full items-center justify-center bg-[#F5F6F7]">
        <div className="w-6 h-6 border-2 border-[#D50C2D] border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  if (!session) return null;
  return <>{children(session)}</>;
}
