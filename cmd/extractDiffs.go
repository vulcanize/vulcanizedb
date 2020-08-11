package cmd

import (
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/streamer"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func getContractAddresses() []string {
	LogWithCommand.Info("Getting contract addresses from config file")
	contracts := viper.GetStringMap("contract")
	var addresses []string
	for contractName := range contracts {
		address := viper.GetStringMapString("contract." + contractName)["address"]
		addresses = append(addresses, address)
	}
	return addresses
}

func extractDiffs() {
	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	addressesToWatch := getContractAddresses()

	healthCheckFile := "/tmp/connection"
	msg := []byte("geth storage fetcher connection established\n")
	gethStatusWriter := fs.NewStatusWriter(healthCheckFile, msg)

	// initialize fetcher
	var storageFetcher fetcher.IStorageFetcher
	logrus.Debug("fetching storage diffs from geth")
	switch storageDiffsSource {
	case "geth":
		logrus.Info("Using new geth patch with filters event system")
		_, ethClient := getClients()
		filterQuery := createFilterQuery(addressesToWatch)
		stateDiffStreamer := streamer.NewEthStateChangeStreamer(ethClient, filterQuery)
		payloadChan := make(chan filters.Payload)
		storageFetcher = fetcher.NewGethRpcStorageFetcher(&stateDiffStreamer, payloadChan, gethStatusWriter)
	default:
		logrus.Debug("fetching storage diffs from csv")
		tailer := fs.FileTailer{Path: storageDiffsPath}
		msg := []byte("csv tail storage fetcher connection established\n")
		statusWriter := fs.NewStatusWriter(healthCheckFile, msg)

		storageFetcher = fetcher.NewCsvTailStorageFetcher(tailer, statusWriter)
	}

	// extract diffs
	extractor := storage.NewDiffExtractor(storageFetcher, &db)
	err := extractor.ExtractDiffs()
	if err != nil {
		LogWithCommand.Fatalf("extracting diffs failed: %s", err.Error())
	}
}

func createFilterQuery(watchedAddresses []string) ethereum.FilterQuery {
	logrus.Infof("Creating a filter query for %d watched addresses", len(watchedAddresses))
	addressesToLog := strings.Join(watchedAddresses[:], ", ")
	logrus.Infof("Watched addresses: %s", addressesToLog)

	var addresses []common.Address
	for _, addressString := range watchedAddresses {
		addresses = append(addresses, common.HexToAddress(addressString))
	}

	return ethereum.FilterQuery{
		Addresses: addresses,
	}
}
