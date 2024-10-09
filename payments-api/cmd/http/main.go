package main

import (
	"log"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"

	/*
		TODO:
		In the hexagonal approach I could use echoRoutes, muxRoutes and others  with  correct  adapters.
		I decided to use Gin because the router engine is less sensitive to changes (but can be changed)
		Common sense indicates that agnostic databases/sources/queues and memory banks are changed more
		frequently, taking better advantage of hexagonal
	*/
	ginRoutes "github.com/jtonynet/go-payments-api/internal/adapter/handler/routes"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatal("cannot initiate app: ", err)
	}

	ginRoutes.GinHandleRequests(
		cfg.API,
		app,
	)
}
