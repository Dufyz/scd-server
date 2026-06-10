package queue

import (
	"time"

	"github.com/hibiken/asynq"
)

const taskRetention = 24 * time.Hour
const maxRetries = 5

type QueueService struct {
	client *asynq.Client
}

func NewQueueService(client *asynq.Client) *QueueService {
	return &QueueService{
		client: client,
	}
}
