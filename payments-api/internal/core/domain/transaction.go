package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	AccountID   uint
	AccountUID  uuid.UUID
	MCC         string
	TotalAmount decimal.Decimal
	Merchant    string
}
