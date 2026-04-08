import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";

const STORAGE_KEY = "lootforge_theme";

export type ThemeMode = "light" | "dark";

interface ThemeSurfaces {
  page: string;
  header: string;
  headerBorder: string;
  nav: string;
  navBorder: string;
  textMuted: string;
  card: string;
  cardBorder: string;
}

function surfacesFor(mode: ThemeMode): ThemeSurfaces {
  if (mode === "light") {
    return {
      page: "bg-slate-100 text-slate-900",
      header: "bg-white/90 backdrop-blur-md",
      headerBorder: "border-slate-200",
      nav: "bg-white/95 backdrop-blur-md",
      navBorder: "border-slate-200",
      textMuted: "text-slate-600",
      card: "bg-white",
      cardBorder: "border-slate-200",
    };
  }
  return {
    page: "bg-slate-950 text-slate-100",
    header: "bg-slate-900/80 backdrop-blur-md",
    headerBorder: "border-white/10",
    nav: "bg-slate-900/95 backdrop-blur-md",
    navBorder: "border-white/10",
    textMuted: "text-slate-400",
    card: "bg-slate-900",
    cardBorder: "border-white/10",
  };
}

interface ThemeContextValue {
  mode: ThemeMode;
  toggleTheme: () => void;
  surfaces: ThemeSurfaces;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setMode] = useState<ThemeMode>(() => {
    if (typeof localStorage === "undefined") return "dark";
    const s = localStorage.getItem(STORAGE_KEY);
    return s === "light" ? "light" : "dark";
  });

  useEffect(() => {
    document.documentElement.dataset.theme = mode;
    localStorage.setItem(STORAGE_KEY, mode);
  }, [mode]);

  const toggleTheme = useCallback(() => {
    setMode((m) => (m === "dark" ? "light" : "dark"));
  }, []);

  const value: ThemeContextValue = {
    mode,
    toggleTheme,
    surfaces: surfacesFor(mode),
  };

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used within ThemeProvider");
  return ctx;
}
