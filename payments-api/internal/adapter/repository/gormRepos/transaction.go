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

type Transaction struct {
	gormConn database.Conn
	db       *gorm.DB
}

func NewTransaction(conn database.Conn) (port.TransactionRepository, error) {
	db, err := conn.GetDB(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("transaction repository failure on conn.GetDB()")
	}

	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("transaction repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Transaction{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (t *Transaction) Save(_ context.Context, tEntity port.TransactionEntity) error {
	transactionModel := &gormModel.Transaction{
		UID:         uuid.New(),
		AccountID:   tEntity.AccountID,
		MCC:         tEntity.MCC,
		Merchant:    tEntity.Merchant,
		TotalAmount: tEntity.TotalAmount,
	}

	if err := t.db.Create(&transactionModel).Error; err != nil {
		return fmt.Errorf("error saving transaction: %s", err)
	}

	return nil
}
