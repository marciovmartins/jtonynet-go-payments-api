package gormRepos

import (
	"fmt"

	"github.com/google/uuid"
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

	return &Account{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

func (gar *Account) FindByUID(accountUID uuid.UUID) (port.AccountEntity, error) {
	return port.AccountEntity{}, nil
}
