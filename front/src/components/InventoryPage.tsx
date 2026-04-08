import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { ArrowLeft, Loader2, Package } from "lucide-react";
import { equipInventoryItem, fetchInventory } from "@/api/shop";
import type { InventoryItemDTO } from "@/api/types";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/authStore";
import { useTheme } from "@/theme/ThemeContext";

interface InventoryPageProps {
  onBack: () => void;
}

function canEquip(type: string) {
  return type === "skin" || type === "mascot_skin";
}

export function InventoryPage({ onBack }: InventoryPageProps) {
  const { t } = useTranslation();
  const { surfaces } = useTheme();
  const authUser = useAuthStore((s) => s.user);
  const refreshProfile = useAuthStore((s) => s.refreshProfile);

  const [items, setItems] = useState<InventoryItemDTO[]>([]);
  const [loading, setLoading] = useState(true);
  const [busyId, setBusyId] = useState<string | null>(null);

  const load = useCallback(() => {
    if (!authUser) return;
    setLoading(true);
    fetchInventory()
      .then(setItems)
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, [authUser]);

  useEffect(() => {
    load();
  }, [load]);

  const onEquip = async (id: string) => {
    setBusyId(id);
    try {
      await equipInventoryItem(id);
      await refreshProfile();
      load();
    } catch {
      /* toast */
    } finally {
      setBusyId(null);
    }
  };

  if (!authUser) {
    return (
      <div className="pb-24 space-y-4">
        <button
          type="button"
          onClick={onBack}
          className={cn(
            "inline-flex items-center gap-2 text-sm font-semibold",
            surfaces.textMuted,
            "hover:text-indigo-400"
          )}
        >
          <ArrowLeft className="w-4 h-4" />
          {t("inventory.back")}
        </button>
        <p className={surfaces.textMuted}>{t("inventory.needLogin")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6 pb-24">
      <button
        type="button"
        onClick={onBack}
        className={cn(
          "inline-flex items-center gap-2 text-sm font-semibold",
          surfaces.textMuted,
          "hover:text-indigo-400"
        )}
      >
        <ArrowLeft className="w-4 h-4" />
        {t("inventory.back")}
      </button>

      <h1 className="text-3xl font-bold bg-gradient-to-r from-emerald-400 to-teal-500 text-transparent bg-clip-text">
        {t("inventory.title")}
      </h1>

      {loading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="w-10 h-10 animate-spin text-indigo-400" />
        </div>
      ) : items.length === 0 ? (
        <p className={surfaces.textMuted}>{t("inventory.empty")}</p>
      ) : (
        <ul className="grid gap-4 md:grid-cols-2">
          {items.map((it) => (
            <li
              key={it.id}
              className={cn(
                "rounded-2xl border p-4 flex gap-4",
                surfaces.card,
                surfaces.cardBorder
              )}
            >
              <div className="w-24 h-24 shrink-0 rounded-xl bg-slate-800 flex items-center justify-center overflow-hidden">
                {it.image_url ? (
                  <img
                    src={it.image_url}
                    alt=""
                    className="max-w-full max-h-full object-contain"
                  />
                ) : (
                  <Package className="w-10 h-10 text-slate-500" />
                )}
              </div>
              <div className="flex-1 min-w-0 space-y-2">
                <p className="font-bold truncate">{it.benefit_name}</p>
                <p className={cn("text-xs", surfaces.textMuted)}>
                  {it.benefit_type}
                </p>
                {canEquip(it.benefit_type) && (
                  <button
                    type="button"
                    disabled={busyId === it.id || it.is_equipped}
                    onClick={() => onEquip(it.id)}
                    className="text-sm font-semibold bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white px-3 py-1.5 rounded-lg"
                  >
                    {busyId === it.id
                      ? t("catalog.loading")
                      : it.is_equipped
                        ? t("inventory.equipped")
                        : t("inventory.equip")}
                  </button>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
