package store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	r *redis.Client
}

func NewRedisClient(r *redis.Client) *RedisClient {
	return &RedisClient{r: r}
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {

	err := c.r.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {

	val, err := c.r.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *RedisClient) Exists(ctx context.Context, key string) bool {
	return true
}
