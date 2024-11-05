package cachedRepository

import (
	"log"

	"github.com/jtonynet/go-payments-api/internal/core/port"

	"github.com/tidwall/gjson"
)

type Merchant struct {
	cacheConn          port.Cache
	merchantRepository port.MerchantRepository
}

func NewMerchant(cacheConn port.Cache, mRepository port.MerchantRepository) *Merchant {
	return &Merchant{
		cacheConn:          cacheConn,
		merchantRepository: mRepository,
	}
}

func (m *Merchant) FindByName(name string) (*port.MerchantEntity, error) {
	var mEntity *port.MerchantEntity

	merchantCached, err := m.cacheConn.Get(name)
	if err != nil {
		mEntity, err = m.merchantRepository.FindByName(name)
		if err != nil {
			return mEntity, err
		}
		m.cacheConn.Set(name, mEntity, m.cacheConn.GetDefaultExpiration())

		log.Printf("MISS CACHE. RETRIEVE FROM DB `%s`", name)
	} else {
		mEntity = &port.MerchantEntity{
			Name:          gjson.Get(merchantCached, "Name").String(),
			MccCode:       gjson.Get(merchantCached, "MccCode").String(),
			MappedMccCode: gjson.Get(merchantCached, "MappedMccCode").String(),
		}

		log.Printf("CACHED! RETRIEVE FROM CACHE  `%s`", name)
	}

	return mEntity, nil
}
