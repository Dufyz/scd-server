"use client";

import Link from "next/link";
import type { ChatResponse } from "@/types/chat";

interface ChatCardProps {
  chat: ChatResponse;
}

export function ChatCard({ chat }: ChatCardProps) {
  return (
    <Link
      href={`/chat/${chat.id}`}
      className="block rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition hover:border-indigo-400 hover:shadow-md"
    >
      <div className="mb-2 flex items-start justify-between gap-2">
        <h2 className="text-base font-semibold text-gray-900 leading-tight line-clamp-2">
          {chat.name}
        </h2>
        <span className="shrink-0 rounded-full bg-indigo-100 px-2 py-0.5 text-xs font-medium text-indigo-700">
          {chat.category}
        </span>
      </div>
      <p className="text-xs text-gray-400">
        Criado em{" "}
        {new Date(chat.created_at).toLocaleDateString("pt-BR", {
          day: "2-digit",
          month: "short",
          year: "numeric",
        })}
      </p>
    </Link>
  );
}
