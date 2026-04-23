"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { getSession } from "@/lib/auth/mockAuth";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const s = getSession();
    if (s) router.replace(`/${s.portal}`);
    else router.replace("/login");
  }, [router]);

  return (
    <div className="flex h-full items-center justify-center bg-[#F5F6F7]">
      <div className="w-6 h-6 border-2 border-[#D50C2D] border-t-transparent rounded-full animate-spin" />
    </div>
  );
}
