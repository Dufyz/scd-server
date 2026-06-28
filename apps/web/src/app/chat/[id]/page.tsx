"use client";

import { useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useMessages } from "@/hooks/useMessages";
import { useSocket } from "@/hooks/useSocket";
import { MessageList } from "@/components/message/MessageList";
import { MessageInput } from "@/components/message/MessageInput";
import { TypingIndicator } from "@/components/message/TypingIndicator";
import type { WsIncomingEvent, MessageResponse } from "@/types/message";

const USER_NAME_KEY = "scd_user_name";

function getUserName(): string {
  if (typeof window === "undefined") return "Anônimo";
  const stored = localStorage.getItem(USER_NAME_KEY);
  if (stored) return stored;
  const generated = `User_${Math.floor(Math.random() * 9000) + 1000}`;
  localStorage.setItem(USER_NAME_KEY, generated);
  return generated;
}

export default function ChatPage() {
  const params = useParams();
  const router = useRouter();
  const chatId = Number(params.id);

  const [userName] = useState<string>(() => getUserName());
  const [typingUsers, setTypingUsers] = useState<string[]>([]);

  const { messages, loading, sendMessage, appendMessage } = useMessages(chatId);

  const handleWsEvent = useCallback(
    (event: WsIncomingEvent) => {
      if (event.type === "send_message") {
        const msg: MessageResponse = {
          id: event.id ?? -Date.now(),
          chat_id: chatId,
          message: event.message,
          user_name: event.user_name,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          language: null,
        };
        appendMessage(msg);
      } else if (event.type === "typing") {
        setTypingUsers((prev: string[]) => {
          if (event.typing) {
            return prev.includes(event.user_name)
              ? prev
              : [...prev, event.user_name];
          }
          return prev.filter((u: string) => u !== event.user_name);
        });
      } else if (event.type === "message_action_update") {
        appendMessage(event);
      }
    },
    [chatId, appendMessage]
  );

  const { sendMessage: wsSend, sendTyping } = useSocket({
    chatId,
    userName,
    onMessage: handleWsEvent,
  });

  const typingOthers = typingUsers.filter((u: string) => u !== userName);

  async function handleSend(text: string) {
    await sendMessage({ chat_id: chatId, message: text, user_name: userName });
  }

  function handleTyping(typing: boolean) {
    sendTyping(typing);
  }

  if (!chatId || isNaN(chatId)) {
    return (
      <div className="flex h-screen items-center justify-center text-gray-500">
        Chat inválido.{" "}
        <Link href="/" className="ml-1 text-indigo-600 underline">
          Voltar
        </Link>
      </div>
    );
  }

  return (
    <div className="flex h-screen flex-col bg-gray-100">
      <header className="flex items-center gap-3 border-b border-gray-200 bg-white px-4 py-3 shadow-sm">
        <button
          onClick={() => router.push("/")}
          className="rounded-lg p-1.5 text-gray-500 hover:bg-gray-100"
          aria-label="Voltar"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2}
            stroke="currentColor"
            className="h-5 w-5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15.75 19.5L8.25 12l7.5-7.5"
            />
          </svg>
        </button>
        <div className="flex-1">
          <h1 className="text-sm font-semibold text-gray-900">
            Sala #{chatId}
          </h1>
          <p className="text-xs text-gray-400">
            Conectado como{" "}
            <span  suppressHydrationWarning className="font-medium text-indigo-600">{userName}</span>
          </p>
        </div>
      </header>

      <MessageList
        messages={messages}
        loading={loading}
        currentUser={userName}
      />

      <TypingIndicator typingUsers={typingOthers} />

      <MessageInput onSend={handleSend} onTyping={handleTyping} />
    </div>
  );
}
