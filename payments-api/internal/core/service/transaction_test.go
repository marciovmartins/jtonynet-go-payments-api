package service

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
)

var (
	accountUIDtoTransact, _ = uuid.FromString("123e4567-e89b-12d3-a456-426614174000")
)

type DBfake struct {
	Account     map[int]port.AccountEntity
	Balance     map[int]port.BalanceByCategoryEntity
	Transaction map[int]port.TransactionEntity
}

type BalanceRegisterInDB struct {
	ID        int
	UID       uuid.UUID
	AccountId int
	Category  string
	Amount    decimal.Decimal
}

func newDBfake() DBfake {
	db := DBfake{}

	db.Transaction = make(map[int]port.TransactionEntity)

	db.Account = map[int]port.AccountEntity{
		1: {
			ID:  1,
			UID: accountUIDtoTransact,
		},
	}

	amountFood, _ := decimal.NewFromString("50.00")
	amountMeal, _ := decimal.NewFromString("75.00")
	amountCash, _ := decimal.NewFromString("100.00")

	db.Balance = map[int]port.BalanceByCategoryEntity{
		1: {
			ID:        1,
			AccountID: 1,
			Amount:    amountFood,
			Category:  port.CategoryFood,
		},
		2: {
			ID:        2,
			AccountID: 1,
			Amount:    amountMeal,
			Category:  port.CategoryMeal,
		},
		3: {
			ID:        3,
			AccountID: 1,
			Amount:    amountCash,
			Category:  port.CategoryCash,
		},
	}

	return db
}

func (dbf *DBfake) AccountRepoFindByUID(uid uuid.UUID) (port.AccountEntity, error) {
	for _, ae := range dbf.Account {
		if ae.UID == uid {
			return ae, nil
		}
	}

	return port.AccountEntity{}, fmt.Errorf("account with AccountUID %s not found", uid.String())
}

type AccountRepoFake struct {
	db DBfake
}

func newAccountRepoFake(db DBfake) port.AccountRepository {
	return &AccountRepoFake{
		db,
	}
}

func (arf *AccountRepoFake) FindByUID(uid uuid.UUID) (port.AccountEntity, error) {
	accountEntity, err := arf.db.AccountRepoFindByUID(uid)
	return accountEntity, err
}

type BalanceRepoFake struct {
	db DBfake
}

func newBalanceRepoFake(db DBfake) port.BalanceRepository {
	return &BalanceRepoFake{
		db,
	}
}

type TransactionRepoFake struct {
	db DBfake
}

func newTransactionRepoFake(db DBfake) port.TransactionRepository {
	return &TransactionRepoFake{
		db,
	}
}

func (dbf *DBfake) BalanceRepoFindByAccountID(accountID int) (port.BalanceEntity, error) {
	categories := make(map[int]port.BalanceByCategoryEntity)
	amountTotal := decimal.NewFromInt(0)

	counter := 1
	for _, be := range dbf.Balance {
		if be.AccountID == accountID {
			amountTotal = amountTotal.Add(be.Amount)
			categories[be.Category.Order] = be
			counter++
		}
	}

	if len(categories) == 0 {
		return port.BalanceEntity{}, fmt.Errorf("balances with AccountID %v not found", accountID)
	}

	b := port.BalanceEntity{
		AccountID:   accountID,
		AmountTotal: amountTotal,
		Categories:  categories,
	}

	return b, nil
}

func (brf *BalanceRepoFake) FindByAccountID(accountID int) (port.BalanceEntity, error) {
	balanceEntity, err := brf.db.BalanceRepoFindByAccountID(accountID)
	return balanceEntity, err
}

func (dbf *DBfake) BalanceRepoUpdate(bCategory port.BalanceByCategoryEntity) bool {
	b, exists := dbf.Balance[bCategory.ID]
	if !exists {
		return false
	}

	b.Amount = bCategory.Amount
	dbf.Balance[bCategory.ID] = b

	return true
}

func (brf *BalanceRepoFake) Update(be port.BalanceEntity) error {
	for _, bCategory := range be.Categories {
		ok := brf.db.BalanceRepoUpdate(bCategory)
		if !ok {
			return fmt.Errorf("balances with ID %v not found to update", bCategory.ID)
		}
	}

	return nil
}

func (dbf *DBfake) TransactionRepoSave(te port.TransactionEntity) bool {
	nextID := len(dbf.Transaction) + 1
	dbf.Transaction[nextID] = te

	return true
}

func (trf *TransactionRepoFake) Save(te port.TransactionEntity) error {
	ok := trf.db.TransactionRepoSave(te)
	if !ok {
		return fmt.Errorf("transaction with AccountUID %v not save", te.AccountUID)
	}

	return nil
}

func TestTransactionService(t *testing.T) {
	//Arrange
	dbFake := newDBfake()

	repos := repository.Repos{}
	repos.Account = newAccountRepoFake(dbFake)
	repos.Balance = newBalanceRepoFake(dbFake)
	repos.Transaction = newTransactionRepoFake(dbFake)

	//Act
	tService := NewTransaction(repos)

	totalAmount, _ := decimal.NewFromString("15.00")
	tRequest := port.TransactionRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: totalAmount,
		MCCcode:     "5411",
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	tService.HandleTransaction(tRequest)

	//Assert
	/*
		- Return Code is equal "00"
		- Balance is updated
		- Transaction was registered
	*/

}
