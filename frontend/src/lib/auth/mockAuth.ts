export type Portal = "admin" | "reseller" | "client";

interface MockUser {
  email: string;
  password: string;
  portal: Portal;
  name: string;
  tenant: string;
}

const MOCK_USERS: MockUser[] = [
  { email: "admin@hanetwork.vn", password: "admin123", portal: "admin", name: "HA Admin", tenant: "HANetwork" },
  { email: "reseller@proxyvn.io", password: "reseller123", portal: "reseller", name: "ProxyVN Admin", tenant: "ProxyVN" },
  { email: "client@example.com", password: "client123", portal: "client", name: "Linh Tran", tenant: "ProxyVN" },
];

const SESSION_KEY = "billing_session";

export interface Session {
  email: string;
  name: string;
  portal: Portal;
  tenant: string;
}

export function login(email: string, password: string): Session | null {
  const user = MOCK_USERS.find((u) => u.email === email && u.password === password);
  if (!user) return null;
  const session: Session = { email: user.email, name: user.name, portal: user.portal, tenant: user.tenant };
  localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  return session;
}

export function logout() {
  localStorage.removeItem(SESSION_KEY);
}

export function getSession(): Session | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(SESSION_KEY);
  if (!raw) return null;
  try { return JSON.parse(raw) as Session; } catch { return null; }
}
