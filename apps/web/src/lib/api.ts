/**
 * Cliente HTTP tipado para o apps/server (Go/Echo)
 * Base URL: NEXT_PUBLIC_API_URL (ex: http://localhost:3000/api)
 */

import type {
  ChatResponse,
  CreateChatPayload,
  UpdateChatPayload,
  ChatFilters,
} from "@/types/chat";
import type {
  MessageResponse,
  CreateMessagePayload,
  UpdateMessagePayload,
} from "@/types/message";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000/api";

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options?.headers ?? {}),
    },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(body?.message ?? `HTTP ${res.status}`);
  }

  // DELETE /chats/:id e DELETE /messages/:id retornam { message }
  const text = await res.text();
  return text ? (JSON.parse(text) as T) : ({} as T);
}

// ─── Chats ────────────────────────────────────────────────────────────────────

export async function listChats(filters?: ChatFilters): Promise<ChatResponse[]> {
  const params = new URLSearchParams();
  if (filters?.name) params.set("name", filters.name);
  if (filters?.category) params.set("category", filters.category);
  const qs = params.toString() ? `?${params.toString()}` : "";
  return request<ChatResponse[]>(`/chats${qs}`);
}

export async function getChatMessages(chatId: number): Promise<MessageResponse[]> {
  return request<MessageResponse[]>(`/chats/${chatId}/messages`);
}

export async function createChat(payload: CreateChatPayload): Promise<ChatResponse> {
  return request<ChatResponse>("/chats", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function updateChat(
  id: number,
  payload: UpdateChatPayload
): Promise<ChatResponse> {
  return request<ChatResponse>(`/chats/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deleteChat(id: number): Promise<void> {
  await request<{ message: string }>(`/chats/${id}`, { method: "DELETE" });
}

// ─── Messages ─────────────────────────────────────────────────────────────────

export async function createMessage(
  payload: CreateMessagePayload
): Promise<MessageResponse> {
  return request<MessageResponse>("/messages", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function updateMessage(
  id: number,
  payload: UpdateMessagePayload
): Promise<MessageResponse> {
  return request<MessageResponse>(`/messages/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deleteMessage(id: number): Promise<void> {
  await request<{ message: string }>(`/messages/${id}`, { method: "DELETE" });
}
