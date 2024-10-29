package database

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database/gormConn"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewConn(cfg config.Database) (port.DBConn, error) {
	switch cfg.Strategy {
	case "gorm":
		return gormConn.New(cfg)
	default:
		return nil, fmt.Errorf("database conn strategy not suported: %s", cfg.Strategy)
	}
}
