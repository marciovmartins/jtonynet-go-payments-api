package pubsub

import (
	"context"
	"log"

	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewRedisClient(c inMemoryDatabase.Client) (PubSub, error) {
	cli, err := c.GetClient(context.TODO())
	if err != nil {
		return nil, err
	}

	client, ok := cli.(*redis.Client)
	if !ok {
		log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
	}

	return RedisPubSub{client: client}, nil
}

func (r RedisPubSub) Subscribe(_ context.Context, topic string) (<-chan string, error) {
	r.pubsub = r.client.Subscribe(context.TODO(), topic)
	channel := make(chan string)

	go func() {
		defer close(channel)
		for msg := range r.pubsub.Channel() {
			channel <- msg.Payload
		}
	}()

	return channel, nil
}

func (r RedisPubSub) Publish(ctx context.Context, topic, message string) error {
	return r.client.Publish(ctx, topic, message).Err()
}

func (r RedisPubSub) Close() error {
	if r.pubsub != nil {
		return r.pubsub.Close()
	}
	return nil
}

func (r RedisPubSub) GetStrategy(_ context.Context) (string, error) {
	return "", nil
}
