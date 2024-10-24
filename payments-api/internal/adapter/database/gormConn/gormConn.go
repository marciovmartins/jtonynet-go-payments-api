package gormConn

import (
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormConn struct {
	DB       *gorm.DB
	strategy string
	driver   string
}

func New(cfg config.Database) (port.DBConn, error) {
	switch cfg.Driver {
	case "postgres":
		strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host,
			cfg.User,
			cfg.Pass,
			cfg.DB,
			cfg.Port,
			cfg.SSLmode)
		db, err := gorm.Open(postgres.Open(strConn), &gorm.Config{})
		if err != nil {
			return GormConn{}, fmt.Errorf("failure on database connection: %w", err)
		}

		/*
			TODO:
			For simplicity, I am using GORM's AutoMigrate. If time permits,
			I should migrate this solution to use the golang-migrate library,
			as it is more robust and scalable.
			See more at: https://github.com/golang-migrate/migrate
		*/
		db.AutoMigrate(&gormModel.Account{})
		db.AutoMigrate(&gormModel.Balance{})
		db.AutoMigrate(&gormModel.Transaction{})
		db.AutoMigrate(&gormModel.MerchantMap{})

		gConn := GormConn{
			DB:       db,
			strategy: cfg.Strategy,
			driver:   cfg.Driver,
		}

		return gConn, nil

	default:
		return GormConn{}, errors.New("database conn driver not suported: " + cfg.Driver)
	}
}

func (gConn GormConn) Readiness() error {
	rawDB, err := gConn.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := rawDB.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}

	return nil
}

func (gConn GormConn) GetDB() interface{} {
	return *gConn.DB
}

func (gConn GormConn) GetStrategy() string {
	return gConn.strategy
}

func (gConn GormConn) GetDriver() string {
	return gConn.driver
}
