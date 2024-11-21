package domain

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	AmountTotal decimal.Decimal
	TransactionByCategories
}
