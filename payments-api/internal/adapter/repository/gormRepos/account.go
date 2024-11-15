package gormRepos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"gorm.io/gorm"
)

type Account struct {
	gormConn database.Conn
	db       *gorm.DB
}

func NewGormAccount(conn database.Conn) (port.AccountRepository, error) {
	db, err := conn.GetDB(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("account repository failure on conn.GetDB()")
	}

	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("account repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Account{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (a *Account) FindByUID(_ context.Context, uid uuid.UUID) (port.AccountEntity, error) {
	accountModel := gormModel.Account{}
	if err := a.db.Where(&gormModel.Account{UID: uid}).First(&accountModel).Error; err != nil {
		return port.AccountEntity{}, fmt.Errorf("account with uid: %s not found", uid)
	}

	return port.AccountEntity{
		ID:  accountModel.ID,
		UID: accountModel.UID,
	}, nil
}
