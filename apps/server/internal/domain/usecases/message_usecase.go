package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"

	redisInfra "github.com/Dufyz/scd-server/infra/redis"
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
	ctx := context.Background()
	key := "messages:list:chat:" + strconv.FormatInt(chatId, 10)

	if exists, _ := redisInfra.Exists(ctx, key); exists {
		if cached, err := redisInfra.Get(ctx, key); err == nil && cached != "" {
			var cachedResp []dtos.MessageResponse
			if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
				return cachedResp, nil
			}
		}
	}

	messages, err := uc.repository.ListByChatId(chatId)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responses = append(responses, uc.buildResponse(message))
	}

	ttl := 60
	if v := os.Getenv("REDIS_TTL_SECONDS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			ttl = parsed
		}
	}

	if b, err := json.Marshal(responses); err == nil {
		_ = redisInfra.Set(ctx, key, string(b), ttl)
	}

	return responses, nil
}

func (uc *MessageUsecase) Create(body dtos.CreateMessage) (dtos.MessageResponse, error) {
	message, err := uc.repository.Create(body)
	if err != nil {
		return dtos.MessageResponse{}, err
	}

	_ = redisInfra.DelByPattern(context.Background(), "messages:list*")

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

	_ = redisInfra.DelByPattern(context.Background(), "messages:list*")

	return uc.buildResponse(message), nil
}

func (uc *MessageUsecase) Delete(id int64) error {
	err := uc.repository.Delete(id)
	if err != nil {
		return err
	}

	_ = redisInfra.DelByPattern(context.Background(), "messages:list*")

	return nil
}
