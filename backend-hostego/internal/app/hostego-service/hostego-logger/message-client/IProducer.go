package messageclient

import (
	"context"
)

type IProducer interface {
	Send(ctx context.Context, req *SendRequest) (string, error)
}
