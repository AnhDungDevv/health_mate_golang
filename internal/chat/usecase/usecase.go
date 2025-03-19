package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"health_backend/config"
	"health_backend/internal/chat"
	"health_backend/internal/chat/delivery/websocket"
	"health_backend/internal/models"
	"health_backend/pkg/logger"
	"sync"
	"time"

	uuid "github.com/gofrs/uuid"
)

type ChatUsecase struct {
	cfg           *config.Config
	pgRepo        chat.Repository
	redisRepo     chat.RedisRepository
	kafkaProducer chat.KafkaProducer
	log           logger.Logger
}

func NewChatUsecase(cfg *config.Config, pg chat.Repository, redis chat.RedisRepository, kafka chat.KafkaProducer, log logger.Logger) chat.UseCase {
	return &ChatUsecase{
		cfg:           cfg,
		pgRepo:        pg,
		redisRepo:     redis,
		kafkaProducer: kafka,
		log:           log,
	}
}

/* ==========================
    CONVERSATION MANAGEMENT
========================== */

// Get conversation ID between two users, create if not exist
func (uc *ChatUsecase) GetConversationID(ctx context.Context, from uuid.UUID, to uuid.UUID) (uuid.UUID, error) {
	redisKey := "conversation:" + from.String() + ":" + to.String()

	var conversationID uuid.UUID
	conversationID, err := uc.redisRepo.GetConversationID(ctx, redisKey)
	if err == nil {
		return conversationID, nil
	}

	if err != nil && err != sql.ErrNoRows {
		uc.log.Error("Error retrieving data from Redis", err)
	}

	conversation, err := uc.pgRepo.GetConversationBetweenUsers(ctx, from, to)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newUUID, err := uuid.NewV4()
			if err != nil {
				return uuid.Nil, fmt.Errorf("failed to generate UUID: %w", err)
			}

			newConversation := models.Conversation{
				ID:        newUUID,
				IsGroup:   false,
				Users:     []models.User{{ID: from}, {ID: to}},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := uc.pgRepo.CreateConversation(ctx, newConversation); err != nil {
				return uuid.Nil, fmt.Errorf("failed to create conversation: %w", err)
			}

			conversationID = newUUID

			cacheDuration := 24 * time.Hour
			if err := uc.redisRepo.SetConversationID(ctx, redisKey, conversationID, cacheDuration); err != nil {
				uc.log.Warn("Failed to cache conversation in Redis", err)
			}

			return conversationID, nil
		}

		return uuid.Nil, err
	}

	conversationID = conversation.ID

	cacheDuration := 24 * time.Hour
	if err := uc.redisRepo.SetConversationID(ctx, redisKey, conversationID, cacheDuration); err != nil {
		uc.log.Warn("Failed to cache conversation in Redis", err)
	}

	return conversationID, nil
}

func (uc *ChatUsecase) GetConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	panic("unimplemented")
}

func (uc *ChatUsecase) GetConversations(ctx context.Context) ([]*models.Conversation, error) {
	panic("unimplemented")
}

func (uc *ChatUsecase) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	panic("unimplemented")
}

func (uc *ChatUsecase) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	panic("unimplemented")
}

/* ==========================
    MESSAGE MANAGEMENT
========================== */

func (uc *ChatUsecase) SendMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error {
	isOnline, err := uc.redisRepo.GetUserOnline(ctx, message.ReceiverID.String())
	if err != nil {
		uc.log.Error("Error checking online status:", err)
		return err
	}

	var wg sync.WaitGroup

	if isOnline {
		msgBytes, err := json.Marshal(message)
		if err != nil {
			uc.log.Error("Error encoding message:", err)
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			websocket.SendToClient(message.ReceiverID.String(), msgBytes)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := uc.redisRepo.StoreMissedMessage(ctx, message.ReceiverID.String(), msgBytes); err != nil {
				uc.log.Error("Error storing message in Redis:", err)
			}
		}()

		wg.Wait()
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := uc.kafkaProducer.ProduceMessageToKafka(ctx, message); err != nil {
			uc.log.Error("Error sending message to Kafka:", err)
		} else {
			uc.log.Info(fmt.Sprintf("Message from %s to %s pushed to Kafka", message.SenderID, message.ReceiverID))
		}
	}()

	wg.Wait()
	return nil
}

func (uc *ChatUsecase) GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	panic("unimplemented")
}

func (uc *ChatUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	panic("unimplemented")
}

func (uc *ChatUsecase) UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *models.Message) error {
	panic("unimplemented")
}

func (uc *ChatUsecase) UnsendMessage(ctx context.Context, messageID uuid.UUID) error {
	panic("unimplemented")
}

/* ==========================
    UNREAD MESSAGES MANAGEMENT
========================== */

func (uc *ChatUsecase) GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*models.Message, error) {
	panic("unimplemented")
}

func (uc *ChatUsecase) SaveUnreadMessage(ctx context.Context, message *models.Message) error {
	panic("unimplemented")
}

/* ==========================
  USER STATUS MANAGEMENT
========================== */

func (uc *ChatUsecase) SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error {
	err := uc.redisRepo.SetUserOnline(ctx, userID.String(), status)
	if err != nil {
		uc.log.Error("Failed to set user online status in Redis", err, "UserID:", userID)
		return err
	}

	uc.log.Info("User online status updated in Redis", "UserID:", userID, "Status:", status)
	return nil
}

func (uc *ChatUsecase) GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error) {
	panic("unimplemented")
}

/* ==========================
    NOTIFICATIONS & REAL-TIME EVENTS
========================== */

func (uc *ChatUsecase) NotifyUserOnline(ctx context.Context, userID string) error {
	return uc.kafkaProducer.NotifyUserOnline(ctx, userID)
}
