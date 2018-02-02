package cmd

import (
	"os"

	"time"

	"log"

	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/utils"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs vulcanizedb with local ethereum node",
	Long: `Syncs vulcanizedb with local ethereum node. 
vulcanizedb sync --startingBlockNumber 0 --config public.toml

Expects ethereum node to be running and requires a .toml config:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/mattkrump/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		sync()
	},
}

const (
	pollingInterval = 7 * time.Second
)

var startingBlockNumber int

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().IntVarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start syncing from")
}

func backFillAllBlocks(blockchain core.Blockchain, repository repositories.Postgres, missingBlocksPopulated chan int, startingBlockNumber int64) {
	go func() {
		missingBlocksPopulated <- history.PopulateMissingBlocks(blockchain, repository, startingBlockNumber)
	}()
}

func sync() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	blockchain := geth.NewBlockchain(ipc)
	if blockchain.LastBlock().Int64() == 0 {
		log.Fatal("geth initial: state sync not finished")
	}
	repository := utils.LoadPostgres(databaseConfig, blockchain.Node())
	validator := history.NewBlockValidator(blockchain, repository, 15)

	missingBlocksPopulated := make(chan int)
	_startingBlockNumber := int64(startingBlockNumber)
	go backFillAllBlocks(blockchain, repository, missingBlocksPopulated, _startingBlockNumber)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateBlocks()
			validator.Log(os.Stdout, window)
		case <-missingBlocksPopulated:
			go backFillAllBlocks(blockchain, repository, missingBlocksPopulated, _startingBlockNumber)
		}
	}
}
