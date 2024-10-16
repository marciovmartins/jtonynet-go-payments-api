package domain

import (
	"github.com/gofrs/uuid"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	AccountUID  uuid.UUID
	MCCcode     string
	TotalAmount decimal.Decimal
	Merchant    string
}

func NewTransactionFromRequest(tRequest port.TransactionRequest) (*Transaction, error) {
	return &Transaction{
		AccountUID:  tRequest.AccountUID,
		MCCcode:     tRequest.MCCcode,
		TotalAmount: tRequest.TotalAmount,
		Merchant:    tRequest.Merchant,
	}, nil
}
