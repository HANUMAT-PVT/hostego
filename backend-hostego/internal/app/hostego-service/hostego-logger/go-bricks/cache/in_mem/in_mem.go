package in_mem

// Get config from local cache(in-memory)
func GetConfigFromLocalCache(cacheKey string) (string, error) {
	inMemCacheSource := NewInMemCacheConfigSource()
	inMemVal, inMemError := inMemCacheSource.Get(cacheKey)

	if inMemError != nil {
		return "", inMemError
	}
	
	return inMemVal, nil
}

// Set config in local cache(in-memory)
func SetConfigInLocalCache(cacheKey string, cacheValue string) {
	inMemCacheSource := NewInMemCacheConfigSource()
	inMemCacheSource.Set(cacheKey, cacheValue)
}
