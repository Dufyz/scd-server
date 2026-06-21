"use client";

import type { ChatResponse } from "@/types/chat";
import { ChatCard } from "./ChatCard";

interface ChatListProps {
  chats: ChatResponse[];
  loading: boolean;
  error: string | null;
}

export function ChatList({ chats, loading, error }: ChatListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-500 border-t-transparent" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-lg bg-red-50 p-4 text-sm text-red-700">
        <strong>Erro:</strong> {error}
      </div>
    );
  }

  if (chats.length === 0) {
    return (
      <div className="py-16 text-center text-gray-400">
        <p className="text-lg">Nenhuma sala encontrada.</p>
        <p className="text-sm">Crie uma nova sala para começar!</p>
      </div>
    );
  }

  return (
    <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {chats.map((chat) => (
        <li key={chat.id}>
          <ChatCard chat={chat} />
        </li>
      ))}
    </ul>
  );
}
