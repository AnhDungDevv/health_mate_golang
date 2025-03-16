package usecase

import (
	"context"
	"health_backend/internal/domain/entity"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ChatUseCase interface {
	// Conversations
	GetConversations(ctx context.Context) ([]*entity.Conversation, error)
	CreateConversation(ctx context.Context, conversation *entity.Conversation) error
	GetConversation(ctx context.Context, conversationID uuid.UUID) (*entity.Conversation, error)
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error

	// Messages
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*entity.Message, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, message *entity.Message) error
	UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *entity.Message) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	UnsendMessage(ctx context.Context, messageID uuid.UUID) error

	// Websocket
	HandleWebSocket(conn *websocket.Conn)

	// Redis
	SaveUnreadMessage(ctx context.Context, message *entity.Message) error
	GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*entity.Message, error)
	SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error
	GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error)
}
