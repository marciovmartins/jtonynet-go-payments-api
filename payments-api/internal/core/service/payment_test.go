package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

var (
	timeoutSLAcfg = 100

	accountUIDtoTransact, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	balanceFoodAmount = decimal.NewFromFloat(205.11)
	balanceMealAmount = decimal.NewFromFloat(310.22)
	balanceCashAmount = decimal.NewFromFloat(415.33)

	correctFoodMCC = "5411"

	amountFoodFundsApproved         = decimal.NewFromFloat(100.10)
	amountFoodFundsRejected         = decimal.NewFromFloat(720.45)
	amountFoodFundsFallbackApproved = decimal.NewFromFloat(320.00)
)

type DBfake struct {
	Account     map[uint]port.AccountEntity
	Balance     map[uint]port.BalanceByCategoryEntity
	Transaction map[uint]port.TransactionEntity
	Merchant    map[uint]port.MerchantEntity
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

	db.Merchant = map[uint]port.MerchantEntity{
		1: {
			Name: "UBER EATS                   SAO PAULO BR",
			MCC:  "5412",
		},
	}

	db.Balance = map[uint]port.BalanceByCategoryEntity{
		1: {
			ID:        1,
			AccountID: 1,
			Amount:    balanceFoodAmount,
			Category:  port.CategoryEntity{Name: "FOOD", MCCs: []string{"5411", "5412"}, Priority: 1},
		},
		2: {
			ID:        2,
			AccountID: 1,
			Amount:    balanceMealAmount,
			Category:  port.CategoryEntity{Name: "MEAL", MCCs: []string{"5811", "5812"}, Priority: 2},
		},
		3: {
			ID:        3,
			AccountID: 1,
			Amount:    balanceCashAmount,
			Category:  port.CategoryEntity{Name: "CASH", Priority: 3},
		},
	}

	return db
}

func (dbf *DBfake) GetDB() *DBfake {
	return dbf
}

