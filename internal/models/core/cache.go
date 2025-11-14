package core

import (
	"context"
	"time"
)

type Cache interface {
	cacheOperations

	Close() error
}

type cacheOperations interface {
	Set		(ctx context.Context, key string, value any, ttl time.Duration)		error
	SetNX	(ctx context.Context, key string, value any, ttl time.Duration)		(bool, error)
	Get		(ctx context.Context, key string) 									(any, error)
	Del		(ctx context.Context, key string) 									error
}