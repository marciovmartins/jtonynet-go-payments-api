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
	ID          uint
	UID         uuid.UUID
	AccountID   uint
	TotalAmount decimal.Decimal
	MCCcode     string
	Merchant    string
}

type TransactionRepository interface {
	Save(TransactionEntity) error
}
