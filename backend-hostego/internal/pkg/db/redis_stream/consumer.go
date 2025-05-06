package redis_stream

// import (
// 	"context"
// 	"strings"
// 	"time"

// 	"github.com/redis/go-redis/v9"

// 	redis_stream "bitbucket.org/coinswitch/message-client/redis-stream"
// )

// type RedisStreamMessageReader struct { // should be same as redisStreamMessagePublisher = RedisStream
// 	Redis        *redis.Client
// 	Stream       string
// 	GroupName    string
// 	MessageCount int64
// 	BlockTime    time.Duration
// }

// func NewStream(redisClient *redis.Client, streamName string, groupName string, msgCount int64, blockTime time.Duration) RedisStreamMessageReader {
// 	return RedisStreamMessageReader{
// 		Redis:        redisClient,
// 		Stream:       streamName,
// 		GroupName:    groupName,
// 		MessageCount: msgCount,
// 		BlockTime:    blockTime,
// 	}
// }

// // BuildMessageConsumer creates a conumer with the given name and will define consumer to check every retryInterval duration if there is any message in pending state for more than idletTimeThreshold duration.
// // Retry batch size is same as stream defined batch size
// func (r *RedisStreamMessageReader) BuildMessageConsumer(ctx context.Context, consumerName string, idleTimeThreshold time.Duration, retryInterval time.Duration) (*redis_stream.RedisStreamConsumer, error) {

// 	consumerBuilder := redis_stream.NewRedisStreamConsumerBuilder()

// 	consumer, er := consumerBuilder.RedisClient(r.Redis).
// 		Stream(r.Stream).
// 		BatchSize(r.MessageCount).
// 		BlockTime(r.BlockTime).
// 		GroupName(r.GroupName).
// 		RetryConfig(idleTimeThreshold, retryInterval).
// 		ConsumerName(consumerName).
// 		Build()

// 	if er != nil {
// 		log.Errorf("unable to create consumer %s", er)
// 		return nil, er
// 	}

// 	res, err := consumer.RegisterConsumerGroup(ctx) // no need to catch error here, just needed for the first time, thereafter will be a redundant step
// 	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
// 		return nil, err
// 	}
// 	log.Infof("Consumer %v registered with stream %v and group %v: %v", consumerName, r.Stream, r.GroupName, res)

// 	return consumer, nil
// }
