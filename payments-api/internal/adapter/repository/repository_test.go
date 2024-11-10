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

	balanceFoodAmount = decimal.NewFromFloat(205.11)
	balanceMealAmount = decimal.NewFromFloat(110.22)
	balanceCashAmount = decimal.NewFromFloat(115.33)

	amountFoodTransaction   = decimal.NewFromFloat(100.10)
	MccCodeFoodTransaction  = "5411"
	merchantFoodTransaction = "PADARIA DO ZE               SAO PAULO BR"

	merchant1NameToMap         = "UBER EATS                   SAO PAULO BR"
	merchant1IncorrectMccToMap = "5555"
	merchant1CorrectMccToMap   = "5412"

	merchant2NameToMap         = "PAG*JoseDaSilva          RIO DE JANEI BR"
	merchant2IncorrectMccToMap = "5555"
	merchant2CorrectMccToMap   = "5812"
)

type RepositoriesSuite struct {
	suite.Suite

	AccountRepo     port.AccountRepository
	BalanceRepo     port.BalanceRepository
	TransactionRepo port.TransactionRepository
	MerchantRepo    port.MerchantRepository

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
	suite.MerchantRepo = repositories.Merchant

	suite.loadDBtestData(conn)
}

func (suite *RepositoriesSuite) loadDBtestData(conn port.DBConn) {
	switch conn.GetStrategy() {
	case "gorm":
		db := conn.GetDB()
		dbGorm, ok := db.(*gorm.DB)
		if !ok {
			log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
		}

		dbGorm.Exec("TRUNCATE TABLE merchants RESTART IDENTITY CASCADE")
		insertMerchantQuery := fmt.Sprintf(`
			INSERT INTO merchants (uid, name, mcc_code, mapped_mcc_code, created_at, updated_at)
			VALUES
				('95abe1ff-6f67-4a17-a4eb-d4842e324f1f', '%s', '%s', '%s', NOW(), NOW()),
				('a53c6a52-8a18-4e7d-8827-7f612233c7ec', '%s', '%s', '%s', NOW(), NOW())`,
			merchant1NameToMap,
			merchant1IncorrectMccToMap,
			merchant1CorrectMccToMap,
			merchant2NameToMap,
			merchant2IncorrectMccToMap,
			merchant2CorrectMccToMap,
		)
		dbGorm.Exec(insertMerchantQuery)

		dbGorm.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")
		insertAccountQuery := fmt.Sprintf(
			`INSERT INTO accounts (uid, name, created_at, updated_at) 
			 VALUES('%s', 'Jonh Doe', NOW(), NOW())`,
			accountUID)
		dbGorm.Exec(insertAccountQuery)

		dbGorm.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE")
		insertCategoryQuery := `
			INSERT INTO categories (uid, name, priority, created_at, updated_at)
			VALUES
				('5681b4b5-6176-498a-a856-8932f79c05cc', 'FOOD', 1, NOW(), NOW()),
				('7bcfcd2a-2fde-4564-916b-92410e794272', 'MEAL', 2, NOW(), NOW()),
				('056de185-bff0-4c4a-93fa-7245f9e72b67', 'CASH', 3, NOW(), NOW())
		`
		dbGorm.Exec(insertCategoryQuery)

		dbGorm.Exec("TRUNCATE TABLE mcc_codes RESTART IDENTITY CASCADE")
		insertMccCodeQuery := `
			INSERT INTO mcc_codes (uid, mcc_code, category_id, created_at, updated_at)
			VALUES
				('11f0c06e-0dff-4643-86bf-998d11e9374f', '5411', 1, NOW(), NOW()),
				('fe5a4c17-a7cd-4072-a793-e99e2642e21a', '5412', 1, NOW(), NOW()),
				('5268ec2b-aa14-4d55-906a-13c91d89826c', '5811', 2, NOW(), NOW()),
				('6179e57c-e630-4e2f-a5db-d153e0cdb9a9', '5812', 2, NOW(), NOW())
		`
		dbGorm.Exec(insertMccCodeQuery)

		dbGorm.Exec("TRUNCATE TABLE balances RESTART IDENTITY CASCADE")
		insertBalancesQuery := fmt.Sprintf(`
			INSERT INTO balances (uid, account_id, amount, category_id, created_at, updated_at)
			VALUES
				('%s', 1, %v, 1, NOW(), NOW()),
				('%s', 1, %v, 2, NOW(), NOW()),
				('%s', 1, %v, 3, NOW(), NOW())`,
			balanceFoodUID, balanceFoodAmount,
			balanceMealUID, balanceMealAmount,
			balanceCashUID, balanceCashAmount,
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

	foodCategoryPriority := 1

	foodBalanceCategory := balanceEntity.Categories[foodCategoryPriority]
	foodBalanceCategory.Amount = foodBalanceCategory.Amount.Sub(amountFoodTransaction)
	balanceEntity.Categories[foodCategoryPriority] = foodBalanceCategory

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

func (suite *RepositoriesSuite) MerchantRepositoryFindByName() {
	merchantEntity, err := suite.MerchantRepo.FindByName(merchant1NameToMap)
	assert.Equal(suite.T(), merchantEntity.MccCode, merchant1IncorrectMccToMap)
	assert.Equal(suite.T(), merchantEntity.MappedMccCode, merchant1CorrectMccToMap)
	assert.NoError(suite.T(), err)
}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}

func (suite *RepositoriesSuite) TestCases() {
	suite.T().Run("TestAccountRepositoryFindByUIDSuccess", func(t *testing.T) {
		suite.AccountRepositoryFindByUIDsuccess()
	})

	suite.T().Run("TestBalanceRepositoryFindByAccountIDSuccess", func(t *testing.T) {
		suite.BalanceRepositoryFindByAccountIDsuccess()
	})

	suite.T().Run("TestBalanceRepositoryUpdateTotalAmountSuccess", func(t *testing.T) {
		suite.BalanceRepositoryUpdateTotalAmountSuccess()
	})

	suite.T().Run("TestTransactionRepositorySaveSuccess", func(t *testing.T) {
		suite.TransactionRepositorySaveSuccess()
	})

	suite.T().Run("TestMerchantRepositoryFindByName", func(t *testing.T) {
		suite.MerchantRepositoryFindByName()
	})
}
