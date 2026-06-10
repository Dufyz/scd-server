package usecases

import (
	"database/sql"

	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/errors"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
)

type ChatUsecase struct {
	repository interfaces.ChatRepositoryInterface
}

func NewChatUsecase(
	repository interfaces.ChatRepositoryInterface,
) ChatUsecase {
	return ChatUsecase{
		repository: repository,
	}
}

func (uc *ChatUsecase) buildResponse(chat entities.Chat) dtos.ChatResponse {
	return dtos.ChatResponse{
		ID:        chat.ID,
		Name:      chat.Name,
		Category:  chat.Category,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}

func (uc *ChatUsecase) List(filters dtos.ChatFilters) ([]dtos.ChatResponse, error) {
	chats, err := uc.repository.List(filters)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		responses = append(responses, uc.buildResponse(chat))
	}

	return responses, nil
}

func (uc *ChatUsecase) Create(body dtos.CreateChat) (dtos.ChatResponse, error) {
	chat, err := uc.repository.Create(body)
	if err != nil {
		return dtos.ChatResponse{}, err
	}

	return uc.buildResponse(chat), nil
}

func (uc *ChatUsecase) Update(id int64, body dtos.UpdateChat) (dtos.ChatResponse, error) {
	chat, err := uc.repository.Update(id, body)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.ChatResponse{}, errors.ErrChatNotFound
		}

		return dtos.ChatResponse{}, err
	}

	return uc.buildResponse(chat), nil
}

func (uc *ChatUsecase) Delete(id int64) error {
	err := uc.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
