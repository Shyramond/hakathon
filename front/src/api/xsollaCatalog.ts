import type { CatalogVirtualItem } from "./types";

interface RawXsollaItem {
  sku?: string;
  name?: string;
  description?: string;
  image_url?: string;
  price?: { amount?: number; currency?: string };
  virtual_prices?: { amount?: number }[];
  groups?: { external_id?: string }[];
  attributes?: unknown;
}

function mapItem(item: RawXsollaItem): CatalogVirtualItem {
  return {
    sku: item.sku ?? "",
    name: item.name ?? "",
    description: item.description ?? "",
    imageUrl: item.image_url ?? "",
    price: item.price?.amount,
    currency: item.price?.currency,
    virtualPrice: item.virtual_prices?.[0]?.amount,
    groups: (item.groups ?? [])
      .map((g) => g.external_id)
      .filter((id): id is string => !!id),
    attributes: item.attributes,
  };
}

function getCatalogBase(): string {
  return (
    import.meta.env.VITE_XSOLLA_CATALOG_BASE_URL || "/xsolla-catalog"
  ).replace(/\/$/, "");
}

function getProjectId(): string {
  return import.meta.env.VITE_XSOLLA_PROJECT_ID || "";
}

/**
 * Каталог виртуальных предметов Xsolla Store API v2.
 * В dev используйте прокси `/xsolla-catalog` (см. vite.config.ts).
 */
export async function fetchXsollaVirtualItems(): Promise<CatalogVirtualItem[]> {
  const projectId = getProjectId();
  if (!projectId) {
    throw new Error("Не задан VITE_XSOLLA_PROJECT_ID");
  }

  const base = getCatalogBase();
  const path = `/api/v2/project/${projectId}/items/virtual_items?locale=ru_RU&country=RU`;
  const url = `${base}${path}`;

  const res = await fetch(url);
  if (!res.ok) {
    const t = await res.text();
    throw new Error(t || `Catalog HTTP ${res.status}`);
  }

  const json = (await res.json()) as { items?: RawXsollaItem[] } | RawXsollaItem[];
  const items = Array.isArray(json) ? json : json.items ?? [];
  return items.map(mapItem);
}
