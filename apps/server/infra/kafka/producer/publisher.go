package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dufyz/scd-server/infra/kafka"
	"go.uber.org/zap"
)

func publishToKafka(ctx context.Context, topic string, action string, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("kafka: failed to marshal event", zap.Error(err))
		return err
	}

	brokers := kafka.BrokersFromEnv()

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = kafka.Produce(ctxWithTimeout, brokers, topic, nil, b)
		cancel()

		if err == nil {
			lastErr = nil
			break
		}

		lastErr = err
		zap.L().Warn("kafka: produce attempt failed, will retry",
			zap.Int("attempt", attempt), zap.String("topic", topic), zap.Error(err))

		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	if lastErr != nil {
		zap.L().Error("kafka: failed to produce message", zap.String("topic", topic), zap.Error(lastErr))
		return fmt.Errorf("kafka produce error: %w", lastErr)
	}

	zap.L().Debug("kafka: produced event", zap.String("topic", topic), zap.String("action", action))
	return nil
}
