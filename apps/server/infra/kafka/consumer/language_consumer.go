package consumer

import (
	"context"
	"encoding/json"

	kafkaProducer "github.com/Dufyz/scd-server/infra/kafka/producer"
	redisInfra "github.com/Dufyz/scd-server/infra/redis"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// LanguageDetectedEvent represents the event payload from ai-server
type LanguageDetectedEvent struct {
	MessageID        int64  `json:"message_id"`
	ChatID           int64  `json:"chat_id"`
	DetectedLanguage string `json:"detected_language"`
	OriginalMessage  string `json:"original_message"`
}

// StartLanguageDetectionConsumer starts consuming language detection events
func StartLanguageDetectionConsumer(
	ctx context.Context,
	reader *kafka.Reader,
	messageRepo interfaces.MessageRepositoryInterface,
) error {
	zap.L().Info("Starting language detection consumer",
		zap.String("topic", reader.Config().Topic),
		zap.String("groupID", reader.Config().GroupID),
	)

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("Language detection consumer stopped")
			return reader.Close()
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				zap.L().Error("Error fetching message from Kafka", zap.Error(err))
				continue
			}

			if err := handleLanguageDetectedEvent(ctx, msg, messageRepo); err != nil {
				zap.L().Error("Error handling language detected event",
					zap.Error(err),
					zap.String("message_key", string(msg.Key)),
				)
				// Continue processing next message even if this one failed
				if err := reader.CommitMessages(ctx, msg); err != nil {
					zap.L().Error("Error committing message offset", zap.Error(err))
				}
				continue
			}

			// Commit message offset after successful processing
			if err := reader.CommitMessages(ctx, msg); err != nil {
				zap.L().Error("Error committing message offset", zap.Error(err))
			}
		}
	}
}

func handleLanguageDetectedEvent(
	ctx context.Context,
	msg kafka.Message,
	messageRepo interfaces.MessageRepositoryInterface,
) error {
	var event LanguageDetectedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		zap.L().Error("Error unmarshaling language detected event", zap.Error(err))
		return err
	}

	zap.L().Info("Processing language detected event",
		zap.Int64("message_id", event.MessageID),
		zap.Int64("chat_id", event.ChatID),
		zap.String("language", event.DetectedLanguage),
	)

	// Update message language in database
	_, err := messageRepo.UpdateLanguage(event.MessageID, event.DetectedLanguage)
	if err != nil {
		zap.L().Error("Error updating message language in database",
			zap.Error(err),
			zap.Int64("message_id", event.MessageID),
			zap.String("language", event.DetectedLanguage),
		)
		return err
	}

	// Publish language update event to Kafka so socket-server can broadcast to clients
	lang := event.DetectedLanguage
	if err := kafkaProducer.ProduceMessageLanguageUpdated(ctx, dtos.MessageResponse{
		ID:       event.MessageID,
		ChatID:   event.ChatID,
		Language: &lang,
	}); err != nil {
		zap.L().Warn("Error publishing language update event to Kafka",
			zap.Error(err),
			zap.Int64("message_id", event.MessageID),
		)
		// Don't return error as this is not critical
	}

	// Invalidate cache for this chat's messages
	if err := redisInfra.DelByPattern(ctx, "messages:list*"); err != nil {
		zap.L().Warn("Error invalidating cache after language update",
			zap.Error(err),
			zap.Int64("chat_id", event.ChatID),
		)
	}

	zap.L().Info("Successfully updated message language",
		zap.Int64("message_id", event.MessageID),
		zap.String("language", event.DetectedLanguage),
	)

	return nil
}
