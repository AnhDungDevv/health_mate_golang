package usecase

import (
	"context"
	"health_backend/internal/domain/entity"
)

type ChatKafkaUseCase interface {
	// Kafka (Streaming Messages)
	ProduceMessageToKafka(ctx context.Context, message *entity.Message) error
	ConsumeMessagesFromKafka(ctx context.Context) error
}
