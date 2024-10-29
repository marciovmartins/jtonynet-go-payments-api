package logger

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/support"
	"github.com/jtonynet/go-payments-api/internal/support/logger/slogStrategy"
)

func New(cfg config.Logger) (support.Logger, error) {
	switch cfg.Strategy {
	case "slog":
		return slogStrategy.New(cfg)
	default:
		return nil, fmt.Errorf("router strategy not suported: %s", cfg.Strategy)
	}
}
