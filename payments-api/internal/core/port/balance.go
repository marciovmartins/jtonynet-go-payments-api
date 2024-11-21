package port

import (
	"github.com/shopspring/decimal"
)

type BalanceEntity struct {
	AmountTotal decimal.Decimal
	Categories  map[int]TransactionByCategoryEntity
}
