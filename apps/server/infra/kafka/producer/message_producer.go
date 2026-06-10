package producer

import (
	"context"

	"github.com/Dufyz/scd-server/internal/shared/dtos"
)

type MessageEventPayload struct {
	Action  string               `json:"action"`
	Message dtos.MessageResponse `json:"message"`
}

type MessageDeleteEventPayload struct {
	Action  string           `json:"action"`
	Message map[string]int64 `json:"message"`
}

type MessageEvent struct {
	Type    string              `json:"type"`
	Payload MessageEventPayload `json:"payload"`
}

type MessageDeleteEvent struct {
	Type    string                    `json:"type"`
	Payload MessageDeleteEventPayload `json:"payload"`
}

func ProduceMessageCreated(ctx context.Context, message dtos.MessageResponse) error {
	event := MessageEvent{
		Type: "message",
		Payload: MessageEventPayload{
			Action:  "create",
			Message: message,
		},
	}
	return publishToKafka(ctx, "message", "create", event)
}

func ProduceMessageUpdated(ctx context.Context, message dtos.MessageResponse) error {
	event := MessageEvent{
		Type: "message",
		Payload: MessageEventPayload{
			Action:  "update",
			Message: message,
		},
	}
	return publishToKafka(ctx, "message", "update", event)
}

func ProduceMessageDeleted(ctx context.Context, id int64) error {
	event := MessageDeleteEvent{
		Type: "message",
		Payload: MessageDeleteEventPayload{
			Action:  "delete",
			Message: map[string]int64{"id": id},
		},
	}
	return publishToKafka(ctx, "message", "delete", event)
}
