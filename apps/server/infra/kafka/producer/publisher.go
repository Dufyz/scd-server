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
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = kafka.Produce(ctxWithTimeout, brokers, topic, nil, b)
	if err != nil {
		zap.L().Error("kafka: failed to produce message", zap.String("topic", topic), zap.Error(err))
		return fmt.Errorf("kafka produce error: %w", err)
	}

	zap.L().Debug("kafka: produced event", zap.String("topic", topic), zap.String("action", action))
	return nil
}
