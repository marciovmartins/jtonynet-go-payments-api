package domain

import (
	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Account struct {
	ID  uint
	UID uuid.UUID
}

func NewAccountFromEntity(aEntity port.AccountEntity) (Account, error) {
	return Account{
		ID:  aEntity.ID,
		UID: aEntity.UID,
	}, nil
}
