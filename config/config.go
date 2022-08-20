package config

import "sync"

var once sync.Once
var config *Config

type Config struct {
	Cache
}

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			Cache: NewCache(),
		}
	})

	return config
}
