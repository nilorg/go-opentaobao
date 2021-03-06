package opentaobao

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// SetRedis 设置RedisCache
func SetRedis(redisClient *redis.Client) {
	GetCache = func(cacheKey string) []byte {
		bytes, err := redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == redis.Nil {
			return nil
		} else if err != nil {
			log.Println(err)
			return nil
		}
		return bytes
	}

	SetCache = func(key string, value []byte, expiration time.Duration) bool {
		err := redisClient.SetNX(context.Background(), key, value, expiration).Err()
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}
