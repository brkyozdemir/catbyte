package cache

import (
	"time"
	"twitch_chat_analysis/cmd/messageProcessor/cache/redis"
	"twitch_chat_analysis/config"
)

type Repository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Cache struct {
	Repository Repository
}

func NewCache(CacheConnection interface{}) *Cache {
	switch CacheConnection.(type) {
	case config.RedisConnection:
		cfg := CacheConnection.(config.RedisConnection)
		return &Cache{
			Repository: redis.NewCache(cfg),
		}
	default:
		panic("Invalid Cache Connection!")
	}
}
