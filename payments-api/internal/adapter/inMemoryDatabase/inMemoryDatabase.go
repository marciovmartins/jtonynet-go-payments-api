package inMemoryDatabase

import (
	"context"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/pubsub"
)

type Client interface {
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetDefaultExpiration(ctx context.Context) (time.Duration, error)
	GetClient(ctx context.Context) (interface{}, error)
}

func NewClient(cfg config.InMemoryDatabase, pubSubLock pubsub.PubSub) (Client, error) {
	switch cfg.Strategy {
	case "redis":
		return NewRedisClient(cfg)
	default:
		return nil, fmt.Errorf("InMemoryDB strategy not suported: %s", cfg.Strategy)
	}
}
