package kafka

import (
	"health_backend/internal/infrastructure/kafka"
)

var ChatProducer = kafka.NewKafkaProducer(kafka.TopicChat)
