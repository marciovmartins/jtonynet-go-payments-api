package cachedRepository

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/cachedRepository/redisRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewMerchant(conn port.Cache, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	var mr port.MerchantRepository

	strategy, err := conn.GetStrategy(context.Background())
	if err != nil {
		return mr, fmt.Errorf("error: dont retrieve cache strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewMerchant(conn, mRepository)
	default:

		return mr, fmt.Errorf("cached repository strategy not suported: %s", strategy)
	}
}
