package redis

import (
	"errors"
	"fmt"

	"github.com/femitubosun/go-sweepline-availability/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	config *config.Config
}

func NewCache(c *config.Config) (*Cache, error) {
	if c == nil {
		return nil, errors.New("missing app configuration")
	}

	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort),
			Password: c.RedisPwd,
			DB:       c.RedisFamily,
		}),
		config: c,
	}, nil
}

func (c *Cache) Client() *redis.Client {
	return c.client
}
