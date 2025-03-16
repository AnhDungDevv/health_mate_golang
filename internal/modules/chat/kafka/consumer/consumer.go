package kafka_chat

import (
	"context"
	"encoding/json"
	"fmt"
	"health_backend/internal/chat/delivery/websocket"
	"health_backend/pkg/logger"
	"io"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	log    logger.Logger
	cancel context.CancelFunc
	ctx    context.Context
}

func NewKafkaConsumer(brokers []string, topic, groupID string, logger logger.Logger) *KafkaConsumer {
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

	// Parse the message
	var chatMsg websocket.ChatMessage
	if err := json.Unmarshal(msg.Value, &chatMsg); err != nil {
		c.log.Error("Failed to unmarshal message:", err)
		return
	}

	// Validate message
	if chatMsg.To == "" {
		c.log.Error("Invalid message: missing recipient")
		return
	}

	// Forward message to recipient
	msgBytes, err := json.Marshal(chatMsg)
	if err != nil {
		c.log.Error("Failed to marshal message for recipient:", err)
		return
	}

	c.log.Info(fmt.Sprintf("Forwarding message from %s to %s", chatMsg.From, chatMsg.To))
	websocket.SendToClient(chatMsg.To, msgBytes)
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
