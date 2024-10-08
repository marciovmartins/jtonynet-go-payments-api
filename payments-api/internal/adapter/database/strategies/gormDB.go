package database

import (
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	model "github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

func NewGormDB(cfg config.Database) (GormDB, error) {
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
			return GormDB{}, fmt.Errorf("failure on database connection: %w", err)
		}

		db.AutoMigrate(&model.GormAccountModel{})

		gdb := GormDB{DB: db}

		// var accountRepo port.AccountRepository = repository.NewGormAccountRepository(gdb)
		// fmt.Println("\nolar:\n")
		// fmt.Println(accountRepo)
		// fmt.Println("\n-----------------\n")

		return gdb, nil

	default:
		return GormDB{}, errors.New("database driver not suported: " + cfg.Driver)
	}
}

func (g GormDB) GetDB() interface{} {
	return g.DB
}

func (g GormDB) Readiness() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}

	return nil
}
