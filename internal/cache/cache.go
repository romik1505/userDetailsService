package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/romik1505/userDetailsService/internal/config"
	"time"
)

type Redis struct {
	Client *redis.Client
}

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

func NewCache() *Redis {
	addr := fmt.Sprintf("%s:%s", config.Config.Cache.Host, config.Config.Cache.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Config.Cache.Password,
	})

	return &Redis{
		Client: client,
	}
}

func NewMockCache() (*Redis, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	return &Redis{
		Client: client,
	}, mock
}

func (c *Redis) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.Client.Get(ctx, key)
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.Client.Set(ctx, key, value, expiration)
}

func (c *Redis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.Client.Del(ctx, keys...)
}
