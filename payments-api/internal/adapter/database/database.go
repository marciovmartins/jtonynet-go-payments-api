package database

import (
	"errors"

	"github.com/jtonynet/go-payments-api/config"
	dbStrategy "github.com/jtonynet/go-payments-api/internal/adapter/database/strategies"
)

type Conn interface {
	GetDB() interface{}
	Readiness() error
}

func NewDB(cfg config.Database) (Conn, error) {
	switch cfg.Strategy {
	case "gorm":
		return dbStrategy.NewGormDB(cfg)
	default:
		return nil, errors.New("database strategy not suported: " + cfg.Strategy)
	}
}
