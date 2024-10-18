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
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	accountUID, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
)

type RepositoriesSuite struct {
	suite.Suite

	Conn        port.DBConn
	AccountRepo port.AccountRepository
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

	suite.dbChargeData(conn)
}

func (suite *RepositoriesSuite) dbChargeData(conn port.DBConn) {
	switch conn.GetStrategy() {
	case "gorm":
		db := conn.GetDB()
		dbGorm, ok := db.(gorm.DB)
		if !ok {
			log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
		}

		dbGorm.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")
		dbGorm.Exec(
			fmt.Sprintf("INSERT INTO accounts (name, uid, created_at, updated_at) VALUES('John Doe', '%s', NOW(), NOW())", accountUID),
		)

	default:
		log.Fatalf("error connecting to database migrate to charge test data")
	}
}

func (suite *RepositoriesSuite) TearDownSuite() {
	os.Unsetenv("ENV")
}

func (suite *RepositoriesSuite) TestAccountRepository() {
	accountEntity, err := suite.AccountRepo.FindByUID(accountUID)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), accountEntity.ID, 1)
}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}
