package queue

import (
	"context"
	"database/sql"

	redisLib "github.com/Dufyz/scd-server/infra/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func NewAsynqServerQueue(redisAddr string, db *sql.DB) error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues:      map[string]int{},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				retries, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				taskID, _ := asynq.GetTaskID(ctx)
				zap.L().Error("Task processing failed",
					zap.String("task_type", task.Type()),
					zap.String("task_id", taskID),
					zap.Int("retry_count", retries),
					zap.Int("max_retry", maxRetry),
					zap.Error(err))
			}),
		},
	)

	_ = redisLib.GetClient()

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	_ = NewQueueService(queueClient)

	mux := asynq.NewServeMux()

	if err := srv.Run(mux); err != nil {
		zap.L().Error("Error on NewAsynqServerQueue", zap.Error(err))
		return err
	}

	return nil
}
