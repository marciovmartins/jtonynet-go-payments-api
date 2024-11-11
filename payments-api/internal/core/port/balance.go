package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type BalanceByCategoryEntity struct {
	ID        uint
	UID       uuid.UUID
	AccountID uint
	Amount    decimal.Decimal
	Category  CategoryEntity
}

type BalanceEntity struct {
	AccountID   uint
	AmountTotal decimal.Decimal
	Categories  map[int]BalanceByCategoryEntity
}

type BalanceRepository interface {
	FindByAccountID(ctx context.Context, id uint) (BalanceEntity, error)
	UpdateTotalAmount(ctx context.Context, be BalanceEntity) error
}
