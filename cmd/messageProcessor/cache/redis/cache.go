package redis

import (
	"context"
	redisCache "github.com/go-redis/redis/v8"
	"time"
	"twitch_chat_analysis/config"
)

type Cache struct {
	Client      *redisCache.Client
	millisecond func()
}

var ctx = context.TODO()

func NewCache(connection config.RedisConnection) *Cache {
	client := redisCache.NewClient(&redisCache.Options{
		Addr: connection.Address,
	})

	return &Cache{
		Client: client,
	}
}

func (c Cache) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c Cache) Get(key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c Cache) Delete(key string) error {
	return c.Client.Del(ctx, key).Err()
}
