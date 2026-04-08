import {
  API_ERROR_TOAST_EVENT,
  AUTH_401_EVENT,
  AUTH_TOKEN_KEY,
} from "./constants";
import type { ApiEnvelope } from "./types";

function getApiBase(): string {
  const base =
    import.meta.env.VITE_API_BASE_URL || "http://127.0.0.1:8080/api/v1";
  return base.replace(/\/$/, "");
}

function extractMessage(json: Record<string, unknown>, res: Response): string {
  if (typeof json.message === "string") return json.message;
  const err = json.error;
  if (err && typeof err === "object" && err !== null && "message" in err) {
    return String((err as { message?: string }).message);
  }
  if (typeof err === "string") return err;
  return res.statusText || "Ошибка запроса";
}

function emitApiError(message: string) {
  window.dispatchEvent(
    new CustomEvent(API_ERROR_TOAST_EVENT, { detail: { message } })
  );
}

function clearAuthAndNotify401() {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  window.dispatchEvent(new CustomEvent(AUTH_401_EVENT));
}

export type RequestOptions = RequestInit & {
  /** Не добавлять Authorization */
  publicRoute?: boolean;
};

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const { publicRoute, headers: initHeaders, ...rest } = options;
  const headers = new Headers(initHeaders);

  if (
    rest.body &&
    typeof rest.body === "string" &&
    !headers.has("Content-Type")
  ) {
    headers.set("Content-Type", "application/json");
  }

  const token = localStorage.getItem(AUTH_TOKEN_KEY);
  if (!publicRoute && token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const rel = path.startsWith("/") ? path : `/${path}`;
  const url = `${getApiBase()}${rel}`;

  const res = await fetch(url, { ...rest, headers });

  const text = await res.text();
  let json: Record<string, unknown> = {};
  if (text) {
    try {
      json = JSON.parse(text) as Record<string, unknown>;
    } catch {
      json = {};
    }
  }

  if (res.status === 401) {
    const msg = extractMessage(json, res);
    if (!publicRoute) {
      clearAuthAndNotify401();
    }
    emitApiError(msg);
    throw new Error(msg);
  }

  if (!res.ok) {
    const msg = extractMessage(json, res);
    emitApiError(msg);
    throw new Error(msg);
  }

  if ("success" in json && json.success === false) {
    const msg = extractMessage(json, res);
    emitApiError(msg);
    throw new Error(msg);
  }

  if ("success" in json && json.success === true && "data" in json) {
    return (json as ApiEnvelope<T>).data as T;
  }

  return json as T;
}

export async function apiRequestWithMeta<T>(
  path: string,
  options: RequestOptions = {}
): Promise<{ data: T; meta?: ApiEnvelope<T>["meta"] }> {
  const { publicRoute, headers: initHeaders, ...rest } = options;
  const headers = new Headers(initHeaders);

  if (
    rest.body &&
    typeof rest.body === "string" &&
    !headers.has("Content-Type")
  ) {
    headers.set("Content-Type", "application/json");
  }

  const token = localStorage.getItem(AUTH_TOKEN_KEY);
  if (!publicRoute && token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const rel = path.startsWith("/") ? path : `/${path}`;
  const url = `${getApiBase()}${rel}`;

  const res = await fetch(url, { ...rest, headers });
  const text = await res.text();
  let json: ApiEnvelope<T> & Record<string, unknown> = { success: false };
  if (text) {
    try {
      json = JSON.parse(text) as ApiEnvelope<T> & Record<string, unknown>;
    } catch {
      json = { success: false };
    }
  }

  if (res.status === 401) {
    const msg = extractMessage(json as Record<string, unknown>, res);
    if (!publicRoute) {
      clearAuthAndNotify401();
    }
    emitApiError(msg);
    throw new Error(msg);
  }

  if (!res.ok) {
    const msg = extractMessage(json as Record<string, unknown>, res);
    emitApiError(msg);
    throw new Error(msg);
  }

  if (json.success === false) {
    const msg = extractMessage(json as Record<string, unknown>, res);
    emitApiError(msg);
    throw new Error(msg);
  }

  return { data: json.data as T, meta: json.meta };
}
