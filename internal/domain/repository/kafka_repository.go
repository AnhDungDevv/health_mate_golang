package repository

import (
	"context"
	"health_backend/internal/domain/entity"
)

type ChatKafkaRepository interface {
	PublishMessage(ctx context.Context, message entity.Message) error
}
