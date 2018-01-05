package main

import (
	"log"

	"flag"

	"math/big"

	"time"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

const (
	windowSize      = 24
	pollingInterval = 10 * time.Second
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	flag.Parse()

	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())

	lastBlockNumber := blockchain.LastBlock().Int64()
	stepSize := int64(1000)

	go func() {
		for i := int64(0); i < lastBlockNumber; i = min(i+stepSize, lastBlockNumber) {
			logs, err := blockchain.GetLogs(core.Contract{Hash: *contractHash}, big.NewInt(i), big.NewInt(i+stepSize))
			log.Println("Backfilling Logs:", i)
			if err != nil {
				log.Println(err)
			}
			repository.CreateLogs(logs)
		}
	}()

	done := make(chan struct{})
	go func() { done <- struct{}{} }()
	for range ticker.C {
		select {
		case <-done:
			go func() {
				z := &big.Int{}
				z.Sub(blockchain.LastBlock(), big.NewInt(25))
				log.Printf("Logs Window: %d - %d", z.Int64(), blockchain.LastBlock().Int64())
				logs, _ := blockchain.GetLogs(core.Contract{Hash: *contractHash}, z, blockchain.LastBlock())
				repository.CreateLogs(logs)
				done <- struct{}{}
			}()
		default:
		}
	}
}
