// package redis

// import (
// 	redismock "github.com/go-redis/redismock/v9"
// )

// type MockClient struct {
// 	RedisClient Redis
// 	Mock        redismock.ClusterClientMock
// }

// func GetMockRedisClient() MockClient {
// 	rc, rm := redismock.NewClusterMock()
// 	return MockClient{
// 		RedisClient: Redis{
// 			client: rc,
// 		},
// 		Mock: rm,
// 	}
// }

package redis

import (
	"context"

	redismock "github.com/go-redis/redismock/v9"
)

// MockClient struct for Redis mock
type MockClient struct {
	RedisClient *Redis
	Mock        redismock.ClientMock
}

// GetMockRedisClient initializes a mock Redis client for testing
func GetMockRedisClient() *MockClient {
	rc, rm := redismock.NewClientMock() // Using single Redis client mock

	mockRedis := &Redis{
		client: rc,
		ctx:    context.TODO(),
	}

	return &MockClient{
		RedisClient: mockRedis,
		Mock:        rm,
	}
}
