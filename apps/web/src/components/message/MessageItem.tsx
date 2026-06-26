"use client";

import type { MessageResponse } from "@/types/message";

interface MessageItemProps {
  message: MessageResponse;
  currentUser: string;
}

export function MessageItem({ message, currentUser }: MessageItemProps) {
  const isOwn = message.user_name === currentUser;

  return (
    <div className={`flex ${isOwn ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[75%] rounded-2xl px-4 py-2 shadow-sm ${
          isOwn
            ? "rounded-br-sm bg-indigo-600 text-white"
            : "rounded-bl-sm bg-white text-gray-900"
        }`}
      >
        {!isOwn && (
          <p className="mb-0.5 text-xs font-semibold text-indigo-500">
            {message.user_name}
          </p>
        )}
        <p className="text-sm leading-relaxed">{message.message}</p>
        <div
          suppressHydrationWarning
          className={`mt-1 flex items-center justify-end gap-2 text-[10px] ${
            isOwn ? "text-indigo-200" : "text-gray-400"
          }`}
        >
          {message.language && (
            <span
              className={`rounded px-1 py-0.5 font-medium uppercase tracking-wide ${
                isOwn ? "bg-indigo-500 text-indigo-100" : "bg-gray-100 text-gray-500"
              }`}
            >
              {message.language}
            </span>
          )}
          <span>
            {new Date(message.created_at).toLocaleTimeString("pt-BR", {
              hour: "2-digit",
              minute: "2-digit",
            })}
          </span>
        </div>
      </div>
    </div>
  );
}
