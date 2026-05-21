import { createContext, useContext, useMemo, useState } from "react";

type StatusType = "info" | "success" | "error";

type ProcessState = {
  statusMessage: string;
  statusType: StatusType;
  cabinetItemsCount: number;
  recommendationsCount: number;
  lastAction: string;
};

type ProcessContextValue = ProcessState & {
  setStatus: (type: StatusType, message: string) => void;
  clearStatus: () => void;
  setCabinetItemsCount: (count: number) => void;
  setRecommendationsCount: (count: number) => void;
  setLastAction: (action: string) => void;
};

const ProcessContext = createContext<ProcessContextValue | null>(null);

const initialState: ProcessState = {
  statusMessage: "",
  statusType: "info",
  cabinetItemsCount: 0,
  recommendationsCount: 0,
  lastAction: ""
};

export function ProcessProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<ProcessState>(initialState);

  const value = useMemo<ProcessContextValue>(
    () => ({
      ...state,
      setStatus: (type, message) => setState((prev) => ({ ...prev, statusType: type, statusMessage: message })),
      clearStatus: () => setState((prev) => ({ ...prev, statusMessage: "" })),
      setCabinetItemsCount: (count) => setState((prev) => ({ ...prev, cabinetItemsCount: count })),
      setRecommendationsCount: (count) => setState((prev) => ({ ...prev, recommendationsCount: count })),
      setLastAction: (action) => setState((prev) => ({ ...prev, lastAction: action }))
    }),
    [state]
  );

  return <ProcessContext.Provider value={value}>{children}</ProcessContext.Provider>;
}

export function useProcess() {
  const ctx = useContext(ProcessContext);
  if (!ctx) {
    throw new Error("ProcessProvider is missing");
  }
  return ctx;
}

