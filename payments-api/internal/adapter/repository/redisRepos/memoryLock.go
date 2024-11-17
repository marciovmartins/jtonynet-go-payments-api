package redisRepos

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type MemoryLock struct {
	lockConn inMemoryDatabase.Client
}

func NewMemoryLock(lockConn inMemoryDatabase.Client) (port.MemoryLockRepository, error) {
	return &MemoryLock{
		lockConn,
	}, nil
}

func (ml *MemoryLock) Lock(
	_ context.Context,
	timeoutSLA port.TimeoutSLA,
	mle port.MemoryLockEntity,
) (port.MemoryLockEntity, error) {
	var lockErr error
	maxElapsedTime := time.Duration(timeoutSLA / 2)

	retry := backoff.NewExponentialBackOff()
	retry.MaxElapsedTime = maxElapsedTime
	retry.InitialInterval = time.Duration(1) * time.Millisecond

	backoff.RetryNotify(func() error {
		_, lockErr = ml.isUnlocked(mle)
		return lockErr
	}, retry, nil)

	if lockErr != nil {
		return port.MemoryLockEntity{}, lockErr
	}

	expiration, err := ml.lockConn.GetDefaultExpiration(context.TODO())
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	err = ml.lockConn.Set(context.TODO(), mle.Key, mle.Timestamp, expiration)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	return mle, nil
}

func (ml *MemoryLock) Unlock(_ context.Context, key string) error {
	return ml.lockConn.Delete(context.TODO(), key)
}

func (ml *MemoryLock) isUnlocked(mle port.MemoryLockEntity) (bool, error) {
	locked, err := ml.get(context.TODO(), mle.Key)
	if err == nil {
		elapsedTime := time.Now().UnixMilli() - locked.Timestamp
		return false, fmt.Errorf("key is already locked by another process, held for %v ms", elapsedTime)
	}

	return true, nil
}

func (ml *MemoryLock) get(_ context.Context, key string) (port.MemoryLockEntity, error) {
	timestampStr, err := ml.lockConn.Get(context.TODO(), key)
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
