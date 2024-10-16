package repository

import (
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository/gormRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Repos struct {
	Account     port.AccountRepository
	Balance     port.BalanceRepository
	Transaction port.TransactionRepository
}

func GetAll(conn port.DBConn) (Repos, error) {
	repos := Repos{}

	strategy := conn.GetStrategy()
	switch strategy {
	case "gorm":
		account, err := gormRepos.NewAccount(conn)
		if err != nil {
			return Repos{}, fmt.Errorf("error when instantiating Gorm Account repository: %v", err)
		}
		repos.Account = account

		return repos, nil
	default:
		return Repos{}, errors.New("repository strategy not suported: " + strategy)
	}
}
