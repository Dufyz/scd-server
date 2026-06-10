package usecases

import (
	"database/sql"

	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/errors"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
)

type MessageUsecase struct {
	repository interfaces.MessageRepositoryInterface
}

func NewMessageUsecase(
	repository interfaces.MessageRepositoryInterface,
) MessageUsecase {
	return MessageUsecase{
		repository: repository,
	}
}

func (uc *MessageUsecase) buildResponse(message entities.Message) dtos.MessageResponse {
	return dtos.MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Message:   message.Message,
		UserName:  message.UserName,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func (uc *MessageUsecase) ListByChatId(chatId int64) ([]dtos.MessageResponse, error) {
	messages, err := uc.repository.ListByChatId(chatId)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responses = append(responses, uc.buildResponse(message))
	}

	return responses, nil
}

func (uc *MessageUsecase) Create(body dtos.CreateMessage) (dtos.MessageResponse, error) {
	message, err := uc.repository.Create(body)
	if err != nil {
		return dtos.MessageResponse{}, err
	}

	return uc.buildResponse(message), nil
}

func (uc *MessageUsecase) Update(id int64, body dtos.UpdateMessage) (dtos.MessageResponse, error) {
	message, err := uc.repository.Update(id, body)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.MessageResponse{}, errors.ErrMessageNotFound
		}

		return dtos.MessageResponse{}, err
	}

	return uc.buildResponse(message), nil
}

func (uc *MessageUsecase) Delete(id int64) error {
	err := uc.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
