import { useEffect } from "react";
import { Navigate, NavLink, Route, Routes } from "react-router-dom";
import { FiActivity, FiSettings, FiUsers } from "react-icons/fi";

import { connectWebSocket } from "@/api/websocket";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import DashboardPage from "@/pages/dashboard";
import SettingsPage from "@/pages/settings";
import UsersPage from "@/pages/users";

const NAVIGATION = [
  { to: "/dashboard", label: "仪表盘", icon: <FiActivity className="h-4 w-4" /> },
  { to: "/users", label: "用户管理", icon: <FiUsers className="h-4 w-4" /> },
  { to: "/settings", label: "系统设置", icon: <FiSettings className="h-4 w-4" /> },
];

function App() {
  useEffect(() => {
    connectWebSocket();
  }, []);

  return (
    <div className="min-h-screen bg-muted/20">
      <header className="border-b bg-background/80 backdrop-blur">
        <div className="container flex h-16 items-center justify-between gap-4">
          <div className="flex items-center gap-6">
            <span className="text-lg font-semibold tracking-tight">Go Desktop Admin</span>
            <nav className="hidden gap-1 md:flex">
              {NAVIGATION.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  className={({ isActive }) =>
                    cn(
                      "flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                      isActive
                        ? "bg-primary text-primary-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground"
                    )
                  }
                >
                  {item.icon}
                  <span>{item.label}</span>
                </NavLink>
              ))}
            </nav>
          </div>
          <div className="flex items-center gap-2">
            <Button asChild variant="outline" size="sm" className="md:hidden">
              <NavLink to="/dashboard">菜单</NavLink>
            </Button>
            <ThemeToggle />
          </div>
        </div>
      </header>
      <main className="container py-6">
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </main>
      <footer className="border-t bg-background">
        <div className="container flex h-14 items-center justify-between text-sm text-muted-foreground">
          <span>© {new Date().getFullYear()} Go Desktop Admin</span>
          <span>构建于 Go + React</span>
        </div>
      </footer>
    </div>
  );
}

export default App;
