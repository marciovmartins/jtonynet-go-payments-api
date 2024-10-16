package domain

import (
	"github.com/gofrs/uuid"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Account struct {
	ID  int
	UID uuid.UUID
}

func NewAccountFromEntity(aEntity port.AccountEntity) (Account, error) {
	return Account{
		ID:  aEntity.ID,
		UID: aEntity.UID,
	}, nil
}
