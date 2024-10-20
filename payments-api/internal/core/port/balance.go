package port

import (
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
	FindByAccountID(uint) (BalanceEntity, error)
	UpdateTotalAmount(BalanceEntity) error
}
