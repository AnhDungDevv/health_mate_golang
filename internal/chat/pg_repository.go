package chat

import (
	"context"
	"health_backend/internal/models"

	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetConversations(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error)
	CreateConversation(ctx context.Context, conversation models.Conversation) error
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error)
	SendMessage(ctx context.Context, message models.Message) error
}
