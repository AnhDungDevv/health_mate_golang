package kafka_chat

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

var (
	instance *KafkaProducer
	once     sync.Once
)

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	once.Do(func() {
		instance = &KafkaProducer{
			Writer: &kafka.Writer{
				Addr:     kafka.TCP(brokers...),
				Topic:    topic,
				Balancer: &kafka.LeastBytes{},
			},
		}
	})
	return instance
}

func (p *KafkaProducer) PublishMessage(key string, message []byte) error {
	return p.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: message,
	})
}

// Close Kafka writer
func (p *KafkaProducer) Close() error {
	return p.Writer.Close()
}

// GetWriter returns the kafka writer
func (p *KafkaProducer) GetWriter() *kafka.Writer {
	return p.Writer
}
