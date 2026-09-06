package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/TimurGaliev44/notifyx/internal/domain"
)

type Producer struct {
	producer sarama.SyncProducer
}

func New(producer sarama.SyncProducer) *Producer {
	return &Producer{producer: producer}
}

func (p *Producer) Publish(ctx context.Context, notification *domain.Notification) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}
	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	msg := &sarama.ProducerMessage{
		Topic:     "notifications.incoming",
		Value:     sarama.ByteEncoder(data),
		Key:       sarama.StringEncoder(notification.IdempotencyKey),
		Partition: -1,
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}
