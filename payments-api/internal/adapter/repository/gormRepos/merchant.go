package gormRepos

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"gorm.io/gorm"
)

type Merchant struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewMerchant(conn port.DBConn) (port.MerchantRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(gorm.DB)
	if !ok {
		return nil, fmt.Errorf("account repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Merchant{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

func (m *Merchant) FindByName(Name string) (port.MerchantEntity, error) {
	MerchantModel := gormModel.Merchant{}
	if err := m.db.Where(&gormModel.Merchant{Name: Name}).First(&MerchantModel).Error; err != nil {
		return port.MerchantEntity{}, fmt.Errorf("Merchant with uid: %s not found", Name)
	}

	return port.MerchantEntity{
		Name:          MerchantModel.Name,
		MccCode:       MerchantModel.MccCode,
		MappedMccCode: MerchantModel.MappedMccCode,
	}, nil
}
