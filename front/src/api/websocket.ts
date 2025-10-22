export type WSMessage<T = unknown> = {
  sender: string;
  receiver: string;
  timestamp: string;
  type: string;
  payload: T;
};

type Listener = (message: WSMessage) => void;

let socket: WebSocket | null = null;
let reconnectTimer: number | null = null;
let reconnectAttempts = 0;
let cachedUrl: string | null = null;
const listeners = new Set<Listener>();

function getSocketUrl(): string {
  if (cachedUrl) {
    return cachedUrl;
  }
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const base = `${protocol}//${window.location.host}`;
  const clientId =
    typeof window.crypto !== "undefined" && "randomUUID" in window.crypto
      ? window.crypto.randomUUID()
      : Math.random().toString(36).slice(2);
  cachedUrl = `${base}/api/ws?clientId=${encodeURIComponent(clientId)}`;
  return cachedUrl;
}

function notify(message: WSMessage) {
  listeners.forEach((listener) => listener(message));
}

function scheduleReconnect() {
  if (reconnectTimer) {
    window.clearTimeout(reconnectTimer);
  }
  const delay = Math.min(1000 * 2 ** reconnectAttempts, 10000);
  reconnectTimer = window.setTimeout(() => {
    reconnectAttempts += 1;
    connectWebSocket();
  }, delay);
}

export function connectWebSocket() {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)
  ) {
    return;
  }

  try {
    socket = new WebSocket(getSocketUrl());
  } catch (error) {
    console.error("Failed to open websocket", error);
    scheduleReconnect();
    return;
  }

  socket.addEventListener("open", () => {
    reconnectAttempts = 0;
  });

  socket.addEventListener("message", (event) => {
    try {
      const data = JSON.parse(event.data) as WSMessage;
      notify(data);
    } catch (error) {
      console.error("Unable to parse websocket message", error);
    }
  });

  socket.addEventListener("close", () => {
    scheduleReconnect();
  });

  socket.addEventListener("error", () => {
    socket?.close();
  });
}

export function subscribeToMessages(listener: Listener) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

export function sendMessage<T>(message: Omit<WSMessage<T>, "timestamp">) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    console.warn("Websocket not ready; skipping send.");
    return;
  }
  const payload = {
    ...message,
    timestamp: new Date().toISOString(),
  } satisfies WSMessage<T>;
  socket.send(JSON.stringify(payload));
}
