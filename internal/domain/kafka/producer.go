package kafka

import "context"

type Producer interface {
	PublishMessage(ctx context.Context, topic string, message interface{}) error
	Close() error
}
