package post

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"social-network-otus/internal/cache"
	"social-network-otus/internal/logger"
	"time"
)

type Cache struct {
	CacheClient *cache.CacheService
	logger      logger.LoggerInterface
}

func NewCacheService(service *cache.CacheService, logger logger.LoggerInterface) *Cache {
	return &Cache{
		CacheClient: service,
		logger:      logger,
	}
}

func (c *Cache) GetFeed(userId string) (*Feed, error) {
	result, err := c.CacheClient.Client.Get(context.Background(), c.buildKey(userId)).Result()

	if err != nil {
		if err == redis.Nil {
			return &Feed{}, nil
		}
		c.logger.Error(err.Error(), err, logger.Fields{})
		return nil, err
	}

	var feed Feed

	err = json.Unmarshal([]byte(result), &feed)
	if err != nil {
		c.logger.Error(err.Error(), err, logger.Fields{})
		return nil, err
	}

	return &feed, nil
}

func (c *Cache) SetFeed(userId string, feed *Feed, expiration time.Duration) error {
	json, err := json.Marshal(feed)
	if err != nil {
		c.logger.Error(err.Error(), err, logger.Fields{})
		return err
	}

	_, err = c.CacheClient.Client.Set(context.Background(), c.buildKey(userId), json, expiration).Result()
	return err
}

func (c *Cache) Clear(userId string) {
	c.CacheClient.Client.Del(context.Background(), c.buildKey(userId))
}

func (c *Cache) Status(userId string) StatusFeed {
	result, err := c.GetFeed(userId)
	if err != nil {
		return NotFound
	}

	return result.Status
}

func (c *Cache) buildKey(userId string) string {
	return fmt.Sprintf("feed:user:%v", userId)
}
