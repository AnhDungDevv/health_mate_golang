package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

var (
	instance *KafkaProducer
	once     sync.Once
)

// Khởi tạo Producer Singleton
func NewKafkaProducer(topic string) *KafkaProducer {
	once.Do(func() {
		instance = &KafkaProducer{
			Writer: &kafka.Writer{
				Addr:     kafka.TCP(KafkaBrokers...),
				Topic:    topic,
				Balancer: &kafka.LeastBytes{},
			},
		}
	})
	return instance
}

// Publish Message
func (p *KafkaProducer) PublishMessage(key string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: message,
	})
}

// Đóng Producer
func (p *KafkaProducer) Close() error {
	return p.Writer.Close()
}
