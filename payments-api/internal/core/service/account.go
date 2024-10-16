package service

import (
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
