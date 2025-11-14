package valkey

import (
	"context"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"time"
)

type Valkey_Mock struct{}

// Del implements core.Cache.
func (v *Valkey_Mock) Del(ctx context.Context, key string) error {
	return nil
}

// SetNX implements core.Cache.
func (v *Valkey_Mock) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return true, nil
}

// Close implements core.Cache.
func (v *Valkey_Mock) Close() error {
	return nil
}

// Get implements core.Cache.
func (v *Valkey_Mock) Get(ctx context.Context, key string) (any, error) {
	return nil, nil
}

// Set implements core.Cache.
func (v *Valkey_Mock) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return nil
}

func NewValkey_Mock() core.Cache {
	return &Valkey_Mock{}
}
