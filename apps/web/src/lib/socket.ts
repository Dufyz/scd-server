/**
 * WebSocket manager para o apps/socket-server (Python)
 * URL: NEXT_PUBLIC_SOCKET_URL (ex: ws://localhost:8765)
 *
 * Protocolo: JSON puro sobre WebSocket nativo
 * Eventos enviados: join_room | send_message | typing
 * Eventos recebidos: send_message | typing
 */

import type { WsIncomingEvent, WsOutgoingPayload } from "@/types/message";

export type WsEventHandler = (event: WsIncomingEvent) => void;

export class SocketManager {
  private ws: WebSocket | null = null;
  private url: string;
  private handlers: WsEventHandler[] = [];
  private reconnectDelay = 2000;
  private shouldReconnect = true;

  constructor(url?: string) {
    this.url = url ?? (typeof window !== "undefined"
      ? (process.env.NEXT_PUBLIC_SOCKET_URL ?? "ws://localhost:8765")
      : "ws://localhost:8765");
  }

  connect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) return;

    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log("[Socket] Connected to", this.url);
    };

    this.ws.onmessage = (ev: MessageEvent) => {
      try {
        const data = JSON.parse(ev.data as string) as WsIncomingEvent;
        this.handlers.forEach((h) => h(data));
      } catch {
        console.warn("[Socket] Failed to parse message:", ev.data);
      }
    };

    this.ws.onclose = () => {
      console.log("[Socket] Disconnected");
      if (this.shouldReconnect) {
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };

    this.ws.onerror = (err) => {
      console.error("[Socket] Error:", err);
    };
  }

  disconnect(): void {
    this.shouldReconnect = false;
    this.ws?.close();
    this.ws = null;
  }

  send(payload: WsOutgoingPayload): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn("[Socket] Not connected, cannot send:", payload);
      return;
    }
    this.ws.send(JSON.stringify(payload));
  }

  onMessage(handler: WsEventHandler): () => void {
    this.handlers.push(handler);
    // Retorna função de cleanup
    return () => {
      this.handlers = this.handlers.filter((h) => h !== handler);
    };
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

// Singleton global para reutilizar a conexão entre componentes
let _instance: SocketManager | null = null;

export function getSocketManager(): SocketManager {
  if (!_instance) {
    _instance = new SocketManager();
  }
  return _instance;
}
