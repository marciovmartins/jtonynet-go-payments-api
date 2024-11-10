package redisRepos

import (
	"context"
	"log"
	"testing"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/cache"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var merchantName = "XYZ*TestCachedRepositoryMerchant                   PIRAPORINHA BR"

type RedisReposSuite struct {
	suite.Suite

	cacheConn          port.Cache
	cachedMerchantRepo port.MerchantRepository
}

type DBfake struct {
	Merchant map[uint]port.MerchantEntity
}

func newDBfake() DBfake {
	db := DBfake{}

	db.Merchant = map[uint]port.MerchantEntity{
		1: {
			Name: merchantName,
			MCC:  "5412",
		},
	}

	return db
}

type MerchantRepoFake struct {
	db DBfake
}

func newMerchantRepoFake(db DBfake) port.MerchantRepository {
	return &MerchantRepoFake{
		db,
	}
}

func (m *MerchantRepoFake) FindByName(Name string) (*port.MerchantEntity, error) {
	MerchantEntity, err := m.db.MerchantRepoFindByName(Name)
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

func (suite *RedisReposSuite) SetupSuite() {
	cfg, err := config.LoadConfig("./../../../../")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	cacheConn, err := cache.New(cfg.Cache)
	if err != nil {
		log.Fatalf("error: dont instantiate cache client: %v", err)
	}

	if cacheConn.Readiness(context.Background()) != nil {
		log.Fatalf("error: dont connecting to cache: %v", err)
	}

	cacheConn.Delete(context.Background(), merchantName)

	dbFake := newDBfake()
	merchantRepo := newMerchantRepoFake(dbFake)

	cachedMerchantRepo, err := NewMerchant(cacheConn, merchantRepo)
	if err != nil {
		log.Fatalf("error: dont instantiate merchant cached repository: %v", err)
	}

	suite.cacheConn = cacheConn
	suite.cachedMerchantRepo = cachedMerchantRepo
}

func (suite *RedisReposSuite) TearDownSuite() {
	suite.cacheConn.Delete(context.Background(), merchantName)
}

func (suite *RedisReposSuite) MerchantRepositoryFindByNameNotCached() {
	_, err := suite.cacheConn.Get(context.Background(), merchantName)
	assert.EqualError(suite.T(), err, "redis: nil")

	merchantEntity, err := suite.cachedMerchantRepo.FindByName(merchantName)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), merchantEntity)

	_, err = suite.cacheConn.Get(context.Background(), merchantName)
	assert.NoError(suite.T(), err)
}

func (suite *RedisReposSuite) MerchantRepositoryFindByNameCached() {
	_, err := suite.cacheConn.Get(context.Background(), merchantName)
	assert.NoError(suite.T(), err)

	merchantEntity, err := suite.cachedMerchantRepo.FindByName(merchantName)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), merchantEntity)
}

func TestRedisReposSuite(t *testing.T) {
	suite.Run(t, new(RedisReposSuite))
}

func (suite *RedisReposSuite) TestCases() {
	suite.T().Run("TestMerchantRepositoryFindByNameNotCached", func(t *testing.T) {
		suite.MerchantRepositoryFindByNameNotCached()
	})

	suite.T().Run("TestMerchantRepositoryFindByNameCached", func(t *testing.T) {
		suite.MerchantRepositoryFindByNameCached()
	})
}
