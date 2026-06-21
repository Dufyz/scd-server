"use client";

import { useState, useEffect, useCallback } from "react";
import { listChats, createChat, updateChat, deleteChat } from "@/lib/api";
import type { ChatResponse, CreateChatPayload, UpdateChatPayload, ChatFilters } from "@/types/chat";

export function useChats(filters?: ChatFilters) {
  const [chats, setChats] = useState<ChatResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchChats = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await listChats(filters);
      setChats(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao carregar chats");
    } finally {
      setLoading(false);
    }
  }, [filters?.name, filters?.category]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    fetchChats();
  }, [fetchChats]);

  const addChat = useCallback(async (payload: CreateChatPayload): Promise<ChatResponse> => {
    const newChat = await createChat(payload);
    setChats((prev) => [newChat, ...prev]);
    return newChat;
  }, []);

  const editChat = useCallback(async (id: number, payload: UpdateChatPayload): Promise<ChatResponse> => {
    const updated = await updateChat(id, payload);
    setChats((prev) => prev.map((c) => (c.id === id ? updated : c)));
    return updated;
  }, []);

  const removeChat = useCallback(async (id: number): Promise<void> => {
    await deleteChat(id);
    setChats((prev) => prev.filter((c) => c.id !== id));
  }, []);

  return { chats, loading, error, refetch: fetchChats, addChat, editChat, removeChat };
}
