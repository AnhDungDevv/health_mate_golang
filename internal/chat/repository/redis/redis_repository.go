package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"health_backend/internal/chat"
	baseredis "health_backend/internal/infrastructure/redis"
	"health_backend/pkg/logger"

	uuid "github.com/gofrs/uuid"
)

const (
	Online  = 1
	Offline = 0
)

type chatRedisRepository struct {
	baseRedisRepository *baseredis.BaseRedisRepository
	logger              logger.Logger
}

func NewChatRedisRepository(base *baseredis.BaseRedisRepository, log logger.Logger) chat.RedisRepository {
	return &chatRedisRepository{
		baseRedisRepository: base,
		logger:              log,
	}
}

// Triển khai phương thức GetConversationID
func (r *chatRedisRepository) GetConversationID(ctx context.Context, key string) (uuid.UUID, error) {
	val, err := r.baseRedisRepository.Get(ctx, key)
	if err != nil || val == "" {
		return uuid.Nil, err
	}

	conversationID, err := uuid.FromString(val)
	if err != nil {
		return uuid.Nil, err
	}

	return conversationID, nil
}

func (r *chatRedisRepository) SetConversationID(ctx context.Context, key string, conversationID uuid.UUID, expiration time.Duration) error {
	err := r.baseRedisRepository.Set(ctx, key, conversationID.String(), expiration)
	if err != nil {
		r.logger.Warn("Unable to cache to Redis", err)
		return err
	}
	return nil
}
func (c *chatRedisRepository) SetUserOnline(ctx context.Context, clientID string, isOnline bool) error {
	status := Offline
	if isOnline {
		status = Online
	}
	// Using the Set method from BaseRedisRepository to set the user's online status
	key := "online:" + clientID
	err := c.baseRedisRepository.Set(ctx, key, status, 10*time.Minute)
	if err != nil {
		c.logger.Error("Error updating user online status", err, clientID)
		return err
	}
	c.logger.Info("User online status updated", clientID, status)
	return nil
}

// GetUserOnline retrieves the user's online status from Redis
// GetUserOnline retrieves the user's online status from Redis
func (c *chatRedisRepository) GetUserOnline(ctx context.Context, clientID string) (bool, error) {
	key := "online:" + clientID

	// Gọi Get từ BaseRedisRepository (trả về string)
	statusStr, err := c.baseRedisRepository.Get(ctx, key)
	if err != nil {
		c.logger.Error("Error retrieving user online status", err, clientID)
		return false, err
	}

	// Kiểm tra nếu Redis trả về chuỗi rỗng hoặc null
	if statusStr == "" {
		c.logger.Warn("User online status is empty, treating as offline", clientID)
		return false, nil // Mặc định là offline
	}

	// Chuyển đổi string -> int
	status, convErr := strconv.Atoi(statusStr)
	if convErr != nil {
		c.logger.Error("Error converting user online status", convErr, clientID, "Value:", statusStr)
		return false, convErr
	}

	// Return true nếu status == Online
	return status == Online, nil
}

func (c *chatRedisRepository) DeleteUser(ctx context.Context, clientID string) error {
	// Construct the key for the user's online status
	key := "online:" + clientID

	// Using the Delete method from BaseRedisRepository to remove the user's online status from Redis
	err := c.baseRedisRepository.Delete(ctx, key)
	if err != nil {
		c.logger.Error("Failed to delete user from Redis", err, clientID)
		return err
	}

	c.logger.Info("User deleted from Redis", clientID)
	return nil
}

func (r *chatRedisRepository) StoreMissedMessage(ctx context.Context, clientID string, message []byte) error {
	// Define a key to store the missed messages for each client
	missedMessagesKey := fmt.Sprintf("missed_messages:%s", clientID)

	// Push the message into a Redis list (using LPUSH to add the message to the left of the list)
	err := r.baseRedisRepository.LPush(ctx, missedMessagesKey, message)
	if err != nil {
		return fmt.Errorf("failed to store missed message in Redis: %v", err)
	}

	// Optionally, set an expiration time for the key to avoid indefinite storage
	err = r.baseRedisRepository.Expire(ctx, missedMessagesKey, time.Hour*24) // Set to expire after 24 hours
	if err != nil {
		return fmt.Errorf("failed to set expiration for missed messages: %v", err)
	}

	return nil
}

func (r *chatRedisRepository) GetMissedMessages(ctx context.Context, clientID string) ([][]byte, error) {
	// Define a key to retrieve missed messages for each client
	missedMessagesKey := fmt.Sprintf("missed_messages:%s", clientID)

	// Get all the messages from the Redis list (LRANGE retrieves the list)
	messages, err := r.baseRedisRepository.LRange(ctx, missedMessagesKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to get missed messages from Redis: %v", err)
	}

	// Convert the string messages to byte arrays
	var result [][]byte
	for _, msg := range messages {
		result = append(result, []byte(msg))
	}

	return result, nil
}

func (r *chatRedisRepository) RemoveMissedMessages(ctx context.Context, clientID string) error {
	// Define a key to remove the missed messages for each client
	missedMessagesKey := fmt.Sprintf("missed_messages:%s", clientID)

	// Delete the Redis list (removes all missed messages for this client)
	err := r.baseRedisRepository.Delete(ctx, missedMessagesKey)
	if err != nil {
		return fmt.Errorf("failed to remove missed messages from Redis: %v", err)
	}

	return nil
}

// GetConversationID
