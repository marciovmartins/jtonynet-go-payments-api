package repository

import (
	"fmt"
	"log"
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

	balanceFoodAmount, _ = decimal.NewFromString("205.11")
	balanceMealAmount, _ = decimal.NewFromString("110.22")
	balanceCashAmount, _ = decimal.NewFromString("115.33")

	amountFoodTransaction, _ = decimal.NewFromString("100.10")
	MccCodeFoodTransaction   = "5411"
	merchantFoodTransaction  = "PADARIA DO ZE               SAO PAULO BR"

	merchantNameToMap         = "UBER EATS                   SAO PAULO BR"
	MerchantIncorrectMccToMap = "5555"
	MerchantCorrectMccToMap   = "5412"
)

type RepositoriesSuite struct {
	suite.Suite

	AccountRepo     port.AccountRepository
	BalanceRepo     port.BalanceRepository
	TransactionRepo port.TransactionRepository
	MerchantMapRepo port.MerchantMaptRepository

	AccountEntity port.AccountEntity
	BalanceEntity port.BalanceEntity
}

func (suite *RepositoriesSuite) SetupSuite() {
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

	repositories, err := GetAll(conn)
	if err != nil {
		log.Fatalf("error when instantiating repositories: %v", err)
	}

	suite.AccountRepo = repositories.Account
	suite.BalanceRepo = repositories.Balance
	suite.TransactionRepo = repositories.Transaction
	suite.MerchantMapRepo = repositories.MerchanMap

	suite.loadDBtestData(conn)
}

func (suite *RepositoriesSuite) loadDBtestData(conn port.DBConn) {
	switch conn.GetStrategy() {
	case "gorm":
		db := conn.GetDB()
		dbGorm, ok := db.(gorm.DB)
		if !ok {
			log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
		}

		dbGorm.Exec("TRUNCATE TABLE merchant_maps RESTART IDENTITY CASCADE")
		insertMerchantMapQuery := fmt.Sprintf(`
			INSERT INTO merchant_maps (uid, merchant_name, mcc_code, mapped_mcc_code, created_at, updated_at)
			VALUES
				('95abe1ff-6f67-4a17-a4eb-d4842e324f1f', '%s', '%s', '%s', NOW(), NOW()),
				('a53c6a52-8a18-4e7d-8827-7f612233c7ec', 'PAG*JoseDaSilva          RIO DE JANEI BR', '5555', '5812', NOW(), NOW())`,
			merchantNameToMap,
			MerchantIncorrectMccToMap,
			MerchantCorrectMccToMap,
		)
		dbGorm.Exec(insertMerchantMapQuery)

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

func (suite *RepositoriesSuite) AccountRepositoryFindByUIDsuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.Equal(suite.T(), accountEntity.ID, uint(1))
	assert.NoError(suite.T(), err)

	suite.AccountEntity = accountEntity
}

func (suite *RepositoriesSuite) BalanceRepositoryFindByAccountIDsuccess() {
	amountTotal := decimal.Sum(
		balanceFoodAmount,
		balanceMealAmount,
		balanceCashAmount,
	)

	accountEntity := suite.AccountEntity

	balanceEntity, err := suite.BalanceRepo.FindByAccountID(accountEntity.ID)
	assert.Equal(suite.T(), balanceEntity.AmountTotal, amountTotal)
	assert.NoError(suite.T(), err)

	suite.BalanceEntity = balanceEntity
}

func (suite *RepositoriesSuite) BalanceRepositoryUpdateTotalAmountSuccess() {
	balanceEntity := suite.BalanceEntity

	foodBalanceCategory := balanceEntity.Categories[port.CategoryFood.Order]
	foodBalanceCategory.Amount = foodBalanceCategory.Amount.Sub(amountFoodTransaction)
	balanceEntity.Categories[port.CategoryFood.Order] = foodBalanceCategory

	balanceEntityUpdateErr := suite.BalanceRepo.UpdateTotalAmount(balanceEntity)
	assert.NoError(suite.T(), balanceEntityUpdateErr)
}

func (suite *RepositoriesSuite) TransactionRepositorySaveSuccess() {
	accountEntity := suite.AccountEntity

	transactionEntity := port.TransactionEntity{
		AccountID:   accountEntity.ID,
		MccCode:     MccCodeFoodTransaction,
		Merchant:    merchantFoodTransaction,
		TotalAmount: amountFoodTransaction,
	}

	transactionEntitySaveErr := suite.TransactionRepo.Save(transactionEntity)
	assert.NoError(suite.T(), transactionEntitySaveErr)

}

func (suite *RepositoriesSuite) MerchantMapRepositoryFindByMerchantName() {
	merchantMapEntity, err := suite.MerchantMapRepo.FindByMerchantName(merchantNameToMap)
	assert.Equal(suite.T(), merchantMapEntity.MccCode, MerchantIncorrectMccToMap)
	assert.Equal(suite.T(), merchantMapEntity.MappedMccCode, MerchantCorrectMccToMap)
	assert.NoError(suite.T(), err)

}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}

func (suite *RepositoriesSuite) TestCases() {
	suite.T().Run("Test AccountRepository Find By UID Success", func(t *testing.T) {
		suite.AccountRepositoryFindByUIDsuccess()
	})
	suite.T().Run("Test Balance Repository Find By AccountID Success", func(t *testing.T) {
		suite.BalanceRepositoryFindByAccountIDsuccess()
	})

	suite.T().Run("Test Balance Repository Update TotalAmount Success", func(t *testing.T) {
		suite.BalanceRepositoryUpdateTotalAmountSuccess()
	})

	suite.T().Run("Test Transaction Repository Save Success", func(t *testing.T) {
		suite.TransactionRepositorySaveSuccess()
	})
}
