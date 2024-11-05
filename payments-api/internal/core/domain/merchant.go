package domain

import "github.com/shopspring/decimal"

type Merchant struct {
	Name          string
	MccCode       string
	MappedMccCode string
}

func (m *Merchant) NewTransaction(
	mcc string,
	totalAmount decimal.Decimal,
	account Account,
) *Transaction {
	correctMccCode := mcc
	if m.MappedMccCode != "" {
		correctMccCode = m.MappedMccCode
	}

	return &Transaction{
		AccountID:   account.ID,
		AccountUID:  account.UID,
		MccCode:     correctMccCode,
		TotalAmount: totalAmount,
		Merchant:    m.Name,
	}
}
