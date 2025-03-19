package kafka_chat

import (
	"context"
	"encoding/json"
	"fmt"
	"health_backend/internal/chat"
	"health_backend/internal/chat/delivery/websocket"
	"health_backend/pkg/logger"
	"io"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader         *kafka.Reader
	log            logger.Logger
	cancel         context.CancelFunc
	ctx            context.Context
	redisClient    chat.RedisRepository
	chatRepository chat.Repository
}

func NewKafkaConsumerChat(brokers []string, topic, groupID string, logger logger.Logger, redisClient chat.RedisRepository, chatRepository chat.Repository) *KafkaConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		ctx:    ctx,
		log:    logger,
		cancel: cancel,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: time.Second,
			MaxWait:        time.Second * 3,
			RetentionTime:  time.Hour * 24,
		}),
		redisClient:    redisClient,
		chatRepository: chatRepository,
	}
}

func (c *KafkaConsumer) StartConsumer() {
	defer c.reader.Close()

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info("Kafka consumer stopping...")
			return

		default:
			msg, err := c.reader.ReadMessage(c.ctx)
			if err != nil {
				c.log.Error("Kafka read error:", err)

				if isFatalKafkaError(err) {
					return
				}

				time.Sleep(2 * time.Second)
				continue
			}

			// Handle message with panic recovery
			go c.handleMessage(msg)
		}
	}
}
func (c *KafkaConsumer) handleMessage(msg kafka.Message) {
	defer func() {
		if r := recover(); r != nil {
			c.log.Error(fmt.Sprintf("Recovered from panic in handleMessage: %v", r))
		}
	}()

	c.log.Info("Received Kafka Message:", string(msg.Value))

	var event struct {
		UserID string `json:"user_id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(msg.Value, &event); err == nil {
		if event.Status == "online" {
			c.log.Info("User is online again:", event.UserID)
			c.SendMissedMessagesFromKafka(c.ctx, event.UserID)
			return
		}
	}

	var chatMsg websocket.ChatMessage
	if err := json.Unmarshal(msg.Value, &chatMsg); err != nil {
		c.log.Error("Failed to unmarshal chat message:", err)
		return
	}

	err := c.redisClient.StoreMissedMessage(c.ctx, chatMsg.To, msg.Value)
	if err != nil {
		c.log.Error("Failed to store missed message in Redis:", err)
	}
}

func isFatalKafkaError(err error) bool {
	// Check for specific Kafka errors
	if err == io.EOF ||
		err == context.Canceled ||
		err == context.DeadlineExceeded {
		return true
	}
	// Check for network errors
	if _, ok := err.(net.Error); ok {
		return true
	}
	return false
}

func (c *KafkaConsumer) Close() {
	c.cancel()
	if err := c.reader.Close(); err != nil {
		c.log.Error("Error closing Kafka reader:", err)
	}
}

// Send missed messages when the user comes online
func (c *KafkaConsumer) SendMissedMessagesFromKafka(ctx context.Context, clientID string) {
	// Get missed messages for the client from Redis or another storage
	messages, err := c.redisClient.GetMissedMessages(c.ctx, clientID)
	if err != nil {
		c.log.Error("Failed to get missed messages from Redis:", err)
		return
	}

	// Send missed messages to WebSocket
	for _, message := range messages {
		websocket.SendToClient(clientID, message)
	}

	// Optionally, remove the missed messages after they've been sent
	err = c.redisClient.RemoveMissedMessages(c.ctx, clientID)
	if err != nil {
		c.log.Error("Failed to remove missed messages from Redis:", err)
	}
}
