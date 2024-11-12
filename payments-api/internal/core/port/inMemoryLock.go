package port

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type InMemoryLockEntity struct {
	Key       uuid.UUID
	Timestamp time.Time
}

type InMemoryLockRepository interface {
	Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Unlock(ctx context.Context, key string) error
	GetLockedValue(ctx context.Context, key string) (string, error)
}
