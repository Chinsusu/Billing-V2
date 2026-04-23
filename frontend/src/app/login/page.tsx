"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { login, getSession } from "@/lib/auth/mockAuth";
import { Eye, EyeOff } from "lucide-react";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPass, setShowPass] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const s = getSession();
    if (s) router.replace(`/${s.portal}`);
  }, [router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!email.trim()) { setError("Email is required"); return; }
    if (!password) { setError("Password is required"); return; }
    setLoading(true);
    await new Promise((r) => setTimeout(r, 600));
    const session = login(email.trim(), password);
    setLoading(false);
    if (!session) { setError("Invalid email or password"); return; }
    router.replace(`/${session.portal}`);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#F5F6F7] px-4">
      <div className="w-full max-w-[360px]">
        {/* Logo */}
        <div className="flex flex-col items-center mb-8">
          <div className="w-10 h-10 bg-[#D50C2D] rounded flex items-center justify-center mb-3">
            <span className="text-white text-[16px] font-bold">H</span>
          </div>
          <h1 className="text-[18px] font-semibold text-gray-900">HANetwork Billing</h1>
          <p className="text-[12px] text-gray-400 mt-1">Sign in to continue</p>
        </div>

        <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-lg p-6 flex flex-col gap-4 shadow-sm">
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 bg-white"
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Password</label>
            <div className="relative">
              <input
                type={showPass ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                autoComplete="current-password"
                className="h-9 w-full px-3 pr-9 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 bg-white"
              />
              <button
                type="button"
                onClick={() => setShowPass((v) => !v)}
                className="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 border-0 bg-transparent cursor-pointer"
              >
                {showPass ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
            </div>
          </div>

          {error && (
            <p className="text-[12px] text-red-600 bg-red-50 border border-red-100 rounded px-3 py-2">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="h-9 w-full text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors disabled:opacity-60 disabled:cursor-not-allowed mt-1"
          >
            {loading ? "Signing in…" : "Sign in"}
          </button>
        </form>

        {/* Demo hint */}
        <div className="mt-4 bg-white border border-gray-200 rounded-lg p-4 text-[11px] text-gray-400">
          <p className="font-medium text-gray-500 mb-1.5">Demo accounts</p>
          <div className="flex flex-col gap-1">
            <span>admin@hanetwork.vn · admin123</span>
            <span>reseller@proxyvn.io · reseller123</span>
            <span>client@example.com · client123</span>
          </div>
        </div>
      </div>
    </div>
  );
}
