package port

import (
	"context"
)

type MemoryLockEntity struct {
	Key       string
	Timestamp int64 //timestamp.UnixMilli
}

type MemoryLockRepository interface {
	Lock(ctx context.Context, mle MemoryLockEntity) error
	Unlock(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (MemoryLockEntity, error)
}
