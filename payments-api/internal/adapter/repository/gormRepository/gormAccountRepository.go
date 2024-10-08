package repository

import (
	"fmt"
	"log"

	database "github.com/jtonynet/go-payments-api/internal/adapter/database"
	dbStrategy "github.com/jtonynet/go-payments-api/internal/adapter/database/strategies"
	model "github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type GormAccountRepository struct {
	Conn dbStrategy.GormDB
}

func NewGormAccountRepository(conn database.Conn) port.AccountRepository {

	gormConn, ok := conn.(dbStrategy.GormDB)
	if !ok {
		log.Fatal("Failed to assert conn as GormDB")
	} else {
		fmt.Println("Deu bom nessa parada!")
	}

	//---
	var accounts []model.GormAccountModel

	db := gormConn.GetDB()

	err := db.Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failure on database connection: %w", err)
	}
	//---

	return GormAccountRepository{Conn: gormConn}
}

// []domain.Account
func (ga GormAccountRepository) FindAll() ([]model.GormAccountModel, error) {
	// var accounts []model.GormAccountModel

	// db := ga.Conn.GetDB()

	// err := db.Find(&accounts).Error
	// if err != nil {
	// 	return nil, fmt.Errorf("failure on database connection: %w", err)
	// }

	// return accounts, err
	return nil, nil
}
