package redis_stream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	messageclient "backend-hostego/internal/app/hostego-service/hostego-logger/message-client"

	"github.com/redis/go-redis/v9"
)

func NewRedisStreamConsumerBuilder() *RedisStreamConsumerBuilder {
	return &RedisStreamConsumerBuilder{
		Consumer: &RedisStreamConsumer{},
	}
}

type RedisStreamConsumerBuilder struct {
	Consumer *RedisStreamConsumer
}

type RetryConfig struct {
	retryInterval     time.Duration
	idleTimeThreshold time.Duration
}
type RedisStreamConsumer struct {
	redisClient   *redis.Client
	stream        string
	batchSize     int64
	groupName     string
	consumerName  string
	retryProducer RedisStreamProducer
	blockTime     time.Duration
	retryConfig   RetryConfig
}

/*
Stream Sets stream name from which messages are to be consumed
*/
func (r *RedisStreamConsumerBuilder) Stream(stream string) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.stream = stream

	return r
}

/*
RedisClient Sets RedisClient, to be used for redis stream operations
*/
func (r *RedisStreamConsumerBuilder) RedisClient(client *redis.Client) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.redisClient = client

	return r
}

/*
BatchSize Sets batch size, consumer will consume messages in batches of size <= batchSize
*/
func (r *RedisStreamConsumerBuilder) BatchSize(batchSize int64) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.batchSize = batchSize

	return r
}

/*
GroupName Sets name of consumer group
*/
func (r *RedisStreamConsumerBuilder) GroupName(groupName string) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.groupName = groupName

	return r
}

/*
ConsumerName Sets name of consumer
*/
func (r *RedisStreamConsumerBuilder) ConsumerName(consumerName string) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.consumerName = consumerName

	return r
}

/*
BlockTime Sets blockTime of consumer
*/
func (r *RedisStreamConsumerBuilder) BlockTime(blockTime time.Duration) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.blockTime = blockTime

	return r
}

/*
RetryConfig Sets retry config, messages pending for time more than or equal to idleTimeThreshold will be consumed again, default value is 5s
Polling of pending messages happens ata frequency of retryInterval, default value is 5s
*/
func (r *RedisStreamConsumerBuilder) RetryConfig(idleTimeThreshold time.Duration, retryInterval time.Duration) *RedisStreamConsumerBuilder {
	consumer := r.Consumer
	consumer.retryConfig = RetryConfig{
		idleTimeThreshold,
		retryInterval,
	}

	return r
}

/*
Build Returns a redis stream consumer, validates attributes of consumer, returns error if attributes are invalid
*/
func (r *RedisStreamConsumerBuilder) Build() (*RedisStreamConsumer, error) {
	consumer := r.Consumer
	if consumer.consumerName == "" {
		return nil, errors.New("consumer name cannot be empty")
	}

	if consumer.stream == "" {
		return nil, errors.New("stream cannot be empty")
	}

	if consumer.batchSize <= 0 {
		return nil, fmt.Errorf("invalid batch size %d", consumer.batchSize)
	}

	if consumer.groupName == "" {
		return nil, errors.New("group name cannot be empty")
	}

	if consumer.redisClient == nil {
		return nil, errors.New("redis client not found")
	}

	if consumer.retryConfig.retryInterval < 0 {
		return nil, errors.New("invalid retry interval")
	}
	if consumer.retryConfig.idleTimeThreshold < 0 {
		return nil, errors.New("invalid retry idle time threshold")
	}

	if consumer.retryConfig.retryInterval == 0 {
		consumer.retryConfig.retryInterval = 5 * time.Second
	}

	if consumer.retryConfig.idleTimeThreshold == 0 {
		consumer.retryConfig.idleTimeThreshold = 5 * time.Second
	}

	consumer.retryProducer = RedisStreamProducer{
		RedisClient: consumer.redisClient,
	}

	return consumer, nil
}

/*
Receive Recieves message from stream in batches of batchSize
*/
func (r *RedisStreamConsumer) Receive(ctx context.Context) ([]messageclient.Message, error) {

	inputParam := &redis.XReadGroupArgs{
		Group:    r.groupName,
		Consumer: r.consumerName,
		Streams:  []string{r.stream, ">"},
		Count:    r.batchSize,
		Block:    r.blockTime,
	}

	streamMessages, err := r.redisClient.XReadGroup(ctx, inputParam).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis stream consume failure, stream - %s error: %w", r.stream, err)
	}

	messages := make([]messageclient.Message, 0)
	for _, streamMsg := range streamMessages {
		for _, msg := range streamMsg.Messages {
			messages = append(messages, messageclient.Message{
				ID:     msg.ID,
				Values: msg.Values,
			})
		}
	}

	return messages, nil
}

