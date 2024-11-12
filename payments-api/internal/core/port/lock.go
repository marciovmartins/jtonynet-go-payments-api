package port

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LockEntity struct {
	Key       uuid.UUID
	Timestamp time.Time
}

type LockRepository interface {
	Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Unlock(ctx context.Context, key string) error
	GetLockedValue(ctx context.Context, key string) (string, error)
}
