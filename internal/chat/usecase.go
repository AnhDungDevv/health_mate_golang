package chat

import (
	"context"
	"health_backend/internal/models"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type UseCase interface {
	// Conversations
	GetConversations(ctx context.Context) ([]*models.Conversation, error)
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error

	// Messages
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error
	UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *models.Message) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	UnsendMessage(ctx context.Context, messageID uuid.UUID) error

	// Websocket
	HandleWebSocket(conn *websocket.Conn)

	// Redis
	SaveUnreadMessage(ctx context.Context, message *models.Message) error
	GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*models.Message, error)
	SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error
	GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error)
}
