package in_mem

import (
	"context"
	"time"

	"backend-hostego/internal/app/hostego-service/hostego-logger/go-bricks/logger"
)

// CacheConfigSource is a struct that represents a cache configuration source.
// Ideally the `cache` client should be an interface to a singleton cache client
// as of now there is no interface created so using the redis client directly
type InMemCacheConfigSource struct {
	SourceType string
}

var (
	log        = logger.GetNewLogger(context.Background(), logger.Logrus, logger.LogLevelInfo)
	inMemCache = GetSyncMapInstance()
)

const InMemSourceType = "IN_MEM"

// NewInMemCacheConfigSource returns a new InMemCacheConfigSource.
func NewInMemCacheConfigSource() *InMemCacheConfigSource {
	return &InMemCacheConfigSource{
		SourceType: InMemSourceType,
	}
}

// Get returns config interface for a given key.
func (ccs *InMemCacheConfigSource) Get(key string) (string, error) {
	// logPrefix := "Get=>"
	dest, err := inMemCache.Get(key)

	if err != nil {
		// log.Warnf("%s received nil value for key: %v", logPrefix, key)
		return "", err
	}

	return dest, nil
}

// Set updates the configuration value for the given key in the source.
func (ccs *InMemCacheConfigSource) Set(key string, value string) error {
	inMemCache.Set(key, value)
	return nil
}

// Set updates the configuration value for the given key in the source.
func (ccs *InMemCacheConfigSource) SetWithExpiry(key string, value string, expiration time.Duration) error {
	inMemCache.SetWithExpiry(key, value, expiration)
	return nil
}

func (ccs *InMemCacheConfigSource) GetSourceType() string {
	return ccs.SourceType
}
