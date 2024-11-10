package redisStrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/config"
	redis "github.com/redis/go-redis/v9"
)

/*
	font: https://github.com/redis/go-redis
*/

type RedisConn struct {
	ctx context.Context

	db         *redis.Client
	strategy   string
	expiration time.Duration
}

func New(cfg config.Cache) (*RedisConn, error) {
	strAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	db := redis.NewClient(&redis.Options{
		Addr:     strAddr,
		Password: cfg.Pass,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	Expiration := time.Duration(cfg.Expiration * int(time.Millisecond))

	return &RedisConn{
		ctx: context.Background(),

		db:         db,
		strategy:   cfg.Strategy,
		expiration: Expiration,
	}, nil
}

func (c *RedisConn) GetDB() interface{} {
	return c.db
}

func (c *RedisConn) Readiness() error {
	_, err := c.db.Ping(c.ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisConn) GetStrategy() string {
	return c.strategy
}

func (c *RedisConn) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.db.Set(c.ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisConn) Get(key string) (string, error) {
	val, err := c.db.Get(c.ctx, key).Result()
	if err != nil {
		// slog.Error("cannot get key: %v, CacheClient error: %v ", key, err)
		return "", err
	}
	if val == "" {
		return "", errors.New("get data empty")
	}

	return val, nil
}

func (c *RedisConn) Delete(key string) error {
	err := c.db.Del(c.ctx, key).Err()
	if err != nil {
		// slog.Error("cannot get key", "key", key, "CacheClient error", err)
		return err
	}
	return nil
}

func (c *RedisConn) GetDefaultExpiration() time.Duration {
	return c.expiration
}
