"use client";

import { FormEvent, useState } from "react";
import { AdminPortal } from "@/modules/admin/AdminPortal";
import { ResellerPortal } from "@/modules/reseller/ResellerPortal";
import { ClientPortal } from "@/modules/client/ClientPortal";
import { AuthSession, TwoFactorSetup, login, logout, setupTwoFactor, verifyTwoFactor } from "@/lib/api/auth";
import { BillingApiError } from "@/lib/api/client";
import { demoPortalModeEnabled } from "@/lib/api/config";

type Portal = "admin" | "reseller" | "client";

const PORTALS: { id: Portal; label: string }[] = [
  { id: "admin", label: "Admin · HANetwork" },
  { id: "reseller", label: "Reseller · ProxyVN" },
  { id: "client", label: "Client · Linh Tran" },
];

function portalsForSession(session: AuthSession): Portal[] {
  if (session.actor_type === "platform_admin" || session.actor_type === "platform_staff") {
    return ["admin"];
  }
  if (session.actor_type === "reseller_owner" || session.actor_type === "reseller_staff") {
    return ["reseller"];
  }
  return ["client"];
}

function firstPortalForSession(session: AuthSession): Portal {
  return portalsForSession(session)[0];
}

function portalLabel(portal: Portal): string {
  return PORTALS.find((candidate) => candidate.id === portal)?.label ?? portal;
}

export function AuthenticatedApp() {
  const [session, setSession] = useState<AuthSession | null>(null);
  const [portal, setPortal] = useState<Portal>("admin");
  const [signingOut, setSigningOut] = useState(false);
  const demoMode = demoPortalModeEnabled();

  function handleAuthenticated(nextSession: AuthSession) {
    setSession(nextSession);
    setPortal(firstPortalForSession(nextSession));
  }

  if (demoMode) {
    return <PortalShell portal={portal} portals={PORTALS.map((item) => item.id)} modeLabel="Demo portal mode" onSelectPortal={setPortal} />;
  }

  if (!session) {
    return <LoginPanel onAuthenticated={handleAuthenticated} />;
  }

  if (session.two_factor_required && !session.two_factor_satisfied) {
    return (
      <TwoFactorPanel
        session={session}
        onVerified={(verified) => setSession({ ...session, ...verified, two_factor_satisfied: true, two_factor_setup_required: false })}
        onSignOut={() => handleSignOut(setSession, setSigningOut)}
        signingOut={signingOut}
      />
    );
  }

  return (
    <PortalShell
      portal={portal}
      portals={portalsForSession(session)}
      modeLabel={`Signed in · ${portalLabel(firstPortalForSession(session))}`}
      onSelectPortal={setPortal}
      onSignOut={() => handleSignOut(setSession, setSigningOut)}
      signingOut={signingOut}
    />
  );
}

function PortalShell({
  portal,
  portals,
  modeLabel,
  onSelectPortal,
  onSignOut,
  signingOut,
}: {
  portal: Portal;
  portals: Portal[];
  modeLabel: string;
  onSelectPortal: (portal: Portal) => void;
  onSignOut?: () => void;
  signingOut?: boolean;
}) {
  return (
    <div className="flex flex-col h-full">
      <div className="shrink-0 h-8 bg-[#0E1116] flex items-center gap-0.5 p-4 z-50">
        <span className="text-[11px] text-gray-500 mr-2">Portal</span>
        {PORTALS.filter((candidate) => portals.includes(candidate.id)).map((p) => (
          <button
            key={p.id}
            onClick={() => onSelectPortal(p.id)}
            className={`h-[22px] p-4.5 text-[11px] font-medium rounded-[3px] border-0 cursor-pointer transition-colors
              ${portal === p.id
                ? "bg-[#D50C2D] text-white"
                : "bg-transparent text-gray-400 hover:bg-[#1F2937] hover:text-gray-200"
              }`}
          >
            {p.label}
          </button>
        ))}
        <span className="ml-auto text-[10px] text-gray-600">{modeLabel}</span>
        {onSignOut && (
          <button
            onClick={onSignOut}
            disabled={signingOut}
            className="ml-3 h-[22px] rounded-[3px] border border-gray-700 px-2 text-[10px] font-medium text-gray-300 hover:border-gray-500 hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
          >
            {signingOut ? "Signing out" : "Sign out"}
          </button>
        )}
      </div>

      <div className="flex-1 min-h-0">
        {portal === "admin" && <AdminPortal />}
        {portal === "reseller" && <ResellerPortal />}
        {portal === "client" && <ClientPortal />}
      </div>
    </div>
  );
}

