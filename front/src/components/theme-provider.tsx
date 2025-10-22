import * as React from "react";

type Theme = "light" | "dark" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeContextValue = {
  theme: Theme;
  resolvedTheme: "light" | "dark";
  setTheme: (theme: Theme) => void;
};

const ThemeContext = React.createContext<ThemeContextValue | undefined>(
  undefined
);

const getSystemTheme = (): "light" | "dark" => {
  if (
    typeof window !== "undefined" &&
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    return "dark";
  }
  return "light";
};

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "ui-theme",
}: ThemeProviderProps) {
  const [theme, setThemeState] = React.useState<Theme>(() => {
    if (typeof window === "undefined") {
      return defaultTheme;
    }
    const stored = window.localStorage.getItem(storageKey) as Theme | null;
    return stored ?? defaultTheme;
  });

  const resolved = React.useMemo(() => {
    if (theme === "system") {
      return getSystemTheme();
    }
    return theme;
  }, [theme]);

  const applyTheme = React.useCallback(
    (next: "light" | "dark") => {
      const root = window.document.documentElement;
      root.classList.remove("light", "dark");
      root.classList.add(next);
    },
    []
  );

  React.useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    applyTheme(resolved);
    window.localStorage.setItem(storageKey, theme);
  }, [resolved, theme, storageKey, applyTheme]);

  React.useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const listener = () => {
      if (theme === "system") {
        applyTheme(media.matches ? "dark" : "light");
      }
    };
    media.addEventListener("change", listener);
    return () => media.removeEventListener("change", listener);
  }, [theme, applyTheme]);

  const handleSetTheme = React.useCallback(
    (next: Theme) => {
      setThemeState(next);
    },
    []
  );

  const value = React.useMemo<ThemeContextValue>(
    () => ({ theme, resolvedTheme: resolved, setTheme: handleSetTheme }),
    [theme, resolved, handleSetTheme]
  );

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const context = React.useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within ThemeProvider");
  }
  return context;
}
