// API client for the 3rike backend.
//
// Responsibilities:
//   - Read base URL from VITE_API_URL.
//   - Inject the JWT from session storage on every request.
//   - Normalize errors into a typed ApiError that carries HTTP status + the
//     backend "error" code so callers can show specific messages.
//   - Be tolerant of Render free-tier cold starts (first hit can take ~60s).
//   - On 401, clear the local session and emit a global event so AuthProvider
//     can react without coupling api.ts to React.

import { clearSession, getToken } from "./session";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

// Cold starts on Render free tier can take ~60s.
const DEFAULT_TIMEOUT_MS = 90_000;

export class ApiError extends Error {
  status: number;
  code: string;
  constructor(status: number, code: string, message?: string) {
    super(message ?? code);
    this.status = status;
    this.code = code;
  }
}

// Fired when the API returns 401, so AuthProvider can navigate to /login.
// Using a window event keeps this module decoupled from React/router.
export const UNAUTHORIZED_EVENT = "3rike:unauthorized";

type RequestOptions = RequestInit & {
  /** Skip JWT injection (used by /auth/login + /auth/register). */
  skipAuth?: boolean;
  /** Override the default timeout. */
  timeoutMs?: number;
};

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { skipAuth, timeoutMs = DEFAULT_TIMEOUT_MS, headers, ...rest } = opts;

  const finalHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...(headers as Record<string, string> | undefined),
  };

  if (!skipAuth) {
    const token = getToken();
    if (token) finalHeaders["Authorization"] = `Bearer ${token}`;
  }

  const controller = new AbortController();
  const timeout = window.setTimeout(() => controller.abort(), timeoutMs);

  let res: Response;
  try {
    res = await fetch(`${API_URL}${path}`, {
      ...rest,
      headers: finalHeaders,
      signal: controller.signal,
    });
  } catch (err) {
    if ((err as Error).name === "AbortError") {
      throw new ApiError(0, "timeout", "Request timed out. Please try again.");
    }
    throw new ApiError(0, "network_error", "Network error. Please check your connection.");
  } finally {
    window.clearTimeout(timeout);
  }

  // Handle 401 globally: drop the session and let listeners react.
  if (res.status === 401 && !skipAuth) {
    clearSession();
    window.dispatchEvent(new CustomEvent(UNAUTHORIZED_EVENT));
  }

  const text = await res.text();
  const data = text ? safeParse(text) : {};

  if (!res.ok) {
    const code = (data as Record<string, unknown>)?.error;
    throw new ApiError(
      res.status,
      typeof code === "string" ? code : "request_failed",
    );
  }

  return data as T;
}

function safeParse(text: string): unknown {
  try {
    return JSON.parse(text);
  } catch {
    return {};
  }
}

// =============================================================================
// Auth
// =============================================================================

export type Role = "driver" | "investor" | "admin";

export type User = {
  id: number;
  email: string;
  role: Role;
  created_at: string;
};

export type LoginResponse = {
  token: string;
  session_id: string;
};

export function register(payload: {
  email: string;
  password: string;
  role: Role;
}): Promise<User> {
  return request<User>("/auth/register", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true,
  });
}

export function login(payload: {
  email: string;
  password: string;
}): Promise<LoginResponse> {
  return request<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true,
  });
}

export function me(): Promise<User> {
  return request<User>("/auth/me");
}

export function logout(): Promise<void> {
  return request<void>("/auth/logout", { method: "POST" });
}

// =============================================================================
// Waitlist (public — pre-launch landing page)
// =============================================================================

export type WaitlistEntry = {
  id: number;
  email: string;
  phone?: string;
  referral_code: string;
  referred_by?: string;
  position: number;
  created_at: string;
};

export type JoinResponse = {
  entry: WaitlistEntry;
  totalCount: number;
};

export type StatsResponse = {
  totalCount: number;
};

export type GetByCodeResponse = {
  entry: WaitlistEntry;
  totalCount: number;
  referralCount: number;
};

export function joinWaitlist(payload: {
  email: string;
  phone?: string;
  referredBy?: string;
}): Promise<JoinResponse> {
  return request<JoinResponse>("/waitlist/join", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true,
  });
}

export function getWaitlistStats(): Promise<StatsResponse> {
  return request<StatsResponse>("/waitlist/stats", { skipAuth: true });
}

export function getWaitlistEntry(code: string): Promise<GetByCodeResponse> {
  return request<GetByCodeResponse>(`/waitlist/${encodeURIComponent(code)}`, {
    skipAuth: true,
  });
}
