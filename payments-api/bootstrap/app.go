package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/jtonynet/go-payments-api/config"

	"github.com/jtonynet/go-payments-api/internal/support"
	"github.com/jtonynet/go-payments-api/internal/support/logger"

	"github.com/jtonynet/go-payments-api/internal/adapter/cache"
	"github.com/jtonynet/go-payments-api/internal/adapter/cachedRepository"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"

	"github.com/jtonynet/go-payments-api/internal/core/service"
)

type App struct {
	Logger support.Logger

	PaymentService *service.Payment
}

func NewApp(cfg *config.Config) (App, error) {
	app := App{}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Printf("warning: dont instantiate logger: %v", err)
	}
	app.Logger = logger

	cacheConn, err := cache.New(cfg.Cache)
	if err != nil {
		return App{}, fmt.Errorf("error: dont instantiate cache client: %v", err)
	}

	if cacheConn.Readiness(context.Background()) != nil {
		return App{}, fmt.Errorf("error: dont connecting to cache: %v", err)
	}

	if logger != nil {
		logger.Debug("successfully: connected to the cache!")
	}

	dbConn, err := database.NewConn(cfg.Database)
	if err != nil {
		return App{}, fmt.Errorf("error: dont instantiate database: %v", err)
	}

	if dbConn.Readiness(context.Background()) != nil {
		return App{}, fmt.Errorf("error: dont connecting to database: %v", err)
	}

	if logger != nil {
		logger.Debug("successfully: connected to the database!")
	}

	allRepos, err := repository.GetAll(dbConn)
	if err != nil {
		return App{}, fmt.Errorf("error: dont instantiate repositories: %v", err)
	}

	cachedMerchantRepo, err := cachedRepository.NewMerchant(
		cacheConn,
		allRepos.Merchant,
	)
	if err != nil {
		return App{}, fmt.Errorf("error: dont instantiate merchant cached repository: %v", err)
	}

	app.PaymentService = service.NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		cachedMerchantRepo,
		logger,
	)

	return app, nil
}
