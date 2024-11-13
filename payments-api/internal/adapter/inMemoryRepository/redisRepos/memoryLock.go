package redisRepos

import (
	"context"
	"strconv"

	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type MemoryLock struct {
	lockConn port.InMemoryDBConn
}

func NewMemoryLock(lockConn port.InMemoryDBConn) (port.MemoryLockRepository, error) {
	return &MemoryLock{
		lockConn,
	}, nil
}

func (ml *MemoryLock) Lock(
	_ context.Context,
	mle port.MemoryLockEntity,
) error {
	expiration, err := ml.lockConn.GetDefaultExpiration(context.Background())
	if err != nil {
		return err
	}

	return ml.lockConn.Set(context.Background(), mle.Key, mle.Timestamp, expiration)
}

func (ml *MemoryLock) Unlock(_ context.Context, key string) error {
	return ml.lockConn.Delete(context.Background(), key)
}

func (ml *MemoryLock) Get(_ context.Context, key string) (port.MemoryLockEntity, error) {
	timestampStr, err := ml.lockConn.Get(context.Background(), key)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	return port.MemoryLockEntity{
		Key:       key,
		Timestamp: timestamp,
	}, nil
}
