package pubSub

import (
	"context"
	"fmt"
	"sync"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub

	subscriptions sync.Map
	strategy      string
}

func NewRedisPubSub(cfg config.PubSub) (*RedisPubSub, error) {
	strAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     strAddr,
		Password: cfg.Pass,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	rps := &RedisPubSub{
		client:        client,
		strategy:      cfg.Strategy,
		subscriptions: sync.Map{},
	}

	err := rps.subscribe(context.Background())
	if err != nil {
		return &RedisPubSub{}, err
	}

	return rps, nil
}

func (r *RedisPubSub) Subscribe(_ context.Context, key Key) (<-chan string, error) {
	listenerBufferSize := 1

	if transactionsSubscriptions, ok := r.subscriptions.Load(key.Account); ok {
		transactionMap, _ := transactionsSubscriptions.(map[string]chan string)

		if subscription, exists := transactionMap[key.Transaction]; exists {
			return subscription, nil
		}

		subscription := make(chan string, listenerBufferSize)
		transactionMap[key.Transaction] = subscription
		r.subscriptions.Store(key.Account, transactionMap)
		return subscription, nil
	}

	subscription := make(chan string, listenerBufferSize)
	transactionMap := map[string]chan string{
		key.Transaction: subscription,
	}

	r.subscriptions.Store(key.Account, transactionMap)
	return subscription, nil
}

func (r *RedisPubSub) UnSubscribe(_ context.Context, key Key) error {
	if transactionsSubscriptions, ok := r.subscriptions.Load(key.Account); ok {
		transactionMap, _ := transactionsSubscriptions.(map[string]chan string)

		if subscription, exists := transactionMap[key.Transaction]; exists {

			delete(transactionMap, key.Transaction)
			close(subscription)

			if len(transactionMap) == 0 {
				r.subscriptions.Delete(key.Account)
			} else {
				r.subscriptions.Store(key.Account, transactionMap)
			}
		}
	}

	return nil
}

func (r *RedisPubSub) subscribe(ctx context.Context) error {
	keyspaceChannel := fmt.Sprintf("__keyevent@%d__:expired", r.client.Options().DB)
	r.pubsub = r.client.Subscribe(ctx, keyspaceChannel)

	go func() {
		defer r.pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-r.pubsub.Channel():
				if !ok {
					return
				}

				accountUID := msg.Payload
				if transactionsSubscriptions, ok := r.subscriptions.Load(accountUID); ok {
					subscriptions, valid := transactionsSubscriptions.(map[string]chan string)
					if !valid {
						continue
					}

					for _, subscription := range subscriptions {
						select {
						case subscription <- accountUID:
						default:
							continue
						}
					}
				}
			}
		}
	}()

	return nil
}

func (r *RedisPubSub) Publish(ctx context.Context, topic, message string) error {
	return r.client.Publish(ctx, topic, message).Err()
}

func (r *RedisPubSub) Close() error {
	if r.pubsub != nil {
		return r.pubsub.Close()
	}
	return nil
}

func (r *RedisPubSub) GetStrategy(_ context.Context) (string, error) {
	return r.strategy, nil
}
