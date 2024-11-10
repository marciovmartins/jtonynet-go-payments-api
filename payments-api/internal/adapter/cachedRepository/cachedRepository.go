package cachedRepository

import (
	"errors"

	"github.com/jtonynet/go-payments-api/internal/adapter/cachedRepository/redisRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewMerchant(conn port.Cache, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	strategy := conn.GetStrategy()
	switch strategy {
	case "redis":
		return redisRepos.NewMerchant(conn, mRepository)
	default:
		var mr port.MerchantRepository
		return mr, errors.New("cached repository strategy not suported: " + strategy)
	}
}
