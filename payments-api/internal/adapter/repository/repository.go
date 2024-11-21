package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository/gormRepos"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository/redisRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type AllRepos struct {
	Account  port.AccountRepository
	Merchant port.MerchantRepository
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

func NewCachedMerchant(cacheConn inMemoryDatabase.Client, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	var mr port.MerchantRepository

	strategy, err := cacheConn.GetStrategy(context.TODO())
	if err != nil {
		return mr, fmt.Errorf("error: dont retrieve cache strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewRedisMerchant(cacheConn, mRepository)
	default:
		return mr, fmt.Errorf("cached repository strategy not suported: %s", strategy)
	}
}

func NewMemoryLock(lockConn inMemoryDatabase.Client) (port.MemoryLockRepository, error) {
	var mlr port.MemoryLockRepository

	strategy, err := lockConn.GetStrategy(context.TODO())
	if err != nil {
		return mlr, fmt.Errorf("error: dont retrieve memory lock strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewMemoryLock(lockConn)
	default:
		return mlr, fmt.Errorf("memory lock repository strategy not suported: %s", strategy)
	}
}
