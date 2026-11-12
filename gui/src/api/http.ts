import { AppError } from "./types";

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";
export const SESSION_EXPIRED_EVENT = "mypills:session-expired";

type RequestOptions = RequestInit & {
  token?: string | null;
};

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Content-Type", "application/json");

  if (options.token) {
    headers.set("Authorization", `Bearer ${options.token}`);
  }

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      headers
    });
  } catch (error) {
    throw toAppError(error);
  }

  if (!response.ok) {
    const error = await parseError(response);
    if (response.status === 401 && options.token) {
      window.dispatchEvent(new CustomEvent(SESSION_EXPIRED_EVENT, {
        detail: options.token
      }));
    }
    throw error;
  }

  if (response.status === 204) {
    return {} as T;
  }

  return (await response.json()) as T;
}

async function parseError(response: Response): Promise<AppError> {
  try {
    const data = (await response.json()) as AppError;
    if (data?.code && data?.message) {
      return data;
    }
  } catch {
    // ignore JSON parse errors
  }

  return {
    code: String(response.status),
    message: response.statusText || "Запрос не выполнен"
  };
}

function toAppError(error: unknown): AppError {
  if (error instanceof Error) {
    return {
      code: "NETWORK",
      message: `Ошибка сети: ${error.message}`
    };
  }

  return {
    code: "NETWORK",
    message: "Ошибка сети: запрос не был отправлен"
  };
}
