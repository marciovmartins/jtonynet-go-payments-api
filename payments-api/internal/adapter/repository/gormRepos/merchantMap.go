package gormRepos

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"gorm.io/gorm"
)

type MerchantMap struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewMerchantMap(conn port.DBConn) (port.MerchantMaptRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(gorm.DB)
	if !ok {
		return nil, fmt.Errorf("account repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &MerchantMap{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

func (m *MerchantMap) FindByMerchantName(merchantName string) (port.MerchantMapEntity, error) {
	merchantMapModel := gormModel.MerchantMap{}
	if err := m.db.Where(&gormModel.MerchantMap{MerchantName: merchantName}).First(&merchantMapModel).Error; err != nil {
		return port.MerchantMapEntity{}, fmt.Errorf("merchantMap with uid: %s not found", merchantName)
	}

	return port.MerchantMapEntity{
		MerchantName:  merchantMapModel.MerchantName,
		MccCode:       merchantMapModel.MccCode,
		MappedMccCode: merchantMapModel.MappedMccCode,
	}, nil
}
