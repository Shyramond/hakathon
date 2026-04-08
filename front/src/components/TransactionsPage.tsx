import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { format } from "date-fns";
import { enUS, ru } from "date-fns/locale";
import { ArrowLeft, Loader2 } from "lucide-react";
import { fetchTransactionsPage } from "@/api/shop";
import type { TransactionDTO } from "@/api/types";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/authStore";
import { useTheme } from "@/theme/ThemeContext";

interface TransactionsPageProps {
  onBack: () => void;
}

export function TransactionsPage({ onBack }: TransactionsPageProps) {
  const { t, i18n } = useTranslation();
  const { surfaces } = useTheme();
  const authUser = useAuthStore((s) => s.user);

  const [rows, setRows] = useState<TransactionDTO[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(() => {
    if (!authUser) return;
    setLoading(true);
    fetchTransactionsPage(50, 0)
      .then((r) => setRows(r.items))
      .catch(() => setRows([]))
      .finally(() => setLoading(false));
  }, [authUser]);

  useEffect(() => {
    load();
  }, [load]);

  const dateLocale = i18n.language === "en" ? enUS : ru;

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
          {t("transactions.back")}
        </button>
        <p className={surfaces.textMuted}>{t("transactions.needLogin")}</p>
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
        {t("transactions.back")}
      </button>

      <h1 className="text-3xl font-bold bg-gradient-to-r from-amber-400 to-orange-500 text-transparent bg-clip-text">
        {t("transactions.title")}
      </h1>

      {loading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="w-10 h-10 animate-spin text-indigo-400" />
        </div>
      ) : rows.length === 0 ? (
        <p className={surfaces.textMuted}>{t("transactions.empty")}</p>
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-white/10">
          <table className="w-full text-sm">
            <thead>
              <tr className={cn("border-b", surfaces.cardBorder)}>
                <th className="text-left p-3">{t("transactions.amount")}</th>
                <th className="text-left p-3">{t("transactions.type")}</th>
                <th className="text-left p-3">{t("transactions.date")}</th>
                <th className="text-left p-3 hidden md:table-cell">
                  {t("transactions.description")}
                </th>
              </tr>
            </thead>
            <tbody>
              {rows.map((tx) => (
                <tr
                  key={tx.id}
                  className={cn("border-b border-white/5", surfaces.card)}
                >
                  <td className="p-3 font-mono font-semibold">{tx.amount}</td>
                  <td className="p-3">{tx.type}</td>
                  <td className="p-3 whitespace-nowrap">
                    {format(new Date(tx.created_at), "Pp", {
                      locale: dateLocale,
                    })}
                  </td>
                  <td className="p-3 hidden md:table-cell max-w-xs truncate">
                    {tx.description}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
