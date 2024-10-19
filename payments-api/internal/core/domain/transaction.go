package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	AccountID   uint
	AccountUID  uuid.UUID
	MCCcode     string
	TotalAmount decimal.Decimal
	Merchant    string
}

func NewTransaction(
	accountID uint,
	accountUID uuid.UUID,
	mccCode string,
	totalAmount decimal.Decimal,
	merchant string) (*Transaction, error) {
	return &Transaction{
		AccountID:   accountID,
		AccountUID:  accountUID,
		MCCcode:     mccCode,
		TotalAmount: totalAmount,
		Merchant:    merchant,
	}, nil
}
