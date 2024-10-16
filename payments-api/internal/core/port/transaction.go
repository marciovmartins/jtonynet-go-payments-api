package port

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	AccountUID  uuid.UUID
	TotalAmount decimal.Decimal
	MCCcode     string
	Merchant    string
}

type TransactionEntity struct {
	AccountUID  uuid.UUID
	TotalAmount decimal.Decimal
	MCCcode     string
	Merchant    string
}

type TransactionRepository interface {
	Save(TransactionEntity) error
}
