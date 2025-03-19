package chat

import (
	"context"
	"health_backend/internal/models"

	uuid "github.com/gofrs/uuid"
)

type Repository interface {
	// Conversation
	GetConversations(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error)
	CreateConversation(ctx context.Context, conversation models.Conversation) error
	GetConversationByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (models.Conversation, error)

	// Message
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error)
	SaveMessage(ctx context.Context, message models.Message) error
	MarkMessageAsRead(ctx context.Context, messageID uuid.UUID, readerID uuid.UUID) error
	GetUnreadMessages(ctx context.Context, receiverID uuid.UUID) ([]models.Message, error)

	// Lấy tin nhắn cuối cùng (Hiển thị preview ở danh sách chat)
	GetLastMessage(ctx context.Context, conversationID uuid.UUID) (models.Message, error)
}
