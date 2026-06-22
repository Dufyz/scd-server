package dtos

import "time"

type MessageResponse struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Message   string    `json:"message"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Language  *string 	`json:"language"`
}

type CreateMessage struct {
	ChatID   int64  `json:"chat_id" validate:"required"`
	Message  string `json:"message" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
}

type UpdateMessage struct {
	Message string `json:"message" validate:"required"`
}
