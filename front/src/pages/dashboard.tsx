import { useEffect, useState } from "react";
import { FiActivity, FiMessageSquare, FiUsers } from "react-icons/fi";

import { type WSMessage, subscribeToMessages } from "@/api/websocket";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const METRICS = [
  {
    title: "在线用户",
    value: "128",
    description: "当前保持 websocket 连接的总数",
    icon: <FiUsers className="h-5 w-5 text-primary" />,
  },
  {
    title: "今日请求",
    value: "4,621",
    description: "后台 API 已处理的 HTTP 请求",
    icon: <FiActivity className="h-5 w-5 text-primary" />,
  },
  {
    title: "消息吞吐",
    value: "892/min",
    description: "近一分钟的 websocket 消息",
    icon: <FiMessageSquare className="h-5 w-5 text-primary" />,
  },
];

function DashboardPage() {
  const [messages, setMessages] = useState<WSMessage[]>([]);

  useEffect(() => {
    const unsubscribe = subscribeToMessages((message) => {
      setMessages((prev) => [message, ...prev].slice(0, 5));
    });
    return unsubscribe;
  }, []);

  return (
    <div className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {METRICS.map((metric) => (
          <Card key={metric.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{metric.title}</CardTitle>
              {metric.icon}
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metric.value}</div>
              <p className="text-xs text-muted-foreground">{metric.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>实时消息</CardTitle>
          <CardDescription>
            订阅服务器推送的 websocket 消息，仅展示最近五条。
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {messages.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              暂无消息，等待服务器推送。
            </p>
          ) : (
            <ul className="space-y-3">
              {messages.map((message, index) => (
                <li
                  key={`${message.timestamp}-${index}`}
                  className="rounded-md border bg-card/50 p-3"
                >
                  <div className="flex items-center justify-between text-sm">
                    <span className="font-medium text-foreground">
                      {message.type}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(message.timestamp).toLocaleTimeString()}
                    </span>
                  </div>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {JSON.stringify(message.payload)}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    {message.sender} → {message.receiver || "*"}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export default DashboardPage;
