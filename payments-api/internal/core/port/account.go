package port

import (
	"context"

	"github.com/google/uuid"
)

type AccountEntity struct {
	ID      uint
	UID     uuid.UUID
	Balance BalanceEntity
}

type AccountRepository interface {
	FindByUID(ctx context.Context, uid uuid.UUID) (AccountEntity, error)
	SaveTransactions(ctx context.Context, transactions map[int]TransactionEntity) error
}
