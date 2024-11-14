package inMemoryDatabase

import (
	"context"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/config"
)

type DBConn interface {
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetDefaultExpiration(ctx context.Context) (time.Duration, error)
}

func NewConn(cfg config.InMemoryDatabase) (DBConn, error) {
	switch cfg.Strategy {
	case "redis":
		return NewRedisConn(cfg)
	default:
		return nil, fmt.Errorf("InMemoryDB strategy not suported: %s", cfg.Strategy)
	}
}
