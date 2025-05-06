package in_mem

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type SyncMapSingleton struct {
	syncMap sync.Map
}

var instance *SyncMapSingleton
var once sync.Once

var (
	KeyExpiredError  = "Key has expired"
	KeyNotFoundError = "Key not found"
)

// GetSyncMapInstance returns the singleton instance of the sync.Map.
func GetSyncMapInstance() *SyncMapSingleton {
	once.Do(func() {
		instance = &SyncMapSingleton{
			syncMap: sync.Map{},
		}

	})
	return instance
}

// ValueWithExpiration represents the value along with its expiration time
type ValueWithExpiration struct {
	Value      interface{}
	Expiration *time.Time
}

func GetSyncMap() *SyncMapSingleton {
	return instance
}

// Get retrieves the value associated with the key and check if it is expired from the sync.Map.
func (s *SyncMapSingleton) Get(key string) (string, error) {
	value, found := s.syncMap.Load(key)
	if !found {
		err := errors.New(KeyNotFoundError)
		return "", err
	}

	now := time.Now()
	valueWithExpiration, ok := value.(ValueWithExpiration)
	if !ok {
		err := errors.New(KeyNotFoundError)
		return "", err
	}

	exp := valueWithExpiration.Expiration
	if exp != nil && now.After(*exp) {
		err := errors.New(KeyExpiredError)
		return "", err
	}

	return valueWithExpiration.Value.(string), nil
}

func (s *SyncMapSingleton) GetWithExpiry(key string, data interface{}) error {
	value, found := s.syncMap.Load(key)
	if !found {
		err := errors.New(KeyNotFoundError)
		return err
	}

	now := time.Now()
	valueWithExpiration, ok := value.(ValueWithExpiration)
	if !ok {
		err := errors.New(KeyNotFoundError)
		return err
	}

	exp := valueWithExpiration.Expiration
	if exp != nil && now.After(*exp) {
		err := errors.New(KeyExpiredError)
		return err
	}

	valBytes, err := json.Marshal(valueWithExpiration.Value)
	if err != nil {
		return err
	}

	err = json.Unmarshal(valBytes, data)
	if err != nil {
		return err
	}

	return nil
}

// Set sets the value associated with the key in the sync.Map.
func (s *SyncMapSingleton) Set(key, value string) {
	valueWithExpiration := ValueWithExpiration{
		Value:      value,
		Expiration: nil,
	}

	s.syncMap.Store(key, valueWithExpiration)
}

// Set sets the value associated with the key in the sync.Map.
func (s *SyncMapSingleton) SetWithExpiry(key string, value interface{}, expiration time.Duration) {
	expTime := time.Now().Add(expiration)

	valueWithExpiration := ValueWithExpiration{
		Value:      value,
		Expiration: &expTime,
	}

	s.syncMap.Store(key, valueWithExpiration)
}
