package cache

import (
	"github.com/redis/go-redis/v9"
	"social-network-otus/internal/config"
)

type CacheService struct {
	Client *redis.Client
}

func NewCacheService(config *config.Config) *CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "", // no password set
		DB:       config.RedisCacheDb,
	})
	return &CacheService{
		Client: rdb,
	}
}
