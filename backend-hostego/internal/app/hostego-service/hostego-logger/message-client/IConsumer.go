package messageclient

import (
	"context"
)

type IConsumer interface {
	// Receive messages from stream
	Receive(ctx context.Context) ([]Message, error)

	// ReceivePendingMessages COnsumes pending messages from stream
	ReceivePendingMessages(ctx context.Context) ([]Message, error)

	// Acknowledge message to stream
	Acknowledge(ctx context.Context, messageId string) error

	// RegisterConsumerGroup Register consumer group
	RegisterConsumerGroup(ctx context.Context) (string, error)

	// Start a consumer to consume messages
	// Stops consumption on receiving ctx.Done()
	Start(ctx context.Context, msgChan chan<- []Message, errChan chan<- error) error
}
