import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { FiMoon, FiSun } from "react-icons/fi";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();

  const handleToggle = () => {
    setTheme(resolvedTheme === "dark" ? "light" : "dark");
  };

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={handleToggle}
      aria-label="Toggle theme"
    >
      {resolvedTheme === "dark" ? (
        <FiMoon className="h-5 w-5" />
      ) : (
        <FiSun className="h-5 w-5" />
      )}
    </Button>
  );
}
