package chat

import (
	"context"
	"health_backend/internal/models"

	uuid "github.com/gofrs/uuid"
)

type UseCase interface {
	// Conversation management
	GetConversationID(ctx context.Context, from uuid.UUID, to uuid.UUID) (uuid.UUID, error)
	GetConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	GetConversations(ctx context.Context) ([]*models.Conversation, error)
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error

	// Message management
	SendMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *models.Message) error
	UnsendMessage(ctx context.Context, messageID uuid.UUID) error

	// Unread message management
	GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*models.Message, error)
	SaveUnreadMessage(ctx context.Context, message *models.Message) error

	// User status management
	SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error
	GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error)
	NotifyUserOnline(ctx context.Context, userID string) error
}