function LoginPanel({ onAuthenticated }: { onAuthenticated: (session: AuthSession) => void }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      onAuthenticated(await login(email, password));
    } catch (err) {
      setError(errorMessage(err, "Login failed."));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-full bg-[#0B0E13] text-white">
      <div className="mx-auto flex min-h-full max-w-[1180px] items-center px-6 py-12">
        <div className="grid w-full gap-8 lg:grid-cols-[1.1fr_0.9fr]">
          <section className="rounded-[28px] border border-white/10 bg-[radial-gradient(circle_at_20%_20%,rgba(213,12,45,0.24),transparent_32%),linear-gradient(135deg,#151A22,#07090D)] p-8 shadow-2xl">
            <p className="text-xs font-semibold uppercase tracking-[0.36em] text-[#F05870]">Billing Control</p>
            <h1 className="mt-8 max-w-[720px] text-5xl font-semibold leading-[0.95] tracking-[-0.05em] text-white">
              Session-gated billing operations for production use.
            </h1>
            <p className="mt-6 max-w-[560px] text-sm leading-6 text-gray-300">
              Use backend auth sessions for admin, reseller, and client portals. Demo switching is disabled unless explicitly configured.
            </p>
          </section>

          <form onSubmit={handleSubmit} className="rounded-[24px] border border-white/10 bg-white p-6 text-[#10131A] shadow-2xl">
            <div>
              <p className="text-[11px] font-bold uppercase tracking-[0.28em] text-[#D50C2D]">Secure login</p>
              <h2 className="mt-3 text-2xl font-semibold tracking-[-0.03em]">Sign in to continue</h2>
              <p className="mt-2 text-sm text-gray-500">Your portal is selected from the authenticated actor type.</p>
            </div>

            <label className="mt-8 block text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">
              Email
              <input
                type="email"
                autoComplete="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                className="mt-2 w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm font-medium text-gray-900 outline-none transition focus:border-[#D50C2D] focus:bg-white"
                required
              />
            </label>

            <label className="mt-5 block text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">
              Password
              <input
                type="password"
                autoComplete="current-password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                className="mt-2 w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm font-medium text-gray-900 outline-none transition focus:border-[#D50C2D] focus:bg-white"
                required
              />
            </label>

            {error && <p className="mt-4 rounded-xl bg-red-50 px-4 py-3 text-sm font-medium text-red-700">{error}</p>}

            <button
              type="submit"
              disabled={submitting}
              className="mt-6 w-full rounded-xl bg-[#D50C2D] px-4 py-3 text-sm font-bold text-white shadow-lg shadow-red-900/20 transition hover:bg-[#B3082A] disabled:cursor-not-allowed disabled:opacity-60"
            >
              {submitting ? "Signing in" : "Sign in"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

function TwoFactorPanel({
  session,
  onVerified,
  onSignOut,
  signingOut,
}: {
  session: AuthSession;
  onVerified: (verified: Partial<AuthSession>) => void;
  onSignOut: () => void;
  signingOut: boolean;
}) {
  const [setup, setSetup] = useState<TwoFactorSetup | null>(null);
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [loadingSetup, setLoadingSetup] = useState(false);
  const [verifying, setVerifying] = useState(false);

  async function handleSetup() {
    setError("");
    setLoadingSetup(true);
    try {
      setSetup(await setupTwoFactor());
    } catch (err) {
      setError(errorMessage(err, "Unable to start two-factor setup."));
    } finally {
      setLoadingSetup(false);
    }
  }

  async function handleVerify(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setVerifying(true);
    try {
      onVerified(await verifyTwoFactor(code));
    } catch (err) {
      setError(errorMessage(err, "Two-factor verification failed."));
    } finally {
      setVerifying(false);
    }
  }

  return (
    <div className="grid min-h-full place-items-center bg-[#0B0E13] px-6 py-12 text-white">
      <div className="w-full max-w-[520px] rounded-[24px] border border-white/10 bg-white p-6 text-[#10131A] shadow-2xl">
        <p className="text-[11px] font-bold uppercase tracking-[0.28em] text-[#D50C2D]">Two-factor required</p>
        <h1 className="mt-3 text-2xl font-semibold tracking-[-0.03em]">Verify this session</h1>
        <p className="mt-2 text-sm text-gray-500">Admin and privileged sessions must complete TOTP before accessing protected screens.</p>

        {session.two_factor_setup_required && (
          <div className="mt-6 rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <button
              type="button"
              onClick={handleSetup}
              disabled={loadingSetup || Boolean(setup)}
              className="rounded-xl bg-[#10131A] px-4 py-2 text-sm font-bold text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loadingSetup ? "Preparing setup" : setup ? "Setup prepared" : "Prepare TOTP setup"}
            </button>
            {setup && (
              <div className="mt-4 space-y-3 text-sm">
                <p className="font-medium text-gray-700">Add this setup key to your authenticator app, then enter the generated code.</p>
                <code className="block rounded-xl bg-white p-3 text-xs font-semibold text-gray-800 break-all">{setup.secret}</code> {/* sensitive-text-allowlist: user-facing TOTP setup value. */}
                <code className="block rounded-xl bg-white p-3 text-xs text-gray-500 break-all">{setup.provision_uri}</code>
              </div>
            )}
          </div>
        )}

        <form onSubmit={handleVerify} className="mt-6">
          <label className="block text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">
            Authentication code
            <input
              inputMode="numeric"
              autoComplete="one-time-code"
              value={code}
              onChange={(event) => setCode(event.target.value)}
              className="mt-2 w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm font-medium text-gray-900 outline-none transition focus:border-[#D50C2D] focus:bg-white"
              required
            />
          </label>

          {error && <p className="mt-4 rounded-xl bg-red-50 px-4 py-3 text-sm font-medium text-red-700">{error}</p>}

          <button
            type="submit"
            disabled={verifying}
            className="mt-6 w-full rounded-xl bg-[#D50C2D] px-4 py-3 text-sm font-bold text-white shadow-lg shadow-red-900/20 transition hover:bg-[#B3082A] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {verifying ? "Verifying" : "Verify and continue"}
          </button>
        </form>

        <button
          type="button"
          onClick={onSignOut}
          disabled={signingOut}
          className="mt-4 w-full rounded-xl border border-gray-200 px-4 py-3 text-sm font-bold text-gray-700 transition hover:border-gray-400 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {signingOut ? "Signing out" : "Sign out"}
        </button>
      </div>
    </div>
  );
}

async function handleSignOut(setSession: (session: AuthSession | null) => void, setSigningOut: (value: boolean) => void) {
  setSigningOut(true);
  try {
    await logout();
  } finally {
    setSession(null);
    setSigningOut(false);
  }
}

function errorMessage(err: unknown, fallback: string): string {
  if (err instanceof BillingApiError) {
    return err.message || fallback;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return fallback;
}
