package cache

import "time"

// GetCacheFunc 获取缓存委托
type GetCacheFunc func(cacheKey string) []byte

// SetCacheFunc 设置缓存委托
type SetCacheFunc func(cacheKey string, value []byte, expiration time.Duration) bool
