package chat

import (
	"context"
	"health_backend/internal/models"
)

type KafkaProducer interface {
	ProduceMessageToKafka(ctx context.Context, message *models.Message) error
	NotifyUserOnline(ctx context.Context, userID string) error
	Close() error
}

type KafkaConsumer interface {
}
