package repository

import (
	"errors"
	"fmt"

	gormStrategy "github.com/jtonynet/go-payments-api/internal/adapter/repository/strategies/gormRepos"
	port "github.com/jtonynet/go-payments-api/internal/core/port"
)

type Repos struct {
	Account port.AccountRepository
}

func GetRepos(conn port.DBConn) (Repos, error) {
	repos := Repos{}

	strategy := conn.GetStrategy()
	switch strategy {
	case "gorm":
		accountRepo, err := gormStrategy.NewGormAccount(conn)
		if err != nil {
			return Repos{}, fmt.Errorf("error when instantiating Gorm Account repository: %v", err)
		}

		repos.Account = accountRepo

		return repos, nil
	default:
		return Repos{}, errors.New("repository strategy not suported: " + strategy)
	}
}
