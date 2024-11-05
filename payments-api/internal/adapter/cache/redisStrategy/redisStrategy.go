package redisStrategy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jtonynet/go-payments-api/config"
	redis "github.com/redis/go-redis/v9"
)

/*
	font: https://github.com/redis/go-redis
*/

type RedisConn struct {
	db  *redis.Client
	ctx context.Context

	Expiration time.Duration
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
		db:  db,
		ctx: context.Background(),

		Expiration: Expiration,
	}, nil
}

func (c *RedisConn) Set(key string, value interface{}, expiration time.Duration) error {
	err := c.db.Set(c.ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisConn) Get(key string) (string, error) {
	val, err := c.db.Get(c.ctx, key).Result()
	if err != nil {
		slog.Error("cannot get key: %v, CacheClient error: %v ", key, err)
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
		slog.Error("cannot delete key: %v, cache client error: %v", key, err)
		return err
	}
	return nil
}

func (c *RedisConn) Readiness() bool {
	_, err := c.db.Ping(c.ctx).Result()
	return err == nil
}

func (c *RedisConn) GetDefaultExpiration() time.Duration {
	return c.Expiration
}
