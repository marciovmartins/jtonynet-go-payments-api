package repository

import (
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository/gormRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type AllRepos struct {
	Account     port.AccountRepository
	Balance     port.BalanceRepository
	Transaction port.TransactionRepository
	MerchanMap  port.MerchantMaptRepository
}

func GetAll(conn port.DBConn) (AllRepos, error) {
	repos := AllRepos{}

	strategy := conn.GetStrategy()
	switch strategy {
	case "gorm":
		account, err := gormRepos.NewAccount(conn)
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

		merchantMap, err := gormRepos.NewMerchantMap(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating merchantMap repository: %v", err)
		}
		repos.MerchanMap = merchantMap

		return repos, nil
	default:
		return AllRepos{}, errors.New("repository strategy not suported: " + strategy)
	}
}
