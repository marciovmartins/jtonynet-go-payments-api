package domain

import "github.com/shopspring/decimal"

type Merchant struct {
	Name string
	MCC  string
}

func (m *Merchant) NewTransaction(
	mcc string,
	totalAmount decimal.Decimal,
	name string,
	account Account,
) *Transaction {
	correctMCC := mcc
	if m.MCC != "" {
		correctMCC = m.MCC
	}

	return &Transaction{
		AccountID:   account.ID,
		AccountUID:  account.UID,
		MCC:         correctMCC,
		TotalAmount: totalAmount,
		Merchant:    name,
	}
}
