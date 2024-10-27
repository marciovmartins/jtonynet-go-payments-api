package gormRepos

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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

type BalanceResult struct {
	AccountID  uint
	BalanceID  uint
	BalanceUID uuid.UUID
	Amount     decimal.Decimal
	Name       string
	Priority   int
	Codes      string
}

func (b *Balance) FindByAccountID(accountID uint) (port.BalanceEntity, error) {
	var be port.BalanceEntity
	var bResults []BalanceResult

	err := b.db.Debug().Table("balances AS b").
		Select("b.account_id, b.id AS balance_id, b.uid AS balance_uid, b.amount, c.name AS category_name, c.priority, STRING_AGG(mc.mcc_code, ',') AS codes").
		Joins("JOIN categories AS c ON b.category_id = c.id").
		Joins("LEFT JOIN mcc_codes AS mc ON c.id = mc.category_id").
		Where("b.account_id = ?", accountID).
		Group("b.account_id, b.id, b.uid, b.amount, c.name, c.priority").
		Scan(&bResults).Error
	if err != nil {
		return port.BalanceEntity{}, fmt.Errorf("error retrying balance with id: %v", accountID)
	}

	if len(bResults) > 0 {
		amountTotal := decimal.NewFromInt(0)
		balanceCategories := make(map[int]port.BalanceByCategoryEntity)

		for _, bResult := range bResults {
			mccCodes := []string{}
			if bResult.Codes != "" {
				mccCodes = strings.Split(bResult.Codes, ",")
			}

			balanceCategories[bResult.Priority] = port.BalanceByCategoryEntity{
				ID:     bResult.BalanceID,
				UID:    bResult.BalanceUID,
				Amount: bResult.Amount,
				Category: port.CategoryEntity{
					Name:     bResult.Name,
					MccCodes: mccCodes,

					Order:    bResult.Priority,
					Priority: bResult.Priority,
				},
			}

			amountTotal = amountTotal.Add(bResult.Amount)
		}

		be.AccountID = accountID
		be.AmountTotal = amountTotal
		be.Categories = balanceCategories

		if len(be.Categories) > 0 {
			return be, nil
		}

		return port.BalanceEntity{}, fmt.Errorf("balance with id: %v not found", accountID)
	}

	return port.BalanceEntity{}, nil
}

func (b *Balance) UpdateTotalAmount(be port.BalanceEntity) error {
	var bindParameters []string
	var bindValues []interface{}

	for _, balanceCategory := range be.Categories {
		bindParameters = append(bindParameters, "(?::int, ?::numeric)")
		bindValues = append(bindValues, balanceCategory.ID, balanceCategory.Amount)
	}

	query := fmt.Sprintf(
		`UPDATE balances AS b
		 SET amount = v.new_amount
		 FROM (VALUES %s) AS v(id, new_amount)
		 WHERE b.id = v.id`,
		strings.Join(bindParameters, ","),
	)

	if err := b.db.Exec(query, bindValues...).Error; err != nil {
		return fmt.Errorf("error performing update balances: %w query: %s", err, query)
	}

	return nil
}
