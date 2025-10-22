import { useMemo } from "react";
import { FiMonitor, FiMoon, FiSun } from "react-icons/fi";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useTheme } from "@/components/theme-provider";

const THEME_OPTIONS = [
  {
    label: "跟随系统",
    value: "system" as const,
    icon: <FiMonitor className="h-4 w-4" />,
  },
  { label: "亮色", value: "light" as const, icon: <FiSun className="h-4 w-4" /> },
  { label: "暗色", value: "dark" as const, icon: <FiMoon className="h-4 w-4" /> },
];

function SettingsPage() {
  const { theme, resolvedTheme, setTheme } = useTheme();

  const activeLabel = useMemo(() => {
    const current = THEME_OPTIONS.find((option) => option.value === theme);
    if (current) {
      return current.label;
    }
    return resolvedTheme === "dark" ? "暗色" : "亮色";
  }, [theme, resolvedTheme]);

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>外观设置</CardTitle>
          <CardDescription>
            支持亮暗主题切换，也可以跟随操作系统自动调整。
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground">
            当前主题：<span className="font-medium text-foreground">{activeLabel}</span>
          </p>
          <div className="flex flex-wrap gap-3">
            {THEME_OPTIONS.map((option) => (
              <Button
                key={option.value}
                variant={option.value === theme ? "default" : "outline"}
                onClick={() => setTheme(option.value)}
              >
                <span className="mr-2">{option.icon}</span>
                {option.label}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>关于模板</CardTitle>
          <CardDescription>
            一套基于 Go + React 的桌面应用后台快速开发起始模板。
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          <p>• Gin 提供 HTTP API，并与前端构建产物做静态资源整合。</p>
          <p>• Gorm 负责数据库访问，默认使用 SQLite，可一行切换其它驱动。</p>
          <p>• WebSocket 统一入口，前端可订阅实时数据。</p>
        </CardContent>
      </Card>
    </div>
  );
}

export default SettingsPage;
