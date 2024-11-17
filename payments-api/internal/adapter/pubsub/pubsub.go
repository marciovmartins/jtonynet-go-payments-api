package pubsub

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase"
)

type PubSub interface {
	GetStrategy(ctx context.Context) (string, error)
	Subscribe(ctx context.Context, topic string) (<-chan string, error)
	Publish(ctx context.Context, topic, message string) error
	Close() error
}

/*
func NewCLient(cfg config.InMemoryDatabase) (Conn, error) {
	switch cfg.Strategy {
	case "redis":
		return NewRedisClient(cfg)
	default:
		return nil, fmt.Errorf("pubsub strategy not suported: %s", cfg.Strategy)
	}
}
*/

func NewPubSubFromClient(client inMemoryDatabase.Client) (PubSub, error) {
	strategy, err := client.GetStrategy(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error when instantiating pubsub from client: %v", err)
	}

	switch strategy {
	case "redis":
		return NewRedisClient(client)
	default:
		return nil, fmt.Errorf("pubsub strategy not suported: %s", strategy)
	}
}
