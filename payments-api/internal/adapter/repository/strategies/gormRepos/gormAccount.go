package repository

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"gorm.io/gorm"
)

type GormAccountRepository struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewGormAccount(conn port.DBConn) (port.AccountRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failure to cast conn.GetDB() as gorm.DB")
	}

	return GormAccountRepository{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

// []domain.Account || []dto.Account
func (gar GormAccountRepository) FindAll() ([]gormModel.Account, error) {
	var accounts []gormModel.Account
	err := gar.db.Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failure on fetch account list: %w", err)
	}

	return accounts, nil
}

// domain.Account  || dto.Account
func (gar GormAccountRepository) FindByUUID(accountUUID uuid.UUID) (gormModel.Account, error) {
	var account gormModel.Account
	if err := gar.db.Where(&gormModel.Account{UUID: accountUUID}).First(&account).Error; err != nil {
		return account, fmt.Errorf("failure on fetch account: %w", err)
	}

	return account, nil
}
