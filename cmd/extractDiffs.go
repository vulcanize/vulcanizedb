package cmd

import (
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/streamer"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// extractDiffsCmd represents the extractDiffs command
var extractDiffsCmd = &cobra.Command{
	Use:   "extractDiffs",
	Short: "Extract storage diffs from a node and write them to postgres",
	Long: `Reads storage diffs from either a CSV or JSON RPC subscription.
	Configure which with the STORAGEDIFFS_SOURCE flag. Received diffs are
	written to public.storage_diff.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		extractDiffs()
	},
}

func init() {
	rootCmd.AddCommand(extractDiffsCmd)
}

func extractDiffs() {
	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	// initialize fetcher
	var storageFetcher fetcher.IStorageFetcher
	switch storageDiffsSource {
	case "geth":
		logrus.Debug("fetching storage diffs from geth pub sub")
		rpcClient, _ := getClients()
		stateDiffStreamer := streamer.NewStateDiffStreamer(rpcClient)
		payloadChan := make(chan statediff.Payload)
		storageFetcher = fetcher.NewGethRpcStorageFetcher(&stateDiffStreamer, payloadChan)
	default:
		logrus.Debug("fetching storage diffs from csv")
		tailer := fs.FileTailer{Path: storageDiffsPath}
		storageFetcher = fetcher.NewCsvTailStorageFetcher(tailer)
	}

	// extract diffs
	extractor := storage.NewDiffExtractor(storageFetcher, &db)
	err := extractor.ExtractDiffs()
	if err != nil {
		LogWithCommand.Fatalf("extracting diffs failed: %s", err.Error())
	}
}
