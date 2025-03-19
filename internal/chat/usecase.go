package chat

import (
	"context"
	"health_backend/internal/models"

	uuid "github.com/gofrs/uuid"
)

type UseCase interface {
	// Conversations
	GetConversations(ctx context.Context) ([]*models.Conversation, error)
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationID(ctx context.Context, from, to uuid.UUID) (uuid.UUID, error)
	GetConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error

	// Messages
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error
	UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *models.Message) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	UnsendMessage(ctx context.Context, messageID uuid.UUID) error

	NotifyUserOnline(ctx context.Context, userID string) error
	SaveUnreadMessage(ctx context.Context, message *models.Message) error
	GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*models.Message, error)
	SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error
	GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error)
}
