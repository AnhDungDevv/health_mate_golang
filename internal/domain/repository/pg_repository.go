package repository

import (
	"context"
	"health_backend/internal/domain/entity"
	"health_backend/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*entity.UsersList, error)
	FindByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
}

type ChatRepository interface {
	GetConversations(ctx context.Context, userID uuid.UUID) ([]entity.Conversation, error)
	CreateConversation(ctx context.Context, conversation entity.Conversation) error
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]entity.Message, error)
	SendMessage(ctx context.Context, message entity.Message) error
}
