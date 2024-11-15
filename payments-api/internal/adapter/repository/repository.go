package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository/gormRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type AllRepos struct {
	Account     port.AccountRepository
	Balance     port.BalanceRepository
	Transaction port.TransactionRepository
	Merchant    port.MerchantRepository
}

func GetAll(conn database.Conn) (AllRepos, error) {
	repos := AllRepos{}

	strategy, err := conn.GetStrategy(context.TODO())
	if err != nil {
		return AllRepos{}, fmt.Errorf("error when instantiating merchant repository: %v", err)
	}

	switch strategy {
	case "gorm":
		account, err := gormRepos.NewGormAccount(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating account repository: %v", err)
		}
		repos.Account = account

		balance, err := gormRepos.NewBalance(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating balance repository: %v", err)
		}
		repos.Balance = balance

		transaction, err := gormRepos.NewTransaction(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating transaction repository: %v", err)
		}
		repos.Transaction = transaction

		Merchant, err := gormRepos.NewMerchant(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating merchant repository: %v", err)
		}
		repos.Merchant = Merchant

		return repos, nil
	default:
		return AllRepos{}, errors.New("repository strategy not suported: " + strategy)
	}
}
