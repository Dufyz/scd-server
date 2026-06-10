package dtos

import "time"

type ChatResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatFilters struct {
	Name     *string
	Category *string
}

type CreateChat struct {
	Name     string `json:"name" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type UpdateChat struct {
	Name     string `json:"name" validate:"required"`
	Category string `json:"category" validate:"required"`
}
