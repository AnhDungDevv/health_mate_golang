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

func (c *chatRepo) CreateConversation(ctx context.Context, conversation models.Conversation) error {
	// Create ID if not have
	if conversation.ID == uuid.Nil {
		conversation.ID = uuid.Must(uuid.NewV4())
	}

	// Create new conversation
	if err := c.db.WithContext(ctx).Create(conversation).Error; err != nil {
		return err
	}

	return nil
}

// GetConversationByUsers implements chat.Repository.
func (c *chatRepo) GetConversationBetweenUsers(ctx context.Context, user1ID uuid.UUID, user2ID uuid.UUID) (models.Conversation, error) {
	var conversation models.Conversation

	err := c.db.Joins("JOIN conversation_users cu1 ON cu1.conversation_id = conversations.id").
		Joins("JOIN conversation_users cu2 ON cu2.conversation_id = conversations.id").
		Where("cu1.user_id = ? AND cu2.user_id = ? AND conversations.is_group = FALSE", user1ID, user2ID).
		Group("conversations.id").
		Having("COUNT(cu1.user_id) = 2").
		First(&conversation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Conversation{}, nil
		}
		return models.Conversation{}, err
	}

	return conversation, nil

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
