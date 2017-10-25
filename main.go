package main

import (
	"flag"

	"log"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/jmoiron/sqlx"
)

func main() {
	ipcPath := flag.String("ipcPath", "", "location geth.ipc")
	flag.Parse()

	var blockchain core.Blockchain = core.NewGethBlockchain(*ipcPath)
	blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
	pgConfig := "host=localhost port=5432 dbname=vulcanize sslmode=disable"
	db, err := sqlx.Connect("postgres", pgConfig)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	blockchain.RegisterObserver(core.BlockchainDBObserver{Db: db})
	blockchain.SubscribeToEvents()
}
