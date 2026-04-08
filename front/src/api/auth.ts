import { apiRequest } from "./client";
import { AUTH_TOKEN_KEY } from "./constants";
import type { LoginResponseData, UserDTO } from "./types";

/**
 * POST /auth/login: бэкенд принимает Xsolla JWT в поле `token`
 * (в ТЗ также упоминается xsolla_token — дублируем для совместимости).
 */
export async function loginWithXsollaToken(xsollaToken: string): Promise<LoginResponseData> {
  const data = await apiRequest<LoginResponseData>("/auth/login", {
    method: "POST",
    publicRoute: true,
    body: JSON.stringify({
      token: xsollaToken,
      xsolla_token: xsollaToken,
    }),
  });
  localStorage.setItem(AUTH_TOKEN_KEY, xsollaToken);
  return data;
}

export function clearAuthToken() {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}

export function getStoredToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}
