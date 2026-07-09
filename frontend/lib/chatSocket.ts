"use client";

import { ChatMessage } from "./types";
import { API_BASE_URL } from "./api";

type MessageCallback = (msg: ChatMessage) => void;

const getWebSocketUrl = () => {
  if (typeof window !== "undefined") {
    const url = API_BASE_URL;
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
    try {
      if (url.startsWith("http://") || url.startsWith("https://")) {
        const parsedUrl = new URL(url);
        const host = (parsedUrl.hostname === "localhost" || parsedUrl.hostname === "127.0.0.1")
          ? (window.location.hostname === "localhost" ? "127.0.0.1" : window.location.hostname)
          : parsedUrl.hostname;

        const isLocal = parsedUrl.hostname === "localhost" || parsedUrl.hostname === "127.0.0.1";
        const finalPort = isLocal ? ":8080" : (parsedUrl.port ? `:${parsedUrl.port}` : "");
        return `${proto}//${host}${finalPort}/chat`;
      }
    } catch (e) {
      console.error("[WebSocket] Failed to parse API_BASE_URL, falling back:", e);
    }
    const fallbackHost = window.location.hostname === "localhost" ? "127.0.0.1" : window.location.hostname;
    return `${proto}//${fallbackHost}:8080/chat`;
  }
  return "ws://localhost:8080/chat";
};

class ChatSocketService {
  private socket: WebSocket | null = null;
  private channel: BroadcastChannel | null = null;
  private listeners: Set<MessageCallback> = new Set();
  private currentUserId: string | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      this.connect();
    }
  }

  public authenticate(userId: string) {
    this.currentUserId = userId;
    console.log("[WebSocket] Authenticating with userId:", userId);
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ type: "auth", userId }));
    }
  }

  private connect() {
    try {
      const wsUrl = getWebSocketUrl();
      console.log("[WebSocket] Connecting to:", wsUrl);
      this.socket = new WebSocket(wsUrl);
      
      this.socket.onopen = () => {
        console.log("[WebSocket] Connection established!");
        if (this.currentUserId) {
          console.log("[WebSocket] Sending auth token on connect for:", this.currentUserId);
          this.socket?.send(JSON.stringify({ type: "auth", userId: this.currentUserId }));
        }
      };

      this.socket.onmessage = (event) => {
        try {
          const msg: ChatMessage = JSON.parse(event.data);
          console.log("[WebSocket] Message received:", msg);
          this.notify(msg);
        } catch (e) {
          console.error("[WebSocket] Failed to parse message", e);
        }
      };

      this.socket.onerror = (e) => {
        console.error("[WebSocket] Error occurred:", e);
        this.setupBroadcastChannelFallback();
      };

      this.socket.onclose = (e) => {
        console.warn("[WebSocket] Connection closed:", e);
        this.setupBroadcastChannelFallback();
      };
    } catch (e) {
      this.setupBroadcastChannelFallback();
    }
  }

  private setupBroadcastChannelFallback() {
    if (this.channel) return;
    this.channel = new BroadcastChannel("rent_ws_chat");
    this.channel.onmessage = (event) => {
      const msg: ChatMessage = event.data;
      this.notify(msg);
    };
  }

  public sendMessage(msg: ChatMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(msg));
    }
    if (this.channel) {
      this.channel.postMessage(msg);
    }
  }

  public subscribe(callback: MessageCallback) {
    this.listeners.add(callback);
    return () => {
      this.listeners.delete(callback);
    };
  }

  private notify(msg: ChatMessage) {
    this.listeners.forEach((listener) => listener(msg));
  }
}

export const chatSocket = new ChatSocketService();
