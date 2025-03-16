package kafka

import "github.com/IBM/sarama"

type kafkaConsumer struct {
	consumer sarama.Consumer
	group    sarama.ConsumerGroup
}

func NewKafkaConsumer(brokers []string, groupID string) (domain.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{group: group}, nil
}
