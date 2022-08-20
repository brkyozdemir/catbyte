package config

import (
	"os"
)

type Cache struct {
	Default RedisConnection
}

type RedisConnection struct {
	Address string
}

func NewCache() Cache {
	return Cache{
		Default: RedisConnection{
			Address: os.Getenv("REDIS_ADDRESS"),
		},
	}
}
