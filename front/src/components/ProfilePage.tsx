import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { ArrowLeft, Loader2, UserCircle } from "lucide-react";
import { fetchProfile } from "@/api/profile";
import type { ProfileResponseData } from "@/api/types";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/authStore";
import { useTheme } from "@/theme/ThemeContext";

interface ProfilePageProps {
  onBack: () => void;
}

export function ProfilePage({ onBack }: ProfilePageProps) {
  const { t } = useTranslation();
  const { surfaces } = useTheme();
  const authUser = useAuthStore((s) => s.user);

  const [data, setData] = useState<ProfileResponseData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!authUser) {
      setLoading(false);
      return;
    }
    setLoading(true);
    fetchProfile()
      .then(setData)
      .catch(() => setData(null))
      .finally(() => setLoading(false));
  }, [authUser]);

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
          {t("profile.back")}
        </button>
        <p className={surfaces.textMuted}>{t("profile.needLogin")}</p>
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
        {t("profile.back")}
      </button>

      <h1 className="text-3xl font-bold bg-gradient-to-r from-cyan-400 to-blue-500 text-transparent bg-clip-text">
        {t("profile.title")}
      </h1>

      {loading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="w-10 h-10 animate-spin text-indigo-400" />
        </div>
      ) : data ? (
        <div
          className={cn(
            "rounded-2xl border p-6 max-w-lg space-y-4",
            surfaces.card,
            surfaces.cardBorder
          )}
        >
          <div className="flex items-center gap-4">
            <UserCircle className="w-16 h-16 text-indigo-400" />
            <div>
              <p className="text-xl font-bold">{data.user.username}</p>
              <p className={cn("text-sm", surfaces.textMuted)}>
                {data.user.email}
              </p>
            </div>
          </div>
          <div className="border-t border-white/10 pt-4 space-y-2">
            <p className={surfaces.textMuted}>{t("profile.balance")}</p>
            <p className="text-3xl font-mono font-bold text-amber-400">
              {data.user.token_balance}
            </p>
          </div>
          <div className={cn("text-sm", surfaces.textMuted)}>
            {t("profile.inventoryPreview")}: {data.inventory.length}
          </div>
        </div>
      ) : (
        <p className={surfaces.textMuted}>{t("catalog.error")}</p>
      )}
    </div>
  );
}
