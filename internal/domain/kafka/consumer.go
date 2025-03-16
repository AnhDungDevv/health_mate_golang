package kafka

import "context"

type Consumer interface {
	Subscribe(topics []string) error
	Consume(ctx context.Context, handler func(message []byte) error)
	Close() error
}
