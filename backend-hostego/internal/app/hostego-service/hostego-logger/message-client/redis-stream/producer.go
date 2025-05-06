package redis_stream

import (
	"context"
	"errors"
	"fmt"

	messageclient "backend-hostego/internal/app/hostego-service/hostego-logger/message-client"

	"github.com/redis/go-redis/v9"
)

func NewRedisStreamProducerBuilder() *RedisStreamProducerBuilder {
	return &RedisStreamProducerBuilder{
		Producer: &RedisStreamProducer{},
	}
}

type RedisStreamProducerBuilder struct {
	Producer *RedisStreamProducer
}

type RedisStreamProducer struct {
	RedisClient *redis.Client
}

func (r *RedisStreamProducerBuilder) RedisClient(client *redis.Client) *RedisStreamProducerBuilder {
	producer := r.Producer
	producer.RedisClient = client

	return r
}

func (r *RedisStreamProducerBuilder) Build() (*RedisStreamProducer, error) {
	producer := r.Producer

	if producer.RedisClient == nil {
		return nil, errors.New("redis client not found")
	}

	return producer, nil
}

func (r *RedisStreamProducer) Send(ctx context.Context, req *messageclient.SendRequest) (string, error) {
	attrs := make(map[string]interface{})

	attrs["Body"] = req.Body
	for _, attr := range req.Attributes {
		attrs[attr.Key] = attr.Value
	}

	if req.StreamLength == 0 {
		req.StreamLength = 50000
	}

	inputParam := &redis.XAddArgs{
		Stream: req.Stream,
		ID:     "*",
		Values: attrs,
		MaxLen: req.StreamLength,
	}
	msgID, err := r.RedisClient.XAdd(ctx, inputParam).Result()
	if err != nil {
		return "", fmt.Errorf("redis stream publish failure, error : %w", err)
	}
	return msgID, nil
}
