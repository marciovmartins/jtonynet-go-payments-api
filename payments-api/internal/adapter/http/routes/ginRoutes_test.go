package ginRoutes

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/config"
	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
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
)

type GinRoutesSuite struct {
	suite.Suite

	router   *gin.Engine
	apiGroup *gin.RouterGroup
}

func (suite *GinRoutesSuite) SetupSuite() {
	cfg, err := config.LoadConfig("./../../../../", ".env.TEST")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatal("cannot initiate app: ", err)
	}

	suite.loadDBtestData(app.Conn)
	suite.router, suite.apiGroup = setupRouterAndGroup(cfg.API, app)

	suite.apiGroup.POST("/payment", ginHandler.PaymentExecution)
}

func (suite *GinRoutesSuite) loadDBtestData(conn port.DBConn) {
	switch conn.GetStrategy() {
	case "gorm":
		db := conn.GetDB()
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

func setupRouterAndGroup(cfg config.API, app bootstrap.App) (*gin.Engine, *gin.RouterGroup) {
	basePath := "/"

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(ginMiddleware.ConfigInject(cfg))
	router.Use(ginMiddleware.AppInject(app))

	return router, router.Group(basePath)
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(GinRoutesSuite))
}

func (suite *GinRoutesSuite) TestPaymentExecuteTransactionApproved() {
	codeApproved := "00" // constants.CODE_APPROVED

	transactionJSON := fmt.Sprintf(
		`{
  			"account": "%s",
  			"mcc": "5411",
  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
  			"totalAmount": 100.09
		}`,
		accountUID,
	)
	reqPaymentExecution, err := http.NewRequest("POST", "/payment", bytes.NewBuffer([]byte(transactionJSON)))
	assert.NoError(suite.T(), err)

	respPaymentExecution := httptest.NewRecorder()
	suite.router.ServeHTTP(respPaymentExecution, reqPaymentExecution)
	assert.Equal(suite.T(), http.StatusOK, respPaymentExecution.Code)

	bodyRespPaymentExecution := respPaymentExecution.Body.String()
	assert.Equal(suite.T(), gjson.Get(bodyRespPaymentExecution, "code").String(), codeApproved)
}

func (suite *GinRoutesSuite) TestPaymentExecuteTransactionRejectedInsufficientFunds() {
	codeRejectedInsufficientFunds := "51" // constants.CODE_REJECTED_INSUFICIENT_FUNDS

	transactionJSON := fmt.Sprintf(
		`{
  			"account": "%s",
  			"mcc": "5411",
  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
  			"totalAmount": 9999.99
		}`,
		accountUID,
	)
	reqPaymentExecution, err := http.NewRequest("POST", "/payment", bytes.NewBuffer([]byte(transactionJSON)))
	assert.NoError(suite.T(), err)

	respPaymentExecution := httptest.NewRecorder()
	suite.router.ServeHTTP(respPaymentExecution, reqPaymentExecution)
	assert.Equal(suite.T(), http.StatusOK, respPaymentExecution.Code)

	bodyRespPaymentExecution := respPaymentExecution.Body.String()
	assert.Equal(suite.T(), gjson.Get(bodyRespPaymentExecution, "code").String(), codeRejectedInsufficientFunds)
}
