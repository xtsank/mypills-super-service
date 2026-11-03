import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { ProfileResDto } from "../api/types";
import { SESSION_EXPIRED_EVENT } from "../api/http";

type AuthState = {
  token: string | null;
  user: ProfileResDto | null;
};

type AuthContextValue = AuthState & {
  isAuthenticated: boolean;
  isAdmin: boolean;
  setAuth: (token: string, user: ProfileResDto) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

const STORAGE_KEY = "mypills-auth";

function tokenClaims(token: string | null): Record<string, unknown> | null {
  if (!token) {
    return null;
  }
  try {
    const encoded = token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(atob(encoded)) as Record<string, unknown>;
  } catch {
    return null;
  }
}

function tokenIsValid(token: string | null): boolean {
  const claims = tokenClaims(token);
  return typeof claims?.exp === "number" && claims.exp > Date.now() / 1000;
}

function loadAuthState(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return { token: null, user: null };
    }
    const saved = JSON.parse(raw) as AuthState;
    if (typeof saved?.token !== "string" || !tokenIsValid(saved.token) ||
      typeof saved.user?.login !== "string") {
      return { token: null, user: null };
    }
    return saved;
  } catch {
    return { token: null, user: null };
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>(() => loadAuthState());

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  }, [state]);

  useEffect(() => {
    function handleExpired(event: Event) {
      const rejectedToken = (event as CustomEvent<string>).detail;
      setState((current) => current.token === rejectedToken
        ? { token: null, user: null }
        : current);
    }
    window.addEventListener(SESSION_EXPIRED_EVENT, handleExpired);
    return () => window.removeEventListener(SESSION_EXPIRED_EVENT, handleExpired);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      token: state.token,
      user: state.user,
      isAuthenticated: tokenIsValid(state.token),
      isAdmin: tokenClaims(state.token)?.is_admin === true,
      setAuth: (token, user) => setState({ token, user }),
      logout: () => setState({ token: null, user: null })
    }),
    [state.token, state.user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("AuthProvider is missing");
  }
  return ctx;
}
