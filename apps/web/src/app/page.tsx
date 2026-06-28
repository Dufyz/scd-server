"use client";

import { useState } from "react";
import { useChats } from "@/hooks/useChats";
import { ChatList } from "@/components/chat/ChatList";
import { CreateChatModal } from "@/components/chat/CreateChatModal";
import type { CreateChatPayload } from "@/types/chat";

export default function HomePage() {
  const [showModal, setShowModal] = useState(false);
  const [nameFilter, setNameFilter] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("");

  const { chats, loading, error, addChat } = useChats({
    name: nameFilter || undefined,
    category: categoryFilter || undefined,
  });

  async function handleCreate(payload: CreateChatPayload) {
    await addChat(payload);
  }

  return (
    <main className="mx-auto max-w-5xl px-4 py-8">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">💬 SCD Chat</h1>
          <p className="text-sm text-gray-500">
            Sistemas Concorrentes e Distribuídos
          </p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700"
        >
          + Nova Sala
        </button>
      </div>

      <div className="mb-6 flex gap-3">
        <input
          type="text"
          value={nameFilter}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setNameFilter(e.target.value)
          }
          placeholder="Buscar por nome..."
          className="flex-1 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm shadow-sm outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
        <input
          type="text"
          value={categoryFilter}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setCategoryFilter(e.target.value)
          }
          placeholder="Buscar por categoria..."
          className="flex-1 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm shadow-sm outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
      </div>

      <ChatList chats={chats} loading={loading} error={error} />

      {showModal && (
        <CreateChatModal
          onClose={() => setShowModal(false)}
          onSubmit={handleCreate}
        />
      )}
    </main>
  );
}
