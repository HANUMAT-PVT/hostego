package redis_stream

// import (
// 	"context"
// 	"backend-hostego/internal/pkg/logger"
// 	"encoding/json"
// 	"fmt"

// 	messageclient "bitbucket.org/coinswitch/message-client"
// 	redis_stream "bitbucket.org/coinswitch/message-client/redis-stream"
// 	"github.com/redis/go-redis/v9"
// )

// var log = logger.GetLogger()

// type MessagePublisher interface {
// 	Publish(payload interface{}) error
// }

// type RedisStreamPublisher struct {
// 	Redis     *redis.Client
// 	Stream    string
// 	GroupName string
// }

// func NewStreamPublisher(redisHost string, streamName string, groupName string) RedisStreamPublisher {
// 	streamOptions := &redis.Options{
// 		Addr:     redisHost,
// 		Password: "",
// 		DB:       0,
// 	}
// 	return RedisStreamPublisher{
// 		Redis:     redis.NewClient(streamOptions),
// 		Stream:    streamName,
// 		GroupName: groupName,
// 	}
// }

// // Publish is a generic function to push messages to Redis Stream
// func (r *RedisStreamPublisher) Publish(payload interface{}) (e error) {

// 	//log.Infof("publishing message, publisher: %+v, --> redis: %+v", r, *r.Redis)
// 	producerBuilder := redis_stream.NewRedisStreamProducerBuilder()
// 	producer, err := producerBuilder.RedisClient(r.Redis).Build()
// 	if err != nil {
// 		log.Errorf("Unable to create producer %s", err)
// 		return err
// 	}
// 	////	CHECK if the Client is active
// 	// if _, err := r.Redis.Ping(context.Background()).Result(); err != nil {
// 	//	log.Error(err)
// 	//}
// 	marshaledPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Errorf("error in marshaling payload %v", payload)
// 		return err
// 	}
// 	// Producing message, use context.Background() to avoid nil pointer error arising due to ctx
// 	msgId, err := producer.Send(context.Background(), &messageclient.SendRequest{
// 		Stream:     r.Stream,
// 		Body:       string(marshaledPayload),
// 		Attributes: []messageclient.Attribute{},
// 	})
// 	if err != nil {
// 		log.Error("Error : %v", err)
// 		panic(any("error occurred while publishing message"))
// 	}
// 	log.Debug(fmt.Sprintf("debug: Data published to stream %v successfully: messageID %s", r.Stream, msgId))
// 	return
// }
