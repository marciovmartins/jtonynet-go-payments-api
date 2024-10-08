package main

import (
	"fmt"
	"log"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
)

func main() {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	if conn.Readiness() != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// var accountRepo port.AccountRepository = repository.NewGormAccountRepository(conn)
	// var accountRepo port.AccountRepository = repository.NewGormAccountRepository(conn) // its OK on gorm strategy
	// fmt.Println(accountRepo)

	fmt.Println("Successfully Readiness database!")
}
