package main

import (
	"log"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/jmoiron/sqlx"
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {

	p.Task("run", nil, func(context *do.Context) {
		ipcPath := context.Args.MayString("", "ipc-path", "i")

		port := "5432"
		host := "localhost"
		databaseName := "vulcanize"

		var blockchain core.Blockchain = core.NewGethBlockchain(ipcPath)
		blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
		pgConfig := "host=" + host + " port=" + port + " dbname=" + databaseName + " sslmode=disable"
		db, err := sqlx.Connect("postgres", pgConfig)
		if err != nil {
			log.Fatalf("Error connecting to DB: %v\n", err)
		}
		blockchain.RegisterObserver(core.BlockchainDBObserver{Db: db})
		blockchain.SubscribeToEvents()
	})

}

func main() {
	do.Godo(tasks)
}
