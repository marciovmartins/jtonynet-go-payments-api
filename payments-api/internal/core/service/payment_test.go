package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

var (
	accountUIDtoTransact, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	balanceFoodAmount, _ = decimal.NewFromString("205.11")
	balanceMealAmount, _ = decimal.NewFromString("310.22")
	balanceCashAmount, _ = decimal.NewFromString("415.33")

	correctFoodMCC = "5411"

	amountFoodFundsApproved, _         = decimal.NewFromString("100.10")
	amountFoodFundsRejected, _         = decimal.NewFromString("720.45")
	amountFoodFundsFallbackApproved, _ = decimal.NewFromString("320.00")
)

type DBfake struct {
	Account     map[uint]port.AccountEntity
	Balance     map[uint]port.BalanceByCategoryEntity
	Transaction map[uint]port.TransactionEntity
}

func newDBfake() DBfake {
	db := DBfake{}

	db.Transaction = make(map[uint]port.TransactionEntity)

	db.Account = map[uint]port.AccountEntity{
		1: {
			ID:  1,
			UID: accountUIDtoTransact,
		},
	}

	db.Balance = map[uint]port.BalanceByCategoryEntity{
		1: {
			ID:        1,
			AccountID: 1,
			Amount:    balanceFoodAmount,
			Category:  port.CategoryFood,
		},
		2: {
			ID:        2,
			AccountID: 1,
			Amount:    balanceMealAmount,
			Category:  port.CategoryMeal,
		},
		3: {
			ID:        3,
			AccountID: 1,
			Amount:    balanceCashAmount,
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

func (dbf *DBfake) BalanceRepoFindByAccountID(accountID uint) (port.BalanceEntity, error) {
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

func (brf *BalanceRepoFake) FindByAccountID(accountID uint) (port.BalanceEntity, error) {
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

func (brf *BalanceRepoFake) UpdateTotalAmount(be port.BalanceEntity) error {
	for _, bCategory := range be.Categories {
		ok := brf.db.BalanceRepoUpdate(bCategory)
		if !ok {
			return fmt.Errorf("balances with ID %v not found to update", bCategory.ID)
		}
	}

	return nil
}

func (dbf *DBfake) TransactionRepoSave(te port.TransactionEntity) bool {
	nextID := uint(len(dbf.Transaction) + 1)
	te.ID = nextID
	te.UID, _ = uuid.NewUUID()
	dbf.Transaction[nextID] = te

	return true
}

func (dbf *DBfake) TransactionRepoFindLastByAcountId(accountID uint) (port.TransactionEntity, error) {
	var lastTransaction port.TransactionEntity
	found := false
	maxKey := uint(0)

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

type PaymentSuite struct {
	suite.Suite
}

func (suite *PaymentSuite) SetupSuite() {

}

func (suite *PaymentSuite) getDBfake() *DBfake {
	dbFake := newDBfake()
	return &dbFake
}

func (suite *PaymentSuite) getAllRepositories(dbFake *DBfake) *repository.AllRepos {
	allRepos := repository.AllRepos{}
	allRepos.Account = newAccountRepoFake(*dbFake)
	allRepos.Balance = newBalanceRepoFake(*dbFake)
	allRepos.Transaction = newTransactionRepoFake(*dbFake)

	return &allRepos
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsRejected() {
	//Arrange
	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsRejected,
		MccCode:     correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
	)

	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	// - Payment execution with received transaction has been rejected
	codeRejected := "51" // domain.CODE_REJECTED_INSUFICIENT_FUNDS
	insuficientFundsError := "failed to approve balance domain: balance category has insuficient funds"

	assert.Equal(suite.T(), returnCode, codeRejected)
	assert.Equal(suite.T(), err.Error(), insuficientFundsError)
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsApproved() {
	//Arrange
	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MccCode:     correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	accountEntity, err := allRepos.Account.FindByUID(tRequest.AccountUID)
	assert.NoError(suite.T(), err)

	//Act
	paymentService := NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	// - Balance is updated
	expectedAmountTotal, _ := decimal.NewFromString("830.56")
	balanceEntity, err := allRepos.Balance.FindByAccountID(accountEntity.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, expectedAmountTotal)

	balanceDomain, err := mapBalanceEntityToDomain(balanceEntity)
	assert.NoError(suite.T(), err)

	expectedAmountCategory, _ := decimal.NewFromString("105.01")
	balanceCategoryMCC, err := balanceDomain.GetByMCC(correctFoodMCC)
	assert.Equal(suite.T(), balanceCategoryMCC.Amount, expectedAmountCategory)
	assert.NoError(suite.T(), err)

	expectedAmountFallback, _ := decimal.NewFromString("415.33")
	balanceCategoryFallback, err := balanceDomain.GetFallback()
	assert.Equal(suite.T(), balanceCategoryFallback.Amount, expectedAmountFallback)
	assert.NoError(suite.T(), err)

	// - Transaction was registered
	transactionByAcountId, _ := dbFake.TransactionRepoFindLastByAcountId(accountEntity.ID)
	assert.Equal(suite.T(), transactionByAcountId.TotalAmount, amountFoodFundsApproved)
}

func (suite *PaymentSuite) TestL2PaymentExecuteCorrectMCCFallbackApproved() {
	//Arrange
	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsFallbackApproved,
		MccCode:     correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	accountEntity, err := allRepos.Account.FindByUID(tRequest.AccountUID)
	assert.NoError(suite.T(), err)

	//Act
	paymentService := NewPayment(
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	// - Balance is updated
	expectedAmountTotal := decimal.NewFromFloat(610.66)
	balanceEntity, err := allRepos.Balance.FindByAccountID(accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, expectedAmountTotal)
	assert.NoError(suite.T(), err)

	balanceDomain, err := mapBalanceEntityToDomain(balanceEntity)
	assert.NoError(suite.T(), err)

	expectedAmountCategory := decimal.NewFromFloat(0)
	balanceCategoryMCC, err := balanceDomain.GetByMCC(correctFoodMCC)
	assert.Equal(suite.T(), balanceCategoryMCC.Amount, expectedAmountCategory)
	assert.NoError(suite.T(), err)

	expectedAmountFallback := decimal.NewFromFloat(300.44)
	balanceCategoryFallback, err := balanceDomain.GetFallback()
	assert.Equal(suite.T(), balanceCategoryFallback.Amount, expectedAmountFallback)
	assert.NoError(suite.T(), err)

	// - Transaction was registered
	transactionByAcountId, _ := dbFake.TransactionRepoFindLastByAcountId(accountEntity.ID)
	assert.Equal(suite.T(), transactionByAcountId.TotalAmount, amountFoodFundsFallbackApproved)
}

func TestPaymentSuite(t *testing.T) {
	suite.Run(t, new(PaymentSuite))
}
