"use client";

import { useState, useEffect, useCallback } from "react";
import { getChatMessages, createMessage } from "@/lib/api";
import type { MessageResponse, CreateMessagePayload } from "@/types/message";

export function useMessages(chatId: number) {
  const [messages, setMessages] = useState<MessageResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchMessages = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getChatMessages(chatId);
      setMessages(data ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao carregar mensagens");
    } finally {
      setLoading(false);
    }
  }, [chatId]);

  useEffect(() => {
    if (chatId) fetchMessages();
  }, [fetchMessages, chatId]);

  const sendMessage = useCallback(
    async (payload: CreateMessagePayload): Promise<MessageResponse> => {
      const apiMsg = await createMessage(payload);
      
      setMessages((prev: MessageResponse[]) => {
        // Se a mensagem já entrou pelo REST (milagre), ignora
        if (prev.some((m) => String(m.id) === String(apiMsg.id))) return prev;
        
        // A FUSÃO: Se a última mensagem da tela tem o mesmo TEXTO e o mesmo USUÁRIO, 
        // e está SEM id (veio do WebSocket), nós substituímos pela oficial do Go!
        const lastMsg = prev[prev.length - 1];
        if (lastMsg && lastMsg.message === apiMsg.message && lastMsg.user_name === apiMsg.user_name && !lastMsg.id) {
            const newArray = [...prev];
            newArray[newArray.length - 1] = apiMsg;
            return newArray;
        }

        return [...prev, apiMsg];
      });
      return apiMsg;
    },
    []
  );

  const appendMessage = useCallback((incomingData: any) => {
    const wsMsg = incomingData.data || incomingData.payload || incomingData;
    
    // Ignora eventos de digitação, queremos apenas as mensagens de texto
    if (wsMsg.type === "typing") return;

    setMessages((prev: MessageResponse[]) => {
      // Se a última mensagem oficial da API já está na tela com o mesmo texto, ignora o WebSocket
      const lastMsg = prev[prev.length - 1];
      if (lastMsg && lastMsg.message === wsMsg.message && lastMsg.user_name === wsMsg.user_name) {
         return prev;
      }
      
      // Coloca a mensagem na tela temporariamente (até o sendMessage a substituir pela oficial com ID)
      return [...prev, wsMsg as MessageResponse];
    });
  }, []);

  return { messages, loading, error, refetch: fetchMessages, sendMessage, appendMessage };
}
