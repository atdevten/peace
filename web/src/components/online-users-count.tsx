"use client";

import React, { useEffect, useRef, useState } from "react";
import Cookies from "js-cookie";
import { useAuth } from "@/contexts/auth-context";

type Props = {
  className?: string;
  wsUrl?: string;
};

export default function OnlineUsersCount({ className, wsUrl }: Props) {
  const [count, setCount] = useState<number>(0);
  const [connected, setConnected] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const { isAuthenticated } = useAuth();

  const wsRef = useRef<WebSocket | null>(null);
  const pingTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const countTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const retryAttemptRef = useRef<number>(0);
  const lastTokenRef = useRef<string | null>(null);
  const tokenCheckTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  function applyTimers(ws: WebSocket, opts: { hidden: boolean }) {
    // Ping interval
    if (pingTimerRef.current) clearInterval(pingTimerRef.current);
    const pingMs = opts.hidden ? 30000 : 10000;
    pingTimerRef.current = setInterval(() => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "ping" }));
      }
    }, pingMs);

    // Count interval (skip if hidden)
    if (countTimerRef.current) clearInterval(countTimerRef.current);
    if (!opts.hidden) {
      countTimerRef.current = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: "get_amount_online_users" }));
        }
      }, 30000);
    }
  }

  function isAmountMsg(msg: unknown): msg is { type: "amount_online_users"; data: { count: number } } {
    if (typeof msg !== "object" || msg === null) return false;
    const obj = msg as Record<string, unknown>;
    if (obj["type"] !== "amount_online_users") return false;
    const data = obj["data"];
    if (typeof data !== "object" || data === null) return false;
    const dataObj = data as Record<string, unknown>;
    return typeof dataObj["count"] === "number";
  }

  useEffect(() => {
    const url = wsUrl || process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8081/ws";
    const token = Cookies.get("access_token");
    if (!token) {
      setError("Missing access token");
      return;
    }

    function connect(authToken: string) {
      try {
        const protocols: string[] = ["bearer", authToken];
        const ws = new WebSocket(url, protocols);
        wsRef.current = ws;

        ws.onopen = () => {
          setConnected(true);
          setError(null);
          retryAttemptRef.current = 0; // reset backoff
          // yêu cầu server trả về số lượng online ngay khi kết nối
          ws.send(JSON.stringify({ type: "get_amount_online_users" }));
          // Thiết lập timers theo visibility hiện tại
          applyTimers(ws, { hidden: document.hidden });
        };

        ws.onmessage = (evt) => {
          try {
            const parsed: unknown = JSON.parse(evt.data);
            if (isAmountMsg(parsed)) setCount(parsed.data.count);
          } catch {
            // ignore parse errors
          }
        };

        ws.onerror = () => {
          setError("WebSocket error");
        };

        ws.onclose = () => {
          setConnected(false);
          if (pingTimerRef.current) {
            clearInterval(pingTimerRef.current);
            pingTimerRef.current = null;
          }
          if (countTimerRef.current) {
            clearInterval(countTimerRef.current);
            countTimerRef.current = null;
          }
          // reconnect sau 2s
          if (!reconnectTimerRef.current) {
            const attempt = retryAttemptRef.current;
            const baseDelay = 2000; // 2s
            const maxDelay = 30000; // 30s cap
            const jitter = Math.floor(Math.random() * 1000); // 0-1s jitter
            const delay = Math.min(maxDelay, baseDelay * Math.pow(2, attempt)) + jitter;
            reconnectTimerRef.current = setTimeout(() => {
              reconnectTimerRef.current = null;
              retryAttemptRef.current = Math.min(attempt + 1, 10);
              connect(authToken);
            }, delay);
          }
        };
      } catch {
        setError("Failed to open WebSocket");
      }
    }

    connect(token);
    lastTokenRef.current = token;

    // Poll token thay đổi để reconnect chủ động khi token refresh
    if (tokenCheckTimerRef.current) clearInterval(tokenCheckTimerRef.current);
    tokenCheckTimerRef.current = setInterval(() => {
      const current = Cookies.get("access_token") || null;
      const prev = lastTokenRef.current;
      if (current && prev && current !== prev) {
        lastTokenRef.current = current;
        // Đóng WS để kết nối lại với token mới
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
          wsRef.current.close();
        }
      }
    }, 10000);

    return () => {
      if (pingTimerRef.current) clearInterval(pingTimerRef.current);
      if (countTimerRef.current) clearInterval(countTimerRef.current);
      if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
      if (tokenCheckTimerRef.current) clearInterval(tokenCheckTimerRef.current);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
  }, [wsUrl, isAuthenticated]);

  // Refresh count khi tab lấy lại focus
  useEffect(() => {
    const onFocus = () => {
      const ws = wsRef.current;
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "get_amount_online_users" }));
      }
    };
    window.addEventListener("focus", onFocus);
    return () => window.removeEventListener("focus", onFocus);
  }, []);

  // Tối ưu visibility: tăng ping interval và tắt count khi tab ẩn; khôi phục khi hiện
  useEffect(() => {
    const onVisibility = () => {
      const ws = wsRef.current;
      if (!ws) return;
      // Áp dụng timers mới
      applyTimers(ws, { hidden: document.hidden });
      if (!document.hidden && ws.readyState === WebSocket.OPEN) {
        // Khi hiện lại: cập nhật tức thì
        ws.send(JSON.stringify({ type: "ping" }));
        ws.send(JSON.stringify({ type: "get_amount_online_users" }));
      }
    };
    document.addEventListener("visibilitychange", onVisibility);
    return () => document.removeEventListener("visibilitychange", onVisibility);
  }, []);

  return (
    <span
      className={
        className ||
        "inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700"
      }
      title={connected ? "WebSocket connected" : error || "Connecting..."}
    >
      {count}
    </span>
  );
}


