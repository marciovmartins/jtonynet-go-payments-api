package gormRepos

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BalanceResult struct {
	AccountID  uint
	BalanceID  uint
	BalanceUID uuid.UUID
	Amount     decimal.Decimal
	Name       string
	Priority   int
	Codes      sql.NullString
}

type Balance struct {
	gormConn database.DBConn
	db       *gorm.DB
}

func NewBalance(conn database.DBConn) (port.BalanceRepository, error) {
	db, err := conn.GetDB(context.Background())
	if err != nil {
		return nil, fmt.Errorf("balance repository failure on conn.GetDB()")
	}

	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("balance repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Balance{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (b *Balance) FindByAccountID(_ context.Context, accountID uint) (port.BalanceEntity, error) {
	var be port.BalanceEntity
	var bResults []BalanceResult
	firstFound := false

	err := b.db.Table("balances AS b").
		Select("b.id AS balance_id, b.uid AS balance_uid, b.amount, c.name, c.priority, STRING_AGG(mc.mcc, ',') AS codes").
		Joins("JOIN categories AS c ON b.category_id = c.id").
		Joins("LEFT JOIN mccs AS mc ON c.id = mc.category_id").
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
			mccs := []string{}
			if bResult.Codes.Valid {
				mccs = strings.Split(bResult.Codes.String, ",")
			}

			if !firstFound {
				firstFound = true
				be.AccountID = bResult.AccountID
			}

			balanceCategories[bResult.Priority] = port.BalanceByCategoryEntity{
				ID:     bResult.BalanceID,
				UID:    bResult.BalanceUID,
				Amount: bResult.Amount,
				Category: port.CategoryEntity{
					Name:     bResult.Name,
					MCCs:     mccs,
					Priority: bResult.Priority,
				},
			}

			amountTotal = amountTotal.Add(bResult.Amount)
		}

		if len(balanceCategories) > 0 {
			be.AmountTotal = amountTotal
			be.Categories = balanceCategories

			return be, nil
		}

		return port.BalanceEntity{}, fmt.Errorf("balance with id: %v not found", accountID)
	}

	return port.BalanceEntity{}, nil
}

func (b *Balance) UpdateTotalAmount(_ context.Context, be port.BalanceEntity) error {
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
