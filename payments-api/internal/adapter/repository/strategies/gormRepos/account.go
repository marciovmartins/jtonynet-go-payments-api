package gormRepos

import (
	"fmt"

	"github.com/gofrs/uuid"
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
	dbGorm, ok := db.(gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failure to cast conn.GetDB() as gorm.DB")
	}

	return Account{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

func (gar Account) FindByUUID(accountUUID uuid.UUID) (port.AccountDTORepository, error) {
	return port.AccountDTORepository{}, nil
}

// TODO: to remove, test purpose
func (gar Account) FindAll() ([]gormModel.Account, error) {
	var accounts []gormModel.Account
	err := gar.db.Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failure on fetch account list: %w", err)
	}

	return accounts, nil
}
