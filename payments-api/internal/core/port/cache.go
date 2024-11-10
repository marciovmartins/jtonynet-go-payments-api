package port

import (
	"context"
	"time"
)

type Cache interface {
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetDefaultExpiration(ctx context.Context) (time.Duration, error)
}
