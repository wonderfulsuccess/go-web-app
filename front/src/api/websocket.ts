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
const pendingMessages: Array<Omit<WSMessage, "timestamp">> = [];

function getSocketUrl(): string {
  if (cachedUrl) {
    return cachedUrl;
  }
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const base = `${protocol}://${window.location.host}`;
  const clientId =
    typeof window.crypto !== "undefined" && "randomUUID" in window.crypto
      ? window.crypto.randomUUID()
      : Math.random().toString(36).slice(2);
  cachedUrl = `${base}/api/ws?clientId=${encodeURIComponent(clientId)}`;
  if (import.meta.env.DEV) {
    console.debug("[websocket] resolved socket url:", cachedUrl);
  }
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
    return socket;
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
    flushPending();
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
    socket = null;
    scheduleReconnect();
  });

  socket.addEventListener("error", () => {
    socket?.close();
  });

  return socket;
}

export function subscribeToMessages(listener: Listener) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

export function sendMessage<T>(message: Omit<WSMessage<T>, "timestamp">) {
  const currentSocket = connectWebSocket();
  if (!currentSocket || currentSocket.readyState !== WebSocket.OPEN) {
    pendingMessages.push(message as Omit<WSMessage, "timestamp">);
    return;
  }

  sendNow(currentSocket, message);
}

function sendNow(socketInstance: WebSocket, message: Omit<WSMessage, "timestamp">) {
  const payload = {
    ...message,
    timestamp: new Date().toISOString(),
  };
  socketInstance.send(JSON.stringify(payload));
}

function flushPending() {
  const activeSocket = socket;
  if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) {
    return;
  }

  while (pendingMessages.length > 0) {
    const message = pendingMessages.shift();
    if (!message) {
      continue;
    }
    sendNow(activeSocket, message);
  }
}
