package gormRepos

import (
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type GormBalanceRepository struct{}

func NewGormBalanceRepository() (GormBalanceRepository, error) {
	return GormBalanceRepository{}, nil
}

func (gbr *GormBalanceRepository) FindByAccountUUID(accountUUID string) (port.BalanceDTORepository, error) {
	return port.BalanceDTORepository{}, nil
}

func (gbr *GormBalanceRepository) Update(port.BalanceDTORepository) error {
	return nil
}
