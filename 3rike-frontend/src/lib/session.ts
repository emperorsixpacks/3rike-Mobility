// Session storage for the authenticated user. Keeps JWT + session ID in
// localStorage behind a single API so the rest of the app never reaches in
// directly. Switching to httpOnly cookies later only requires changing this
// file.
//
// We intentionally avoid storing the User object here — call `/auth/me` on
// boot and let AuthProvider hold the resolved user in memory.

const TOKEN_KEY = "3rike.token";
const SESSION_ID_KEY = "3rike.sessionId";

export type StoredSession = {
  token: string;
  sessionId: string;
};

export function getSession(): StoredSession | null {
  const token = localStorage.getItem(TOKEN_KEY);
  const sessionId = localStorage.getItem(SESSION_ID_KEY);
  if (!token || !sessionId) return null;
  return { token, sessionId };
}

export function setSession(s: StoredSession): void {
  localStorage.setItem(TOKEN_KEY, s.token);
  localStorage.setItem(SESSION_ID_KEY, s.sessionId);
}

export function clearSession(): void {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(SESSION_ID_KEY);
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}
