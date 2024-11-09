package gormRepos

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"gorm.io/gorm"
)

type Account struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewAccount(conn port.DBConn) (port.AccountRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("account repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Account{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (a *Account) FindByUID(uid uuid.UUID) (port.AccountEntity, error) {
	accountModel := gormModel.Account{}
	if err := a.db.Where(&gormModel.Account{UID: uid}).First(&accountModel).Error; err != nil {
		return port.AccountEntity{}, fmt.Errorf("account with uid: %s not found", uid)
	}

	return port.AccountEntity{
		ID:  accountModel.ID,
		UID: accountModel.UID,
	}, nil
}
