// Espelha dtos.ChatResponse do apps/server
export interface ChatResponse {
  id: number;
  name: string;
  category: string;
  created_at: string;
  updated_at: string;
}

// Espelha dtos.CreateChat do apps/server
export interface CreateChatPayload {
  name: string;
  category: string;
}

// Espelha dtos.UpdateChat do apps/server
export interface UpdateChatPayload {
  name: string;
  category: string;
}

// Filtros de query para GET /api/chats
export interface ChatFilters {
  name?: string;
  category?: string;
}
