package database

import (
	"errors"

	"github.com/jtonynet/go-payments-api/config"
	dbStrategy "github.com/jtonynet/go-payments-api/internal/adapter/database/strategies"
	port "github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewConn(cfg config.Database) (port.DBConn, error) {
	switch cfg.Strategy {
	case "gorm":
		return dbStrategy.NewGormConn(cfg)
	default:
		return nil, errors.New("database conn strategy not suported: " + cfg.Strategy)
	}
}
