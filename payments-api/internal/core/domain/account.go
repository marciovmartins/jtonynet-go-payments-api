package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID  uint
	UID uuid.UUID
}

func (a *Account) NewTransaction(
	mccCode string,
	totalAmount decimal.Decimal,
	mDomain Merchant,
) *Transaction {
	correctMccCode := mccCode
	if mDomain.MappedMccCode != "" {
		correctMccCode = mDomain.MappedMccCode
	}

	return &Transaction{
		AccountID:   a.ID,
		AccountUID:  a.UID,
		MccCode:     correctMccCode,
		TotalAmount: totalAmount,
		Merchant:    mDomain.Name,
	}
}
