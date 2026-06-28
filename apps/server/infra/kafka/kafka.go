package kafka

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	writers   = make(map[string]*kafka.Writer)
	writersMu sync.Mutex
)

func getWriter(brokers []string, topic string) *kafka.Writer {
	key := strings.Join(brokers, ",") + "/" + topic
	writersMu.Lock()
	defer writersMu.Unlock()
	if w, ok := writers[key]; ok {
		return w
	}
	w := NewWriter(brokers, topic)
	writers[key] = w
	return w
}

func NewWriter(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func NewReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
}

func Produce(ctx context.Context, brokers []string, topic string, key, value []byte) error {
	w := getWriter(brokers, topic)

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	return w.WriteMessages(ctx, msg)
}

func Consume(ctx context.Context, brokers []string, topic, groupID string, handler func(kafka.Message) error) error {
	r := NewReader(brokers, topic, groupID)
	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if handler != nil {
			if err := handler(m); err != nil {
				continue
			}
		}
	}
}

func BrokersFromEnv() []string {
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return []string{"localhost:9094"}
}
