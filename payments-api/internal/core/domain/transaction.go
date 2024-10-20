package domain

import (
	"fmt"

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

	if !totalAmount.GreaterThan(decimal.NewFromInt(0)) {
		return &Transaction{}, fmt.Errorf("transaction totalAmount must be a positive value %v", totalAmount)
	}

	return &Transaction{
		AccountID:   accountID,
		AccountUID:  accountUID,
		MCCcode:     mccCode,
		TotalAmount: totalAmount,
		Merchant:    merchant,
	}, nil
}
