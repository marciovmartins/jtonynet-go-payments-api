package gormRepos

import (
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type GormBalanceRepository struct{}

func NewGormBalanceRepository() (GormBalanceRepository, error) {
	return GormBalanceRepository{}, nil
}

func (gbr *GormBalanceRepository) FindByAccountID(id int) ([]port.BalanceEntity, error) {
	return nil, nil
}

func (gbr *GormBalanceRepository) Update(port.BalanceEntity) error {
	return nil
}
