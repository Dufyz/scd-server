// Espelha dtos.MessageResponse do apps/server
export interface MessageResponse {
  id: number;
  chat_id: number;
  message: string;
  user_name: string;
  created_at: string;
  updated_at: string;
  language: string | null;
}

// Espelha dtos.CreateMessage do apps/server
export interface CreateMessagePayload {
  chat_id: number;
  message: string;
  user_name: string;
}

// Espelha dtos.UpdateMessage do apps/server
export interface UpdateMessagePayload {
  message: string;
}

// Eventos WebSocket recebidos do socket-server
export interface WsMessageEvent {
  type: "send_message";
  room_id: string;
  id?: number;
  message: string;
  user_name: string;
}

export interface WsTypingEvent {
  type: "typing";
  room_id: string;
  user_name: string;
  typing: boolean;
}

export interface WsMessageUpdateEvent {
  type: "message_action_update";
  room_id: string;
  id: number;
  language: string;
}

export type WsIncomingEvent = WsMessageEvent | WsTypingEvent | WsMessageUpdateEvent;

// Eventos WebSocket enviados ao socket-server
export interface WsJoinRoomPayload {
  type: "join_room";
  room_id: number;
}

export interface WsSendMessagePayload {
  type: "send_message";
  room_id: number;
  message: string;
  user_name: string;
}

export interface WsTypingPayload {
  type: "typing";
  room_id: number;
  user_name: string;
  typing: boolean;
}

export type WsOutgoingPayload =
  | WsJoinRoomPayload
  | WsSendMessagePayload
  | WsTypingPayload;
