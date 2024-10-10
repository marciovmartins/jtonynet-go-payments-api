package gormRepos

import "github.com/jtonynet/go-payments-api/internal/core/port"

type GormTransactionRepository struct{}

func NewGormTransactionRepository() (port.TransactionRepository, error) {
	return GormTransactionRepository{}, nil
}

func (gtr GormTransactionRepository) Save(tEntity port.TransactionEntity) error {
	return nil
}
