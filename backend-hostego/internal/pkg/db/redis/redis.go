package redis

import (
	"backend-hostego/internal/app/hostego-service/constants/config_constants"
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/pkg/logger"
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/spf13/viper"
)

// var Client = &Redis{}

// type Redis struct {
// 	redis *redis.ClusterClient
// 	ctx   context.Context
// }

// var log = logger.GetLogger()

// func InitRedis() {
// 	ctx := context.TODO()

// 	GlobalRedisClient := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs:    viper.GetStringSlice(config_constants.VKEYS_REDIS_CLUSTERS_HOST_URL),
// 		PoolSize: viper.GetInt(config_constants.VKEYS_REDIS_CLUSTERS_POOL_SIZE),
// 	})

// 	_, err := GlobalRedisClient.Ping(ctx).Result()
// 	if err != nil {
// 		log.Fatal(errConstants.CRITICAL_ISSUE + " failed to initialize Redis error:-" + err.Error())
// 		return
// 	}
// 	Client.redis = GlobalRedisClient
// 	Client.ctx = ctx

// 	log.Info("redis initialized successfully")
// }

// func GetClient() *Redis {
// 	if Client.redis == nil {
// 		InitRedis()
// 	}

// 	return Client
// }

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

var (
	instance *Redis
	once     sync.Once
	log      = logger.GetLogger()
)

// InitRedis initializes the Redis Sentinel client (Singleton)
func InitRedis() {
	once.Do(func() {
		ctx := context.TODO()

		// Fetch Redis host and configurations
		redisAddr := viper.GetString(config_constants.VKEYS_REDIS_HOST) // Default: 0.0.0.0:6379
		password := viper.GetString("")                                 // Default: ""
		poolSize := viper.GetInt("10")                                  // Default: 10

		// Set default values if not configured
		if redisAddr == "" {
			redisAddr = "127.0.0.1:6379"
		}
		if poolSize == 0 {
			poolSize = 10
		}

		// Create Redis client
		client := redis.NewClient(&redis.Options{
			Addr:         redisAddr,
			Password:     password,
			DB:           0, // Default Redis database
			PoolSize:     poolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		// Ping Redis to verify connection
		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatal("Failed to initialize Redis: " + err.Error())
			return
		}

		instance = &Redis{
			client: client,
			ctx:    ctx,
		}

		log.Info("Redis initialized successfully")
	})
}

// GetClient returns the singleton Redis client
func GetClient() *Redis {
	if instance == nil {
		InitRedis()
	}
	return instance
}

// GetCtx returns the Redis context
func (r *Redis) GetCtx() context.Context {
	return r.ctx
}

// GetRedis returns the underlying Redis client
func (r *Redis) GetRedis() *redis.Client {
	return r.client
}

func (r *Redis) Get(ctx dto.ReqCtx, key string) (string, error) {
	// defer setupNewRelicSegment(&ctx, "get", key)()
	result, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *Redis) Set(ctx dto.ReqCtx, key, value string, timeout time.Duration) error {
	// defer setupNewRelicSegment(&ctx, "set", key)()
	err := r.client.Set(r.ctx, key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) MGet(ctx dto.ReqCtx, keys ...string) ([]string, error) {
	// defer setupNewRelicSegment(&ctx, "mget", "")()
	result, err := r.client.MGet(r.ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	values := make([]string, len(result))
	for i, val := range result {
		if val == nil {
			values[i] = ""
		} else {
			values[i] = val.(string)
		}
	}
	return values, nil
}

func (r *Redis) HGet(ctx dto.ReqCtx, key, field string) (string, error) {
	// defer setupNewRelicSegment(&ctx, "hget", key)()
	result, err := r.client.HGet(r.ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *Redis) HSet(ctx dto.ReqCtx, key string, field string, value interface{}) error {
	// defer setupNewRelicSegment(&ctx, "hset", key)()
	err := r.client.HSet(r.ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) HMGet(ctx dto.ReqCtx, key string, fields ...string) (map[string]string, error) {
	// defer setupNewRelicSegment(&ctx, "hmget", key)()
	result, err := r.client.HMGet(r.ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for i, val := range result {
		if val != nil {
			values[fields[i]] = val.(string)
		}
	}

	return values, nil
}

func (r *Redis) HGetAll(ctx dto.ReqCtx, key string) (map[string]string, error) {
	// defer setupNewRelicSegment(&ctx, "hgetall", key)()
	result, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for field, val := range result {
		values[field] = val
	}
	return values, nil
}

func (r *Redis) Delete(ctx dto.ReqCtx, key string) error {
	// defer setupNewRelicSegment(&ctx, "delete", key)()
	_, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetAllKeys(ctx dto.ReqCtx) ([]string, error) {
	res, err := r.client.Keys(r.ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// setupNewRelicSegment start a datastore segment for the redis call
// it starts the segment and return the End function which can be deferred
// func setupNewRelicSegment(rCtx *dto.ReqCtx, operation string, key string) func() {
// 	newRelicSegment := newrelic.DatastoreSegment{
// 		Product:            newrelic.DatastoreRedis,
// 		Host:               viper.GetString(config_constants.VKEYS_REDIS_CLUSTERS_HOST_URL),
// 		Operation:          operation,
// 		ParameterizedQuery: key,
// 		RawQuery:           key,
// 	}
// 	newRelicSegment.StartTime = rCtx.NewRelicTxn.StartSegmentNow()
// 	return newRelicSegment.End

// }
