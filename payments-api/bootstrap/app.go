package bootstrap

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/support"
	"github.com/jtonynet/go-payments-api/internal/support/logger"

	"github.com/jtonynet/go-payments-api/internal/core/service"
)

type App struct {
	Logger         support.Logger
	PaymentService *service.Payment
}

func NewApp(cfg *config.Config) (App, error) {
	app := App{}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		return App{}, fmt.Errorf("error when instantiate logger: %v", err)
	}
	app.Logger = logger

	conn, err := database.NewConn(cfg.Database)
	if err != nil {
		return App{}, fmt.Errorf("error connecting to database: %v", err)
	}

	if conn.Readiness() != nil {
		return App{}, fmt.Errorf("error connecting to database: %v", err)
	}

	logger.Debug("successfully connected to the database!")

	allRepos, err := repository.GetAll(conn)
	if err != nil {
		return App{}, fmt.Errorf("error when instantiate repositories: %v", err)
	}

	app.PaymentService = service.NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.MerchanMap,
		logger,
	)

	return app, nil
}
