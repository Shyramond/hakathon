/**
 * Забрать access_token из hash/query после редиректа Xsolla Login (response_type=token).
 */
export function consumeOAuthTokenFromUrl(): string | null {
  const hash = window.location.hash.replace(/^#/, "");
  const hp = new URLSearchParams(hash);
  let token = hp.get("access_token");

  if (!token) {
    const sp = new URLSearchParams(window.location.search);
    token = sp.get("access_token") || sp.get("token");
  }

  if (token) {
    const url = new URL(window.location.href);
    url.hash = "";
    url.searchParams.delete("access_token");
    url.searchParams.delete("token");
    window.history.replaceState(null, "", `${url.pathname}${url.search}`);
  }

  return token;
}

export function buildXsollaLoginUrl(): string | null {
  const clientId = import.meta.env.VITE_XSOLLA_LOGIN_CLIENT_ID?.trim();
  const redirectUri =
    import.meta.env.VITE_XSOLLA_LOGIN_REDIRECT_URI?.trim() ||
    `${window.location.origin}${window.location.pathname}`;

  if (!clientId) return null;

  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: redirectUri,
    response_type: "token",
    state: crypto.randomUUID(),
  });

  return `https://login.xsolla.com/api/oauth2/auth?${params.toString()}`;
}
