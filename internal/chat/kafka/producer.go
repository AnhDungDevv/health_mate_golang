package kafka_chat

import (
	"context"
	"encoding/json"
	"fmt"
	"health_backend/internal/chat"
	kafka_config "health_backend/internal/infrastructure/kafka"
	"health_backend/internal/models"
	"health_backend/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaProducer(logger logger.Logger) chat.KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafka_config.KafkaBrokers...),
		Topic:    kafka_config.TopicChat,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{kafkaWriter: writer, log: logger}
}

func (p *KafkaProducer) ProduceMessageToKafka(ctx context.Context, message *models.Message) error {
	keyBytes := message.ReceiverID.Bytes()

	// Chuyển đổi message thành JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		p.log.Errorf("failed to marshal message: %v", err)
		return err
	}

	err = p.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   keyBytes,
		Value: messageJSON, // Gửi JSON thay vì string raw
	})
	if err != nil {
		p.log.Errorf("failed to write message to kafka: %v", err)
	}

	return err
}

func (p *KafkaProducer) NotifyUserOnline(ctx context.Context, userID string) error {
	// Tạo thông điệp JSON
	messageData := fmt.Sprintf(`{"user_id": "%s", "status": "online"}`, userID)

	message := kafka.Message{
		Key:   []byte(userID),
		Value: []byte(messageData),
	}

	p.log.Info("Notify user online:", messageData)
	err := p.kafkaWriter.WriteMessages(ctx, message)
	if err != nil {
		p.log.Errorf("Failed to notify user online: %v", err)
	}
	return err
}

func (p *KafkaProducer) Close() error {
	return p.kafkaWriter.Close()
}
