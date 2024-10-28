package gormRepos

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"gorm.io/gorm"
)

type Transaction struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewTransaction(conn port.DBConn) (port.TransactionRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("transaction repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Transaction{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (t *Transaction) Save(tEntity port.TransactionEntity) error {
	transactionModel := &gormModel.Transaction{
		UID:         uuid.New(),
		AccountID:   tEntity.AccountID,
		MccCode:     tEntity.MccCode,
		Merchant:    tEntity.Merchant,
		TotalAmount: tEntity.TotalAmount,
	}

	if err := t.db.Create(&transactionModel).Error; err != nil {
		return fmt.Errorf("error saving transaction: %s", err)
	}

	return nil
}
