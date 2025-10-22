import { useEffect, useMemo, useState } from "react";

import {
  connectWebSocket,
  sendMessage,
  subscribeToMessages,
} from "@/api/websocket";
import type { WSMessage } from "@/api/websocket";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

type DemoEvent = {
  id: string;
  direction: "in" | "out";
  summary: string;
  timestamp: string;
};

const COMPONENT_ID = "websocket-demo";

function formatPayload(payload: unknown): string {
  if (payload && typeof payload === "object") {
    const maybePayload = payload as { message?: unknown };
    if (typeof maybePayload.message === "string") {
      return maybePayload.message;
    }
  }
  if (typeof payload === "string") {
    return payload;
  }
  try {
    return JSON.stringify(payload);
  } catch {
    return String(payload);
  }
}

export default function WebsocketDemo() {
  const [events, setEvents] = useState<DemoEvent[]>([]);
  const [hasRequested, setHasRequested] = useState(false);
  const [isActive, setIsActive] = useState(false);

  useEffect(() => {
    connectWebSocket();
  }, []);

  useEffect(() => {
    const unsubscribe = subscribeToMessages((message: WSMessage) => {
      const incomingEvent: DemoEvent = {
        id: `${message.timestamp}-${message.type}-in`,
        direction: "in",
        summary: `[${message.type}] ${formatPayload(message.payload)}`,
        timestamp: message.timestamp,
      };

      let outgoingEvent: DemoEvent | null = null;

      if (message.type === "server-tick") {
        setIsActive(true);

        const replyPayload = {
          message: `Ack: ${formatPayload(message.payload)}`,
          repliedAt: new Date().toISOString(),
        };

        sendMessage({
          sender: COMPONENT_ID,
          receiver: "server",
          type: "client-ack",
          payload: replyPayload,
        });

        outgoingEvent = {
          id: `${replyPayload.repliedAt}-client-ack-out`,
          direction: "out",
          summary: `[client-ack] ${replyPayload.message}`,
          timestamp: replyPayload.repliedAt,
        };
      }

      setEvents((prev) => {
        const next = [...prev, incomingEvent];
        if (outgoingEvent) {
          return [...next, outgoingEvent];
        }
        return next;
      });
    });

    return unsubscribe;
  }, []);

  const handleStart = () => {
    connectWebSocket();
    sendMessage({
      sender: COMPONENT_ID,
      receiver: "server",
      type: "demo-start",
      payload: {
        message: "start websocket demo broadcast",
        requestedAt: new Date().toISOString(),
      },
    });
    setHasRequested(true);
  };

  const orderedEvents = useMemo(
    () => [...events].sort((a, b) => a.timestamp.localeCompare(b.timestamp)),
    [events]
  );

  const statusHint = isActive
    ? "服务端已经开始推送消息。"
    : hasRequested
      ? "指令已发送，等待服务端首条消息…"
      : "点击按钮后才会请求后端开始推送。";

  return (
    <Card className="max-w-2xl">
      <CardHeader className="space-y-4">
        <div className="space-y-1">
          <CardTitle>WebSocket 演示</CardTitle>
          <CardDescription>
            服务端每秒推送一条消息，组件收到后立即回传确认。
          </CardDescription>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <Button onClick={handleStart} disabled={isActive}>
            {isActive ? "演示进行中" : "开始请求"}
          </Button>
          <p className="text-sm text-muted-foreground">{statusHint}</p>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {orderedEvents.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            等待服务端推送消息…
          </p>
        ) : (
          <ul className="space-y-2 text-sm">
            {orderedEvents.map((event) => (
              <li
                key={event.id}
                className={
                  event.direction === "in"
                    ? "text-emerald-600"
                    : "text-sky-600"
                }
              >
                <span className="font-medium">
                  {event.direction === "in" ? "⬇ 收到" : "⬆ 发送"}
                </span>{" "}
                <span className="text-muted-foreground">
                  {new Date(event.timestamp).toLocaleTimeString()}
                </span>{" "}
                {event.summary}
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
