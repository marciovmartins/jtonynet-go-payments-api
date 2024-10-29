package bootstrap

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"

	"github.com/jtonynet/go-payments-api/internal/core/service"
)

type App struct {
	PaymentService *service.Payment
}

func NewApp(cfg *config.Config) (App, error) {
	app := App{}

	conn, err := database.NewConn(cfg.Database)
	if err != nil {
		return App{}, fmt.Errorf("error connecting to database: %v", err)
	}

	if conn.Readiness() != nil {
		return App{}, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("successfully connected to the database!")

	allRepos, err := repository.GetAll(conn)
	if err != nil {
		return App{}, fmt.Errorf("error when instantiating Account repository: %v", err)
	}

	app.PaymentService = service.NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.MerchanMap,
	)

	return app, nil
}
