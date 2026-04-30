// AuthProvider holds the resolved User in memory and exposes login/register/
// logout helpers. On boot, if a JWT exists in storage, it calls /auth/me to
// resolve the user. Listens for global 401s (from api.ts) and clears state.

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";
import {
  ApiError,
  UNAUTHORIZED_EVENT,
  login as apiLogin,
  logout as apiLogout,
  me as apiMe,
  register as apiRegister,
  type Role,
  type User,
} from "./api";
import { clearSession, getSession, setSession } from "./session";

type AuthState =
  | { status: "loading"; user: null }
  | { status: "authenticated"; user: User }
  | { status: "anonymous"; user: null };

type AuthContextValue = {
  state: AuthState;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<User>;
  register: (email: string, password: string, role: Role) => Promise<User>;
  logout: () => Promise<void>;
  refresh: () => Promise<User | null>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const [state, setState] = useState<AuthState>(() =>
    getSession() ? { status: "loading", user: null } : { status: "anonymous", user: null },
  );

  // Hold the latest navigate fn in a ref so the 401 listener can use it
  // without resubscribing on every render.
  const navigateRef = useRef(navigate);
  navigateRef.current = navigate;

  const refresh = useCallback(async (): Promise<User | null> => {
    if (!getSession()) {
      setState({ status: "anonymous", user: null });
      return null;
    }
    try {
      const user = await apiMe();
      setState({ status: "authenticated", user });
      return user;
    } catch (err) {
      // Any failure on /me means the token is no good.
      clearSession();
      setState({ status: "anonymous", user: null });
      if (err instanceof ApiError && err.status === 401) {
        // Already cleared by api.ts, but keep flow consistent.
      }
      return null;
    }
  }, []);

  // Boot: resolve the user if we have a stored token.
  useEffect(() => {
    if (state.status === "loading") {
      void refresh();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Global 401 listener: drop user state and bounce to /login.
  useEffect(() => {
    const handler = () => {
      setState({ status: "anonymous", user: null });
      navigateRef.current("/login", { replace: true });
    };
    window.addEventListener(UNAUTHORIZED_EVENT, handler);
    return () => window.removeEventListener(UNAUTHORIZED_EVENT, handler);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const { token, session_id } = await apiLogin({ email, password });
    setSession({ token, sessionId: session_id });
    const user = await apiMe();
    setState({ status: "authenticated", user });
    return user;
  }, []);

  const register = useCallback(
    async (email: string, password: string, role: Role) => {
      const user = await apiRegister({ email, password, role });
      // Backend's register returns the user but no token — log the user in
      // immediately so the app is in an authenticated state on success.
      const { token, session_id } = await apiLogin({ email, password });
      setSession({ token, sessionId: session_id });
      setState({ status: "authenticated", user });
      return user;
    },
    [],
  );

  const logout = useCallback(async () => {
    try {
      await apiLogout();
    } catch {
      // Best-effort — even if the server call fails, drop local state.
    }
    clearSession();
    setState({ status: "anonymous", user: null });
    navigate("/login", { replace: true });
  }, [navigate]);

  const value = useMemo<AuthContextValue>(
    () => ({
      state,
      user: state.user,
      isAuthenticated: state.status === "authenticated",
      isLoading: state.status === "loading",
      login,
      register,
      logout,
      refresh,
    }),
    [state, login, register, logout, refresh],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside <AuthProvider>");
  }
  return ctx;
}
