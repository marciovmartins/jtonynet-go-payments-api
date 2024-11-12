package inMemoryDatabase

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"

	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase/redisConn"

	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewConn(cfg config.InMemoryDatabase) (port.InMemoryDBConn, error) {
	switch cfg.Strategy {
	case "redis":
		return redisConn.New(cfg)
	default:
		return nil, fmt.Errorf("InMemoryDB strategy not suported: %s", cfg.Strategy)
	}
}
