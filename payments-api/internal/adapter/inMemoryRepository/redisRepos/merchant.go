package redisRepos

import (
	"context"

	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryDatabase"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"github.com/tidwall/gjson"
)

type Merchant struct {
	cacheConn inMemoryDatabase.Conn

	merchantRepository port.MerchantRepository
}

func NewMerchant(cacheConn inMemoryDatabase.Conn, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	return &Merchant{
		cacheConn:          cacheConn,
		merchantRepository: mRepository,
	}, nil
}

func (m *Merchant) FindByName(_ context.Context, name string) (*port.MerchantEntity, error) {
	var mEntity *port.MerchantEntity

	merchantCached, err := m.cacheConn.Get(context.TODO(), name)
	if err != nil {
		mEntity, err = m.merchantRepository.FindByName(context.TODO(), name)
		if err != nil {
			return mEntity, err
		}

		defaultExpiration, err := m.cacheConn.GetDefaultExpiration(context.TODO())
		if err != nil {
			return mEntity, err
		}

		err = m.cacheConn.Set(context.TODO(), name, mEntity, defaultExpiration)
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
