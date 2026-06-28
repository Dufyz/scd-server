package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

type ChatUsecase struct {
	repository interfaces.ChatRepositoryInterface
	cacheTTL   int
}

func NewChatUsecase(
	repository interfaces.ChatRepositoryInterface,
) ChatUsecase {
	ttl := 60
	if v := os.Getenv("REDIS_TTL_SECONDS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			ttl = parsed
		}
	}
	return ChatUsecase{
		repository: repository,
		cacheTTL:   ttl,
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
	ctx := context.Background()
	name := ""
	category := ""
	if filters.Name != nil {
		name = *filters.Name
	}
	if filters.Category != nil {
		category = *filters.Category
	}

	key := fmt.Sprintf("chats:list:name=%s:category=%s", name, category)
	if cached, err := redisInfra.Get(ctx, key); err == nil && cached != "" {
		var cachedResp []dtos.ChatResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return cachedResp, nil
		}
	}

	chats, err := uc.repository.List(filters)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		responses = append(responses, uc.buildResponse(chat))
	}

	if b, err := json.Marshal(responses); err == nil {
		_ = redisInfra.Set(ctx, key, string(b), uc.cacheTTL)
	}

	return responses, nil
}

func (uc *ChatUsecase) Create(body dtos.CreateChat) (dtos.ChatResponse, error) {
	chat, err := uc.repository.Create(body)
	if err != nil {
		return dtos.ChatResponse{}, err
	}

	go redisInfra.DelByPattern(context.Background(), "chats:list*")

	resp := uc.buildResponse(chat)
	go func() {
		if err := kafkaProducer.ProduceChatCreated(context.Background(), resp); err != nil {
			zap.L().Error("failed to publish chat.created event", zap.Error(err))
		}
	}()

	return resp, nil
}

func (uc *ChatUsecase) Update(id int64, body dtos.UpdateChat) (dtos.ChatResponse, error) {
	chat, err := uc.repository.Update(id, body)
	if err != nil {
		if err == sql.ErrNoRows {
			return dtos.ChatResponse{}, errors.ErrChatNotFound
		}

		return dtos.ChatResponse{}, err
	}

	go redisInfra.DelByPattern(context.Background(), "chats:list*")

	resp := uc.buildResponse(chat)
	go func() {
		if err := kafkaProducer.ProduceChatUpdated(context.Background(), resp); err != nil {
			zap.L().Error("failed to publish chat.updated event", zap.Error(err))
		}
	}()

	return resp, nil
}

func (uc *ChatUsecase) Delete(id int64) error {
	err := uc.repository.Delete(id)
	if err != nil {
		return err
	}

	go redisInfra.DelByPattern(context.Background(), "chats:list*")

	go func() {
		if err := kafkaProducer.ProduceChatDeleted(context.Background(), id); err != nil {
			zap.L().Error("failed to publish chat.deleted event", zap.Error(err))
		}
	}()

	return nil
}
