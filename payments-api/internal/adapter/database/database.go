package database

import (
	"errors"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database/strategies/gormConn"
	port "github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewConn(cfg config.Database) (port.DBConn, error) {
	switch cfg.Strategy {
	case "gorm":
		return gormConn.New(cfg)
	default:
		return nil, errors.New("database conn strategy not suported: " + cfg.Strategy)
	}
}
