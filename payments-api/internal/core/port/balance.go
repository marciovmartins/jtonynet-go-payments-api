package port

import (
	"github.com/shopspring/decimal"
)

type BalanceByCategoryEntity struct {
	ID        int
	AccountID int
	Amount    decimal.Decimal
	Category  CategoryEntity
}

type BalanceEntity struct {
	AccountID   int
	AmountTotal decimal.Decimal
	Categories  map[int]BalanceByCategoryEntity
}

type BalanceRepository interface {
	FindByAccountID(int) (BalanceEntity, error)
	Update(BalanceEntity) error
}
