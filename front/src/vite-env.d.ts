/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_XSOLLA_LOGIN_CLIENT_ID: string;
  readonly VITE_XSOLLA_LOGIN_REDIRECT_URI: string;
  readonly VITE_XSOLLA_PROJECT_ID: string;
  readonly VITE_XSOLLA_CATALOG_BASE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
