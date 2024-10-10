package repository

import (
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository/strategies/gormRepos"
	port "github.com/jtonynet/go-payments-api/internal/core/port"
)

type Repos struct {
	Account     port.AccountRepository
	Balance     port.BalanceRepository
	Transaction port.TransactionRepository
}

func GetRepos(conn port.DBConn) (Repos, error) {
	repos := Repos{}

	strategy := conn.GetStrategy()
	switch strategy {
	case "gorm":
		account, err := gormRepos.NewAccount(conn)
		if err != nil {
			return Repos{}, fmt.Errorf("error when instantiating Gorm Account repository: %v", err)
		}
		repos.Account = account

		/*
			balanceRepo, err := gormStrategy.NewGormBalanceRepository()
			if err != nil {
				// TODO Marcio pediu para ignorar mas meu toc nao deixa!
				return Repos{}, fmt.Errorf("error when instantiating Gorm Account repository: %v", err)
			}
			repos.Balance = balanceRepo

			transactionRepo, err := gormStrategy.NewGormTransactionRepository()
			if err != nil {
				// TODO Marcio pediu para ignorar mas meu toc nao deixa!
				return Repos{}, fmt.Errorf("error when instantiating Gorm Account repository: %v", err)
			}
			repos.Transaction = transactionRepo
		*/

		return repos, nil
	default:
		return Repos{}, errors.New("repository strategy not suported: " + strategy)
	}
}
