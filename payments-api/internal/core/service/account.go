package service

import (
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
)

type Account struct {
	repos repository.Repos
}

func NewAccount(repos repository.Repos) Account {
	return Account{
		repos,
	}
}

func (as Account) RetrieveList() ([]gormModel.Account, error) {
	accountRepo := as.repos.Account
	list, err := accountRepo.FindAll()

	return list, err
}
