package port

import (
	"context"

	"github.com/jtonynet/go-payments-api/config"
)

type Router interface {
	HandleRequests(ctx context.Context, cfg config.API) error
}
