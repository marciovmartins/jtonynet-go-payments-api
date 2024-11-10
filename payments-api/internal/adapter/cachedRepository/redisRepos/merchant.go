package redisRepos

import (
	"context"

	"github.com/jtonynet/go-payments-api/internal/core/port"

	"github.com/tidwall/gjson"
)

type Merchant struct {
	redisConn port.Cache

	merchantRepository port.MerchantRepository
}

func NewMerchant(conn port.Cache, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	return &Merchant{
		redisConn:          conn,
		merchantRepository: mRepository,
	}, nil
}

func (m *Merchant) FindByName(name string) (*port.MerchantEntity, error) {
	var mEntity *port.MerchantEntity

	merchantCached, err := m.redisConn.Get(context.Background(), name)
	if err != nil {
		mEntity, err = m.merchantRepository.FindByName(name)
		if err != nil {
			return mEntity, err
		}

		defaultExpiration, err := m.redisConn.GetDefaultExpiration(context.Background())
		if err != nil {
			return mEntity, err
		}

		err = m.redisConn.Set(context.Background(), name, mEntity, defaultExpiration)
		if err != nil {
			return mEntity, err
		}

	} else {
		mEntity = &port.MerchantEntity{
			Name: gjson.Get(merchantCached, "Name").String(),
			MCC:  gjson.Get(merchantCached, "MCC").String(),
		}

	}

	return mEntity, nil
}
