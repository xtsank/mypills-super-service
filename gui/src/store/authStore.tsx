import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { ProfileResDto } from "../api/types";

type AuthState = {
  token: string | null;
  user: ProfileResDto | null;
};

type AuthContextValue = AuthState & {
  isAuthenticated: boolean;
  setAuth: (token: string, user: ProfileResDto) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

const STORAGE_KEY = "mypills-auth";

function loadAuthState(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return { token: null, user: null };
    }
    return JSON.parse(raw) as AuthState;
  } catch {
    return { token: null, user: null };
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>(() => loadAuthState());

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  }, [state]);

  const value = useMemo<AuthContextValue>(
    () => ({
      token: state.token,
      user: state.user,
      isAuthenticated: Boolean(state.token),
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

