package main

import (
	"flag"

	"time"

	"os"
	"text/template"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/blockchain_listener"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/pkg/observers"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

const windowTemplate = `Validating Existing Blocks
|{{.LowerBound}}|-- Validation Window --|{{.UpperBound}}| {{.MaxBlockNumber}}(HEAD)

`

const (
	windowSize      = 24
	pollingInterval = 10 * time.Second
)

func createListener(blockchain *geth.GethBlockchain, repository repositories.Postgres) blockchain_listener.BlockchainListener {
	listener := blockchain_listener.NewBlockchainListener(
		blockchain,
		[]core.BlockchainObserver{
			observers.BlockchainLoggingObserver{},
			observers.NewBlockchainDbObserver(repository),
		},
	)
	return listener
}

func validateBlocks(blockchain *geth.GethBlockchain, repository repositories.Postgres, windowSize int, windowTemplate *template.Template) {
	window := history.UpdateBlocksWindow(blockchain, repository, windowSize)
	repository.SetBlocksStatus(blockchain.LastBlock().Int64())
	windowTemplate.Execute(os.Stdout, window)
}

func main() {
	parsedWindowTemplate := template.Must(template.New("window").Parse(windowTemplate))
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	environment := flag.String("environment", "", "Environment name")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	listner := createListener(blockchain, repository)
	go listner.Start()
	defer listner.Stop()

	missingBlocksPopulated := make(chan int)
	go func() {
		missingBlocksPopulated <- history.PopulateMissingBlocks(blockchain, repository, 0)
	}()

	for range ticker.C {
		validateBlocks(blockchain, repository, windowSize, parsedWindowTemplate)
		select {
		case <-missingBlocksPopulated:
			go func() {
				missingBlocksPopulated <- history.PopulateMissingBlocks(blockchain, repository, 0)
			}()
		default:
		}
	}
}
