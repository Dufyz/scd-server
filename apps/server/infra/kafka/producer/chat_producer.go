package producer

import (
	"context"

	"github.com/Dufyz/scd-server/internal/shared/dtos"
)

type ChatEventPayload struct {
	Action string            `json:"action"`
	Chat   dtos.ChatResponse `json:"chat"`
}

type ChatDeleteEventPayload struct {
	Action string           `json:"action"`
	Chat   map[string]int64 `json:"chat"`
}

type ChatEvent struct {
	Type    string           `json:"type"`
	Payload ChatEventPayload `json:"payload"`
}

type ChatDeleteEvent struct {
	Type    string                 `json:"type"`
	Payload ChatDeleteEventPayload `json:"payload"`
}

func ProduceChatCreated(ctx context.Context, chat dtos.ChatResponse) error {
	event := ChatEvent{
		Type: "chat",
		Payload: ChatEventPayload{
			Action: "create",
			Chat:   chat,
		},
	}
	return publishToKafka(ctx, "chat", "create", event)
}

func ProduceChatUpdated(ctx context.Context, chat dtos.ChatResponse) error {
	event := ChatEvent{
		Type: "chat",
		Payload: ChatEventPayload{
			Action: "update",
			Chat:   chat,
		},
	}
	return publishToKafka(ctx, "chat", "update", event)
}

func ProduceChatDeleted(ctx context.Context, id int64) error {
	event := ChatDeleteEvent{
		Type: "chat",
		Payload: ChatDeleteEventPayload{
			Action: "delete",
			Chat:   map[string]int64{"id": id},
		},
	}
	return publishToKafka(ctx, "chat", "delete", event)
}
