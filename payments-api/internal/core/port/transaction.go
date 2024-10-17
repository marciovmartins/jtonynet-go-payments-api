package port

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	AccountUID  uuid.UUID
	TotalAmount decimal.Decimal
	MCCcode     string
	Merchant    string
}

type TransactionEntity struct {
	ID          int
	UID         uuid.UUID
	AccountID   int
	TotalAmount decimal.Decimal
	MCCcode     string
	Merchant    string
}

type TransactionRepository interface {
	Save(TransactionEntity) error
}
