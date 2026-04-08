import { useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import { ArrowLeft, Loader2, Package, ShoppingCart, Tag } from "lucide-react";
import { fetchBenefitById, fetchBenefits } from "@/api/benefits";
import { purchaseBenefit } from "@/api/shop";
import { fetchXsollaVirtualItems } from "@/api/xsollaCatalog";
import type { BenefitDTO, CatalogVirtualItem } from "@/api/types";
import { useToast } from "@/context/ToastContext";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/authStore";
import { useTheme } from "@/theme/ThemeContext";

export function Catalog() {
  const { t } = useTranslation();
  const { surfaces } = useTheme();
  const showToast = useToast();
  const user = useAuthStore((s) => s.user);
  const refreshProfile = useAuthStore((s) => s.refreshProfile);

  const [benefits, setBenefits] = useState<BenefitDTO[]>([]);
  const [benLoading, setBenLoading] = useState(true);
  const [benError, setBenError] = useState<string | null>(null);

  const [xsollaItems, setXsollaItems] = useState<CatalogVirtualItem[]>([]);
  const [xsLoading, setXsLoading] = useState(true);
  const [xsError, setXsError] = useState<string | null>(null);

  const [detailId, setDetailId] = useState<string | null>(null);
  const [detail, setDetail] = useState<BenefitDTO | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [purchaseBusy, setPurchaseBusy] = useState(false);

  const [groupFilter, setGroupFilter] = useState<string>("all");

  const loadBenefits = useCallback(() => {
    setBenLoading(true);
    setBenError(null);
    fetchBenefits()
      .then(setBenefits)
      .catch(() => setBenError(t("catalog.error")))
      .finally(() => setBenLoading(false));
  }, [t]);

  const loadXsolla = useCallback(() => {
    setXsLoading(true);
    setXsError(null);
    const pid = import.meta.env.VITE_XSOLLA_PROJECT_ID?.trim();
    if (!pid) {
      setXsollaItems([]);
      setXsError(t("catalog.catalogNotConfigured"));
      setXsLoading(false);
      return;
    }
    fetchXsollaVirtualItems()
      .then(setXsollaItems)
      .catch((e: Error) => setXsError(e.message || t("catalog.error")))
      .finally(() => setXsLoading(false));
  }, [t]);

  useEffect(() => {
    loadBenefits();
    loadXsolla();
  }, [loadBenefits, loadXsolla]);

  useEffect(() => {
    if (!detailId) {
      setDetail(null);
      return;
    }
    setDetailLoading(true);
    fetchBenefitById(detailId)
      .then(setDetail)
      .catch(() => setDetail(null))
      .finally(() => setDetailLoading(false));
  }, [detailId]);

  const allGroups = useMemo(() => {
    const s = new Set<string>();
    xsollaItems.forEach((i) => i.groups.forEach((g) => s.add(g)));
    return [...s].sort();
  }, [xsollaItems]);

  const filteredXsolla = useMemo(() => {
    if (groupFilter === "all") return xsollaItems;
    return xsollaItems.filter((i) => i.groups.includes(groupFilter));
  }, [xsollaItems, groupFilter]);

  const handlePurchase = async (id: string) => {
    if (!user) {
      showToast(t("catalog.needLogin"), "info");
      return;
    }
    setPurchaseBusy(true);
    try {
      await purchaseBenefit(id);
      await refreshProfile();
      showToast(t("catalog.purchaseSuccess"), "info");
      setDetailId(null);
      loadBenefits();
    } catch {
      /* toast из api */
    } finally {
      setPurchaseBusy(false);
    }
  };

  const handleXsollaPurchase = (item: CatalogVirtualItem) => {
    if (!user) {
      showToast(t("catalog.needLogin"), "info");
      return;
    }

    const projectId = import.meta.env.VITE_XSOLLA_PROJECT_ID?.trim();
    if (!projectId) {
      showToast(t("catalog.catalogNotConfigured"), "error");
      return;
    }

    const catalogBase = import.meta.env.VITE_XSOLLA_CATALOG_BASE_URL?.trim() || "/xsolla-catalog";
    const url = `${catalogBase}/paystation2/?projectId=${encodeURIComponent(projectId)}&sku=${encodeURIComponent(item.sku)}`;
    window.open(url, "_blank", "noopener,noreferrer");
  };

  const skeletonCard = (
    <div
      className={cn(
        "rounded-2xl border overflow-hidden animate-pulse",
        surfaces.card,
        surfaces.cardBorder
      )}
    >
      <div className="aspect-video bg-slate-800/50" />
      <div className="p-5 space-y-3">
        <div className="h-5 bg-slate-800/80 rounded w-2/3" />
        <div className="h-3 bg-slate-800/60 rounded w-full" />
      </div>
    </div>
  );

  if (detailId) {
    return (
      <div className="space-y-6 pb-24">
        <button
          type="button"
          onClick={() => setDetailId(null)}
          className={cn(
            "inline-flex items-center gap-2 text-sm font-semibold",
            surfaces.textMuted,
            "hover:text-indigo-400"
          )}
        >
          <ArrowLeft className="w-4 h-4" />
          {t("catalog.back")}
        </button>

        {detailLoading ? (
          <div className="flex justify-center py-16">
            <Loader2 className="w-10 h-10 animate-spin text-indigo-400" />
          </div>
        ) : detail ? (
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            className={cn(
              "rounded-2xl border overflow-hidden max-w-lg mx-auto",
              surfaces.card,
              surfaces.cardBorder
            )}
          >
            <div className="aspect-video bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center">
              {detail.image_url ? (
                <img
                  src={detail.image_url}
                  alt=""
                  className="max-h-full max-w-full object-contain"
                />
              ) : (
                <Package className="w-20 h-20 text-slate-500" />
              )}
            </div>
            <div className="p-6 space-y-4">
              <div>
                <h2 className="text-2xl font-bold">{detail.name}</h2>
                <p className={cn("text-sm mt-1", surfaces.textMuted)}>
                  {t("common.type")}: {detail.type}
                </p>
              </div>
              <p className={surfaces.textMuted}>{detail.description}</p>
              {!detail.is_active && (
                <p className="text-amber-500 text-sm font-medium">
                  {t("catalog.inactive")}
                </p>
              )}
              <div className="flex items-center justify-between pt-2">
                <span className="text-xl font-bold text-indigo-400">
                  {detail.price_tokens} {t("catalog.tokens")}
                </span>
                <button
                  type="button"
                  disabled={!detail.is_active || purchaseBusy}
                  onClick={() => handlePurchase(detail.id)}
                  className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white py-3 px-5 rounded-xl font-semibold"
                >
                  {purchaseBusy ? (
                    <Loader2 className="w-5 h-5 animate-spin" />
                  ) : (
                    <ShoppingCart className="w-5 h-5" />
                  )}
                  {t("catalog.buy")}
                </button>
              </div>
            </div>
          </motion.div>
        ) : (
          <p className={surfaces.textMuted}>{t("catalog.error")}</p>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-10 pb-24">
      <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-violet-400 to-fuchsia-400 text-transparent bg-clip-text">
        {t("catalog.title")}
      </h1>

      <section className="space-y-4">
        <h2 className="text-xl font-bold flex items-center gap-2">
          <Tag className="w-6 h-6 text-violet-400" />
          {t("catalog.benefits")}
        </h2>
        {benError && (
          <div className="flex items-center gap-3 text-red-400 text-sm">
            {benError}
            <button
              type="button"
              onClick={loadBenefits}
              className="underline font-semibold"
            >
              {t("catalog.retry")}
            </button>
          </div>
        )}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {benLoading
            ? [1, 2, 3].map((k) => <div key={k}>{skeletonCard}</div>)
            : benefits.length === 0
              ? (
                  <p className={surfaces.textMuted}>{t("catalog.noBenefits")}</p>
                )
              : benefits.map((b) => (
                  <motion.div
                    key={b.id}
                    initial={{ opacity: 0, y: 16 }}
                    animate={{ opacity: 1, y: 0 }}
                    onClick={() => setDetailId(b.id)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter" || event.key === " ") {
                        event.preventDefault();
                        setDetailId(b.id);
                      }
                    }}
                    tabIndex={0}
                    role="button"
                    className={cn(
                      "cursor-pointer text-left rounded-2xl border overflow-hidden transition-all hover:border-indigo-500/50 focus:outline-none focus:ring-2 focus:ring-indigo-500/50",
                      surfaces.card,
                      surfaces.cardBorder
                    )}
                  >
                    <div className="aspect-video bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center">
                      {b.image_url ? (
                        <img
                          src={b.image_url}
                          alt=""
                          className="max-h-full max-w-full object-contain"
                        />
                      ) : (
                        <Package className="w-16 h-16 text-slate-500" />
                      )}
                    </div>
                    <div className="p-4 space-y-3">
                      <div>
                        <h3 className="font-bold text-lg">{b.name}</h3>
                        <p className={cn("text-sm line-clamp-2", surfaces.textMuted)}>
                          {b.description}
                        </p>
                      </div>
                      <div className="flex items-center justify-between gap-3">
                        <p className="text-indigo-400 font-semibold">
                          {b.price_tokens} {t("catalog.tokens")}
                        </p>
                        <button
                          type="button"
                          onClick={(event) => {
                            event.stopPropagation();
                            handlePurchase(b.id);
                          }}
                          disabled={!b.is_active || purchaseBusy}
                          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white py-2 px-4 rounded-xl font-semibold"
                        >
                          <ShoppingCart className="w-4 h-4" />
                          {t("catalog.buy")}
                        </button>
                      </div>
                      {!b.is_active && (
                        <p className="text-amber-500 text-sm font-medium">
                          {t("catalog.inactive")}
                        </p>
                      )}
                    </div>
                  </motion.div>
                ))}
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-xl font-bold flex items-center gap-2">
          <Package className="w-6 h-6 text-fuchsia-400" />
          {t("catalog.xsollaItems")}
        </h2>

        {allGroups.length > 0 && (
          <div className="flex flex-wrap gap-2 items-center">
            <span className={cn("text-sm font-medium", surfaces.textMuted)}>
              {t("catalog.filterGroups")}:
            </span>
            <button
              type="button"
              onClick={() => setGroupFilter("all")}
              className={cn(
                "px-3 py-1 rounded-lg text-sm font-semibold border",
                groupFilter === "all"
                  ? "bg-indigo-600 border-indigo-500 text-white"
                  : cn(surfaces.card, surfaces.cardBorder)
              )}
            >
              {t("catalog.allGroups")}
            </button>
            {allGroups.map((g) => (
              <button
                type="button"
                key={g}
                onClick={() => setGroupFilter(g)}
                className={cn(
                  "px-3 py-1 rounded-lg text-sm font-semibold border",
                  groupFilter === g
                    ? "bg-indigo-600 border-indigo-500 text-white"
                    : cn(surfaces.card, surfaces.cardBorder)
                )}
              >
                {g}
              </button>
            ))}
          </div>
        )}

        {xsError && !xsLoading && (
          <div className="flex items-center gap-3 text-amber-400/90 text-sm">
            {xsError}
            {import.meta.env.VITE_XSOLLA_PROJECT_ID?.trim() && (
              <button
                type="button"
                onClick={loadXsolla}
                className="underline font-semibold"
              >
                {t("catalog.retry")}
              </button>
            )}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {xsLoading
            ? [1, 2, 3].map((k) => <div key={k}>{skeletonCard}</div>)
            : filteredXsolla.length === 0
              ? (
                  <p className={surfaces.textMuted}>{t("catalog.noItems")}</p>
                )
              : filteredXsolla.map((item) => (
                  <motion.div
                    key={item.sku}
                    initial={{ opacity: 0, y: 16 }}
                    animate={{ opacity: 1, y: 0 }}
                    className={cn(
                      "rounded-2xl border overflow-hidden",
                      surfaces.card,
                      surfaces.cardBorder
                    )}
                  >
                    <div className="aspect-video bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center">
                      {item.imageUrl ? (
                        <img
                          src={item.imageUrl}
                          alt=""
                          className="max-h-full max-w-full object-contain"
                        />
                      ) : (
                        <Package className="w-16 h-16 text-slate-500" />
                      )}
                    </div>
                    <div className="p-4 space-y-4">
                      <div className="space-y-2">
                        <h3 className="font-bold">{item.name}</h3>
                        <p className={cn("text-sm line-clamp-3", surfaces.textMuted)}>
                          {item.description}
                        </p>
                      </div>
                      <div className="space-y-2 text-sm">
                        {item.price != null && (
                          <p>
                            {t("common.price")}: {item.price}{" "}
                            {item.currency ?? ""}
                          </p>
                        )}
                        {item.virtualPrice != null && (
                          <p className="text-indigo-400">
                            {t("layout.tokens")}: {item.virtualPrice}
                          </p>
                        )}
                        {item.groups.length > 0 && (
                          <p className={surfaces.textMuted}>
                            {item.groups.join(", ")}
                          </p>
                        )}
                      </div>
                      <button
                        type="button"
                        onClick={() => handleXsollaPurchase(item)}
                        className="w-full inline-flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-500 text-white py-2 rounded-xl font-semibold"
                      >
                        <ShoppingCart className="w-4 h-4" />
                        {t("catalog.buy")}
                      </button>
                    </div>
                  </motion.div>
                ))}
        </div>
      </section>
    </div>
  );
}