func (dbf *DBfake) AccountRepoFindByUID(_ context.Context, uid uuid.UUID) (port.AccountEntity, error) {
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

func (arf *AccountRepoFake) FindByUID(_ context.Context, uid uuid.UUID) (port.AccountEntity, error) {
	accountEntity, err := arf.db.AccountRepoFindByUID(context.TODO(), uid)
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

type MerchantRepoFake struct {
	db DBfake
}

func newMerchantRepoFake(db DBfake) port.MerchantRepository {
	return &MerchantRepoFake{
		db,
	}
}

func (m *MerchantRepoFake) FindByName(_ context.Context, name string) (*port.MerchantEntity, error) {
	MerchantEntity, err := m.db.MerchantRepoFindByName(name)
	return MerchantEntity, err
}

func (dbf *DBfake) MerchantRepoFindByName(Name string) (*port.MerchantEntity, error) {

	for _, m := range dbf.Merchant {
		if m.Name == Name {
			return &m, nil
		}
	}

	return nil, nil
}

func (dbf *DBfake) BalanceRepoFindByAccountID(accountID uint) (port.BalanceEntity, error) {
	categories := make(map[int]port.BalanceByCategoryEntity)
	amountTotal := decimal.NewFromInt(0)

	for _, be := range dbf.Balance {
		if be.AccountID == accountID {
			amountTotal = amountTotal.Add(be.Amount)
			categories[be.Category.Priority] = be
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

func (brf *BalanceRepoFake) FindByAccountID(_ context.Context, accountID uint) (port.BalanceEntity, error) {
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

func (brf *BalanceRepoFake) UpdateTotalAmount(_ context.Context, be port.BalanceEntity) error {
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

func (trf *TransactionRepoFake) Save(_ context.Context, te port.TransactionEntity) error {
	ok := trf.db.TransactionRepoSave(te)
	if !ok {
		return fmt.Errorf("transaction with AccountID %v not save", te.AccountID)
	}

	return nil
}

type InMemoryDBfake struct {
	Lock map[string]string
}

func newInMemoryDBfake() InMemoryDBfake {
	imdbf := InMemoryDBfake{}
	imdbf.Lock = make(map[string]string)

	return imdbf
}

func (suite *PaymentSuite) getInMemoryDBfake() InMemoryDBfake {
	inMemorydbFake := newInMemoryDBfake()
	return inMemorydbFake
}

func (imdbf *InMemoryDBfake) MemoryLockRepoLock(_ context.Context, mle port.MemoryLockEntity) (port.MemoryLockEntity, error) {
	imdbf.Lock[mle.Key] = strconv.FormatInt(mle.Timestamp, 10)
	return mle, nil
}

func (imdbf *InMemoryDBfake) MemoryLockRepoUnlock(_ context.Context, key string) error {
	if _, exists := imdbf.Lock[key]; exists {
		delete(imdbf.Lock, key)
		return nil
	}
	return nil
}

type MemoryLockRepoFake struct {
	memoryDB InMemoryDBfake
}

func newMemoryLockRepoFake(memoryDB InMemoryDBfake) port.MemoryLockRepository {
	return &MemoryLockRepoFake{
		memoryDB,
	}
}

func (m *MemoryLockRepoFake) Lock(_ context.Context, timeoutSLA port.TimeoutSLA, mle port.MemoryLockEntity) (port.MemoryLockEntity, error) {
	return m.memoryDB.MemoryLockRepoLock(context.TODO(), mle)
}

func (m *MemoryLockRepoFake) Unlock(_ context.Context, key string) error {
	return m.memoryDB.MemoryLockRepoUnlock(context.TODO(), key)
}

type PaymentSuite struct {
	suite.Suite
}

func (suite *PaymentSuite) getDBfake() *DBfake {
	dbFake := newDBfake()
	return &dbFake
}

func (suite *PaymentSuite) getMemoryLockRepoFake(memoryDB InMemoryDBfake) port.MemoryLockRepository {
	return newMemoryLockRepoFake(memoryDB)
}

func (suite *PaymentSuite) getAllRepositories(dbFake *DBfake) *repository.AllRepos {
	allRepos := repository.AllRepos{}
	allRepos.Account = newAccountRepoFake(*dbFake)
	allRepos.Balance = newBalanceRepoFake(*dbFake)
	allRepos.Transaction = newTransactionRepoFake(*dbFake)
	allRepos.Merchant = newMerchantRepoFake(*dbFake)

	return &allRepos
}

func (suite *PaymentSuite) TestL1PaymentExecuteGenericRejected() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := DBfake{}
	allRepos := suite.getAllRepositories(&dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)

	returnCode, _ := paymentService.Execute(tRequest)

	//Assert
	codeRejected := "07" // domain.CODE_REJECTED_GENERIC
	assert.Equal(suite.T(), returnCode, codeRejected)
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsRejected() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsRejected,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)

	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeRejected := "51" // domain.CODE_REJECTED_INSUFICIENT_FUNDS
	insuficientFundsError := "failed to approve balance domain: balance category has insuficient funds"

	assert.Equal(suite.T(), returnCode, codeRejected)
	assert.Equal(suite.T(), err.Error(), insuficientFundsError)
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	expectedAmountTotal := decimal.NewFromFloat(830.56)
	expectedAmountCategory := decimal.NewFromFloat(105.01)
	expectedAmountFallback := decimal.NewFromFloat(415.33)
	transactionAmount := amountFoodFundsApproved
	suite.assertBalanceAmounts(
		dbFake,
		allRepos,
		tRequest,
		expectedAmountTotal,
		expectedAmountCategory,
		expectedAmountFallback,
		transactionAmount,
	)
}

func (suite *PaymentSuite) TestL2PaymentExecuteCorrectMCCFallbackApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsFallbackApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	expectedAmountTotal := decimal.NewFromFloat(610.66)
	expectedAmountCategory := decimal.NewFromFloat(0)
	expectedAmountFallback := decimal.NewFromFloat(300.44)
	transactionAmount := amountFoodFundsFallbackApproved
	suite.assertBalanceAmounts(
		dbFake,
		allRepos,
		tRequest,
		expectedAmountTotal,
		expectedAmountCategory,
		expectedAmountFallback,
		transactionAmount,
	)
}

func (suite *PaymentSuite) TestL3PaymentExecuteNameMCCWithFundsApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         "5555",
		Merchant:    "UBER EATS                   SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	expectedAmountTotal := decimal.NewFromFloat(830.56)
	expectedAmountCategory := decimal.NewFromFloat(105.01)
	expectedAmountFallback := decimal.NewFromFloat(415.33)
	transactionAmount := amountFoodFundsApproved
	suite.assertBalanceAmounts(
		dbFake,
		allRepos,
		tRequest,
		expectedAmountTotal,
		expectedAmountCategory,
		expectedAmountFallback,
		transactionAmount,
	)
}

func (suite *PaymentSuite) TestL3PaymentExecuteNameMCCFallbackApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsFallbackApproved,
		MCC:         "5555",
		Merchant:    "UBER EATS                   SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Balance,
		allRepos.Transaction,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.NoError(suite.T(), err)

	expectedAmountTotal := decimal.NewFromFloat(610.66)
	expectedAmountCategory := decimal.NewFromFloat(0)
	expectedAmountFallback := decimal.NewFromFloat(300.44)
	transactionAmount := amountFoodFundsFallbackApproved
	suite.assertBalanceAmounts(
		dbFake,
		allRepos,
		tRequest,
		expectedAmountTotal,
		expectedAmountCategory,
		expectedAmountFallback,
		transactionAmount,
	)
}

func (suite *PaymentSuite) assertBalanceAmounts(
	dbFake *DBfake,
	allRepos *repository.AllRepos,
	tRequest port.TransactionPaymentRequest,
	expectedAmountTotal decimal.Decimal,
	expectedAmountCategory decimal.Decimal,
	expectedAmountFallback decimal.Decimal,
	transactionAmount decimal.Decimal,
) {
	accountEntity, err := allRepos.Account.FindByUID(context.TODO(), tRequest.AccountUID)
	assert.NoError(suite.T(), err)

	// - Balance is updated
	balanceEntity, err := allRepos.Balance.FindByAccountID(context.TODO(), accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, expectedAmountTotal)
	assert.NoError(suite.T(), err)

	balanceDomain, err := mapBalanceEntityToDomain(balanceEntity)
	assert.NoError(suite.T(), err)

	balanceCategoryMCC, err := balanceDomain.GetByMCC(correctFoodMCC)
	assert.Equal(suite.T(), balanceCategoryMCC.Amount, expectedAmountCategory)
	assert.NoError(suite.T(), err)

	balanceCategoryFallback, err := balanceDomain.GetFallback()
	assert.Equal(suite.T(), balanceCategoryFallback.Amount, expectedAmountFallback)
	assert.NoError(suite.T(), err)

	// - Transaction was registered
	transactionByAcountId, _ := dbFake.TransactionRepoFindLastByAcountId(accountEntity.ID)
	assert.Equal(suite.T(), transactionByAcountId.TotalAmount, transactionAmount)
}

func TestPaymentSuite(t *testing.T) {
	suite.Run(t, new(PaymentSuite))
}
