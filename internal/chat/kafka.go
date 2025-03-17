package chat

import (
	"context"
	"health_backend/internal/models"
)

type ChatKafkaUseCase interface {
	// Kafka (Streaming Messages)
	ProduceMessageToKafka(ctx context.Context, message *models.Message) error
	ConsumeMessagesFromKafka(ctx context.Context) error
}