/*
ReceivePendingMessages Receives pending messages from stream for retry
*/
func (r *RedisStreamConsumer) ReceivePendingMessages(ctx context.Context) ([]messageclient.Message, error) {
	pipeline := r.redisClient.Pipeline()

	inputPendingMessages := &redis.XPendingExtArgs{
		Stream:   r.stream,
		Group:    r.groupName,
		Idle:     r.retryConfig.idleTimeThreshold,
		Start:    "-",
		Consumer: r.consumerName,
		Count:    r.batchSize,
		End:      "+",
	}
	pendingMessages := pipeline.XPendingExt(ctx, inputPendingMessages)
	inputParam := &redis.XAutoClaimArgs{
		MinIdle:  r.retryConfig.idleTimeThreshold,
		Group:    r.groupName,
		Stream:   r.stream,
		Count:    r.batchSize,
		Start:    "-",
		Consumer: r.consumerName,
	}
	claimedMessages := pipeline.XAutoClaim(ctx, inputParam)
	_, err := pipeline.Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("redis stream auto claim failure, stream - %s error: %w", r.stream, err)
	}

	pendingMessagesVal := pendingMessages.Val()
	claimedMessagesVal, _ := claimedMessages.Val()

	pendingMap := getRetryCountMap(pendingMessagesVal)
	messagesToBeRetried := r.handleClaimedMsg(ctx, claimedMessagesVal, pendingMap)

	messages := make([]messageclient.Message, 0)
	for _, msg := range messagesToBeRetried {
		messages = append(messages, messageclient.Message{
			ID:     msg.ID,
			Values: msg.Values,
		})
	}

	return messages, nil
}

/*
Acknowledge Acknowledges message with messageId to stream
*/
func (r *RedisStreamConsumer) Acknowledge(ctx context.Context, messageId string) error {

	_, err := r.redisClient.XAck(ctx, r.stream, r.groupName, messageId).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("Could not acknowledge  messageId : %q. Err: %w", messageId, err)
	}
	return nil

}

/*
RegisterConsumerGroup Registers consumer group
*/
func (r *RedisStreamConsumer) RegisterConsumerGroup(ctx context.Context) (string, error) {

	subResponse, err := r.redisClient.XGroupCreateMkStream(ctx, r.stream, r.groupName, "0").Result()
	if err != nil {
		return "", fmt.Errorf("Cannot register Redis Consumer, error: %w", err)
	}
	return subResponse, nil
}

/*
Start Starts consuming from stream
Sends batches of consumed messages to channel
*/
func (r *RedisStreamConsumer) Start(ctx context.Context, msgChan chan<- messageclient.MessagesWithError) error {
	if msgChan == nil {
		return errors.New("Unable to start consumer, No message channel found")
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go r.consumeMessages(ctx, msgChan, wg)
	go r.consumePendingMessages(ctx, msgChan, wg)

	wg.Wait()
	return nil
}

/*
consumeMessages continuously consumes messages from stream
*/
func (r *RedisStreamConsumer) consumeMessages(ctx context.Context, msgChan chan<- messageclient.MessagesWithError, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// Consume from stream
		messages, err := r.Receive(ctx)
		pushToChannel(ctx, messages, err, msgChan)
	}
}

/*
consumePendingMessages continuously consumes pending messages from stream i.e messages which are already consumed >= 1 times but are not acknowledged
*/
func (r *RedisStreamConsumer) consumePendingMessages(ctx context.Context, msgChan chan<- messageclient.MessagesWithError, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// Consume pending messages from stream
		messages, err := r.ReceivePendingMessages(ctx)
		pushToChannel(ctx, messages, err, msgChan)
		time.Sleep(r.retryConfig.retryInterval)
	}

}

/*
handleClaimedMsg If message has been retried more than 5 times, message is pushed to DLQ snd is acknowledged in stream
*/
func (r *RedisStreamConsumer) handleClaimedMsg(ctx context.Context, claimedMessages []redis.XMessage, pendingMsgMap map[string]int64) []redis.XMessage {
	var filteredMsg []redis.XMessage
	for _, val := range claimedMessages {
		retryCount, ok := pendingMsgMap[val.ID]
		if ok && retryCount > 5 {
			v, _ := json.Marshal(val)
			_, err := r.retryProducer.Send(ctx, &messageclient.SendRequest{
				Stream: r.stream + "_DLQ",
				Body:   string(v),
			})
			if err == nil {
				r.Acknowledge(ctx, val.ID)
			}
		} else {
			filteredMsg = append(filteredMsg, val)
		}
	}
	return filteredMsg
}

func getRetryCountMap(msg []redis.XPendingExt) map[string]int64 {
	retryCountMap := map[string]int64{}
	for _, val := range msg {
		retryCountMap[val.ID] = val.RetryCount
	}
	return retryCountMap
}

func pushToChannel(ctx context.Context, messages []messageclient.Message, err error, msgChan chan<- messageclient.MessagesWithError) {
	if err == nil {
		select {
		case <-ctx.Done():
			return
		case msgChan <- messageclient.MessagesWithError{
			Messages: messages,
			Err:      nil,
		}:
		}
	} else {
		select {
		case <-ctx.Done():
			return
		case msgChan <- messageclient.MessagesWithError{
			Messages: nil,
			Err:      fmt.Errorf("error occured while consuming messages, err: %w", err)}:
		}
	}
}
