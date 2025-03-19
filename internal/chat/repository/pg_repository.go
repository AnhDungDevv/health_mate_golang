package repository

import (
	"context"
	"health_backend/internal/chat"
	"health_backend/internal/models"

	uuid "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type chatRepo struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) chat.Repository {
	return &chatRepo{
		db: db,
	}
}

// CreateConversation implements chat.Repository.
func (c *chatRepo) CreateConversation(ctx context.Context, conversation models.Conversation) error {
	panic("unimplemented")
}

// GetConversationByUsers implements chat.Repository.
func (c *chatRepo) GetConversationByUsers(ctx context.Context, user1ID uuid.UUID, user2ID uuid.UUID) (models.Conversation, error) {
	panic("unimplemented")
}

// GetConversations implements chat.Repository.
func (c *chatRepo) GetConversations(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	panic("unimplemented")
}

// GetLastMessage implements chat.Repository.
func (c *chatRepo) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (models.Message, error) {
	panic("unimplemented")
}

// GetMessages implements chat.Repository.
func (c *chatRepo) GetMessages(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error) {
	panic("unimplemented")
}

// GetUnreadMessages implements chat.Repository.
func (c *chatRepo) GetUnreadMessages(ctx context.Context, receiverID uuid.UUID) ([]models.Message, error) {
	panic("unimplemented")
}

// MarkMessageAsRead implements chat.Repository.
func (c *chatRepo) MarkMessageAsRead(ctx context.Context, messageID uuid.UUID, readerID uuid.UUID) error {
	panic("unimplemented")
}

// SaveMessage implements chat.Repository.
func (c *chatRepo) SaveMessage(ctx context.Context, message models.Message) error {
	panic("unimplemented")
}
