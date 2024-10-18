package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	accountUIDtoTransact, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
)

type PaymentSuite struct {
	suite.Suite
}

type DBfake struct {
	Account     map[int]port.AccountEntity
	Balance     map[int]port.BalanceByCategoryEntity
	Transaction map[int]port.TransactionEntity
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

	amountFood, _ := decimal.NewFromString("105.00")
	amountMeal, _ := decimal.NewFromString("110.00")
	amountCash, _ := decimal.NewFromString("115.00")

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

func (dbf *DBfake) GetDB() *DBfake {
	return dbf
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

	for _, be := range dbf.Balance {
		if be.AccountID == accountID {
			amountTotal = amountTotal.Add(be.Amount)
			categories[be.Category.Order] = be
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
	te.ID = nextID
	te.UID, _ = uuid.NewUUID()
	dbf.Transaction[nextID] = te

	return true
}

func (dbf *DBfake) TransactionRepoFindLastByAcountId(accountID int) (port.TransactionEntity, error) {
	var lastTransaction port.TransactionEntity
	found := false
	maxKey := 0

	for key, t := range dbf.Transaction {
		if t.AccountID == accountID && key > maxKey {
			lastTransaction = t
			maxKey = key
			found = true
		}
	}

	if !found {
		return port.TransactionEntity{}, fmt.Errorf("transaction with AccountID %v not found", accountID)
	}

	return lastTransaction, nil
}

func (trf *TransactionRepoFake) Save(te port.TransactionEntity) error {
	ok := trf.db.TransactionRepoSave(te)
	if !ok {
		return fmt.Errorf("transaction with AccountID %v not save", te.AccountID)
	}

	return nil
}

func (suite *PaymentSuite) TestPaymentExecuteSuccess() {
	//Arrange
	dbFake := newDBfake()

	repos := repository.Repos{}
	repos.Account = newAccountRepoFake(dbFake)
	repos.Balance = newBalanceRepoFake(dbFake)
	repos.Transaction = newTransactionRepoFake(dbFake)

	//Act
	paymentService := NewPayment(repos)

	amountTransaction, _ := decimal.NewFromString("100.00")
	tRequest := port.TransactionRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountTransaction,
		MCCcode:     "5411",
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	accountEntity, _ := repos.Account.FindByUID(tRequest.AccountUID)
	balanceEntityAfterTransact, _ := repos.Balance.FindByAccountID(accountEntity.ID)
	amountAfterTransact := balanceEntityAfterTransact.AmountTotal

	returnCode, cErr := paymentService.Execute(tRequest)

	//Assert
	// - Payment execution with received transaction has been approved
	returnCodeApproved := "00" // constants.CODE_APPROVED
	assert.Equal(suite.T(), returnCode, returnCodeApproved)
	assert.Equal(suite.T(), cErr, nil)

	// - Balance is updated
	balanceEntityBeforeTransact, _ := repos.Balance.FindByAccountID(accountEntity.ID)
	amountBeforeTransact := balanceEntityBeforeTransact.AmountTotal
	assert.Equal(suite.T(), amountBeforeTransact, amountAfterTransact.Sub(amountTransaction))

	// - Transaction was registered
	transactionByAcountId, _ := dbFake.TransactionRepoFindLastByAcountId(accountEntity.ID)
	assert.Equal(suite.T(), transactionByAcountId.TotalAmount, amountTransaction)
}

func TestPaymentSuite(t *testing.T) {
	suite.Run(t, new(PaymentSuite))
}
