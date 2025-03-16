package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (domain.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{producer: producer}, nil
}

func (p *kafkaProducer) PublishMessage(ctx context.Context, topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}
