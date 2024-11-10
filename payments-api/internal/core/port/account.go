package port

import (
	"context"

	"github.com/google/uuid"
)

type AccountEntity struct {
	ID  uint
	UID uuid.UUID
}

type AccountRepository interface {
	FindByUID(ctx context.Context, uid uuid.UUID) (AccountEntity, error)
}
