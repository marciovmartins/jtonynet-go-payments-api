package cache

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"

	"github.com/jtonynet/go-payments-api/internal/adapter/cache/redisStrategy"

	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func New(cfg config.Cache) (port.Cache, error) {
	switch cfg.Strategy {
	case "redis":
		return redisStrategy.New(cfg)
	default:
		return nil, fmt.Errorf("cache strategy not suported: %s", cfg.Strategy)
	}
}
