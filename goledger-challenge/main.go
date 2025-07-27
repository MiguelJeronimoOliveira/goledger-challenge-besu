package main

import (
	"fmt"
	"log"

	"goledger-challenge/config"
	blockchain "goledger-challenge/contract"
	"goledger-challenge/db"
	"goledger-challenge/handler"
	"goledger-challenge/router"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	fmt.Println("Starting the server...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	cfg := config.LoadConfig()
	blockchainClient, err := blockchain.InitClient(cfg.RPCURL, cfg.ContractAddress, cfg.PrivateKey, cfg.ABIPath)
	if err != nil {
		log.Fatalf("Blockchain init error: %v", err)
	}

	database, err := db.InitPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer database.Close()

	h := &handler.Handler{
		Blockchain: blockchainClient,
		DB:         database,
	}

	r := router.SetupRouter(h)
	r.Run(":8080")

}
