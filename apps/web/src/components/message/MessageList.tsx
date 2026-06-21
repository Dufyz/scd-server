"use client";

import { useEffect, useRef } from "react";
import type { MessageResponse } from "@/types/message";
import { MessageItem } from "./MessageItem";

interface MessageListProps {
  messages: MessageResponse[];
  loading: boolean;
  currentUser: string;
}

export function MessageList({ messages, loading, currentUser }: MessageListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll para a última mensagem
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-500 border-t-transparent" />
      </div>
    );
  }

  if (messages.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center text-gray-400">
        <p>Nenhuma mensagem ainda. Seja o primeiro a escrever!</p>
      </div>
    );
  }

  return (
    <div className="flex flex-1 flex-col gap-3 overflow-y-auto px-4 py-4">
      {messages.map((msg) => (
        <MessageItem key={msg.id} message={msg} currentUser={currentUser} />
      ))}
      <div ref={bottomRef} />
    </div>
  );
}
