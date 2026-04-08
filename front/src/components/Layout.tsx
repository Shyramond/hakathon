import React, { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  ChevronDown,
  Coins,
  Globe,
  Hammer,
  LogIn,
  LogOut,
  Moon,
  Package,
  RefreshCw,
  Store,
  Sun,
  User,
  UserCircle,
} from "lucide-react";
import { useStore, GEM_TIERS } from "../store/useStore";
import { cn } from "../lib/utils";
import { buildXsollaLoginUrl } from "@/lib/xsollaLogin";
import { setAppLanguage } from "@/i18n";
import { useAuthStore } from "@/store/authStore";
import { useTheme } from "@/theme/ThemeContext";
import { useToast } from "@/context/ToastContext";
import { AUTH_401_EVENT } from "@/api/constants";
import type { AppTab } from "@/types/app";

interface LayoutProps {
  children: React.ReactNode;
  activeTab: AppTab;
  setActiveTab: (tab: AppTab) => void;
  onAuth401?: () => void;
}

const gemColors: Record<string, string> = {
  grey: "text-gray-400",
  green: "text-green-400",
  blue: "text-blue-400",
  purple: "text-purple-400",
  gold: "text-yellow-400",
};

export function Layout({
  children,
  activeTab,
  setActiveTab,
  onAuth401,
}: LayoutProps) {
  const { t, i18n } = useTranslation();
  const { mode, toggleTheme, surfaces } = useTheme();
  const gems = useStore((state) => state.gems);
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);
  const showToast = useToast();

  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const onDoc = (e: MouseEvent) => {
      if (!menuRef.current?.contains(e.target as Node)) setMenuOpen(false);
    };
    document.addEventListener("click", onDoc);
    return () => document.removeEventListener("click", onDoc);
  }, []);

  useEffect(() => {
    const h = () => {
      logout();
      onAuth401?.();
    };
    window.addEventListener(AUTH_401_EVENT, h);
    return () => window.removeEventListener(AUTH_401_EVENT, h);
  }, [logout, onAuth401]);

  const tabs: { id: AppTab; label: string; icon: typeof Store }[] = [
    { id: "store", label: t("layout.tabStore"), icon: Store },
    { id: "forge", label: t("layout.tabForge"), icon: Hammer },
    { id: "exchange", label: t("layout.tabExchange"), icon: RefreshCw },
    { id: "daily", label: t("layout.tabDaily"), icon: LogIn },
  ];

  const goXsollaLogin = () => {
    const url = buildXsollaLoginUrl();
    if (!url) {
      showToast(t("auth.xsollaNotConfigured"), "error");
      return;
    }
    window.location.href = url;
  };

  const switchLang = () => {
    const next = i18n.language === "en" ? "ru" : "en";
    setAppLanguage(next);
  };

  return (
    <div className={cn("min-h-screen font-sans selection:bg-indigo-500/30", surfaces.page)}>
      <header
        className={cn(
          "sticky top-0 z-50 border-b backdrop-blur-md",
          surfaces.header,
          surfaces.headerBorder
        )}
      >
        <div className="max-w-5xl mx-auto px-4 min-h-16 py-2 flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center space-x-2">
            <Coins className="w-8 h-8 text-yellow-500 shrink-0" />
            <span className="text-xl font-bold bg-gradient-to-r from-yellow-400 to-yellow-600 text-transparent bg-clip-text">
              {t("layout.brand")}
            </span>
          </div>

          <div className="flex flex-wrap items-center gap-3 md:gap-4">
            <div className="flex space-x-3 md:space-x-4">
              {GEM_TIERS.map((tier) => (
                <div
                  key={tier}
                  className="flex items-center space-x-1"
                  title={tier}
                >
                  <div
                    className={cn(
                      "w-3 h-3 rotate-45",
                      tier === "grey" && "bg-gray-400",
                      tier === "green" &&
                        "bg-green-400 shadow-[0_0_8px_rgba(74,222,128,0.5)]",
                      tier === "blue" &&
                        "bg-blue-400 shadow-[0_0_8px_rgba(96,165,250,0.5)]",
                      tier === "purple" &&
                        "bg-purple-400 shadow-[0_0_8px_rgba(192,132,252,0.5)]",
                      tier === "gold" &&
                        "bg-yellow-400 shadow-[0_0_8px_rgba(250,204,21,0.5)]"
                    )}
                  />
                  <span className={cn("font-medium text-sm", gemColors[tier])}>
                    {gems[tier]}
                  </span>
                </div>
              ))}
            </div>

            {user && (
              <div
                className={cn(
                  "hidden sm:flex items-center gap-1.5 px-2 py-1 rounded-lg border text-sm font-mono font-bold text-amber-400",
                  surfaces.cardBorder,
                  surfaces.card
                )}
              >
                <Coins className="w-4 h-4" />
                {user.token_balance}
              </div>
            )}

            <div
              className={cn(
                "flex items-center gap-1 border-l pl-3 ml-1",
                mode === "light" ? "border-slate-200" : "border-white/10"
              )}
            >
              <button
                type="button"
                onClick={toggleTheme}
                className={cn(
                  "p-2 rounded-lg transition-colors",
                  surfaces.textMuted,
                  "hover:bg-white/10"
                )}
                title={
                  mode === "dark" ? t("layout.themeLight") : t("layout.themeDark")
                }
              >
                {mode === "dark" ? (
                  <Sun className="w-5 h-5" />
                ) : (
                  <Moon className="w-5 h-5" />
                )}
              </button>
              <button
                type="button"
                onClick={switchLang}
                className={cn(
                  "p-2 rounded-lg font-bold text-xs flex items-center gap-0.5",
                  surfaces.textMuted,
                  "hover:bg-white/10"
                )}
                title="Language"
              >
                <Globe className="w-4 h-4" />
                {i18n.language === "en" ? t("layout.langEn") : t("layout.langRu")}
              </button>
            </div>

            {!user ? (
              <button
                type="button"
                onClick={goXsollaLogin}
                className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm font-semibold px-4 py-2 rounded-xl"
              >
                <LogIn className="w-4 h-4" />
                {t("layout.login")}
              </button>
            ) : (
              <div className="relative" ref={menuRef}>
                <button
                  type="button"
                  onClick={() => setMenuOpen((o) => !o)}
                  className={cn(
                    "flex items-center gap-2 rounded-xl px-3 py-2 border text-sm font-semibold",
                    surfaces.card,
                    surfaces.cardBorder
                  )}
                >
                  <UserCircle className="w-8 h-8 text-indigo-400 shrink-0" />
                  <span className="max-w-[120px] truncate hidden sm:inline">
                    {user.username}
                  </span>
                  <ChevronDown className="w-4 h-4 opacity-60" />
                </button>
                {menuOpen && (
                  <div
                    className={cn(
                      "absolute right-0 mt-2 w-52 rounded-xl border shadow-xl py-1 z-[60]",
                      surfaces.card,
                      surfaces.cardBorder
                    )}
                  >
                    <button
                      type="button"
                      className="w-full text-left px-4 py-2.5 text-sm hover:bg-white/10 flex items-center gap-2"
                      onClick={() => {
                        setMenuOpen(false);
                        setActiveTab("profile");
                      }}
                    >
                      <User className="w-4 h-4" />
                      {t("layout.profile")}
                    </button>
                    <button
                      type="button"
                      className="w-full text-left px-4 py-2.5 text-sm hover:bg-white/10 flex items-center gap-2"
                      onClick={() => {
                        setMenuOpen(false);
                        setActiveTab("inventory");
                      }}
                    >
                      <Package className="w-4 h-4" />
                      {t("layout.inventory")}
                    </button>
                    <button
                      type="button"
                      className="w-full text-left px-4 py-2.5 text-sm hover:bg-white/10 flex items-center gap-2"
                      onClick={() => {
                        setMenuOpen(false);
                        setActiveTab("transactions");
                      }}
                    >
                      <RefreshCw className="w-4 h-4" />
                      {t("layout.transactions")}
                    </button>
                    <button
                      type="button"
                      className="w-full text-left px-4 py-2.5 text-sm text-red-400 hover:bg-red-500/10 flex items-center gap-2"
                      onClick={() => {
                        setMenuOpen(false);
                        logout();
                        setActiveTab("store");
                      }}
                    >
                      <LogOut className="w-4 h-4" />
                      {t("layout.logout")}
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 py-8">{children}</main>

      {['store', 'forge', 'exchange', 'daily'].includes(activeTab) && (
        <nav
          className={cn(
            "fixed bottom-0 w-full backdrop-blur-md border-t pb-safe z-40",
            surfaces.nav,
            surfaces.navBorder
          )}
        >
          <div className="max-w-md mx-auto flex justify-around p-2">
            {tabs.map(({ id, label, icon: Icon }) => (
              <button
                key={id}
                type="button"
                onClick={() => setActiveTab(id)}
                className={cn(
                  "flex flex-col items-center p-2 rounded-xl transition-all duration-200",
                  activeTab === id
                    ? "text-indigo-400 bg-indigo-500/10 scale-110"
                    : cn(surfaces.textMuted, "hover:text-slate-200 hover:bg-white/5")
                )}
              >
                <Icon className="w-6 h-6 mb-1" />
                <span className="text-[10px] uppercase tracking-wider font-semibold">
                  {label}
                </span>
              </button>
            ))}
          </div>
        </nav>
      )}
    </div>
  );
}
