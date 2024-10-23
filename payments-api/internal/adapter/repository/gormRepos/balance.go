package gormRepos

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Balance struct {
	gormConn port.DBConn
	db       *gorm.DB
}

func NewBalance(conn port.DBConn) (port.BalanceRepository, error) {
	db := conn.GetDB()
	dbGorm, ok := db.(gorm.DB)
	if !ok {
		return nil, fmt.Errorf("balance repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Balance{
		gormConn: conn,
		db:       &dbGorm,
	}, nil
}

func (b *Balance) FindByAccountID(accountID uint) (port.BalanceEntity, error) {
	balanceModelList := []gormModel.Balance{}
	if err := b.db.Find(&balanceModelList).Where(&gormModel.Balance{AccountID: accountID}).Error; err != nil {
		return port.BalanceEntity{}, fmt.Errorf("error retrying balance with id: %v", accountID)
	}

	be := port.BalanceEntity{}

	if len(balanceModelList) > 0 {
		amountTotal := decimal.NewFromInt(0)
		balanceCategories := make(map[int]port.BalanceByCategoryEntity)

		for _, balanceModel := range balanceModelList {
			balanceCategory, exists := port.Categories[balanceModel.CategoryName]
			if !exists {
				continue
			}

			balanceCategories[balanceCategory.Order] = port.BalanceByCategoryEntity{
				ID:       balanceModel.ID,
				Amount:   balanceModel.Amount,
				Category: balanceCategory,
			}

			amountTotal = amountTotal.Add(balanceModel.Amount)
		}

		be.AccountID = accountID
		be.AmountTotal = amountTotal
		be.Categories = balanceCategories

		if len(be.Categories) > 0 {
			return be, nil
		}
	}

	return port.BalanceEntity{}, fmt.Errorf("balance with id: %v not found", accountID)
}

func (b *Balance) UpdateTotalAmount(be port.BalanceEntity) error {
	query := "UPDATE balances SET amount = CASE"
	for _, balanceCategory := range be.Categories {
		query += fmt.Sprintf(" WHEN id = %v THEN %v", uint(balanceCategory.ID), decimal.Decimal(balanceCategory.Amount))
	}
	query += " END "

	if err := b.db.Exec(query).Error; err != nil {
		return fmt.Errorf("error performing update balances: %w query: %s", err, query)
	}

	return nil
}
