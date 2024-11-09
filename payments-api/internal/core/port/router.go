package port

import (
	"github.com/jtonynet/go-payments-api/config"
)

type Router interface {
	HandleRequests(cfg config.API)
}
