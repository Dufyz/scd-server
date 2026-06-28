package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"

	kafkaProducer "github.com/Dufyz/scd-server/infra/kafka/producer"
	redisInfra "github.com/Dufyz/scd-server/infra/redis"
	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/errors"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
	"go.uber.org/zap"
)

type MessageUsecase struct {
	repository interfaces.MessageRepositoryInterface
	cacheTTL   int
}

func NewMessageUsecase(
	repository interfaces.MessageRepositoryInterface,
) MessageUsecase {
	ttl := 60
	if v := os.Getenv("REDIS_TTL_SECONDS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			ttl = parsed
		}
	}
	return MessageUsecase{
		repository: repository,
		cacheTTL:   ttl,
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
		Language:  message.Language,
	}
}

func (uc *MessageUsecase) ListByChatId(chatId int64) ([]dtos.MessageResponse, error) {
	ctx := context.Background()
	key := "messages:list:chat:" + strconv.FormatInt(chatId, 10)

	if cached, err := redisInfra.Get(ctx, key); err == nil && cached != "" {
		var cachedResp []dtos.MessageResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return cachedResp, nil
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

	if b, err := json.Marshal(responses); err == nil {
		_ = redisInfra.Set(ctx, key, string(b), uc.cacheTTL)
	}

	return responses, nil
}

func (uc *MessageUsecase) Create(body dtos.CreateMessage) (dtos.MessageResponse, error) {
	message, err := uc.repository.Create(body)
	if err != nil {
		if err.Error() == "pq: insert or update on table \"messages\" violates foreign key constraint \"messages_chat_id_fkey\"" {
			return dtos.MessageResponse{}, errors.ErrMessageFKChatId
		}

		return dtos.MessageResponse{}, err
	}

	go redisInfra.DelByPattern(context.Background(), "messages:list*")

	resp := uc.buildResponse(message)
	go func() {
		if err := kafkaProducer.ProduceMessageCreated(context.Background(), resp); err != nil {
			zap.L().Error("failed to publish message.created event", zap.Error(err))
		}
	}()

	return resp, nil
}

func (uc *MessageUsecase) Update(id int64, body dtos.UpdateMessage) (dtos.MessageResponse, error) {
	message, err := uc.repository.Update(id, body)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.MessageResponse{}, errors.ErrMessageNotFound
		}

		return dtos.MessageResponse{}, err
	}

	go redisInfra.DelByPattern(context.Background(), "messages:list*")

	resp := uc.buildResponse(message)
	go func() {
		if err := kafkaProducer.ProduceMessageUpdated(context.Background(), resp); err != nil {
			zap.L().Error("failed to publish message.updated event", zap.Error(err))
		}
	}()

	return resp, nil
}

func (uc *MessageUsecase) Delete(id int64) error {
	err := uc.repository.Delete(id)
	if err != nil {
		return err
	}

	go redisInfra.DelByPattern(context.Background(), "messages:list*")

	go func() {
		if err := kafkaProducer.ProduceMessageDeleted(context.Background(), id); err != nil {
			zap.L().Error("failed to publish message.deleted event", zap.Error(err))
		}
	}()

	return nil
}
