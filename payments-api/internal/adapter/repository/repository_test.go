package repository

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	accountUID, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	balanceFoodUID, _ = uuid.Parse("bf64460c-efe0-4ffb-bd1f-54b136c2b2ac")
	balanceMealUID, _ = uuid.Parse("19475f4b-ee4c-4bce-add1-df0db5908201")
	balanceCashUID, _ = uuid.Parse("389e9316-ce28-478e-b14e-f971812de22d")

	balanceFoodAmount, _ = decimal.NewFromString("105.11")
	balanceMealAmount, _ = decimal.NewFromString("110.22")
	balanceCashAmount, _ = decimal.NewFromString("115.33")

	amountFoodTransaction, _ = decimal.NewFromString("100.10")
	mccCodeFoodTransaction   = "5411"
	merchantFoodTransaction  = "PADARIA DO ZE               SAO PAULO BR"
)

type RepositoriesSuite struct {
	suite.Suite

	Conn            port.DBConn
	AccountRepo     port.AccountRepository
	BalanceRepo     port.BalanceRepository
	TransactionRepo port.TransactionRepository
}

func (suite *RepositoriesSuite) SetupSuite() {
	os.Setenv("ENV", "test")

	cfg, err := config.LoadConfig("./../../../")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	conn, err := database.NewConn(cfg.Database)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	if conn.Readiness() != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	suite.Conn = conn

	repositories, err := GetAll(conn)
	if err != nil {
		log.Fatalf("error when instantiating repositories: %v", err)
	}

	suite.AccountRepo = repositories.Account
	suite.BalanceRepo = repositories.Balance
	suite.TransactionRepo = repositories.Transaction
}

func (suite *RepositoriesSuite) TearDownSuite() {
	os.Unsetenv("ENV")
}

func (suite *RepositoriesSuite) SetupTest() {
	suite.dbChargeData()
}

func (suite *RepositoriesSuite) dbChargeData() {
	switch suite.Conn.GetStrategy() {
	case "gorm":
		db := suite.Conn.GetDB()
		dbGorm, ok := db.(gorm.DB)
		if !ok {
			log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
		}

		dbGorm.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")
		insertAccountQuery := fmt.Sprintf(
			`INSERT INTO accounts (uid, name, created_at, updated_at) 
			 VALUES('%s', 'Jonh Doe', NOW(), NOW())`,
			accountUID)
		dbGorm.Exec(insertAccountQuery)

		dbGorm.Exec("TRUNCATE TABLE balances RESTART IDENTITY CASCADE")
		insertBalancesQuery := fmt.Sprintf(`
			INSERT INTO balances (uid, account_id, amount, category_name, created_at, updated_at)
			VALUES
				('%s', 1, %v, '%s', NOW(), NOW()),
				('%s', 1, %v, '%s', NOW(), NOW()),
				('%s', 1, %v, '%s', NOW(), NOW())`,
			balanceFoodUID, balanceFoodAmount, port.CATEGORY_FOOD_NAME,
			balanceMealUID, balanceMealAmount, port.CATEGORY_MEAL_NAME,
			balanceCashUID, balanceCashAmount, port.CATEGORY_CASH_NAME,
		)
		dbGorm.Exec(insertBalancesQuery)

		dbGorm.Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE")

	default:
		log.Fatalf("error connecting to database migrate to charge test data")
	}
}

func (suite *RepositoriesSuite) TestAccountRepositoryFindByUIDsuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.Equal(suite.T(), accountEntity.ID, uint(1))
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesSuite) TestBalanceRepositoryFindByAccountIDsuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.Equal(suite.T(), accountEntity.ID, uint(1))
	assert.NoError(suite.T(), err)

	amountTotal := decimal.Sum(
		balanceFoodAmount,
		balanceMealAmount,
		balanceCashAmount,
	)

	balanceEntity, err := suite.BalanceRepo.FindByAccountID(accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, amountTotal)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesSuite) TestBalanceRepositoryUpdateTotalAmountsuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.Equal(suite.T(), accountEntity.ID, uint(1))
	assert.NoError(suite.T(), err)

	amountTotal := decimal.Sum(
		balanceFoodAmount,
		balanceMealAmount,
		balanceCashAmount,
	)

	balanceEntity, err := suite.BalanceRepo.FindByAccountID(accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, amountTotal)
	assert.NoError(suite.T(), err)

	foodBalanceCategory := balanceEntity.Categories[port.CategoryFood.Order]
	foodBalanceCategory.Amount = foodBalanceCategory.Amount.Sub(amountFoodTransaction)
	balanceEntity.Categories[port.CategoryFood.Order] = foodBalanceCategory

	balanceEntityUpdateErr := suite.BalanceRepo.UpdateTotalAmount(balanceEntity)
	assert.NoError(suite.T(), balanceEntityUpdateErr)
}

func (suite *RepositoriesSuite) TestTransactionRepositorySaveSuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.Equal(suite.T(), accountEntity.ID, uint(1))
	assert.NoError(suite.T(), err)

	amountTotal := decimal.Sum(
		balanceFoodAmount,
		balanceMealAmount,
		balanceCashAmount,
	)

	balanceEntity, err := suite.BalanceRepo.FindByAccountID(accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, amountTotal)
	assert.NoError(suite.T(), err)

	foodBalanceCategory := balanceEntity.Categories[port.CategoryFood.Order]
	foodBalanceCategory.Amount = foodBalanceCategory.Amount.Sub(amountFoodTransaction)
	balanceEntity.Categories[port.CategoryFood.Order] = foodBalanceCategory

	balanceEntityUpdateErr := suite.BalanceRepo.UpdateTotalAmount(balanceEntity)
	assert.NoError(suite.T(), balanceEntityUpdateErr)

	transactionEntity := port.TransactionEntity{
		AccountID:   accountEntity.ID,
		MCCcode:     mccCodeFoodTransaction,
		Merchant:    merchantFoodTransaction,
		TotalAmount: amountFoodTransaction,
	}

	transactionEntitySaveErr := suite.TransactionRepo.Save(transactionEntity)
	assert.NoError(suite.T(), transactionEntitySaveErr)

}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}
