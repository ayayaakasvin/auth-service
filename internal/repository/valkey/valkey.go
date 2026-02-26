package valkey

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/redis/go-redis/v9"
)

const origin = "Redis/Valkey"

// for storing methods of storing and retrieving session_id
type Cache struct {
	connection *redis.Client
}

func NewValkeyClient(cfg config.ValkeyConfig) (core.Cache, error) {
	ctx := context.Background()
	opt, err := redis.ParseURL(cfg.URL)
	log.Printf("URL: %s", cfg.URL)
	if err != nil {
		msg := fmt.Sprintf("failed to parse Redis URL: %v", err)
		log.Println(msg)
		return nil, err
	}

	conn := redis.NewClient(opt)
	if err := conn.Ping(ctx).Err(); err != nil {
		msg := fmt.Sprintf("failed to connect to db: %v", err)
		log.Println(msg)
		return nil, err
	}

	return &Cache{
		connection: conn,
	}, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.connection.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	return c.connection.Get(ctx, key).Result()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.connection.Del(ctx, key).Err()
}

func (c *Cache) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	set := c.connection.SetNX(ctx, key, value, ttl)
	return set.Val(), set.Err()
}

func (c *Cache) Close() error {
	return c.connection.Close()
}
