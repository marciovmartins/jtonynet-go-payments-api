package database

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
)

type DBConn interface {
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	GetDB(ctx context.Context) (interface{}, error)
	GetDriver(ctx context.Context) (string, error)
}

func NewConn(cfg config.Database) (DBConn, error) {
	switch cfg.Strategy {
	case "gorm":
		return NewGormConn(cfg)
	default:
		return nil, fmt.Errorf("database conn strategy not suported: %s", cfg.Strategy)
	}
}
