"use client";

import { useState } from "react";
import { useChats } from "@/hooks/useChats";
import { ChatList } from "@/components/chat/ChatList";
import { CreateChatModal } from "@/components/chat/CreateChatModal";
import type { CreateChatPayload, ChatResponse } from "@/types/chat";

export default function HomePage() {
  const [showModal, setShowModal] = useState(false);
  const [search, setSearch] = useState("");
  const { chats, loading, error, addChat } = useChats();

  const filtered = search.trim()
    ? chats.filter((c: ChatResponse) =>
        c.name.toLowerCase().includes(search.toLowerCase()) ||
        c.category.toLowerCase().includes(search.toLowerCase())
      )
    : chats;

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

      <div className="mb-6">
        <input
          type="text"
          value={search}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSearch(e.target.value)
          }
          placeholder="Buscar salas por nome ou categoria..."
          className="w-full rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm shadow-sm outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
      </div>

      <ChatList chats={filtered} loading={loading} error={error} />

      {showModal && (
        <CreateChatModal
          onClose={() => setShowModal(false)}
          onSubmit={handleCreate}
        />
      )}
    </main>
  );
}
