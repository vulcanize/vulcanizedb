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
)

const windowTemplate = `Validating Existing Blocks
|{{.LowerBound}}|-- Validation Window --|{{.UpperBound}}| {{.MaxBlockNumber}}(HEAD)

`

func main() {
	environment := flag.String("environment", "", "Environment name")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	listener := blockchain_listener.NewBlockchainListener(
		blockchain,
		[]core.BlockchainObserver{
			observers.BlockchainLoggingObserver{},
			observers.NewBlockchainDbObserver(repository),
		},
	)
	go listener.Start()

	windowSize := 10
	ticker := time.NewTicker(10 * time.Second)
	t := template.Must(template.New("window").Parse(windowTemplate))
	for _ = range ticker.C {
		window := history.UpdateBlocksWindow(blockchain, repository, windowSize)
		t.Execute(os.Stdout, window)
	}
	ticker.Stop()
}
