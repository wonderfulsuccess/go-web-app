import { useCallback, useEffect, useState } from "react";
import { FiPlus, FiRefreshCcw } from "react-icons/fi";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  createdAt: string;
}

function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch("/api/users");
      if (!response.ok) {
        throw new Error(`请求失败: ${response.status}`);
      }
      const data = (await response.json()) as User[];
      setUsers(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "无法加载用户，请确认后端服务已启动。"
      );
    } finally {
      setLoading(false);
    }
  }, []);

  const createUser = useCallback(async () => {
    try {
      const payload = {
        name: `访客 ${Math.random().toString(36).slice(2, 6)}`,
        email: `${Date.now()}@example.com`,
        role: "viewer",
      };
      const response = await fetch("/api/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        throw new Error(`创建失败: ${response.status}`);
      }
      const created = (await response.json()) as User;
      setUsers((prev) => [created, ...prev]);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "无法创建用户，请确认后端服务已启动。"
      );
    }
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>用户管理</CardTitle>
            <CardDescription>查看和维护后台数据库中的用户。</CardDescription>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={fetchUsers}
              disabled={loading}
            >
              <FiRefreshCcw className="mr-2 h-4 w-4" /> 刷新
            </Button>
            <Button size="sm" onClick={createUser}>
              <FiPlus className="mr-2 h-4 w-4" /> 新建用户
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-sm text-muted-foreground">加载中...</p>
          ) : error ? (
            <p className="text-sm text-destructive">{error}</p>
          ) : users.length === 0 ? (
            <p className="text-sm text-muted-foreground">暂无用户数据。</p>
          ) : (
            <div className="overflow-hidden rounded-md border">
              <table className="min-w-full divide-y divide-border text-sm">
                <thead className="bg-muted/60">
                  <tr>
                    <th className="px-4 py-2 text-left font-medium">姓名</th>
                    <th className="px-4 py-2 text-left font-medium">邮箱</th>
                    <th className="px-4 py-2 text-left font-medium">角色</th>
                    <th className="px-4 py-2 text-left font-medium">
                      创建时间
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {users.map((user) => (
                    <tr key={user.id}>
                      <td className="px-4 py-2 font-medium text-foreground">
                        {user.name}
                      </td>
                      <td className="px-4 py-2 text-muted-foreground">
                        {user.email}
                      </td>
                      <td className="px-4 py-2 text-muted-foreground">
                        {user.role}
                      </td>
                      <td className="px-4 py-2 text-muted-foreground">
                        {new Date(user.createdAt).toLocaleString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export default UsersPage;
