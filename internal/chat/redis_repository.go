package chat

import (
	"context"
	"time"

	uuid "github.com/gofrs/uuid"
)

type RedisRepository interface {
	// Status online
	SetUserOnline(ctx context.Context, clientID string, isOnline bool) error
	GetUserOnline(ctx context.Context, clientID string) (bool, error)
	DeleteUser(ctx context.Context, clientID string) error
	GetConversationID(ctx context.Context, key string) (uuid.UUID, error)
	SetConversationID(ctx context.Context, key string, conversationID uuid.UUID, expiration time.Duration) error
	// Miss message
	StoreMissedMessage(ctx context.Context, clientID string, message []byte) error
	GetMissedMessages(ctx context.Context, clientID string) ([][]byte, error)
	RemoveMissedMessages(ctx context.Context, clientID string) error
}
