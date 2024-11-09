package router

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/http/router/ginStrategy"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func New(cfg config.Router, app bootstrap.App) (port.Router, error) {
	switch cfg.Strategy {
	case "gin":
		return ginStrategy.New(cfg, app)
	default:
		return nil, fmt.Errorf("router strategy not suported: %s", cfg.Strategy)
	}
}
