"use client";

import { useEffect, useRef, useCallback } from "react";
import { getSocketManager } from "@/lib/socket";
import type { WsIncomingEvent, WsOutgoingPayload } from "@/types/message";

interface UseSocketOptions {
  chatId: number;
  userName: string;
  onMessage: (event: WsIncomingEvent) => void;
}

export function useSocket({ chatId, userName, onMessage }: UseSocketOptions) {
  const managerRef = useRef(getSocketManager());

  useEffect(() => {
    const manager = managerRef.current;

    // Conecta ao socket-server
    manager.connect();

    // Entra na sala correspondente ao chat
    const joinRoom = () => {
      manager.send({ type: "join_room", room_id: chatId });
    };

    // Aguarda conexão aberta antes de entrar na sala
    const interval = setInterval(() => {
      if (manager.isConnected) {
        joinRoom();
        clearInterval(interval);
      }
    }, 200);

    // Registra handler de mensagens recebidas
    const unsubscribe = manager.onMessage(onMessage);

    return () => {
      clearInterval(interval);
      unsubscribe();
    };
  }, [chatId, onMessage]);

  const sendMessage = useCallback(
    (message: string) => {
      managerRef.current.send({
        type: "send_message",
        room_id: chatId,
        message,
        user_name: userName,
      });
    },
    [chatId, userName]
  );

  const sendTyping = useCallback(
    (typing: boolean) => {
      managerRef.current.send({
        type: "typing",
        room_id: chatId,
        user_name: userName,
        typing,
      });
    },
    [chatId, userName]
  );

  const sendPayload = useCallback((payload: WsOutgoingPayload) => {
    managerRef.current.send(payload);
  }, []);

  return { sendMessage, sendTyping, sendPayload };
}
