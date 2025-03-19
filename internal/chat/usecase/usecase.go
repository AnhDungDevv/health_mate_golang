package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"health_backend/config"
	"health_backend/internal/chat"
	"health_backend/internal/chat/delivery/websocket"
	"health_backend/internal/models"
	"health_backend/pkg/logger"
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
		pgRepo:        pg, // Không cần &
		redisRepo:     redis,
		kafkaProducer: kafka, // Không cần &
		log:           log,
	}
}

// GetConversationID implements chat.UseCase.
func (uc *ChatUsecase) GetConversationID(ctx context.Context, from uuid.UUID, to uuid.UUID) (uuid.UUID, error) {
	// Tạo key Redis dựa trên hai user
	redisKey := "conversation:" + from.String() + ":" + to.String()

	var conversationID uuid.UUID

	// 1️⃣ Kiểm tra Redis trước
	conversationID, err := uc.redisRepo.GetConversationID(ctx, redisKey)
	if err == nil { // Nếu tìm thấy trong Redis
		return conversationID, nil
	}

	// Nếu Redis không có dữ liệu hoặc lỗi không phải ErrNil, tiếp tục truy vấn DB
	if err != nil && err != sql.ErrNoRows {
		uc.log.Error("Lỗi khi lấy dữ liệu từ Redis", err)
	}

	// 2️⃣ Truy vấn PostgreSQL nếu Redis không có dữ liệu
	// conversationID, err = uc.pgRepo.GetConversationID(ctx, from, to)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return uuid.Nil, ErrConversationNotFound
	// 	}
	// 	return uuid.Nil, err
	// }

	// 3️⃣ Cache lại vào Redis với thời gian hết hạn (ví dụ: 1 ngày)
	cacheDuration := 24 * time.Hour
	err = uc.redisRepo.SetConversationID(ctx, redisKey, conversationID, cacheDuration)
	if err != nil {
		uc.log.Warn("Không thể lưu cache vào Redis", err)
	}

	return conversationID, nil
}
func (uc *ChatUsecase) NotifyUserOnline(ctx context.Context, userID string) error {
	return uc.kafkaProducer.NotifyUserOnline(ctx, userID)
}

func (uc *ChatUsecase) SendMessage(ctx context.Context, conversationID uuid.UUID, message *models.Message) error {
	// Kiểm tra người nhận có online không

	isOnline, err := uc.redisRepo.GetUserOnline(ctx, message.ReceiverID.String())

	if err != nil {
		uc.log.Error("Lỗi kiểm tra trạng thái online:", err)
		return err
	}
	// Nếu người nhận online, gửi qua WebSocket
	if isOnline {
		msgBytes, err := json.Marshal(message)
		if err != nil {
			uc.log.Error("Lỗi mã hóa tin nhắn:", err)
			return err
		}
		websocket.SendToClient(message.ReceiverID.String(), msgBytes)

		// Lưu vào Redis để hỗ trợ việc gửi lại tin nhắn nếu cần
		err = uc.redisRepo.StoreMissedMessage(ctx, message.ReceiverID.String(), msgBytes)
		if err != nil {
			uc.log.Error("Lỗi lưu tin nhắn vào Redis:", err)
		}
		return nil
	}

	// Nếu người nhận offline, gửi tin nhắn vào Kafka
	err = uc.kafkaProducer.ProduceMessageToKafka(ctx, message)
	if err != nil {
		uc.log.Error("Lỗi gửi tin nhắn vào Kafka:", err)
		return err
	}

	uc.log.Info(fmt.Sprintf("Tin nhắn từ %s đến %s đã được đẩy vào Kafka", message.SenderID, message.ReceiverID))
	return nil
}

// CreateConversation implements chat.UseCase.
func (uc *ChatUsecase) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	panic("unimplemented")
}

// DeleteConversation implements chat.UseCase.
func (uc *ChatUsecase) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	panic("unimplemented")
}

// DeleteMessage implements chat.UseCase.
func (uc *ChatUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	panic("unimplemented")
}

// GetConversation implements chat.UseCase.
func (uc *ChatUsecase) GetConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	panic("unimplemented")
}

// GetConversations implements chat.UseCase.
func (uc *ChatUsecase) GetConversations(ctx context.Context) ([]*models.Conversation, error) {
	panic("unimplemented")
}

// GetMessages implements chat.UseCase.
func (uc *ChatUsecase) GetMessages(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	panic("unimplemented")
}

// GetOnlineUsers implements chat.UseCase.
func (uc *ChatUsecase) GetOnlineUsers(ctx context.Context) ([]uuid.UUID, error) {
	panic("unimplemented")
}

// GetUnreadMessages implements chat.UseCase.
func (uc *ChatUsecase) GetUnreadMessages(ctx context.Context, userID uuid.UUID) ([]*models.Message, error) {
	panic("unimplemented")
}

// SaveUnreadMessage implements chat.UseCase.
func (uc *ChatUsecase) SaveUnreadMessage(ctx context.Context, message *models.Message) error {
	panic("unimplemented")
}

func (uc *ChatUsecase) SetUserOnlineStatus(ctx context.Context, userID uuid.UUID, status bool) error {
	// Gọi hàm SetUserOnline từ chatRedisRepository để lưu vào Redis
	err := uc.redisRepo.SetUserOnline(ctx, userID.String(), status)
	if err != nil {
		uc.log.Error("Failed to set user online status in Redis", err, "UserID:", userID)
		return err
	}

	uc.log.Info("User online status updated in Redis", "UserID:", userID, "Status:", status)
	return nil
}

// UnsendMessage implements chat.UseCase.
func (uc *ChatUsecase) UnsendMessage(ctx context.Context, messageID uuid.UUID) error {
	panic("unimplemented")
}

// UpdateMessage implements chat.UseCase.
func (uc *ChatUsecase) UpdateMessage(ctx context.Context, messageID uuid.UUID, updateData *models.Message) error {
	panic("unimplemented")
}
